package v1

import (
	"fmt"

	v1 "github.com/metal-stack/api/go/metalstack/api/v2"
	clitypes "github.com/metal-stack/metal-lib/pkg/commands/types"
	"github.com/metal-stack/v"
	"github.com/spf13/cobra"
)

type version struct {
	Client string
	Server *v1.Version
}

func newVersionCmd(c *clitypes.Config) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "print the client and server version information",
		Long:  "print the client and server version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewRequestContext()
			defer cancel()

			v := version{
				Client: v.V.String(),
			}

			resp, err := c.Client.Apiv2().Version().Get(ctx, &v1.VersionServiceGetRequest{})
			if err == nil {
				v.Server = resp.Version
			}

			if err := c.DescribePrinter.Print(v); err != nil {
				return err
			}

			if err != nil {
				return fmt.Errorf("failed to get server info: %w", err)
			}

			return nil
		},
	}

	return versionCmd
}
