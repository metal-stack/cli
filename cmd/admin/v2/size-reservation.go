package v2

import (
	"fmt"

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

type sizeReservation struct {
	c *config.Config
}

func newSizeReservationCmd(c *config.Config) *cobra.Command {
	w := &sizeReservation{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.SizeReservationServiceCreateRequest, *adminv2.SizeReservationServiceUpdateRequest, *apiv2.SizeReservation]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "size-reservation",
		Plural:          "size-reservations",
		Description:     "manage size reservations which reserve machine capacity for a project",
		Sorter:          sorters.SizeReservationSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.AdminSizeReservation,
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "size reservation name to filter for")
			cmd.Flags().String("description", "", "size reservation description to filter for")
			cmd.Flags().String("size", "", "size id to filter for")
			cmd.Flags().String("project", "", "project id to filter for")
			cmd.Flags().String("partition", "", "partition id to filter for")
			cmd.Flags().StringSlice("labels", nil, "labels to filter for")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("size", c.Completion.Size))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.Partition))
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the size reservation to create")
			cmd.Flags().String("description", "", "the description of the size reservation to create")
			cmd.Flags().String("project", "", "the project id of the size reservation to create")
			cmd.Flags().String("size", "", "the size id of the size reservation to create")
			cmd.Flags().StringSlice("partitions", nil, "the partition ids the size reservation applies to")
			cmd.Flags().Int32("amount", 0, "the amount of machines to reserve")
			cmd.Flags().StringSlice("labels", nil, "labels to add to the size reservation")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("size", c.Completion.Size))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partitions", c.Completion.Partition))
		},
		CreateRequestFromCLI: w.createRequestFromCLI,
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the size reservation")
			cmd.Flags().String("description", "", "the description of the size reservation")
			cmd.Flags().StringSlice("partitions", nil, "the partition ids the size reservation applies to")
			cmd.Flags().Int32("amount", 0, "the amount of machines to reserve")
			cmd.Flags().StringSlice("labels", nil, "labels to replace for the size reservation")
			cmd.Flags().StringSlice("add-labels", nil, "labels to add to the size reservation")
			cmd.Flags().StringSlice("remove-labels", nil, "labels to remove from the size reservation")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("partitions", c.Completion.Partition))
		},
		UpdateRequestFromCLI: w.updateRequestFromCLI,
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *sizeReservation) Get(id string) (*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	// the admin api does not offer a dedicated get call for size reservations,
	// so we filter the list by the given id.
	resp, err := c.c.Client.Adminv2().SizeReservation().List(ctx, &adminv2.SizeReservationServiceListRequest{
		Query: &apiv2.SizeReservationQuery{
			Id: new(id),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get size reservation: %w", err)
	}

	reservations := resp.GetSizeReservations()
	if len(reservations) == 0 {
		return nil, fmt.Errorf("size reservation %q not found", id)
	}

	return reservations[0], nil
}

func (c *sizeReservation) List() ([]*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.SizeReservationServiceListRequest{
		Query: &apiv2.SizeReservationQuery{
			Name:        pointer.PointerOrNil(viper.GetString("name")),
			Description: pointer.PointerOrNil(viper.GetString("description")),
			Size:        pointer.PointerOrNil(viper.GetString("size")),
			Project:     pointer.PointerOrNil(viper.GetString("project")),
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

	resp, err := c.c.Client.Adminv2().SizeReservation().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list size reservations: %w", err)
	}

	return resp.GetSizeReservations(), nil
}

func (c *sizeReservation) Create(rq *adminv2.SizeReservationServiceCreateRequest) (*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().SizeReservation().Create(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, fmt.Errorf("failed to create size reservation: %w", err)
	}

	return resp.GetSizeReservation(), nil
}

func (c *sizeReservation) Update(rq *adminv2.SizeReservationServiceUpdateRequest) (*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().SizeReservation().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update size reservation: %w", err)
	}

	return resp.GetSizeReservation(), nil
}

func (c *sizeReservation) Delete(id string) (*apiv2.SizeReservation, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().SizeReservation().Delete(ctx, &adminv2.SizeReservationServiceDeleteRequest{
		Id: id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete size reservation: %w", err)
	}

	return resp.GetSizeReservation(), nil
}

func (c *sizeReservation) Convert(r *apiv2.SizeReservation) (string, *adminv2.SizeReservationServiceCreateRequest, *adminv2.SizeReservationServiceUpdateRequest, error) {
	return r.Id,
		&adminv2.SizeReservationServiceCreateRequest{
			SizeReservation: r,
		},
		&adminv2.SizeReservationServiceUpdateRequest{
			Id:          r.Id,
			UpdateMeta:  helpers.UpdateMetaFromMeta(r.Meta),
			Name:        pointer.PointerOrNil(r.Name),
			Description: pointer.PointerOrNil(r.Description),
			Partitions:  r.Partitions,
			Amount:      pointer.PointerOrNil(r.Amount),
			Labels:      helpers.UpdateLabelsFromMeta(r.Meta),
		}, nil
}

func (c *sizeReservation) createRequestFromCLI() (*adminv2.SizeReservationServiceCreateRequest, error) {
	reservation := &apiv2.SizeReservation{
		Name:        viper.GetString("name"),
		Description: viper.GetString("description"),
		Project:     viper.GetString("project"),
		Size:        viper.GetString("size"),
		Partitions:  viper.GetStringSlice("partitions"),
		Amount:      viper.GetInt32("amount"),
	}

	labels, err := helpers.LabelsFromSlice(viper.GetStringSlice("labels"))
	if err != nil {
		return nil, err
	}
	if labels != nil {
		reservation.Meta = &apiv2.Meta{Labels: labels}
	}

	return &adminv2.SizeReservationServiceCreateRequest{
		SizeReservation: reservation,
	}, nil
}

func (c *sizeReservation) updateRequestFromCLI(args []string) (*adminv2.SizeReservationServiceUpdateRequest, error) {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	updateLabels, err := helpers.UpdateLabelsFromCLI()
	if err != nil {
		return nil, err
	}

	req := &adminv2.SizeReservationServiceUpdateRequest{
		Id: id,
		UpdateMeta: &apiv2.UpdateMeta{
			LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
		},
		Name:        pointer.PointerOrNil(viper.GetString("name")),
		Description: pointer.PointerOrNil(viper.GetString("description")),
		Amount:      amountFromCLI(),
		Labels:      updateLabels,
	}

	if partitions := viper.GetStringSlice("partitions"); len(partitions) > 0 {
		req.Partitions = partitions
	}

	return req, nil
}

func amountFromCLI() *int32 {
	amount := viper.GetInt32("amount")
	if amount <= 0 {
		return nil
	}

	return new(amount)
}
