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

type adminMachine struct {
	c *config.Config
}

func newAdminMachineCmd(c *config.Config) *cobra.Command {
	w := &adminMachine{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.Machine]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "machine",
		Plural:          "machines",
		Description:     "manage machines",
		Sorter:          sorters.MachineSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.MachineListCompletion,
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd, genericcli.DeleteCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("partition", "", "", "partition to filter for")
			cmd.Flags().String("id", "", "machine id to filter for")
			cmd.Flags().String("size", "", "size to filter for")
			cmd.Flags().String("rack", "", "rack to filter for")
		},
	}

	setStateCmd := &cobra.Command{
		Use:   "set-state <id>",
		Short: "set the state of a machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.setState(args)
		},
		ValidArgsFunction: c.Completion.MachineListCompletion,
	}

	setStateCmd.Flags().String("state", "", "the state to set (e.g. AVAILABLE, LOCKED, TAINTED)")
	setStateCmd.Flags().String("description", "", "description why this machine state was set")

	consolePasswordCmd := &cobra.Command{
		Use:   "console-password <id>",
		Short: "get the console password of a machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.consolePassword(args)
		},
		ValidArgsFunction: c.Completion.MachineListCompletion,
	}

	issuesCmd := &cobra.Command{
		Use:   "issues",
		Short: "list machines with issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.issues()
		},
	}

	issuesCmd.Flags().String("partition", "", "partition to filter for")
	issuesCmd.Flags().String("machine-id", "", "machine id to filter for")

	bmcCmd := &cobra.Command{
		Use:   "bmc",
		Short: "manage machine BMC",
	}

	bmcGetCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "get BMC details of a machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.bmcGet(args)
		},
		ValidArgsFunction: c.Completion.MachineListCompletion,
	}

	bmcListCmd := &cobra.Command{
		Use:   "list",
		Short: "list BMC details of many machines",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.bmcList()
		},
	}

	bmcCmd.AddCommand(bmcGetCmd, bmcListCmd)

	return genericcli.NewCmds(cmdsConfig, setStateCmd, consolePasswordCmd, issuesCmd, bmcCmd)
}

func (c *adminMachine) Get(id string) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.MachineServiceGetRequest{Uuid: id}

	resp, err := c.c.Client.Adminv2().Machine().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine: %w", err)
	}

	return resp.Machine, nil
}

func (c *adminMachine) List() ([]*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.MachineServiceListRequest{
		Query: &apiv2.MachineQuery{
			Uuid:      pointer.PointerOrNil(viper.GetString("id")),
			Partition: pointer.PointerOrNil(viper.GetString("partition")),
			Size:      pointer.PointerOrNil(viper.GetString("size")),
			Rack:      pointer.PointerOrNil(viper.GetString("rack")),
		},
	}

	resp, err := c.c.Client.Adminv2().Machine().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list machines: %w", err)
	}

	return resp.Machines, nil
}

func (c *adminMachine) Delete(id string) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Machine().Delete(ctx, &adminv2.MachineServiceDeleteRequest{Uuid: id})
	if err != nil {
		return nil, fmt.Errorf("failed to delete machine: %w", err)
	}

	return resp.Machine, nil
}

func (c *adminMachine) setState(args []string) error {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	stateStr := viper.GetString("state")
	state, ok := apiv2.MachineState_value[stateStr]
	if !ok {
		return fmt.Errorf("invalid machine state: %s", stateStr)
	}

	_, err = c.c.Client.Adminv2().Machine().SetState(ctx, &adminv2.MachineServiceSetStateRequest{
		Uuid:        id,
		State:       apiv2.MachineState(state),
		Description: viper.GetString("description"),
	})
	if err != nil {
		return fmt.Errorf("failed to set machine state: %w", err)
	}

	return nil
}

func (c *adminMachine) consolePassword(args []string) error {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Machine().ConsolePassword(ctx, &adminv2.MachineServiceConsolePasswordRequest{
		Uuid: id,
	})
	if err != nil {
		return fmt.Errorf("failed to get console password: %w", err)
	}

	return c.c.DescribePrinter.Print(resp)
}

func (c *adminMachine) issues() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Machine().Issues(ctx, &adminv2.MachineServiceIssuesRequest{
		Query: &apiv2.MachineIssuesQuery{
			MachineQuery: &apiv2.MachineQuery{
				Uuid:      pointer.PointerOrNil(viper.GetString("machine-id")),
				Partition: pointer.PointerOrNil(viper.GetString("partition")),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to list machine issues: %w", err)
	}

	return c.c.ListPrinter.Print(resp.Issues)
}

func (c *adminMachine) bmcGet(args []string) error {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Machine().GetBMC(ctx, &adminv2.MachineServiceGetBMCRequest{
		Uuid: id,
	})
	if err != nil {
		return fmt.Errorf("failed to get BMC: %w", err)
	}

	return c.c.DescribePrinter.Print(resp)
}

func (c *adminMachine) bmcList() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Machine().ListBMC(ctx, &adminv2.MachineServiceListBMCRequest{})
	if err != nil {
		return fmt.Errorf("failed to list BMC: %w", err)
	}

	return c.c.DescribePrinter.Print(resp)
}

func (c *adminMachine) Create(rq any) (*apiv2.Machine, error) {
	panic("unimplemented")
}

func (c *adminMachine) Convert(r *apiv2.Machine) (string, any, any, error) {
	panic("unimplemented")
}

func (c *adminMachine) Update(rq any) (*apiv2.Machine, error) {
	panic("unimplemented")
}
