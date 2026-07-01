package v2

import (
	"fmt"

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

type size struct {
	c *config.Config
}

func newSizeCmd(c *config.Config) *cobra.Command {
	w := &size{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.SizeServiceCreateRequest, *adminv2.SizeServiceUpdateRequest, *apiv2.Size]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "size",
		Plural:          "sizes",
		Description:     "manage sizes which defines the cpu, gpu, memory and storage properties of machines",
		Sorter:          sorters.SizeSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.SizeListCompletion,
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("id", "", "", "size id to filter for")
			cmd.Flags().StringP("name", "", "", "size name to filter for")
			cmd.Flags().StringP("description", "", "", "size description to filter for")
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("id", "", "", "size id to update")
			cmd.Flags().StringP("name", "", "", "size name to update")
			cmd.Flags().StringP("description", "", "", "size description to update")
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {

		},
		// TODO: create from CLI might be nice to have?
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *size) Get(id string) (*apiv2.Size, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.SizeServiceGetRequest{Id: id}

	resp, err := c.c.Client.Apiv2().Size().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}

	return resp.Size, nil
}

func (c *size) Create(rq *adminv2.SizeServiceCreateRequest) (*apiv2.Size, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Size().Create(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to create size: %w", err)
	}

	return resp.Size, nil
}

func (c *size) Delete(id string) (*apiv2.Size, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Size().Delete(ctx, &adminv2.SizeServiceDeleteRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("failed to delete size: %w", err)
	}

	return resp.Size, nil
}

func (c *size) List() ([]*apiv2.Size, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.SizeServiceListRequest{Query: &apiv2.SizeQuery{
		Id:          pointer.PointerOrNil(viper.GetString("id")),
		Name:        pointer.PointerOrNil(viper.GetString("name")),
		Description: pointer.PointerOrNil(viper.GetString("description")),
	}}

	resp, err := c.c.Client.Apiv2().Size().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get sizes: %w", err)
	}

	return resp.Sizes, nil
}

func (c *size) Update(rq *adminv2.SizeServiceUpdateRequest) (*apiv2.Size, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Size().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update size: %w", err)
	}

	return resp.Size, nil
}

func (c *size) Convert(r *apiv2.Size) (string, *adminv2.SizeServiceCreateRequest, *adminv2.SizeServiceUpdateRequest, error) {
	return r.Id, &adminv2.SizeServiceCreateRequest{
			Size: &apiv2.Size{
				Id:          r.Id,
				Name:        r.Name,
				Description: r.Description,
				Meta:        r.Meta,
				Constraints: r.Constraints,
			},
		}, &adminv2.SizeServiceUpdateRequest{
			Id:          r.Id,
			Name:        r.Name,
			Description: r.Description,
			Constraints: r.Constraints,
			UpdateMeta: &apiv2.UpdateMeta{
				LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_CLIENT,
				UpdatedAt:       r.Meta.UpdatedAt,
			},
			Labels: &apiv2.UpdateLabels{},
		}, nil
}
