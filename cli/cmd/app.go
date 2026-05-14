package main

import (
	"fmt"
	"os"

	resolver "github.com/Friend-zva/dns-application/cli/internal/adapter/resolver"
	cobraCLI "github.com/Friend-zva/dns-application/cli/internal/collector/cobra"
	usecase "github.com/Friend-zva/dns-application/cli/internal/usecase"
)

func run() error {
	clientResolver, err := resolver.NewClient(":8080")
	if err != nil {
		return fmt.Errorf("cannot init resolver client: %w", err)
	}
	defer func() {
		if err := clientResolver.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: cannot close resolver connection: %v\n", err)
		}
	}()

	usecase := usecase.NewDNSServerUseCase(clientResolver)

	handlerCLI := cobraCLI.NewHandler(usecase)
	if err := handlerCLI.Execute(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
