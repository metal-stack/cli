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
	"github.com/metal-stack/metal-lib/pkg/tag"
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
			cmd.Flags().StringP("project", "p", "", "project from where machines should be listed")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the machine")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		ValidArgsFn: c.Completion.MachineListCompletion,
	}

	bmcCommandCmd := &cobra.Command{
		Use:   "bmc-command",
		Short: "send a command to the bmc of a machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.bmcCommand()
		},
	}
	bmcCommandCmd.Flags().String("id", "", "id of the machine where the command should be sent to")
	bmcCommandCmd.Flags().String("command", "", "the actual command to send to the machine")
	bmcCommandCmd.RegisterFlagCompletionFunc("id", c.Completion.MachineListCompletion)
	bmcCommandCmd.RegisterFlagCompletionFunc("command", c.Completion.BMCCommandListCompletion)
	genericcli.Must(bmcCommandCmd.MarkFlagRequired("id"))
	genericcli.Must(bmcCommandCmd.MarkFlagRequired("command"))

	return genericcli.NewCmds(cmdsConfig, bmcCommandCmd)
}

func (c *machine) bmcCommand() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	commandString := viper.GetString("command")

	cmd, ok := apiv2.MachineBMCCommand_value[commandString]
	if !ok {
		return fmt.Errorf("unknown command: %s", commandString)
	}
	_, err := c.c.Client.Adminv2().Machine().BMCCommand(ctx, &adminv2.MachineServiceBMCCommandRequest{
		Uuid:    viper.GetString("id"),
		Command: apiv2.MachineBMCCommand(cmd),
	})
	if err != nil {
		return err
	}
	return err
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
			Uuid:      pointer.PointerOrNil(viper.GetString("id")),
			Name:      pointer.PointerOrNil(viper.GetString("name")),
			Partition: pointer.PointerOrNil(viper.GetString("partition")),
			Size:      pointer.PointerOrNil(viper.GetString("size")),
			Rack:      pointer.PointerOrNil(viper.GetString("rack")),
			Labels: &apiv2.Labels{
				Labels: tag.NewTagMap(viper.GetStringSlice("labels")),
			},
			Allocation: &apiv2.MachineAllocationQuery{
				Uuid:             pointer.PointerOrNil(viper.GetString("allocation-uuid")),
				Name:             pointer.PointerOrNil(viper.GetString("allocation-name")),
				Project:          pointer.PointerOrNil(viper.GetString("project")),
				Image:            pointer.PointerOrNil(viper.GetString("image")),
				FilesystemLayout: pointer.PointerOrNil(viper.GetString("file-system-layout-id")),
				Hostname:         pointer.PointerOrNil(viper.GetString("hostname")),
				// AllocationType:   &0,
				Labels: &apiv2.Labels{
					Labels: tag.NewTagMap(viper.GetStringSlice("allocation-labels")),
				},
				// Vpn: &apiv2.MachineVPN{}, these query fields are no pointers and somehow seem wrong? how to search for vpn key?
			},
			Network: &apiv2.MachineNetworkQuery{},
			Nic:     &apiv2.MachineNicQuery{},
			Disk: &apiv2.MachineDiskQuery{
				Names: viper.GetStringSlice("disk-names"),
				// Sizes:
			},
			Bmc: &apiv2.MachineBMCQuery{
				Address:   pointer.PointerOrNil(viper.GetString("bmc-address")),
				Mac:       pointer.PointerOrNil(viper.GetString("bmc-mac")),
				User:      pointer.PointerOrNil(viper.GetString("bmc-user")),
				Interface: pointer.PointerOrNil(viper.GetString("bmc-interface")),
			},
			Fru: &apiv2.MachineFRUQuery{
				ChassisPartNumber:   pointer.PointerOrNil(viper.GetString("chassis-part-number")),
				ChassisPartSerial:   pointer.PointerOrNil(viper.GetString("chassis-part-serial")),
				BoardMfg:            pointer.PointerOrNil(viper.GetString("board-mfg")),
				BoardSerial:         pointer.PointerOrNil(viper.GetString("board-serial")),
				BoardPartNumber:     pointer.PointerOrNil(viper.GetString("board-part-number")),
				ProductManufacturer: pointer.PointerOrNil(viper.GetString("product-manufacturer")),
				ProductPartNumber:   pointer.PointerOrNil(viper.GetString("product-part-number")),
				ProductSerial:       pointer.PointerOrNil(viper.GetString("product-serial")),
			},
			Hardware: &apiv2.MachineHardwareQuery{
				Memory:   pointer.PointerOrNil(viper.GetUint64("memory")),
				CpuCores: pointer.PointerOrNil(viper.GetUint32("cpu-cores")),
			},
			// State:    &0,
		},
		Partition: nil, // again partition?
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
