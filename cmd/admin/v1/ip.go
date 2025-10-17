package v1

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type ip struct {
	c *config.Config
}

func newIPCmd(c *config.Config) *cobra.Command {
	w := &ip{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.IP]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[any, any, *apiv2.IP](w).WithFS(c.Fs),
		Singular:        "ip",
		Plural:          "ips",
		Description:     "an ip address of metal-stack.io",
		Sorter:          sorters.IPSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("ip", "", "ipaddress to filter [optional]")
			cmd.Flags().String("name", "", "name to filter [optional]")
			cmd.Flags().String("network", "", "network to filter [optional]")
			cmd.Flags().String("description", "", "description to filter [optional]")
			genericcli.Must(cmd.RegisterFlagCompletionFunc("network", c.Completion.NetworkListCompletion))
		},
		ValidArgsFn: c.Completion.IpListCompletion,
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *ip) List() ([]*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().IP().List(ctx, &adminv2.IPServiceListRequest{
		Query: &apiv2.IPQuery{},
	})
	if err != nil {
		return nil, err
	}

	return resp.Ips, nil
}

func (t *ip) Get(id string) (*apiv2.IP, error) {
	panic("unimplemented")
}
func (c *ip) Delete(id string) (*apiv2.IP, error) {
	panic("unimplemented")
}
func (t *ip) Create(rq any) (*apiv2.IP, error) {
	panic("unimplemented")
}
func (t *ip) Convert(r *apiv2.IP) (string, any, any, error) {
	panic("unimplemented")
}

func (t *ip) Update(rq any) (*apiv2.IP, error) {
	panic("unimplemented")
}
