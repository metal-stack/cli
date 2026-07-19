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

type adminSizeReservation struct {
	c *config.Config
}

func newAdminSizeReservationCmd(c *config.Config) *cobra.Command {
	w := &adminSizeReservation{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, *adminv2.SizeReservationServiceUpdateRequest, *apiv2.SizeReservation]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "size-reservation",
		Plural:          "size-reservations",
		Description:     "manage size reservations",
		Sorter:          sorters.SizeReservationSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.AdminSizeReservationListCompletion,
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd, genericcli.DeleteCmd),
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *adminSizeReservation) Get(id string) (*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.SizeReservationServiceGetRequest{Id: id}

	resp, err := c.c.Client.Apiv2().SizeReservation().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get size reservation: %w", err)
	}

	return resp.SizeReservation, nil
}

func (c *adminSizeReservation) List() ([]*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.SizeReservationServiceListRequest{}

	resp, err := c.c.Client.Adminv2().SizeReservation().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list size reservations: %w", err)
	}

	return resp.SizeReservations, nil
}

func (c *adminSizeReservation) Create(rq any) (*apiv2.SizeReservation, error) {
	panic("unimplemented")
}

func (c *adminSizeReservation) Delete(id string) (*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().SizeReservation().Delete(ctx, &adminv2.SizeReservationServiceDeleteRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("failed to delete size reservation: %w", err)
	}

	return resp.SizeReservation, nil
}

func (c *adminSizeReservation) Convert(r *apiv2.SizeReservation) (string, any, *adminv2.SizeReservationServiceUpdateRequest, error) {
	return r.Id, nil, &adminv2.SizeReservationServiceUpdateRequest{
		Id:          r.Id,
		Name:        &r.Name,
		Description: &r.Description,
		Partitions:  r.Partitions,
		Amount:      &r.Amount,
	}, nil
}

func (c *adminSizeReservation) Update(rq *adminv2.SizeReservationServiceUpdateRequest) (*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().SizeReservation().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update size reservation: %w", err)
	}

	return resp.SizeReservation, nil
}
