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
		Description:     "manage machines",
		Sorter:          sorters.MachineSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("uuid", "", "allocation uuid of machine which should be listed")
			cmd.Flags().String("name", "", "name from machines which should be listed")
			cmd.Flags().String("hostname", "", "hostname from machines which should be listed")
			cmd.Flags().String("size", "", "size from machines which should be listed")
			cmd.Flags().String("image", "", "image")
			cmd.Flags().StringP("project", "p", "", "project from where machines should be listed")
			cmd.Flags().StringP("partition", "", "", "partition from where machines should be listed")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("size", c.Completion.Size))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("image", c.Completion.Image))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.Partition))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the machine")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			// w.addMachineCreateFlags(cmd, "machine")
			cmd.Aliases = []string{"allocate"}
			cmd.Example = `machine create can be done in two different ways:

- default with automatic allocation:

	metalctl machine create \
		--hostname worker01 \
		--name worker \
		--image ubuntu-18.04 \ # query available with: metalctl image list
		--size t1-small-x86 \  # query available with: metalctl size list
		--partition test \     # query available with: metalctl partition list
		--project cluster01 \
		--sshpublickey "@~/.ssh/id_rsa.pub"

- for metal administration with reserved machines:

	reserve a machine you want to allocate:

	metalctl machine reserve 00000000-0000-0000-0000-0cc47ae54694 --description "blocked for maintenance"

	allocate this machine:

	metalctl machine create \
		--hostname worker01 \
		--name worker \
		--image ubuntu-18.04 \ # query available with: metalctl image list
		--project cluster01 \
		--sshpublickey "@~/.ssh/id_rsa.pub" \
		--id 00000000-0000-0000-0000-0cc47ae54694

after you do not want to use this machine exclusive, remove the reservation:

metalctl machine reserve 00000000-0000-0000-0000-0cc47ae54694 --remove

Once created the machine installation can not be modified anymore.
`
		},
		ValidArgsFn: c.Completion.Machine,
	}

	bmcCommandCmd := &cobra.Command{
		Use:   "bmc-command",
		Short: "send a command to the bmc of a machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.bmcCommand(args)
		},
		ValidArgsFunction: c.Completion.Machine,
	}
	bmcCommandCmd.Flags().String("command", "", "the actual command to send to the machine")
	genericcli.Must(bmcCommandCmd.RegisterFlagCompletionFunc("command", c.Completion.BMCCommands))
	genericcli.Must(bmcCommandCmd.MarkFlagRequired("command"))

	lockCmd := &cobra.Command{
		Use:   "lock",
		Short: "lock or unlock a machine, e.g. machine cannot be used",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.lockOrTaint(args, apiv2.MachineState_MACHINE_STATE_LOCKED)
		},
		ValidArgsFunction: c.Completion.Machine,
	}
	lockCmd.Flags().String("description", "", "description of why the machine was locked")
	lockCmd.Flags().Bool("remove", false, "if set to true, machine will be unlocked")

	taintCmd := &cobra.Command{
		Use:   "taint",
		Short: "taint or untaint a machine, e.g. machine will not be automatically selected on machine create, only admins can create them",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.lockOrTaint(args, apiv2.MachineState_MACHINE_STATE_TAINTED)
		},
		ValidArgsFunction: c.Completion.Machine,
	}
	taintCmd.Flags().String("description", "", "description of why the machine was tainted")
	taintCmd.Flags().Bool("remove", false, "if set to true, machine will be untainted")

	consoleCmd := &cobra.Command{
		Use:   "console",
		Short: "establishes a connection to the serial console of a machine. for authentication at the metal-console it uses the token such that no machine ssh key is required for access (unlike the corresponding user API command).",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.console(args)
		},
		ValidArgsFunction: c.Completion.Machine,
	}
	consoleCmd.Flags().Bool("ipmi", false, "if set to true, the serial console will be opened using ipmitool (requires ipmitool to be present)")

	return genericcli.NewCmds(cmdsConfig, bmcCommandCmd, lockCmd, taintCmd, consoleCmd)
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

	var allocation *apiv2.MachineAllocationQuery

	if viper.IsSet("hostname") || viper.IsSet("project") || viper.IsSet("image") {
		allocation = &apiv2.MachineAllocationQuery{
			Hostname: pointer.PointerOrNil(viper.GetString("hostname")),
			Project:  pointer.PointerOrNil(viper.GetString("project")),
			Image:    pointer.PointerOrNil(viper.GetString("image")),
		}
	}

	resp, err := c.c.Client.Adminv2().Machine().List(ctx, &adminv2.MachineServiceListRequest{
		Query: &apiv2.MachineQuery{
			Uuid:       pointer.PointerOrNil(viper.GetString("id")),
			Name:       pointer.PointerOrNil(viper.GetString("name")),
			Partition:  pointer.PointerOrNil(viper.GetString("partition")),
			Size:       pointer.PointerOrNil(viper.GetString("size")),
			Allocation: allocation,
			// 	Rack:      pointer.PointerOrNil(viper.GetString("rack")),
			// 	Labels: &apiv2.Labels{
			// 		Labels: tag.NewTagMap(viper.GetStringSlice("labels")),
			// 	},
			// 	Bmc: &apiv2.MachineBMCQuery{
			// 		Address:   pointer.PointerOrNil(viper.GetString("bmc-address")),
			// 		Mac:       pointer.PointerOrNil(viper.GetString("bmc-mac")),
			// 		User:      pointer.PointerOrNil(viper.GetString("bmc-user")),
			// 		Interface: pointer.PointerOrNil(viper.GetString("bmc-interface")),
			// 	},
			// 	Fru: &apiv2.MachineFRUQuery{
			// 		ChassisPartNumber:   pointer.PointerOrNil(viper.GetString("chassis-part-number")),
			// 		ChassisPartSerial:   pointer.PointerOrNil(viper.GetString("chassis-part-serial")),
			// 		BoardMfg:            pointer.PointerOrNil(viper.GetString("board-mfg")),
			// 		BoardSerial:         pointer.PointerOrNil(viper.GetString("board-serial")),
			// 		BoardPartNumber:     pointer.PointerOrNil(viper.GetString("board-part-number")),
			// 		ProductManufacturer: pointer.PointerOrNil(viper.GetString("product-manufacturer")),
			// 		ProductPartNumber:   pointer.PointerOrNil(viper.GetString("product-part-number")),
			// 		ProductSerial:       pointer.PointerOrNil(viper.GetString("product-serial")),
			// 	},
			// 	Hardware: &apiv2.MachineHardwareQuery{
			// 		Memory:   pointer.PointerOrNil(viper.GetUint64("memory")),
			// 		CpuCores: pointer.PointerOrNil(viper.GetUint32("cpu-cores")),
			// 	},
			// State:    &0,
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

func (c *machine) lockOrTaint(args []string, state apiv2.MachineState) error {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	if viper.GetBool("remove") {
		state = apiv2.MachineState_MACHINE_STATE_AVAILABLE
	}

	resp, err := c.c.Client.Adminv2().Machine().SetState(ctx, &adminv2.MachineServiceSetStateRequest{
		Uuid:        id,
		Description: viper.GetString("description"),
		State:       state,
	})
	if err != nil {
		return err
	}

	return c.c.ListPrinter.Print(resp.Machine)
}

func (c *machine) bmcCommand(args []string) error {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	commandString := viper.GetString("command")

	cmd, ok := apiv2.MachineBMCCommand_value[commandString]
	if !ok {
		return fmt.Errorf("unknown bmc command: %s", commandString)
	}

	_, err = c.c.Client.Adminv2().Machine().BMCCommand(ctx, &adminv2.MachineServiceBMCCommandRequest{
		Uuid:    id,
		Command: apiv2.MachineBMCCommand(cmd),
	})
	if err != nil {
		return err
	}

	return err
}

func (c *machine) console(args []string) error {
	panic("unimplemented")
}
