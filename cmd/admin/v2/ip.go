package v2

import (
	"fmt"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/cli/pkg/helpers"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "ip",
		Plural:          "ips",
		Description:     "an ip address of metal-stack.io",
		Sorter:          sorters.IPSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("ip", "", "ip which should be listed")
			cmd.Flags().String("uuid", "", "allocation uuid of ip which should be listed")
			cmd.Flags().String("project", "", "project from where ips should be listed")
			cmd.Flags().String("name", "", "name from ips which should be listed")
			cmd.Flags().String("network", "", "network from where ips should be listed")
			cmd.Flags().String("machine", "", "machine where ips are attached to")
			cmd.Flags().StringSlice("labels", nil, "lists only ips with the given labels")
			cmd.Flags().String("addressfamily", "", "addressfamily of ips which should be listed")
			cmd.Flags().String("type", "", "type of ips which should be listed")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("network", c.Completion.Network))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("machine", c.Completion.Machine))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("addressfamily", c.Completion.AddressFamily))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("type", c.Completion.IPType))

		},
		ValidArgsFn: c.Completion.Ip,
		OnlyCmds: map[genericcli.DefaultCmd]bool{
			genericcli.ListCmd:     true,
			genericcli.DescribeCmd: true,
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *ip) Create(_ any) (*apiv2.IP, error) {
	panic("unimplemented")
}

func (c *ip) Delete(id string) (*apiv2.IP, error) {
	panic("unimplemented")
}

func (c *ip) Get(id string) (*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().IP().List(ctx, &adminv2.IPServiceListRequest{
		Query: &apiv2.IPQuery{
			Ip: &id,
		},
	})
	if err != nil {
		return nil, err
	}
	switch len(resp.Ips) {
	case 0:
		return nil, fmt.Errorf("no ip found for ip:%s", id)
	case 1:
		return resp.Ips[0], nil
	default:
		return nil, fmt.Errorf("more than one ip found for ip:%s", id)
	}
}

func (c *ip) List() ([]*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	var labels *apiv2.Labels
	if labelSlice := viper.GetStringSlice("labels"); len(labelSlice) > 0 {
		var err error

		labels, err = helpers.LabelsFromSlice(labelSlice)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.c.Client.Adminv2().IP().List(ctx, &adminv2.IPServiceListRequest{
		Query: &apiv2.IPQuery{
			Ip:               pointer.PointerOrNil(viper.GetString("ip")),
			Uuid:             pointer.PointerOrNil(viper.GetString("uuid")),
			Network:          pointer.PointerOrNil(viper.GetString("network")),
			Project:          pointer.PointerOrNil(viper.GetString("project")),
			Name:             pointer.PointerOrNil(viper.GetString("name")),
			Machine:          pointer.PointerOrNil(viper.GetString("machine")),
			ParentPrefixCidr: pointer.PointerOrNil(viper.GetString("parent-prefix")),
			Labels:           labels,
			Type:             helpers.IPTypeToType(viper.GetString("type")),
			AddressFamily:    helpers.IPAddressFamilyToType(viper.GetString("addressfamily")),
			Namespace:        pointer.PointerOrNil(viper.GetString("namespace")),
		},
	})
	if err != nil {
		return nil, err
	}

	return resp.Ips, nil
}

func (c *ip) Update(_ any) (*apiv2.IP, error) {
	panic("unimplemented")
}

func (c *ip) Convert(r *apiv2.IP) (string, any, any, error) {
	panic("unimplemented")
}
