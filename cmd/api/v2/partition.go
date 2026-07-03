package v2

import (
	"fmt"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
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
		Description:     "list and get partitions",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("id", "", "", "image id to filter for")
		},
	}

	return genericcli.NewCmds(cmdsConfig)
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
