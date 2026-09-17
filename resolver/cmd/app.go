package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	grpcserver "github.com/Friend-zva/dns-application/platform/grpcserver"
	logger "github.com/Friend-zva/dns-application/platform/logger"
	resolverpb "github.com/Friend-zva/dns-application/proto/resolver"
	system "github.com/Friend-zva/dns-application/resolver/internal/adapter/system"
	grpcH "github.com/Friend-zva/dns-application/resolver/internal/controller/grpc"
	usecase "github.com/Friend-zva/dns-application/resolver/internal/usecase"
)

func run(ctx context.Context) error {
	log := logger.MustMakeLogger("INFO")
	log.Debug("debug messages are enabled")

	log.Info("starting resolver server...")

	osRepo := system.NewOSDNSServerRepo(log)
	usecase := usecase.NewDNSServerUseCase(osRepo)
	handler := grpcH.NewHandler(log, usecase)

	resolver, err := grpcserver.New(":8080")
	if err != nil {
		return fmt.Errorf("cannot create grpc resolver: %w", err)
	}

	resolverpb.RegisterResolverServer(resolver.GRPC(), handler)

	if err := resolver.Run(ctx); err != nil {
		return fmt.Errorf("cannot run grpc resolver: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("cannot launch resolver server: %v\n", err)
		}
		cancel()
		os.Exit(1)
	}

	cancel()
}
