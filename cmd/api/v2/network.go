package v2

import (
	"fmt"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type network struct {
	c *config.Config
}

func newNetworkCmd(c *config.Config) *cobra.Command {
	w := &network{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.Network]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "network",
		Plural:          "networks",
		Description:     "read project networks",
		Sorter:          sorters.NetworkSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.NetworkListCompletion,
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project to list networks for")
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *network) Get(id string) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.NetworkServiceGetRequest{
		Id:      id,
		Project: c.c.GetProject(),
	}

	resp, err := c.c.Client.Apiv2().Network().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get network: %w", err)
	}

	return resp.Network, nil
}

func (c *network) List() ([]*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.NetworkServiceListRequest{
		Project: c.c.GetProject(),
	}

	resp, err := c.c.Client.Apiv2().Network().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	return resp.Networks, nil
}

func (c *network) Create(rq any) (*apiv2.Network, error) {
	panic("unimplemented")
}

func (c *network) Delete(id string) (*apiv2.Network, error) {
	panic("unimplemented")
}

func (c *network) Convert(r *apiv2.Network) (string, any, any, error) {
	panic("unimplemented")
}

func (c *network) Update(rq any) (*apiv2.Network, error) {
	panic("unimplemented")
}
