package v2

import (
	"fmt"
	"strings"

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

type partition struct {
	c *config.Config
}

func newPartitionCmd(c *config.Config) *cobra.Command {
	w := &partition{
		c: c,
	}

	gcli := genericcli.NewGenericCLI(w).WithFS(c.Fs)

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.Partition]{
		BinaryName:      config.BinaryName,
		GenericCLI:      gcli,
		Singular:        "partition",
		Plural:          "partitions",
		Description:     "manage partitions",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.RunE = func(cmd *cobra.Command, args []string) error {
				return gcli.DescribeAndPrint("", w.c.DescribePrinter)
			}
		},
	}

	capacityCmd := &cobra.Command{
		Use:   "capacity",
		Short: "show partition capacity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.capacity()
		},
	}

	capacityCmd.Flags().StringP("id", "", "", "filter on partition id.")
	capacityCmd.Flags().StringP("size", "", "", "filter on size id.")
	capacityCmd.Flags().StringP("project", "", "", "consider project-specific counts, e.g. size reservations.")
	capacityCmd.Flags().StringSlice("sort-by", []string{}, fmt.Sprintf("order by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: %s", strings.Join(sorters.PartitionCapacitySorter().AvailableKeys(), "|")))
	genericcli.Must(capacityCmd.RegisterFlagCompletionFunc("id", c.Completion.PartitionListCompletion))
	genericcli.Must(capacityCmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
	genericcli.Must(capacityCmd.RegisterFlagCompletionFunc("size", c.Completion.SizeListCompletion))
	genericcli.Must(capacityCmd.RegisterFlagCompletionFunc("sort-by", cobra.FixedCompletions(sorters.PartitionCapacitySorter().AvailableKeys(), cobra.ShellCompDirectiveNoFileComp)))

	return genericcli.NewCmds(cmdsConfig, capacityCmd)
}

func (c *partition) capacity() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.PartitionServiceCapacityRequest{}

	if viper.IsSet("id") {
		req.Id = new(viper.GetString("id"))
	}
	if viper.IsSet("size") {
		req.Size = new(viper.GetString("size"))
	}
	if viper.IsSet("project") {
		req.Project = new(viper.GetString("project"))
	}
	resp, err := c.c.Client.Adminv2().Partition().Capacity(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get partition capacity: %w", err)
	}

	err = sorters.PartitionCapacitySorter().SortBy(resp.PartitionCapacity)
	if err != nil {
		return err
	}

	return c.c.ListPrinter.Print(resp.PartitionCapacity)
}

func (c *partition) Get(id string) (*apiv2.Partition, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.PartitionServiceGetRequest{Id: id}

	resp, err := c.c.Client.Apiv2().Partition().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get partition: %w", err)
	}

	return resp.Partition, nil
}

func (c *partition) List() ([]*apiv2.Partition, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.PartitionServiceListRequest{Query: &apiv2.PartitionQuery{
		Id: pointer.PointerOrNil(viper.GetString("id")),
	}}

	resp, err := c.c.Client.Apiv2().Partition().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get partitions: %w", err)
	}

	return resp.Partitions, nil
}

func (c *partition) Create(rq any) (*apiv2.Partition, error) {
	panic("unimplemented")
}

func (c *partition) Delete(id string) (*apiv2.Partition, error) {
	panic("unimplemented")
}

func (t *partition) Convert(r *apiv2.Partition) (string, any, any, error) {
	panic("unimplemented")
}

func (t *partition) Update(rq any) (*apiv2.Partition, error) {
	panic("unimplemented")
}
