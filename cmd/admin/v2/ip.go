package v2

import (
	"strings"

	"github.com/metal-stack/api/go/enum"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
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
		Description:     "manage ip addresses",
		Sorter:          sorters.IPSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("ip", "", "ipaddress to filter [optional]")
			cmd.Flags().String("name", "", "name to filter [optional]")
			cmd.Flags().String("network", "", "network to filter [optional]")
			cmd.Flags().String("project", "", "project to filter [optional]")
			cmd.Flags().String("uuid", "", "allocation uuid to filter [optional]")
			cmd.Flags().String("machine", "", "machine to filter [optional]")
			cmd.Flags().String("namespace", "", "namespace to filter [optional]")
			cmd.Flags().String("parent-prefix", "", "parent-prefix to filter [optional]")
			cmd.Flags().String("type", "", "type to filter [optional] can be either ephemeral|static")
			cmd.Flags().String("addressfamily", "", "addressfamily to filter [optional] can be either ipv6|ipv6")
			cmd.Flags().StringSlice("label", nil, "label to filter, must be in the form of key=value, can be either specified multiple times, or comma seperated [optional]")
			genericcli.Must(cmd.RegisterFlagCompletionFunc("network", c.Completion.NetworkListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("type", c.Completion.IpTypeCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("addressfamily", c.Completion.IpAddressFamilyCompletion))
		},
		ValidArgsFn: c.Completion.IpListCompletion,
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *ip) List() ([]*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	var (
		labels *apiv2.Labels
		ipType *apiv2.IPType
		af     *apiv2.IPAddressFamily
	)

	if len(viper.GetStringSlice("label")) > 0 {
		labels = &apiv2.Labels{
			Labels: map[string]string{},
		}
		for _, label := range viper.GetStringSlice("label") {
			key, value, _ := strings.Cut(label, "=")
			labels.Labels[key] = value
		}
	}

	if viper.IsSet("type") {
		ipt, err := enum.GetEnum[apiv2.IPType](viper.GetString("type"))
		if err != nil {
			return nil, err
		}
		ipType = &ipt
	}

	if viper.IsSet("addressfamily") {
		ipaf, err := enum.GetEnum[apiv2.IPAddressFamily](viper.GetString("addressfamily"))
		if err != nil {
			return nil, err
		}
		af = &ipaf
	}

	resp, err := c.c.Client.Adminv2().IP().List(ctx, &adminv2.IPServiceListRequest{
		Query: &apiv2.IPQuery{
			Ip:               pointer.PointerOrNil(viper.GetString("ip")),
			Network:          pointer.PointerOrNil(viper.GetString("network")),
			Name:             pointer.PointerOrNil(viper.GetString("name")),
			Project:          pointer.PointerOrNil(viper.GetString("project")),
			Uuid:             pointer.PointerOrNil(viper.GetString("uuid")),
			Machine:          pointer.PointerOrNil(viper.GetString("machine")),
			ParentPrefixCidr: pointer.PointerOrNil(viper.GetString("parent-prefix")),
			Namespace:        pointer.PointerOrNil(viper.GetString("namespace")),
			Labels:           labels,
			Type:             ipType,
			AddressFamily:    af,
		},
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
