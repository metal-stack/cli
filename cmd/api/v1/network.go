package v1

import (
	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/cli/pkg/common"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/metal-stack/metal-lib/pkg/tag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type networkCmd struct {
	c *config.Config
}

func newNetworkCmd(c *config.Config) *cobra.Command {
	w := &networkCmd{
		c: c,
	}

	listFlags := func(cmd *cobra.Command) {
		cmd.Flags().String("id", "", "ID to filter [optional]")
		cmd.Flags().String("name", "", "name to filter [optional]")
		cmd.Flags().String("description", "", "description to filter [optional]")
		cmd.Flags().String("partition", "", "partition to filter [optional]")
		cmd.Flags().String("project", "", "project to filter [optional]")
		cmd.Flags().StringSlice("prefixes", []string{}, "prefixes to filter")
		cmd.Flags().StringSlice("destination-prefixes", []string{}, "destination prefixes to filter")
		cmd.Flags().String("addressfamily", "", "addressfamily to filter, either ipv4 or ipv6 [optional]")
		cmd.Flags().Uint32("vrf", 0, "vrf to filter [optional]")
		cmd.Flags().StringSlice("labels", nil, "labels to filter [optional]")
		cmd.Flags().StringP("type", "t", "", "type of the network. [optional]")

		genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionListCompletion))
		genericcli.Must(cmd.RegisterFlagCompletionFunc("addressfamily", c.Completion.NetworkAddressFamilyCompletion))
		genericcli.Must(cmd.RegisterFlagCompletionFunc("type", c.Completion.NetworkTypeCompletion))
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv2.NetworkServiceCreateRequest, *apiv2.NetworkServiceUpdateRequest, *apiv2.Network]{
		BinaryName:           config.BinaryName,
		GenericCLI:           genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:             "network",
		Plural:               "networks",
		Description:          "networks can be attached to a machine or firewall such that they can communicate with each other.",
		CreateRequestFromCLI: w.createRequestFromCLI,
		UpdateRequestFromCLI: w.updateRequestFromCLI,
		Sorter:               sorters.NetworkSorter(),
		ValidArgsFn:          c.Completion.NetworkListCompletion,
		DescribePrinter:      func() printers.Printer { return c.DescribePrinter },
		ListPrinter:          func() printers.Printer { return c.ListPrinter },
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "name of the network to create. [required]")
			cmd.Flags().String("partition", "", "partition where this network should exist. [required]")
			cmd.Flags().String("project", "", "partition where this network should exist (alternative to parent-network-id). [optional]")
			cmd.Flags().String("parent-network-id", "", "the parent of the network (alternative to partition). [optional]")
			cmd.Flags().String("description", "", "description of the network to create. [optional]")
			cmd.Flags().StringSlice("labels", nil, "labels for this network. [optional]")
			cmd.Flags().String("addressfamily", "", "addressfamily of the network to acquire, if not specified the network inherits the address families from the parent [optional]")
			cmd.Flags().Uint32("ipv4-prefix-length", 0, "ipv4 prefix bit length of the network to create, defaults to default child prefix length of the parent network. [optional]")
			cmd.Flags().Uint32("ipv6-prefix-length", 0, "ipv6 prefix bit length of the network to create, defaults to default child prefix length of the parent network. [optional]")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("addressfamily", c.Completion.NetworkAddressFamilyCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("parent-network-id", c.Completion.NetworkListCompletion))
		},
		ListCmdMutateFn: func(cmd *cobra.Command) {
			listFlags(cmd)
			cmd.Flags().String("parent-network-id", "", "parent network to filter [optional]")
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the network [optional]")
			cmd.Flags().String("description", "", "the description of the network [optional]")
			cmd.Flags().StringSlice("labels", nil, "the labels of the network, must be in the form of key=value, use it like: --labels \"key1=value1,key2=value2\". [optional]")
			cmd.Flags().String("project", "", "project to filter [optional]")
		},
	}

	listBaseNetworksCmd := &cobra.Command{
		Use:   "list-base-networks",
		Short: "lists base networks that can be used for network creation",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return w.listBaseNetworks()
		},
		ValidArgsFunction: c.Completion.TenantMemberListCompletion,
	}
	listFlags(listBaseNetworksCmd)

	return genericcli.NewCmds(cmdsConfig, listBaseNetworksCmd)
}

func (c *networkCmd) Get(id string) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Network().Get(ctx, &apiv2.NetworkServiceGetRequest{
		Id:      id,
		Project: c.c.GetProject(),
	})
	if err != nil {
		return nil, err
	}

	return resp.Network, nil
}

func (c *networkCmd) List() ([]*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Network().List(ctx, &apiv2.NetworkServiceListRequest{
		Project: c.c.GetProject(),
		Query: &apiv2.NetworkQuery{
			Id:                  pointer.PointerOrNil(viper.GetString("id")),
			Name:                pointer.PointerOrNil(viper.GetString("name")),
			Description:         pointer.PointerOrNil(viper.GetString("description")),
			Partition:           pointer.PointerOrNil(viper.GetString("partition")),
			Project:             pointer.Pointer(c.c.GetProject()),
			Prefixes:            viper.GetStringSlice("prefixes"),
			DestinationPrefixes: viper.GetStringSlice("destination-prefixes"),
			Vrf:                 pointer.PointerOrNil(viper.GetUint32("vrf")),
			ParentNetwork:       pointer.PointerOrNil(viper.GetString("parent-network-id")),
			AddressFamily:       common.NetworkAddressFamilyToType(viper.GetString("addressfamily")),
			Labels: &apiv2.Labels{
				Labels: tag.NewTagMap(viper.GetStringSlice("labels")),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return resp.Networks, nil
}

func (c *networkCmd) Delete(id string) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Network().Delete(ctx, &apiv2.NetworkServiceDeleteRequest{
		Id:      id,
		Project: c.c.GetProject(),
	})
	if err != nil {
		return nil, err
	}

	return resp.Network, nil
}

func (c *networkCmd) Create(rq *apiv2.NetworkServiceCreateRequest) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Network().Create(ctx, rq)
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.AlreadyExists {
			return nil, genericcli.AlreadyExistsError()
		}
		return nil, err
	}

	return resp.Network, nil
}

func (c *networkCmd) Update(rq *apiv2.NetworkServiceUpdateRequest) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Network().Update(ctx, rq)
	if err != nil {
		return nil, err
	}

	return resp.Network, nil
}

func (c *networkCmd) Convert(r *apiv2.Network) (string, *apiv2.NetworkServiceCreateRequest, *apiv2.NetworkServiceUpdateRequest, error) {
	return r.Id, networkResponseToCreate(r), networkResponseToUpdate(r), nil
}

func networkResponseToCreate(r *apiv2.Network) *apiv2.NetworkServiceCreateRequest {
	meta := pointer.SafeDeref(r.Meta)

	return &apiv2.NetworkServiceCreateRequest{
		Project:     pointer.SafeDeref(r.Project),
		Name:        r.Name,
		Description: r.Description,
		Partition:   r.Partition,
		Labels: &apiv2.Labels{
			Labels: pointer.SafeDeref(meta.Labels).Labels,
		},
		ParentNetwork: r.ParentNetwork,
		// TODO: allow defining length and addressfamilies somehow?
	}
}

func networkResponseToUpdate(r *apiv2.Network) *apiv2.NetworkServiceUpdateRequest {
	meta := pointer.SafeDeref(r.Meta)

	return &apiv2.NetworkServiceUpdateRequest{
		Id:          r.Id,
		Project:     pointer.SafeDeref(r.Project),
		Name:        r.Name,
		Description: r.Description,
		Labels: &apiv2.UpdateLabels{
			Update: meta.Labels, // TODO: this only ensures that the labels are present but it does not cleanup old one's, which would require fetching the current state and calculating the diff
		}}
}

func (c *networkCmd) createRequestFromCLI() (*apiv2.NetworkServiceCreateRequest, error) {
	labels, err := genericcli.LabelsToMap(viper.GetStringSlice("labels"))
	if err != nil {
		return nil, err
	}

	var (
		cpl = &apiv2.ChildPrefixLength{}
	)
	if viper.IsSet("ipv4-prefix-length") {
		cpl.Ipv4 = pointer.Pointer(viper.GetUint32("ipv4-prefix-length"))
	}
	if viper.IsSet("ipv6-prefix-length") {
		cpl.Ipv6 = pointer.Pointer(viper.GetUint32("ipv6-prefix-length"))
	}

	return &apiv2.NetworkServiceCreateRequest{
		Description: pointer.PointerOrNil(viper.GetString("description")),
		Name:        pointer.PointerOrNil(viper.GetString("name")),
		Project:     c.c.GetProject(),
		Partition:   pointer.PointerOrNil(viper.GetString("partition")),
		Labels: &apiv2.Labels{
			Labels: labels,
		},
		ParentNetwork: pointer.PointerOrNil(viper.GetString("parent-network-id")),
		Length:        cpl,
		AddressFamily: common.NetworkAddressFamilyToType(viper.GetString("addressfamily")),
	}, nil
}

func (c *networkCmd) updateRequestFromCLI(args []string) (*apiv2.NetworkServiceUpdateRequest, error) {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	var labels *apiv2.UpdateLabels
	if viper.IsSet("labels") {
		lbls, err := genericcli.LabelsToMap(viper.GetStringSlice("labels"))
		if err != nil {
			return nil, err
		}

		labels = &apiv2.UpdateLabels{
			Update: &apiv2.Labels{
				Labels: lbls,
			},
		}
	}

	var (
		ur = &apiv2.NetworkServiceUpdateRequest{
			Id:          id,
			Project:     c.c.GetProject(),
			Description: pointer.PointerOrNil(viper.GetString("description")),
			Name:        pointer.PointerOrNil(viper.GetString("name")),
			Labels:      labels,
		}
	)

	return ur, nil
}

func (c *networkCmd) listBaseNetworks() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	var nwType *apiv2.NetworkType
	if viper.IsSet("type") {
		nt, err := enum.GetEnum[apiv2.NetworkType](viper.GetString("type"))
		if err != nil {
			return err
		}
		nwType = &nt
	}

	resp, err := c.c.Client.Apiv2().Network().ListBaseNetworks(ctx, &apiv2.NetworkServiceListBaseNetworksRequest{
		Project: c.c.GetProject(),
		Query: &apiv2.NetworkQuery{
			Id:                  pointer.PointerOrNil(viper.GetString("id")),
			Name:                pointer.PointerOrNil(viper.GetString("name")),
			Description:         pointer.PointerOrNil(viper.GetString("description")),
			Partition:           pointer.PointerOrNil(viper.GetString("partition")),
			Project:             pointer.PointerOrNil(viper.GetString("project")),
			Prefixes:            viper.GetStringSlice("prefixes"),
			DestinationPrefixes: viper.GetStringSlice("destination-prefixes"),
			Vrf:                 pointer.PointerOrNil(viper.GetUint32("vrf")),
			AddressFamily:       common.NetworkAddressFamilyToType(viper.GetString("addressfamily")),
			Labels: &apiv2.Labels{
				Labels: tag.NewTagMap(viper.GetStringSlice("labels")),
			},
			Type: nwType,
		},
	})

	if err != nil {
		return err
	}

	return c.c.ListPrinter.Print(resp.Networks)
}
