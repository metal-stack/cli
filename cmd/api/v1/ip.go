package v1

import (
	"fmt"

	"connectrpc.com/connect"
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

	cmdsConfig := &genericcli.CmdsConfig[*apiv2.IPServiceCreateRequest, *apiv2.IPServiceUpdateRequest, *apiv2.IP]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "ip",
		Plural:          "ips",
		Description:     "an ip address of metal-stack.io",
		Sorter:          sorters.IPSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project from where ips should be listed")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the ip")
			cmd.Flags().StringP("network", "n", "", "network from which the ip should get created")
			cmd.Flags().StringP("name", "", "", "name of the ip")
			cmd.Flags().StringP("description", "", "", "description of the ip")
			cmd.Flags().StringSliceP("tags", "", nil, "tags to add to the ip")
			cmd.Flags().BoolP("static", "", false, "make this ip static")
			cmd.Flags().StringP("addressfamily", "", "", "addressfamily, can be either IPv4|IPv6, defaults to IPv4 (optional)")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the ip")
			cmd.Flags().String("name", "", "name of the ip")
			cmd.Flags().String("description", "", "description of the ip")
			cmd.Flags().StringSlice("tags", nil, "tags of the ip")
			cmd.Flags().Bool("static", false, "make this ip static")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the ip")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the ip")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		CreateRequestFromCLI: func() (*apiv2.IPServiceCreateRequest, error) {
			return &apiv2.IPServiceCreateRequest{
				Project:     c.GetProject(),
				Name:        pointer.Pointer(viper.GetString("name")),
				Description: pointer.Pointer(viper.GetString("description")),
				Network:     viper.GetString("network"),
				// Labels:        viper.GetStringSlice("tags"), // FIXME implement
				Type:          pointer.Pointer(ipStaticToType(viper.GetBool("static"))),
				AddressFamily: addressFamilyToType(viper.GetString("addressfamily")),
			}, nil
		},
		UpdateRequestFromCLI: w.updateFromCLI,
		ValidArgsFn:          c.Completion.IpListCompletion,
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *ip) updateFromCLI(args []string) (*apiv2.IPServiceUpdateRequest, error) {
	uuid, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	ipToUpdate, err := c.Get(uuid)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve ip: %w", err)
	}

	if viper.IsSet("name") {
		ipToUpdate.Name = viper.GetString("name")
	}
	if viper.IsSet("description") {
		ipToUpdate.Description = viper.GetString("description")
	}
	if viper.IsSet("static") {
		ipToUpdate.Type = ipStaticToType(viper.GetBool("static"))
	}
	// if viper.IsSet("tags") {
	// if ipToUpdate.Meta == nil {
	// 	ipToUpdate.Meta = &apiv2.Meta{
	// 		Labels: &apiv2.Labels{},
	// 	}
	// }
	// for _, l :=

	// ipToUpdate.Meta.Labels = viper.GetStringSlice("tags")
	// FIXME implement
	// }

	return IpResponseToUpdate(ipToUpdate), nil
}

func (c *ip) Create(rq *apiv2.IPServiceCreateRequest) (*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().IP().Create(ctx, connect.NewRequest(rq))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Ip, nil
}

func (c *ip) Delete(id string) (*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.IPServiceDeleteRequest{
		Project: c.c.GetProject(),
		Ip:      id,
	}

	if viper.IsSet("file") {
		var err error
		req.Ip, req.Project, err = helpers.DecodeProject(id)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.c.Client.Apiv2().IP().Delete(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Ip, nil
}

func (c *ip) Get(id string) (*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().IP().Get(ctx, connect.NewRequest(&apiv2.IPServiceGetRequest{
		Project: c.c.GetProject(),
		Ip:      id,
	}))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Ip, nil
}

func (c *ip) List() ([]*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().IP().List(ctx, connect.NewRequest(&apiv2.IPServiceListRequest{
		Project: c.c.GetProject(),
	}))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Ips, nil
}

func (c *ip) Update(rq *apiv2.IPServiceUpdateRequest) (*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().IP().Update(ctx, connect.NewRequest(rq))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Ip, nil
}

func (*ip) Convert(r *apiv2.IP) (string, *apiv2.IPServiceCreateRequest, *apiv2.IPServiceUpdateRequest, error) {
	return helpers.EncodeProject(r.Uuid, r.Project), IpResponseToCreate(r), IpResponseToUpdate(r), nil
}

func IpResponseToCreate(r *apiv2.IP) *apiv2.IPServiceCreateRequest {
	return &apiv2.IPServiceCreateRequest{
		Project:     r.Project,
		Name:        &r.Name,
		Description: &r.Description,
		Labels:      r.Meta.Labels,
		Type:        &r.Type,
	}
}

func IpResponseToUpdate(r *apiv2.IP) *apiv2.IPServiceUpdateRequest {
	meta := pointer.SafeDeref(r.Meta)

	return &apiv2.IPServiceUpdateRequest{
		Project:     r.Project,
		Ip:          r.Ip,
		Name:        &r.Name,
		Description: &r.Description,
		Type:        &r.Type,
		Labels: &apiv2.UpdateLabels{
			Update: meta.Labels, // TODO: this only ensures that the labels are present but it does not cleanup old one's, which would require fetching the current state and calculating the diff
		},
	}
}

func ipStaticToType(b bool) apiv2.IPType {
	if b {
		return apiv2.IPType_IP_TYPE_STATIC
	}
	return apiv2.IPType_IP_TYPE_EPHEMERAL
}

func addressFamilyToType(af string) *apiv2.IPAddressFamily {
	switch af {
	case "":
		return nil
	case "ipv4", "IPv4":
		return apiv2.IPAddressFamily_IP_ADDRESS_FAMILY_V4.Enum()
	case "ipv6", "IPv6":
		return apiv2.IPAddressFamily_IP_ADDRESS_FAMILY_V6.Enum()
	default:
		return apiv2.IPAddressFamily_IP_ADDRESS_FAMILY_UNSPECIFIED.Enum()
	}
}
