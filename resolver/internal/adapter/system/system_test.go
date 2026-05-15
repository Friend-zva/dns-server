package system

import (
	"log/slog"
	"os"
	"reflect"
	"testing"

	domain "github.com/Friend-zva/dns-application/resolver/internal/domain"
)

func TestOSDNSServerRepo(t *testing.T) {
	fileTmp, err := os.CreateTemp("", "resolv.conf.*")
	if err != nil {
		t.Fatalf("cannot create temp file: %v", err)
	}
	defer func() {
		err = os.Remove(fileTmp.Name())
		if err != nil {
			t.Errorf("cannot remove temp file: %v", err)
		}
	}()

	pathResolvOld := pathResolv
	pathResolv = fileTmp.Name()
	defer func() { pathResolv = pathResolvOld }()

	initialContent := []byte("# This is a comment\nnameserver 8.8.8.8\n")
	if err := os.WriteFile(pathResolv, initialContent, 0644); err != nil {
		t.Fatalf("cannot write initial content: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewOSDNSServerRepo(logger)

	t.Run("GetDNSServers", func(t *testing.T) {
		servers, err := repo.GetDNSServers()
		if err != nil {
			t.Fatalf("unexpected error, got %v", err)
		}

		expected := []domain.DNSServer{{Name: "8.8.8.8"}}
		if !reflect.DeepEqual(servers, expected) {
			t.Errorf("expected %v, got %v", expected, servers)
		}
	})

	t.Run("AddDNSServer", func(t *testing.T) {
		server := domain.DNSServer{Name: "1.1.1.1"}
		_, err := repo.AddDNSServer(server)
		if err != nil {
			t.Fatalf("unexpected error, got %v", err)
		}

		servers, err := repo.GetDNSServers()
		if err != nil {
			t.Fatalf("unexpected error, got %v", err)
		}

		expected := []domain.DNSServer{{Name: "8.8.8.8"}, {Name: "1.1.1.1"}}
		if !reflect.DeepEqual(servers, expected) {
			t.Errorf("expected %v, got %v", expected, servers)
		}
	})

	t.Run("AddDNSServer_Duplicate", func(t *testing.T) {
		server := domain.DNSServer{Name: "1.1.1.1"}
		_, err := repo.AddDNSServer(server)
		if err != nil {
			t.Fatalf("unexpected error, got %v", err)
		}

		servers, err := repo.GetDNSServers()
		if err != nil {
			t.Fatalf("unexpected error, got %v", err)
		}

		expected := []domain.DNSServer{{Name: "8.8.8.8"}, {Name: "1.1.1.1"}}
		if !reflect.DeepEqual(servers, expected) {
			t.Errorf("expected %v, got %v", expected, servers)
		}
	})

	t.Run("DeleteDNSServer", func(t *testing.T) {
		server := domain.DNSServer{Name: "8.8.8.8"}
		_, err := repo.DeleteDNSServer(server)
		if err != nil {
			t.Fatalf("unexpected error, got %v", err)
		}

		servers, err := repo.GetDNSServers()
		if err != nil {
			t.Fatalf("unexpected error, got %v", err)
		}

		expected := []domain.DNSServer{{Name: "1.1.1.1"}}
		if !reflect.DeepEqual(servers, expected) {
			t.Errorf("expected %v, got %v", expected, servers)
		}
	})
}
