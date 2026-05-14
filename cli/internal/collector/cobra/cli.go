package cobracli

import (
	"context"
	"fmt"

	cobra "github.com/spf13/cobra"

	domain "github.com/Friend-zva/dns-application/cli/internal/domain"
)

type DNSServerUseCase interface {
	AddDNSServer(ctx context.Context, name string) (string, error)
	DeleteDNSServer(ctx context.Context, name string) (string, error)
	GetDNSServers(ctx context.Context) ([]domain.NameServer, error)
}

type handler struct {
	cmdRoot *cobra.Command
	usecase DNSServerUseCase
}

func NewHandler(usecase DNSServerUseCase) *handler {
	h := &handler{
		usecase: usecase,
		cmdRoot: &cobra.Command{
			Use:           "dns-app",
			Short:         "CLI tool for managing DNS servers",
			SilenceErrors: true,
		},
	}

	h.setupCommands()
	return h
}

func (h *handler) Execute() error {
	h.cmdRoot.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd: true,
	}

	return h.cmdRoot.Execute()
}

func (h *handler) setupCommands() {
	cmdAdd := &cobra.Command{
		Use:   "add [ip]",
		Short: "Add a DNS server",
		Args:  cobra.ExactArgs(1),
		RunE:  h.addServer,
	}

	cmdDelete := &cobra.Command{
		Use:   "delete [ip]",
		Short: "Delete a DNS server",
		Args:  cobra.ExactArgs(1),
		RunE:  h.deleteServer,
	}

	cmdGet := &cobra.Command{
		Use:   "get",
		Short: "List all DNS servers",
		Args:  cobra.NoArgs,
		RunE:  h.getServers,
	}

	h.cmdRoot.AddCommand(cmdAdd)
	h.cmdRoot.AddCommand(cmdDelete)
	h.cmdRoot.AddCommand(cmdGet)
}

func (h *handler) addServer(cmd *cobra.Command, args []string) error {
	ip := args[0]

	msg, err := h.usecase.AddDNSServer(context.Background(), ip)
	if err != nil {
		return fmt.Errorf("cannot add server: %w", err)
	}

	fmt.Println(msg)
	return nil
}

func (h *handler) deleteServer(cmd *cobra.Command, args []string) error {
	ip := args[0]

	msg, err := h.usecase.DeleteDNSServer(context.Background(), ip)
	if err != nil {
		return fmt.Errorf("cannot delete server: %w", err)
	}

	fmt.Println(msg)
	return nil
}

func (h *handler) getServers(cmd *cobra.Command, args []string) error {
	servers, err := h.usecase.GetDNSServers(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get servers: %w", err)
	}

	if len(servers) == 0 {
		fmt.Println("No DNS servers found.")
		return nil
	}

	fmt.Println("DNS Servers:")
	for _, s := range servers {
		fmt.Println("-", s)
	}

	return nil
}
