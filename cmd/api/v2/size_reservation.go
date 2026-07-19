package v2

import (
	"fmt"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type sizeReservation struct {
	c *config.Config
}

func newSizeReservationCmd(c *config.Config) *cobra.Command {
	w := &sizeReservation{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.SizeReservation]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "size-reservation",
		Plural:          "size-reservations",
		Description:     "read size reservations which allow to reserve machine capacity",
		Sorter:          sorters.SizeReservationSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.SizeReservationListCompletion,
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *sizeReservation) Get(id string) (*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.SizeReservationServiceGetRequest{Id: id}

	resp, err := c.c.Client.Apiv2().SizeReservation().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get size reservation: %w", err)
	}

	return resp.SizeReservation, nil
}

func (c *sizeReservation) List() ([]*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.SizeReservationServiceListRequest{}

	resp, err := c.c.Client.Apiv2().SizeReservation().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list size reservations: %w", err)
	}

	return resp.SizeReservations, nil
}

func (c *sizeReservation) Create(rq any) (*apiv2.SizeReservation, error) {
	panic("unimplemented")
}

func (c *sizeReservation) Delete(id string) (*apiv2.SizeReservation, error) {
	panic("unimplemented")
}

func (c *sizeReservation) Convert(r *apiv2.SizeReservation) (string, any, any, error) {
	panic("unimplemented")
}

func (c *sizeReservation) Update(rq any) (*apiv2.SizeReservation, error) {
	panic("unimplemented")
}
