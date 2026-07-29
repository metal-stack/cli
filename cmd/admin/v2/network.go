package v2

import (
	"github.com/metal-stack/api/go/enum"
	"github.com/metal-stack/api/go/errorutil"
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

type networkCmd struct {
	c *config.Config
}

func newNetworkCmd(c *config.Config) *cobra.Command {
	var (
		w = &networkCmd{
			c: c,
		}
		listFlags = func(cmd *cobra.Command) {
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
			cmd.Flags().String("nat-type", "", "nat type of the network. [optional]")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("addressfamily", c.Completion.NetworkAddressFamilyCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("type", c.Completion.NetworkTypeCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("nat-type", c.Completion.NetworkTypeCompletion))
		}
	)

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.NetworkServiceCreateRequest, *adminv2.NetworkServiceUpdateRequest, *apiv2.Network]{
		BinaryName:           config.BinaryName,
		GenericCLI:           genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:             "network",
		Plural:               "networks",
		Description:          "networks can be attached to a machine or firewall such that they can communicate with each other.",
		CreateRequestFromCLI: w.createRequestFromCLI,
		UpdateRequestFromCLI: w.updateRequestFromCLI,
		Sorter:               sorters.NetworkSorter(),
		ValidArgsFn:          c.Completion.NetworkAdminListCompletion,
		DescribePrinter:      func() printers.Printer { return c.DescribePrinter },
		ListPrinter:          func() printers.Printer { return c.ListPrinter },
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("id", "", "id of the network to create, defaults to a random uuid if not provided. [optional]")
			cmd.Flags().String("name", "", "name of the network to create. [required]")
			cmd.Flags().StringP("type", "t", "", "type of the network. [required]")
			cmd.Flags().String("nat-type", "", "nat-type of the network. [required]")
			cmd.Flags().String("partition", "", "partition where this network should exist. [required]")
			cmd.Flags().String("project", "", "partition where this network should exist (alternative to parent-network). [optional]")
			cmd.Flags().String("parent-network", "", "the parent of the network (alternative to partition). [optional]")
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
			genericcli.Must(cmd.RegisterFlagCompletionFunc("parent-network", c.Completion.NetworkListCompletion))
		},
		ListCmdMutateFn: func(cmd *cobra.Command) {
			listFlags(cmd)
			cmd.Flags().String("parent-network", "", "parent network to filter [optional]")
			genericcli.Must(cmd.RegisterFlagCompletionFunc("parent-network", c.Completion.NetworkListCompletion))
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the network [optional]")
			cmd.Flags().String("description", "", "the description of the network [optional]")
			cmd.Flags().String("project", "", "project to filter [optional]")
			cmd.Flags().StringSlice("labels", nil, "labels to replace for the network")
			cmd.Flags().StringSlice("add-labels", nil, "labels to add to the network")
			cmd.Flags().StringSlice("remove-labels", nil, "labels to remove to the network")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *networkCmd) Get(id string) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Network().Get(ctx, &adminv2.NetworkServiceGetRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return resp.Network, nil
}

func (c *networkCmd) List() ([]*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	var (
		nwType  *apiv2.NetworkType
		natType *apiv2.NATType
	)

	if viper.IsSet("type") {
		nt, err := enum.GetEnum[apiv2.NetworkType](viper.GetString("type"))
		if err != nil {
			return nil, err
		}

		nwType = &nt
	}

	if viper.IsSet("nat-type") {
		nat, err := enum.GetEnum[apiv2.NATType](viper.GetString("nat-type"))
		if err != nil {
			return nil, err
		}

		natType = &nat
	}

	req := &adminv2.NetworkServiceListRequest{
		Query: &apiv2.NetworkQuery{
			Id:                  pointer.PointerOrNil(viper.GetString("id")),
			Name:                pointer.PointerOrNil(viper.GetString("name")),
			Description:         pointer.PointerOrNil(viper.GetString("description")),
			Partition:           pointer.PointerOrNil(viper.GetString("partition")),
			Project:             pointer.PointerOrNil(viper.GetString("project")),
			Prefixes:            viper.GetStringSlice("prefixes"),
			DestinationPrefixes: viper.GetStringSlice("destination-prefixes"),
			Vrf:                 pointer.PointerOrNil(viper.GetUint32("vrf")),
			ParentNetwork:       pointer.PointerOrNil(viper.GetString("parent-network")),
			AddressFamily:       helpers.NetworkAddressFamilyToType(viper.GetString("addressfamily")),
			Type:                nwType,
			NatType:             natType,
		},
	}

	if labelSlice := viper.GetStringSlice("labels"); len(labelSlice) > 0 {
		var err error

		req.Query.Labels, err = helpers.LabelsFromSlice(labelSlice)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.c.Client.Adminv2().Network().List(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Networks, nil
}

func (c *networkCmd) Delete(id string) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Network().Delete(ctx, &adminv2.NetworkServiceDeleteRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Network, nil
}

func (c *networkCmd) Create(rq *adminv2.NetworkServiceCreateRequest) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Network().Create(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}
		return nil, err
	}

	return resp.Network, nil
}

func (c *networkCmd) Update(rq *adminv2.NetworkServiceUpdateRequest) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Network().Update(ctx, rq)
	if err != nil {
		return nil, err
	}

	return resp.Network, nil
}

func (c *networkCmd) Convert(r *apiv2.Network) (string, *adminv2.NetworkServiceCreateRequest, *adminv2.NetworkServiceUpdateRequest, error) {
	addressFamily, err := helpers.AddressFamilyFromPrefixes(r.Prefixes...)
	if err != nil {
		return "", nil, nil, err
	}

	return r.Id, networkResponseToCreate(r, addressFamily), networkResponseToUpdate(r), nil
}

func networkResponseToCreate(r *apiv2.Network, addressFamily *apiv2.NetworkAddressFamily) *adminv2.NetworkServiceCreateRequest {
	meta := pointer.SafeDeref(r.Meta)

	return &adminv2.NetworkServiceCreateRequest{
		Project:     r.Project,
		Name:        r.Name,
		Description: r.Description,
		Partition:   r.Partition,
		Labels: &apiv2.Labels{
			Labels: pointer.SafeDeref(meta.Labels).Labels,
		},
		ParentNetwork:              r.ParentNetwork,
		Id:                         &r.Id,
		Type:                       r.Type,
		Prefixes:                   r.Prefixes,
		DestinationPrefixes:        r.DestinationPrefixes,
		DefaultChildPrefixLength:   r.DefaultChildPrefixLength,
		MinChildPrefixLength:       r.MinChildPrefixLength,
		NatType:                    &r.NatType,
		Vrf:                        r.Vrf,
		AdditionalAnnouncableCidrs: r.AdditionalAnnouncableCidrs,
		// TODO
		// Length:                     r.Length,
		AddressFamily: addressFamily,
	}
}

func networkResponseToUpdate(r *apiv2.Network) *adminv2.NetworkServiceUpdateRequest {
	return &adminv2.NetworkServiceUpdateRequest{
		UpdateMeta:                 helpers.UpdateMetaFromMeta(r.Meta),
		Labels:                     helpers.UpdateLabelsFromMeta(r.Meta),
		Id:                         r.Id,
		Name:                       r.Name,
		Description:                r.Description,
		Prefixes:                   r.Prefixes,
		DestinationPrefixes:        r.DestinationPrefixes,
		DefaultChildPrefixLength:   r.DefaultChildPrefixLength,
		MinChildPrefixLength:       r.MinChildPrefixLength,
		NatType:                    &r.NatType,
		AdditionalAnnouncableCidrs: r.AdditionalAnnouncableCidrs,
		Force:                      viper.GetBool("FORCE_NETWORK"),
	}
}

func (c *networkCmd) createRequestFromCLI() (*adminv2.NetworkServiceCreateRequest, error) {
	labels, err := helpers.LabelsFromSlice(viper.GetStringSlice("labels"))
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
			Ipv4: new(viper.GetUint32("default-ipv4-prefix-length")),
		}
	}
	if viper.IsSet("default-ipv6-prefix-length") {
		if defaultCPL == nil {
			defaultCPL = &apiv2.ChildPrefixLength{}
		}

		defaultCPL.Ipv6 = new(viper.GetUint32("default-ipv6-prefix-length"))
	}

	if viper.IsSet("min-ipv4-prefix-length") {
		minCPL = &apiv2.ChildPrefixLength{
			Ipv4: new(viper.GetUint32("min-ipv4-prefix-length")),
		}
	}
	if viper.IsSet("min-ipv6-prefix-length") {
		if minCPL == nil {
			minCPL = &apiv2.ChildPrefixLength{}
		}

		minCPL.Ipv6 = new(viper.GetUint32("min-ipv6-prefix-length"))
	}

	if viper.IsSet("ipv4-prefix-length") {
		length = &apiv2.ChildPrefixLength{
			Ipv4: new(viper.GetUint32("ipv4-prefix-length")),
		}
	}
	if viper.IsSet("ipv6-prefix-length") {
		if length == nil {
			length = &apiv2.ChildPrefixLength{}
		}

		length.Ipv6 = new(viper.GetUint32("ipv6-prefix-length"))
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
		vrf = new(viper.GetUint32("vrf"))
	}

	return &adminv2.NetworkServiceCreateRequest{
		Description:                pointer.PointerOrNil(viper.GetString("description")),
		Name:                       pointer.PointerOrNil(viper.GetString("name")),
		Project:                    pointer.PointerOrNil(viper.GetString("project")),
		Partition:                  pointer.PointerOrNil(viper.GetString("partition")),
		Labels:                     labels,
		ParentNetwork:              pointer.PointerOrNil(viper.GetString("parent-network")),
		AddressFamily:              helpers.NetworkAddressFamilyToType(viper.GetString("addressfamily")),
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

	updateLabels, err := helpers.UpdateLabelsFromCLI()
	if err != nil {
		return nil, err
	}

	var (
		natType    = apiv2.NATType_NAT_TYPE_NONE
		defaultCPL *apiv2.ChildPrefixLength
		minCPL     *apiv2.ChildPrefixLength
	)
	if viper.IsSet("default-ipv4-prefix-length") {
		defaultCPL = &apiv2.ChildPrefixLength{
			Ipv4: new(viper.GetUint32("default-ipv4-prefix-length")),
		}
	}
	if viper.IsSet("default-ipv6-prefix-length") {
		if defaultCPL == nil {
			defaultCPL = &apiv2.ChildPrefixLength{}
		}
		defaultCPL.Ipv6 = new(viper.GetUint32("default-ipv6-prefix-length"))
	}
	if viper.IsSet("min-ipv4-prefix-length") {
		minCPL = &apiv2.ChildPrefixLength{
			Ipv4: new(viper.GetUint32("min-ipv4-prefix-length")),
		}
	}
	if viper.IsSet("min-ipv6-prefix-length") {
		if minCPL == nil {
			minCPL = &apiv2.ChildPrefixLength{}
		}
		minCPL.Ipv6 = new(viper.GetUint32("min-ipv6-prefix-length"))
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
			Labels:                     updateLabels,
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
