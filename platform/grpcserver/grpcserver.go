package grpcserver

import (
	"context"
	"net"

	grpc "google.golang.org/grpc"
)

type Server struct {
	srv *grpc.Server
	lis net.Listener
}

func New(address string, opts ...grpc.ServerOption) (*Server, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Server{
		srv: grpc.NewServer(opts...),
		lis: lis,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.srv.GracefulStop()
	}()
	return s.srv.Serve(s.lis)
}

func (s *Server) GRPC() *grpc.Server {
	return s.srv
}
