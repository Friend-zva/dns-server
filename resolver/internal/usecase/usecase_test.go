package usecase

import (
	"reflect"
	"testing"

	domain "github.com/Friend-zva/dns-application/resolver/internal/domain"
)

type mockRepo struct {
	servers []domain.DNSServer
}

func (m *mockRepo) AddDNSServer(s domain.DNSServer) (string, error) {
	for _, server := range m.servers {
		if server.Name == s.Name {
			return "Already exists", nil
		}
	}
	m.servers = append(m.servers, s)
	return "Added", nil
}

func (m *mockRepo) DeleteDNSServer(s domain.DNSServer) (string, error) {
	var result []domain.DNSServer
	reply := "Not found"

	for _, server := range m.servers {
		if server.Name != s.Name {
			result = append(result, server)
		} else {
			reply = "Deleted"
		}
	}

	m.servers = result
	return reply, nil
}

func (m *mockRepo) GetDNSServers() ([]domain.DNSServer, error) {
	return m.servers, nil
}

func TestDNSServerUseCase(t *testing.T) {
	repo := &mockRepo{
		servers: []domain.DNSServer{{Name: "8.8.8.8"}},
	}
	usecase := NewDNSServerUseCase(repo)

	t.Run("GetDNSServers", func(t *testing.T) {
		servers, err := usecase.GetDNSServers()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := []domain.DNSServer{{Name: "8.8.8.8"}}
		if !reflect.DeepEqual(servers, expected) {
			t.Errorf("expected %v, got %v", expected, servers)
		}
	})

	t.Run("AddDNSServer", func(t *testing.T) {
		_, err := usecase.AddDNSServer("1.1.1.1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		servers, _ := usecase.GetDNSServers()
		expected := []domain.DNSServer{{Name: "8.8.8.8"}, {Name: "1.1.1.1"}}
		if !reflect.DeepEqual(servers, expected) {
			t.Errorf("expected %v, got %v", expected, servers)
		}
	})

	t.Run("DeleteDNSServer", func(t *testing.T) {
		_, err := usecase.DeleteDNSServer("8.8.8.8")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		servers, _ := usecase.GetDNSServers()
		expected := []domain.DNSServer{{Name: "1.1.1.1"}}
		if !reflect.DeepEqual(servers, expected) {
			t.Errorf("expected %v, got %v", expected, servers)
		}
	})
}
