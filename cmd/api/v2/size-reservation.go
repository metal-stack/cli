package v2

import (
	"fmt"

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
		Singular:        "reservation",
		Plural:          "reservations",
		Description:     "show size reservations which reserve machine capacity for a project",
		Sorter:          sorters.SizeReservationSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.SizeReservation,
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project from where size reservations should be listed")
			cmd.Flags().String("id", "", "size reservation id to filter for")
			cmd.Flags().String("name", "", "size reservation name to filter for")
			cmd.Flags().String("description", "", "size reservation description to filter for")
			cmd.Flags().String("size", "", "size id to filter for")
			cmd.Flags().String("partition", "", "partition id to filter for")
			cmd.Flags().StringSlice("labels", nil, "labels to filter for")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("size", c.Completion.Size))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.Partition))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the size reservation")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *sizeReservation) Get(id string) (*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().SizeReservation().Get(ctx, &apiv2.SizeReservationServiceGetRequest{
		Id:      id,
		Project: c.c.GetProject(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get size reservation: %w", err)
	}

	return resp.GetSizeReservation(), nil
}

func (c *sizeReservation) List() ([]*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.SizeReservationServiceListRequest{
		Project: c.c.GetProject(),
		Query: &apiv2.SizeReservationQuery{
			Id:          pointer.PointerOrNil(viper.GetString("id")),
			Name:        pointer.PointerOrNil(viper.GetString("name")),
			Description: pointer.PointerOrNil(viper.GetString("description")),
			Size:        pointer.PointerOrNil(viper.GetString("size")),
			Partition:   pointer.PointerOrNil(viper.GetString("partition")),
		},
	}

	if labelSlice := viper.GetStringSlice("labels"); len(labelSlice) > 0 {
		labels, err := helpers.LabelsFromSlice(labelSlice)
		if err != nil {
			return nil, err
		}

		req.Query.Labels = labels
	}

	resp, err := c.c.Client.Apiv2().SizeReservation().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list size reservations: %w", err)
	}

	return resp.GetSizeReservations(), nil
}

func (c *sizeReservation) Create(rq any) (*apiv2.SizeReservation, error) {
	panic("unimplemented")
}

func (c *sizeReservation) Update(rq any) (*apiv2.SizeReservation, error) {
	panic("unimplemented")
}

func (c *sizeReservation) Delete(id string) (*apiv2.SizeReservation, error) {
	panic("unimplemented")
}

func (c *sizeReservation) Convert(r *apiv2.SizeReservation) (string, any, any, error) {
	panic("unimplemented")
}
