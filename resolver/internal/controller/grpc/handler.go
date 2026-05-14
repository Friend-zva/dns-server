package grpc

import (
	"context"
	"log/slog"

	resolverpb "github.com/Friend-zva/dns-application/proto/resolver"
	domain "github.com/Friend-zva/dns-application/resolver/internal/domain"
)

type DNSServerUseCase interface {
	AddDNSServer(name string) (string, error)
	DeleteDNSServer(name string) (string, error)
	GetDNSServers() ([]domain.DNSServer, error)
}

type handler struct {
	resolverpb.UnimplementedResolverServer
	usecase DNSServerUseCase
}

func NewHandler(log *slog.Logger, usecase DNSServerUseCase) *handler {
	return &handler{
		usecase: usecase,
	}
}

func (h *handler) AddDNSServer(ctx context.Context, req *resolverpb.AddDNSServerRequest) (*resolverpb.AddDNSServerResponse, error) {
	reply, err := h.usecase.AddDNSServer(req.Name)
	if err != nil {
		return nil, err
	}

	return &resolverpb.AddDNSServerResponse{
		Reply: reply,
	}, nil
}

func (h *handler) DeleteDNSServer(ctx context.Context, req *resolverpb.DeleteDNSServerRequest) (*resolverpb.DeleteDNSServerResponse, error) {
	reply, err := h.usecase.DeleteDNSServer(req.Name)
	if err != nil {
		return nil, err
	}

	return &resolverpb.DeleteDNSServerResponse{
		Reply: reply,
	}, nil
}

func (h *handler) GetDNSServers(ctx context.Context, req *resolverpb.GetDNSServersRequest) (*resolverpb.GetDNSServersResponse, error) {
	servers, err := h.usecase.GetDNSServers()
	if err != nil {
		return nil, err
	}

	names := make([]string, len(servers))
	for i, server := range servers {
		names[i] = server.Name
	}

	return &resolverpb.GetDNSServersResponse{
		Names: names,
	}, nil
}
