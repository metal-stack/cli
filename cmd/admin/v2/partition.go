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

type adminPartition struct {
	c *config.Config
}

func newAdminPartitionCmd(c *config.Config) *cobra.Command {
	w := &adminPartition{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, *adminv2.PartitionServiceUpdateRequest, *apiv2.Partition]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "partition",
		Plural:          "partitions",
		Description:     "manage partitions (failure domains)",
		Sorter:          sorters.PartitionSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.PartitionListCompletion,
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd, genericcli.DeleteCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("id", "", "", "partition id to filter for")
		},
	}

	capacityCmd := &cobra.Command{
		Use:   "capacity",
		Short: "show partition capacity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.capacity()
		},
	}

	capacityCmd.Flags().String("id", "", "partition id to filter for")
	capacityCmd.Flags().String("size", "", "size to filter for")
	capacityCmd.Flags().String("project", "", "project to filter for")

	return genericcli.NewCmds(cmdsConfig, capacityCmd)
}

func (c *adminPartition) capacity() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.PartitionServiceCapacityRequest{
		Id:      pointer.PointerOrNil(viper.GetString("id")),
		Size:    pointer.PointerOrNil(viper.GetString("size")),
		Project: pointer.PointerOrNil(viper.GetString("project")),
	}

	resp, err := c.c.Client.Adminv2().Partition().Capacity(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get partition capacity: %w", err)
	}

	return c.c.ListPrinter.Print(resp.PartitionCapacity)
}

func (c *adminPartition) Get(id string) (*apiv2.Partition, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.PartitionServiceGetRequest{Id: id}

	resp, err := c.c.Client.Apiv2().Partition().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get partition: %w", err)
	}

	return resp.Partition, nil
}

func (c *adminPartition) List() ([]*apiv2.Partition, error) {
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

func (c *adminPartition) Create(rq any) (*apiv2.Partition, error) {
	panic("unimplemented")
}

func (c *adminPartition) Delete(id string) (*apiv2.Partition, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Partition().Delete(ctx, &adminv2.PartitionServiceDeleteRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("failed to delete partition: %w", err)
	}

	return resp.Partition, nil
}

func (c *adminPartition) Convert(r *apiv2.Partition) (string, any, *adminv2.PartitionServiceUpdateRequest, error) {
	return r.Id, nil, &adminv2.PartitionServiceUpdateRequest{
		Id:                   r.Id,
		Description:          &r.Description,
		BootConfiguration:    r.BootConfiguration,
		DnsServers:           r.DnsServers,
		NtpServers:           r.NtpServers,
		MgmtServiceAddresses: r.MgmtServiceAddresses,
	}, nil
}

func (c *adminPartition) Update(rq *adminv2.PartitionServiceUpdateRequest) (*apiv2.Partition, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Partition().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update partition: %w", err)
	}

	return resp.Partition, nil
}
