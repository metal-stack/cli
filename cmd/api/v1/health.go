package v1

import (
	"fmt"

	"connectrpc.com/connect"
	v1 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/spf13/cobra"
)

func newHealthCmd(c *config.Config) *cobra.Command {
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "print the client and server health information",
		Long:  "print the client and server health information",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewRequestContext()
			defer cancel()

			resp, err := c.Client.Apiv2().Health().Get(ctx, connect.NewRequest(&v1.HealthServiceGetRequest{}))
			if err != nil {
				return fmt.Errorf("failed to get health: %w", err)
			}

			return c.ListPrinter.Print(resp.Msg.Health)
		},
	}

	return healthCmd
}
