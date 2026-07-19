package v2

import (
	"fmt"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type adminImageUsage struct {
	c *config.Config
}

func newAdminImageUsageCmd(c *config.Config) *cobra.Command {
	w := &adminImageUsage{
		c: c,
	}

	usageCmd := &cobra.Command{
		Use:   "image-usage",
		Short: "show image usage information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.imageUsage()
		},
	}

	usageCmd.Flags().String("id", "", "image id to query usage for")

	return usageCmd
}

func (c *adminImageUsage) imageUsage() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.ImageServiceUsageRequest{
		Query: &apiv2.ImageQuery{},
	}

	if viper.IsSet("id") {
		req.Query.Id = pointer.PointerOrNil(viper.GetString("id"))
	}

	resp, err := c.c.Client.Adminv2().Image().Usage(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get image usage: %w", err)
	}

	return c.c.ListPrinter.Print(resp.ImageUsage)
}
