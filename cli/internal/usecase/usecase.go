package usecase

import (
	"context"

	domain "github.com/Friend-zva/dns-application/cli/internal/domain"
)

type DNSServerRepo interface {
	AddDNSServer(ctx context.Context, name domain.NameServer) (string, error)
	DeleteDNSServer(ctx context.Context, name domain.NameServer) (string, error)
	GetDNSServers(ctx context.Context) ([]domain.NameServer, error)
}

type dnsServerUseCase struct {
	repo DNSServerRepo
}

func NewDNSServerUseCase(repo DNSServerRepo) *dnsServerUseCase {
	return &dnsServerUseCase{
		repo: repo,
	}
}

func (u *dnsServerUseCase) AddDNSServer(ctx context.Context, name string) (string, error) {
	reply, err := u.repo.AddDNSServer(ctx, name)
	if err != nil {
		return "", err
	}

	return reply, nil
}

func (u *dnsServerUseCase) DeleteDNSServer(ctx context.Context, name string) (string, error) {
	reply, err := u.repo.DeleteDNSServer(ctx, name)
	if err != nil {
		return "", err
	}

	return reply, nil
}

func (u *dnsServerUseCase) GetDNSServers(ctx context.Context) ([]domain.NameServer, error) {
	servers, err := u.repo.GetDNSServers(ctx)
	if err != nil {
		return []domain.NameServer{}, err
	}

	return servers, nil
}
