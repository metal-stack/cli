package v2

import (
	"fmt"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type machine struct {
	c *config.Config
}

func newMachineCmd(c *config.Config) *cobra.Command {
	w := &machine{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.Machine]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "machine",
		Plural:          "machines",
		Description:     "read machines of the metal cloud",
		Sorter:          sorters.MachineSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.MachineListCompletion,
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project to list machines for")
			cmd.Flags().StringP("partition", "", "", "partition to filter for")
			cmd.Flags().StringP("size", "", "", "size to filter for")
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *machine) Get(id string) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.MachineServiceGetRequest{
		Uuid:    id,
		Project: c.c.GetProject(),
	}

	resp, err := c.c.Client.Apiv2().Machine().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine: %w", err)
	}

	return resp.Machine, nil
}

func (c *machine) List() ([]*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.MachineServiceListRequest{
		Query: &apiv2.MachineQuery{
			Partition: pointer.PointerOrNil(viper.GetString("partition")),
			Size:      pointer.PointerOrNil(viper.GetString("size")),
		},
	}

	resp, err := c.c.Client.Apiv2().Machine().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list machines: %w", err)
	}

	return resp.Machines, nil
}

func (c *machine) Create(rq any) (*apiv2.Machine, error) {
	panic("unimplemented")
}

func (c *machine) Delete(id string) (*apiv2.Machine, error) {
	panic("unimplemented")
}

func (c *machine) Convert(r *apiv2.Machine) (string, any, any, error) {
	panic("unimplemented")
}

func (c *machine) Update(rq any) (*apiv2.Machine, error) {
	panic("unimplemented")
}
