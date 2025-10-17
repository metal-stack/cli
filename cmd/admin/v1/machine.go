package v1

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
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
		Description:     "an machine of metal-stack.io",
		Sorter:          sorters.MachineSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project from where machines should be listed")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the machine")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		ValidArgsFn: c.Completion.MachineListCompletion,
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *machine) updateFromCLI(args []string) (any, error) {
	panic("unimplemented")
}

func (c *machine) Create(rq any) (*apiv2.Machine, error) {
	panic("unimplemented")
}

func (c *machine) Delete(id string) (*apiv2.Machine, error) {
	panic("unimplemented")
}

func (c *machine) Get(id string) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Machine().Get(ctx, &adminv2.MachineServiceGetRequest{
		Uuid: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Machine, nil
}

func (c *machine) List() ([]*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Machine().List(ctx, &adminv2.MachineServiceListRequest{
		Query: &apiv2.MachineQuery{
			// FIXME implement
		},
	})
	if err != nil {
		return nil, err
	}

	return resp.Machines, nil
}

func (c *machine) Update(rq any) (*apiv2.Machine, error) {
	panic("unimplemented")
}

func (c *machine) Convert(r *apiv2.Machine) (string, any, any, error) {
	panic("unimplemented")

}

func (c *machine) MachineResponseToCreate(r *apiv2.Machine) any {
	panic("unimplemented")
}

func (c *machine) MachineResponseToUpdate(desired *apiv2.Machine) (any, error) {
	panic("unimplemented")
}
