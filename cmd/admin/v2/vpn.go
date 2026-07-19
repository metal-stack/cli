package v2

import (
	"fmt"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/durationpb"
)

type adminVPN struct {
	c *config.Config
}

func newAdminVPNCmd(c *config.Config) *cobra.Command {
	w := &adminVPN{
		c: c,
	}

	adminVPNCmd := &cobra.Command{
		Use:   "vpn",
		Short: "manage VPN",
	}

	authKeyCmd := &cobra.Command{
		Use:   "auth-key",
		Short: "generate a VPN authentication key for a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.authKey()
		},
	}

	authKeyCmd.Flags().String("project", "", "project for which to generate the VPN auth key")
	genericcli.Must(authKeyCmd.MarkFlagRequired("project"))
	authKeyCmd.Flags().Bool("ephemeral", false, "whether the auth key should be ephemeral")
	authKeyCmd.Flags().Duration("expires", 0, "duration after which the auth key expires")
	authKeyCmd.Flags().String("reason", "", "reason for requesting VPN access")

	listNodesCmd := &cobra.Command{
		Use:   "list-nodes",
		Short: "list VPN connected machines",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.listNodes()
		},
	}

	adminVPNCmd.AddCommand(authKeyCmd, listNodesCmd)

	return adminVPNCmd
}

func (c *adminVPN) authKey() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.VPNServiceAuthKeyRequest{
		Project:   viper.GetString("project"),
		Ephemeral: viper.GetBool("ephemeral"),
		Reason:    viper.GetString("reason"),
	}

	if viper.IsSet("expires") {
		req.Expires = durationpb.New(viper.GetDuration("expires"))
	}

	resp, err := c.c.Client.Adminv2().VPN().AuthKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to generate VPN auth key: %w", err)
	}

	return c.c.DescribePrinter.Print(resp)
}

func (c *adminVPN) listNodes() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().VPN().ListNodes(ctx, &adminv2.VPNServiceListNodesRequest{})
	if err != nil {
		return fmt.Errorf("failed to list VPN nodes: %w", err)
	}

	return c.c.ListPrinter.Print(resp.Nodes)
}
