package v1

import (
	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/enum"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
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

	// TODO: move to common?
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

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.NetworkServiceCreateRequest, *adminv2.NetworkServiceUpdateRequest, *apiv2.Network]{
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
			cmd.Flags().String("id", "", "id of the network to create, defaults to a random uuid if not provided. [optional]")
			cmd.Flags().String("name", "", "name of the network to create. [required]")
			cmd.Flags().StringP("type", "t", "", "type of the network. [required]")
			cmd.Flags().String("nat-type", "", "nat-type of the network. [required]")
			cmd.Flags().String("partition", "", "partition where this network should exist. [required]")
			cmd.Flags().String("project", "", "partition where this network should exist (alternative to parent-network-id). [optional]")
			cmd.Flags().String("parent-network-id", "", "the parent of the network (alternative to partition). [optional]")
			cmd.Flags().String("description", "", "description of the network to create. [optional]")
			cmd.Flags().StringSlice("labels", nil, "labels for this network. [optional]")
			cmd.Flags().String("addressfamily", "", "addressfamily of the network to acquire, if not specified the network inherits the address families from the parent [optional]")
			cmd.Flags().Uint32("ipv4-prefix-length", 0, "ipv4 prefix bit length of the network to create, defaults to default child prefix length of the parent network. [optional]")
			cmd.Flags().Uint32("ipv6-prefix-length", 0, "ipv6 prefix bit length of the network to create, defaults to default child prefix length of the parent network. [optional]")
			cmd.Flags().Uint32("default-ipv4-prefix-length", 0, "default ipv4 prefix bit length of the network to create. [optional]")
			cmd.Flags().Uint32("default-ipv6-prefix-length", 0, "default ipv6 prefix bit length of the network to create. [optional]")
			cmd.Flags().Uint32("min-ipv4-prefix-length", 0, "min ipv4 prefix bit length of the network to create. [optional]")
			cmd.Flags().Uint32("min-ipv6-prefix-length", 0, "min ipv6 prefix bit length of the network to create. [optional]")
			cmd.Flags().StringSlice("prefixes", nil, "prefixes for this network. [optional]")
			cmd.Flags().StringSlice("destination-prefixes", nil, "destination-prefixes for this network. [optional]")
			cmd.Flags().StringSlice("additional-announcable-cidrs", nil, "additional-announcable-cidrs for this network. [optional]")
			cmd.Flags().Uint32("vrf", 0, "the vrf of the network to create. [optional]")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("addressfamily", c.Completion.NetworkAddressFamilyCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("type", c.Completion.NetworkTypeCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("nat-type", c.Completion.NetworkNatTypeCompletion))
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

	return genericcli.NewCmds(cmdsConfig)
}

func (c *networkCmd) Get(id string) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Network().Get(ctx, connect.NewRequest(&adminv2.NetworkServiceGetRequest{
		Id: id,
	}))

	if err != nil {
		return nil, err
	}

	return resp.Msg.Network, nil
}

func (c *networkCmd) List() ([]*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	var nwType *apiv2.NetworkType
	if viper.IsSet("type") {
		nt, err := enum.GetEnum[apiv2.NetworkType](viper.GetString("type"))
		if err != nil {
			return nil, err
		}
		nwType = &nt
	}

	resp, err := c.c.Client.Adminv2().Network().List(ctx, connect.NewRequest(&adminv2.NetworkServiceListRequest{
		Query: &apiv2.NetworkQuery{
			Id:                  pointer.PointerOrNil(viper.GetString("id")),
			Name:                pointer.PointerOrNil(viper.GetString("name")),
			Description:         pointer.PointerOrNil(viper.GetString("description")),
			Partition:           pointer.PointerOrNil(viper.GetString("partition")),
			Project:             pointer.PointerOrNil(viper.GetString("project")),
			Prefixes:            viper.GetStringSlice("prefixes"),
			DestinationPrefixes: viper.GetStringSlice("destination-prefixes"),
			Vrf:                 pointer.PointerOrNil(viper.GetUint32("vrf")),
			ParentNetworkId:     pointer.PointerOrNil(viper.GetString("parent-network-id")),
			AddressFamily:       common.NetworkAddressFamilyToType(viper.GetString("addressfamily")),
			Labels: &apiv2.Labels{
				Labels: tag.NewTagMap(viper.GetStringSlice("labels")),
			},
			Type: nwType,
			// NatType: (*apiv2.NATType)(nwType),
		},
	}))

	if err != nil {
		return nil, err
	}

	return resp.Msg.Networks, nil
}

func (c *networkCmd) Delete(id string) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Network().Delete(ctx, connect.NewRequest(&adminv2.NetworkServiceDeleteRequest{
		Id: id,
	}))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Network, nil
}

func (c *networkCmd) Create(rq *adminv2.NetworkServiceCreateRequest) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Network().Create(ctx, connect.NewRequest(rq))
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.AlreadyExists {
			return nil, genericcli.AlreadyExistsError()
		}
		return nil, err
	}

	return resp.Msg.Network, nil
}

func (c *networkCmd) Update(rq *adminv2.NetworkServiceUpdateRequest) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Network().Update(ctx, connect.NewRequest(rq))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Network, nil
}

func (c *networkCmd) Convert(r *apiv2.Network) (string, *adminv2.NetworkServiceCreateRequest, *adminv2.NetworkServiceUpdateRequest, error) {
	return r.Id, networkResponseToCreate(r), networkResponseToUpdate(r), nil
}

func networkResponseToCreate(r *apiv2.Network) *adminv2.NetworkServiceCreateRequest {
	meta := pointer.SafeDeref(r.Meta)

	return &adminv2.NetworkServiceCreateRequest{
		Project:     r.Project,
		Name:        r.Name,
		Description: r.Description,
		Partition:   r.Partition,
		Labels: &apiv2.Labels{
			Labels: pointer.SafeDeref(meta.Labels).Labels,
		},
		ParentNetworkId: r.ParentNetworkId,
		// TODO: allow defining length and addressfamilies somehow?
	}
}

func networkResponseToUpdate(r *apiv2.Network) *adminv2.NetworkServiceUpdateRequest {
	meta := pointer.SafeDeref(r.Meta)

	return &adminv2.NetworkServiceUpdateRequest{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Labels: &apiv2.UpdateLabels{
			Update: meta.Labels, // TODO: this only ensures that the labels are present but it does not cleanup old one's, which would require fetching the current state and calculating the diff
		},
		Prefixes:                 []string{},
		DestinationPrefixes:      []string{},
		DefaultChildPrefixLength: &apiv2.ChildPrefixLength{},
		MinChildPrefixLength:     &apiv2.ChildPrefixLength{},
		// NatType:                    &0,
		AdditionalAnnouncableCidrs: []string{},
		Force:                      false,
	}
}

func (c *networkCmd) createRequestFromCLI() (*adminv2.NetworkServiceCreateRequest, error) {
	labels, err := genericcli.LabelsToMap(viper.GetStringSlice("labels"))
	if err != nil {
		return nil, err
	}

	var (
		natType    = apiv2.NATType_NAT_TYPE_NONE
		defaultCPL *apiv2.ChildPrefixLength
		minCPL     *apiv2.ChildPrefixLength
		length     *apiv2.ChildPrefixLength
	)
	if viper.IsSet("default-ipv4-prefix-length") {
		defaultCPL = &apiv2.ChildPrefixLength{
			Ipv4: pointer.Pointer(viper.GetUint32("default-ipv4-prefix-length")),
		}
	}
	if viper.IsSet("default-ipv6-prefix-length") {
		if defaultCPL == nil {
			defaultCPL = &apiv2.ChildPrefixLength{}
		}
		defaultCPL.Ipv6 = pointer.Pointer(viper.GetUint32("default-ipv6-prefix-length"))
	}
	if viper.IsSet("min-ipv4-prefix-length") {
		minCPL = &apiv2.ChildPrefixLength{
			Ipv4: pointer.Pointer(viper.GetUint32("min-ipv4-prefix-length")),
		}
	}
	if viper.IsSet("min-ipv6-prefix-length") {
		if minCPL == nil {
			minCPL = &apiv2.ChildPrefixLength{}
		}
		minCPL.Ipv6 = pointer.Pointer(viper.GetUint32("min-ipv6-prefix-length"))
	}
	if viper.IsSet("ipv4-prefix-length") {
		length = &apiv2.ChildPrefixLength{
			Ipv4: pointer.Pointer(viper.GetUint32("ipv4-prefix-length")),
		}
	}
	if viper.IsSet("ipv6-prefix-length") {
		if length == nil {
			length = &apiv2.ChildPrefixLength{}
		}
		length.Ipv6 = pointer.Pointer(viper.GetUint32("ipv6-prefix-length"))
	}

	nwType, err := enum.GetEnum[apiv2.NetworkType](viper.GetString("type"))
	if err != nil {
		return nil, err
	}

	if viper.IsSet("nat-type") {
		natType, err = enum.GetEnum[apiv2.NATType](viper.GetString("nat-type"))
		if err != nil {
			return nil, err
		}
	}

	var vrf *uint32
	if viper.IsSet("vrf") {
		vrf = pointer.Pointer(viper.GetUint32("vrf"))
	}

	return &adminv2.NetworkServiceCreateRequest{
		Description: pointer.PointerOrNil(viper.GetString("description")),
		Name:        pointer.PointerOrNil(viper.GetString("name")),
		Project:     pointer.PointerOrNil(viper.GetString("project")),
		Partition:   pointer.PointerOrNil(viper.GetString("partition")),
		Labels: &apiv2.Labels{
			Labels: labels,
		},
		ParentNetworkId:            pointer.PointerOrNil(viper.GetString("parent-network-id")),
		AddressFamily:              common.NetworkAddressFamilyToType(viper.GetString("addressfamily")),
		Id:                         pointer.PointerOrNil(viper.GetString("id")),
		Type:                       nwType,
		Prefixes:                   viper.GetStringSlice("prefixes"),
		DestinationPrefixes:        viper.GetStringSlice("destination-prefixes"),
		DefaultChildPrefixLength:   defaultCPL,
		MinChildPrefixLength:       minCPL,
		NatType:                    &natType,
		Vrf:                        vrf,
		AdditionalAnnouncableCidrs: viper.GetStringSlice("additional-announcable-cidrs"),
		Length:                     length,
	}, nil
}

func (c *networkCmd) updateRequestFromCLI(args []string) (*adminv2.NetworkServiceUpdateRequest, error) {
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
		natType    = apiv2.NATType_NAT_TYPE_NONE
		defaultCPL *apiv2.ChildPrefixLength
		minCPL     *apiv2.ChildPrefixLength
	)
	if viper.IsSet("default-ipv4-prefix-length") {
		defaultCPL = &apiv2.ChildPrefixLength{
			Ipv4: pointer.Pointer(viper.GetUint32("default-ipv4-prefix-length")),
		}
	}
	if viper.IsSet("default-ipv6-prefix-length") {
		if defaultCPL == nil {
			defaultCPL = &apiv2.ChildPrefixLength{}
		}
		defaultCPL.Ipv6 = pointer.Pointer(viper.GetUint32("default-ipv6-prefix-length"))
	}
	if viper.IsSet("min-ipv4-prefix-length") {
		minCPL = &apiv2.ChildPrefixLength{
			Ipv4: pointer.Pointer(viper.GetUint32("min-ipv4-prefix-length")),
		}
	}
	if viper.IsSet("min-ipv6-prefix-length") {
		if minCPL == nil {
			minCPL = &apiv2.ChildPrefixLength{}
		}
		minCPL.Ipv6 = pointer.Pointer(viper.GetUint32("min-ipv6-prefix-length"))
	}

	if viper.IsSet("nat-type") {
		natType, err = enum.GetEnum[apiv2.NATType](viper.GetString("nat-type"))
		if err != nil {
			return nil, err
		}
	}

	var (
		ur = &adminv2.NetworkServiceUpdateRequest{
			Id:                         id,
			Description:                pointer.PointerOrNil(viper.GetString("description")),
			Name:                       pointer.PointerOrNil(viper.GetString("name")),
			Labels:                     labels,
			Prefixes:                   viper.GetStringSlice("prefixes"),
			DestinationPrefixes:        viper.GetStringSlice("destination-prefixes"),
			DefaultChildPrefixLength:   defaultCPL,
			MinChildPrefixLength:       minCPL,
			NatType:                    &natType,
			AdditionalAnnouncableCidrs: viper.GetStringSlice("additional-announcable-cidrs"),
			Force:                      viper.GetBool("force"),
		}
	)

	return ur, nil
}
