package system

import (
	"bufio"
	"bytes"
	"log/slog"
	"os"
	"strings"

	domain "github.com/Friend-zva/dns-application/resolver/internal/domain"
	apperror "github.com/Friend-zva/dns-application/resolver/platform/apperror"
)

var pathResolv = "/etc/resolv.conf"

type osDNSServerRepo struct {
	log *slog.Logger
}

func NewOSDNSServerRepo(log *slog.Logger) *osDNSServerRepo {
	return &osDNSServerRepo{
		log: log,
	}
}

func (r *osDNSServerRepo) AddDNSServer(server domain.DNSServer) (string, error) {
	file, err := os.OpenFile(pathResolv, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return "", apperror.ErrInternal.Wrap(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			r.log.Warn("cannot close file", "error", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	serverExists := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		tokens := strings.Fields(line)
		if len(tokens) >= 2 {
			name := tokens[1]
			if name == server.Name {
				serverExists = true
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", apperror.ErrInternal.Wrap(err)
	}

	if !serverExists {
		_, err = file.WriteString("nameserver " + server.Name + "\n")
		if err != nil {
			return "", apperror.ErrInternal.Wrap(err)
		}
		r.log.Info("server added", "name", server.Name)
		return "Added", nil
	}

	r.log.Warn("server already exists", "name", server.Name)
	return "Already exists", nil

}

func (r *osDNSServerRepo) DeleteDNSServer(server domain.DNSServer) (string, error) {
	data, err := os.ReadFile(pathResolv)
	if err != nil {
		return "", apperror.ErrInternal.Wrap(err)
	}

	lines := bytes.Split(data, []byte("\n"))
	var output [][]byte
	serverExists := false

	for _, line := range lines {
		tokens := bytes.Fields(line)

		if len(tokens) == 2 &&
			bytes.Equal(tokens[0], []byte("nameserver")) &&
			bytes.Equal(tokens[1], []byte(server.Name)) {
			serverExists = true
			continue
		}

		output = append(output, line)
	}

	result := bytes.Join(output, []byte("\n"))
	err = os.WriteFile(pathResolv, result, 0644)
	if err != nil {
		return "", apperror.ErrInternal.Wrap(err)
	}

	if serverExists {
		r.log.Info("server deleted", "name", server.Name)
		return "Deleted", nil
	}
	r.log.Info("server not found", "name", server.Name)
	return "Not found", nil
}

func (r *osDNSServerRepo) GetDNSServers() ([]domain.DNSServer, error) {
	data, err := os.ReadFile(pathResolv)
	if err != nil {
		return []domain.DNSServer{}, apperror.ErrInternal.Wrap(err)
	}

	lines := bytes.Split(data, []byte("\n"))
	servers := []domain.DNSServer{}

	for _, text := range lines {
		line := bytes.TrimSpace(text)
		if len(line) == 0 || bytes.HasPrefix(line, []byte("#")) {
			continue
		}
		tokens := bytes.Fields(line)
		if len(tokens) >= 2 && string(tokens[0]) == "nameserver" {
			name := string(tokens[1])
			servers = append(servers, domain.DNSServer{Name: name})
		}
	}

	r.log.Info("get servers", "count", len(servers))
	return servers, nil
}
