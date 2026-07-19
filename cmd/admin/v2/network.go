package v2

import (
	"fmt"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type adminNetwork struct {
	c *config.Config
}

func newAdminNetworkCmd(c *config.Config) *cobra.Command {
	w := &adminNetwork{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.NetworkServiceCreateRequest, *adminv2.NetworkServiceUpdateRequest, *apiv2.Network]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "network",
		Plural:          "networks",
		Description:     "manage networks",
		Sorter:          sorters.NetworkSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.NetworkListCompletion,
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *adminNetwork) Get(id string) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.NetworkServiceGetRequest{Id: id}

	resp, err := c.c.Client.Adminv2().Network().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get network: %w", err)
	}

	return resp.Network, nil
}

func (c *adminNetwork) List() ([]*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.NetworkServiceListRequest{}

	resp, err := c.c.Client.Adminv2().Network().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	return resp.Networks, nil
}

func (c *adminNetwork) Create(rq *adminv2.NetworkServiceCreateRequest) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Network().Create(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	return resp.Network, nil
}

func (c *adminNetwork) Delete(id string) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Network().Delete(ctx, &adminv2.NetworkServiceDeleteRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("failed to delete network: %w", err)
	}

	return resp.Network, nil
}

func (c *adminNetwork) Convert(r *apiv2.Network) (string, *adminv2.NetworkServiceCreateRequest, *adminv2.NetworkServiceUpdateRequest, error) {
	return r.Id, &adminv2.NetworkServiceCreateRequest{
			Id:                         &r.Id,
			Name:                       r.Name,
			Description:                r.Description,
			Partition:                  r.Partition,
			Project:                    r.Project,
			Type:                       r.Type,
			Prefixes:                   r.Prefixes,
			DestinationPrefixes:        r.DestinationPrefixes,
			DefaultChildPrefixLength:   r.DefaultChildPrefixLength,
			MinChildPrefixLength:       r.MinChildPrefixLength,
			Labels:                     r.Meta.Labels,
			NatType:                    r.NatType.Enum(),
			Vrf:                        r.Vrf,
			ParentNetwork:              r.ParentNetwork,
			AdditionalAnnouncableCidrs: r.AdditionalAnnouncableCidrs,
			Length:                     r.DefaultChildPrefixLength,
		}, &adminv2.NetworkServiceUpdateRequest{
			Id:                         r.Id,
			Name:                       r.Name,
			Description:                r.Description,
			Prefixes:                   r.Prefixes,
			DestinationPrefixes:        r.DestinationPrefixes,
			DefaultChildPrefixLength:   r.DefaultChildPrefixLength,
			MinChildPrefixLength:       r.MinChildPrefixLength,
			NatType:                    r.NatType.Enum(),
			AdditionalAnnouncableCidrs: r.AdditionalAnnouncableCidrs,
			Labels: &apiv2.UpdateLabels{
				Strategy: &apiv2.UpdateLabels_Replace{
					Replace: r.Meta.Labels,
				},
			},
			UpdateMeta: &apiv2.UpdateMeta{
				LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_CLIENT,
				UpdatedAt:       r.Meta.UpdatedAt,
			},
		}, nil
}

func (c *adminNetwork) Update(rq *adminv2.NetworkServiceUpdateRequest) (*apiv2.Network, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Network().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update network: %w", err)
	}

	return resp.Network, nil
}
