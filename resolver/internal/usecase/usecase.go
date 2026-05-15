package usecase

import (
	"fmt"
	"net"
	"strings"

	domain "github.com/Friend-zva/dns-application/resolver/internal/domain"
	apperror "github.com/Friend-zva/dns-application/resolver/platform/apperror"
)

type DNSServerRepo interface {
	AddDNSServer(domain.DNSServer) (string, error)
	DeleteDNSServer(domain.DNSServer) (string, error)
	GetDNSServers() ([]domain.DNSServer, error)
}

type dnsServerUseCase struct {
	repo DNSServerRepo
}

func NewDNSServerUseCase(repo DNSServerRepo) *dnsServerUseCase {
	return &dnsServerUseCase{
		repo: repo,
	}
}

func (u *dnsServerUseCase) AddDNSServer(name string) (string, error) {
	if name == "" {
		err := fmt.Errorf("empty name server")
		return "", apperror.ErrInvArgs.Wrap(err)
	}

	if net.ParseIP(name) == nil {
		err := fmt.Errorf("invalid IP address")
		return "", apperror.ErrInvArgs.Wrap(err)
	}

	server := domain.DNSServer{Name: name}

	reply, err := u.repo.AddDNSServer(server)
	if err != nil {
		return "", err
	}

	return reply, nil
}

func (u *dnsServerUseCase) DeleteDNSServer(name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		err := fmt.Errorf("empty name server")
		return "", apperror.ErrInvArgs.Wrap(err)
	}

	if net.ParseIP(name) == nil {
		err := fmt.Errorf("invalid IP address")
		return "", apperror.ErrInvArgs.Wrap(err)
	}

	server := domain.DNSServer{Name: name}

	reply, err := u.repo.DeleteDNSServer(server)
	if err != nil {
		return "", err
	}

	return reply, nil
}

func (u *dnsServerUseCase) GetDNSServers() ([]domain.DNSServer, error) {
	servers, err := u.repo.GetDNSServers()
	if err != nil {
		return []domain.DNSServer{}, err
	}

	return servers, nil
}
