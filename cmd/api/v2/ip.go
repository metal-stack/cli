package v2

import (
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
			cmd.Flags().StringSlice("labels", nil, "labels to add to the ip")
			cmd.Flags().BoolP("static", "", false, "make this ip static")
			cmd.Flags().StringP("addressfamily", "", "", "addressfamily, can be either IPv4|IPv6, defaults to IPv4 (optional)")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the ip")
			cmd.Flags().String("name", "", "name of the ip")
			cmd.Flags().String("description", "", "description of the ip")
			cmd.Flags().StringArray("labels", nil, "adds (or edits) the volume labels in the form of <key>=<value>")
			cmd.Flags().StringArray("remove-labels", nil, "removes the volume labels with the given key")
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
		CreateRequestFromCLI: w.createFromCLI,
		UpdateRequestFromCLI: w.updateFromCLI,
		ValidArgsFn:          c.Completion.IpListCompletion,
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *ip) createFromCLI() (*apiv2.IPServiceCreateRequest, error) {
	var labels *apiv2.Labels = nil

	labelSlice := viper.GetStringSlice("labels")
	if len(labelSlice) > 0 {
		labelsMap, err := genericcli.LabelsToMap(labelSlice)
		if err != nil {
			return nil, err
		}
		labels = &apiv2.Labels{Labels: labelsMap}
	}

	return &apiv2.IPServiceCreateRequest{
		Project:       c.c.GetProject(),
		Network:       viper.GetString("network"),
		Name:          pointer.PointerOrNil(viper.GetString("name")),
		Description:   pointer.PointerOrNil(viper.GetString("description")),
		Labels:        labels,
		Type:          new(ipStaticToType(viper.GetBool("static"))),
		AddressFamily: addressFamilyToType(viper.GetString("addressfamily")),
	}, nil
}

func (c *ip) updateFromCLI(args []string) (*apiv2.IPServiceUpdateRequest, error) {
	uuid, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	req := &apiv2.IPServiceUpdateRequest{
		Ip:      uuid,
		Project: c.c.GetProject(),
		UpdateMeta: &apiv2.UpdateMeta{
			LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
		},
	}

	if viper.IsSet("name") {
		req.Name = pointer.PointerOrNil(viper.GetString("name"))
	}
	if viper.IsSet("description") {
		req.Description = pointer.PointerOrNil(viper.GetString("description"))
	}
	if viper.IsSet("static") {
		req.Type = pointer.PointerOrNil(ipStaticToType(viper.GetBool("static")))
	}
	if viper.IsSet("remove-labels") || viper.IsSet("labels") {
		updates := &apiv2.LabelsPatch{}

		if viper.IsSet("remove-labels") {
			updates.Remove = viper.GetStringSlice("remove-labels")
		}

		if viper.IsSet("labels") {
			labels, err := genericcli.LabelsToMap(viper.GetStringSlice("labels"))
			if err != nil {
				return nil, err
			}
			updates.Update = &apiv2.Labels{Labels: labels}
		}

		req.Labels = &apiv2.UpdateLabels{
			Strategy: &apiv2.UpdateLabels_Patch{
				Patch: updates,
			},
		}
	}

	return req, nil
}

func (c *ip) Create(rq *apiv2.IPServiceCreateRequest) (*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().IP().Create(ctx, rq)
	if err != nil {
		if helpers.IsAlreadyExists(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, err
	}

	return resp.Ip, nil
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

	resp, err := c.c.Client.Apiv2().IP().Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Ip, nil
}

func (c *ip) Get(id string) (*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().IP().Get(ctx, &apiv2.IPServiceGetRequest{
		Project: c.c.GetProject(),
		Ip:      id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Ip, nil
}

func (c *ip) List() ([]*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().IP().List(ctx, &apiv2.IPServiceListRequest{
		Project: c.c.GetProject(),
	})
	if err != nil {
		return nil, err
	}

	return resp.Ips, nil
}

func (c *ip) Update(rq *apiv2.IPServiceUpdateRequest) (*apiv2.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().IP().Update(ctx, rq)
	if err != nil {
		return nil, err
	}

	return resp.Ip, nil
}

func (c *ip) Convert(r *apiv2.IP) (string, *apiv2.IPServiceCreateRequest, *apiv2.IPServiceUpdateRequest, error) {
	responseToUpdate, err := c.IpResponseToUpdate(r)
	return helpers.EncodeProject(r.Ip, r.Project), IpResponseToCreate(r), responseToUpdate, err
}

func IpResponseToCreate(ip *apiv2.IP) *apiv2.IPServiceCreateRequest {
	return &apiv2.IPServiceCreateRequest{
		Ip:          &ip.Ip,
		Project:     ip.Project,
		Network:     ip.Network,
		Name:        &ip.Name,
		Description: &ip.Description,
		Labels:      pointer.SafeDeref(ip.Meta).Labels,
		Type:        &ip.Type,
	}
}

func (c *ip) IpResponseToUpdate(ip *apiv2.IP) (*apiv2.IPServiceUpdateRequest, error) {
	return &apiv2.IPServiceUpdateRequest{
		Project:     ip.Project,
		Ip:          ip.Ip,
		Name:        &ip.Name,
		Description: &ip.Description,
		Type:        &ip.Type,
		Labels: &apiv2.UpdateLabels{
			Strategy: &apiv2.UpdateLabels_Replace{
				Replace: &apiv2.Labels{
					Labels: pointer.SafeDeref(pointer.SafeDeref(ip.Meta).Labels).Labels,
				},
			},
		},
		UpdateMeta: &apiv2.UpdateMeta{
			UpdatedAt:       pointer.SafeDeref(ip.Meta).UpdatedAt,
			LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_CLIENT,
		},
	}, nil
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
