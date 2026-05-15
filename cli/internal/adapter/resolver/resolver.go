package resolver

import (
	"context"

	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"

	domain "github.com/Friend-zva/dns-application/cli/internal/domain"
	resolverpb "github.com/Friend-zva/dns-application/proto/resolver"
	apperror "github.com/Friend-zva/dns-application/resolver/platform/apperror"
)

type client struct {
	conn *grpc.ClientConn
	pb   resolverpb.ResolverClient
}

func NewClient(address string) (*client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, apperror.ErrExternal.Wrap(err)
	}

	return &client{
		conn: conn,
		pb:   resolverpb.NewResolverClient(conn),
	}, nil
}

func (c *client) AddDNSServer(ctx context.Context, name domain.NameServer) (string, error) {
	resp, err := c.pb.AddDNSServer(ctx, &resolverpb.AddDNSServerRequest{Name: name})
	if err != nil {
		return "", err
	}

	return resp.Reply, nil
}

func (c *client) DeleteDNSServer(ctx context.Context, name domain.NameServer) (string, error) {
	resp, err := c.pb.DeleteDNSServer(ctx, &resolverpb.DeleteDNSServerRequest{Name: name})
	if err != nil {
		return "", err
	}

	return resp.Reply, nil
}

func (c *client) GetDNSServers(ctx context.Context) ([]domain.NameServer, error) {
	resp, err := c.pb.GetDNSServers(ctx, &resolverpb.GetDNSServersRequest{})
	if err != nil {
		return []domain.NameServer{}, err
	}

	return resp.Names, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
