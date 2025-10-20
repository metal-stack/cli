package v1

import (
	"fmt"

	v1 "github.com/metal-stack/api/go/metalstack/api/v2"
	clitypes "github.com/metal-stack/metal-lib/pkg/commands/types"
	"github.com/spf13/cobra"
)

func newHealthCmd(c *clitypes.Config) *cobra.Command {
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "print the client and server health information",
		Long:  "print the client and server health information",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewRequestContext()
			defer cancel()

			resp, err := c.Client.Apiv2().Health().Get(ctx, &v1.HealthServiceGetRequest{})
			if err != nil {
				return fmt.Errorf("failed to get health: %w", err)
			}

			return c.ListPrinter.Print(resp.Health)
		},
	}

	return healthCmd
}
