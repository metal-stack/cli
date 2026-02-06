package v2

import (
	"fmt"
	"time"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/durationpb"
)

type vpn struct {
	c *config.Config
}

func newVPNCmd(c *config.Config) *cobra.Command {
	w := &vpn{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.VPNNode]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "vpn",
		Plural:          "vpn",
		Description:     "manage vpn keys and list nodes connected",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("project", "", "the project for which vpn nodes should be listed")
		},
		OnlyCmds:    genericcli.OnlyCmds(genericcli.ListCmd),
		ValidArgsFn: w.c.Completion.ProjectListCompletion,
	}

	authKeyCmd := &cobra.Command{
		Use:   "authkey",
		Short: "generate a authkey to connect to the vpn",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.authKey()
		},
		ValidArgsFunction: c.Completion.ProjectListCompletion,
	}
	authKeyCmd.Flags().String("project", "", "the project for which the authkey should be generated")
	authKeyCmd.Flags().Bool("ephemeral", true, "ephemeral defines if the key can only be used once")
	authKeyCmd.Flags().Duration("expires", 1*time.Hour, "the duration after the generated key is not valid anymore")
	genericcli.Must(authKeyCmd.MarkFlagRequired("project"))

	return genericcli.NewCmds(cmdsConfig, authKeyCmd)
}

func (v *vpn) authKey() error {
	ctx, cancel := v.c.NewRequestContext()
	defer cancel()

	req := &adminv2.VPNServiceAuthKeyRequest{
		Project:   viper.GetString("project"),
		Ephemeral: viper.GetBool("ephemeral"),
		Expires:   durationpb.New(viper.GetDuration("expires")),
	}

	resp, err := v.c.Client.Adminv2().VPN().AuthKey(ctx, req)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(v.c.Out, "authkey: %s ephemeral:%t created at:%s expires at:%s\n", resp.AuthKey, resp.Ephemeral, resp.CreatedAt, resp.ExpiresAt)
	_, _ = fmt.Fprintf(v.c.Out, "vpn endpoint: %s\n", resp.Address)

	return nil
}

func (v *vpn) Get(id string) (*apiv2.VPNNode, error) {
	panic("unimplemented")
}

func (v *vpn) List() ([]*apiv2.VPNNode, error) {
	ctx, cancel := v.c.NewRequestContext()
	defer cancel()

	req := &adminv2.VPNServiceListNodesRequest{}

	if viper.IsSet("project") {
		req.Project = pointer.Pointer(viper.GetString("project"))
	}

	resp, err := v.c.Client.Adminv2().VPN().ListNodes(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list vpn nodes: %w", err)
	}

	return resp.Nodes, nil
}

func (v *vpn) Create(rq any) (*apiv2.VPNNode, error) {
	panic("unimplemented")
}

func (v *vpn) Delete(id string) (*apiv2.VPNNode, error) {
	panic("unimplemented")
}

func (v *vpn) Convert(r *apiv2.VPNNode) (string, any, any, error) {
	panic("unimplemented")
}

func (v *vpn) Update(rq any) (*apiv2.VPNNode, error) {
	panic("unimplemented")
}
