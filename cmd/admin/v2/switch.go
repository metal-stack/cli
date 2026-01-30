package v2

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/metal-stack/api/go/enum"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/cli/cmd/tableprinters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/multisort"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type switchCmd struct {
	c *config.Config
}

func newSwitchCmd(c *config.Config) *cobra.Command {
	sw := &switchCmd{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, *adminv2.SwitchServiceUpdateRequest, *apiv2.Switch]{
		GenericCLI: genericcli.NewGenericCLI(sw).WithFS(c.Fs),
		OnlyCmds: genericcli.OnlyCmds(
			genericcli.DescribeCmd,
			genericcli.ListCmd,
			genericcli.UpdateCmd,
			genericcli.DeleteCmd,
			genericcli.EditCmd,
		),
		BinaryName:      config.BinaryName,
		Singular:        "switch",
		Plural:          "switches",
		Description:     "view and manage network switches",
		Aliases:         []string{"sw"},
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		Sorter:          &multisort.Sorter[*apiv2.Switch]{},
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("id", "", "ID of the switch.")
			cmd.Flags().String("name", "", "Name of the switch.")
			cmd.Flags().String("os-vendor", "", "OS vendor of this switch.")
			cmd.Flags().String("os-version", "", "OS version of this switch.")
			cmd.Flags().String("partition", "", "Partition of this switch.")
			cmd.Flags().String("rack", "", "Rack of this switch.")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("id", c.Completion.SwitchListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("name", c.Completion.SwitchNameListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("rack", c.Completion.SwitchRackListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("os-vendor", c.Completion.SwitchOSVendorListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("os-version", c.Completion.SwitchOSVersionListCompletion))
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().Bool("force", false, "forcefully delete the switch accepting the risk that it still has machines connected to it")
		},
	}

	switchConnectedMachinesCmd := &cobra.Command{
		Use:   "connected-machines",
		Short: "shows switches with their connected machines",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sw.switchConnectedMachines()
		},
		Example: "The command will show the machines connected to the switch ports.",
	}

	switchConnectedMachinesCmd.Flags().String("id", "", "ID of the switch.")
	switchConnectedMachinesCmd.Flags().String("name", "", "Name of the switch.")
	switchConnectedMachinesCmd.Flags().String("os-vendor", "", "OS vendor of this switch.")
	switchConnectedMachinesCmd.Flags().String("os-version", "", "OS version of this switch.")
	switchConnectedMachinesCmd.Flags().String("partition", "", "Partition of this switch.")
	switchConnectedMachinesCmd.Flags().String("rack", "", "Rack of this switch.")

	// TODO: add once size and machine completion are implemented
	// switchMachinesCmd.Flags().String("size", "", "Size of the connected machines.")
	// switchMachinesCmd.Flags().String("machine-id", "", "The id of the connected machine, ignores size flag if set.")

	genericcli.Must(switchConnectedMachinesCmd.RegisterFlagCompletionFunc("id", c.Completion.SwitchListCompletion))
	genericcli.Must(switchConnectedMachinesCmd.RegisterFlagCompletionFunc("name", c.Completion.SwitchNameListCompletion))
	genericcli.Must(switchConnectedMachinesCmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionListCompletion))
	genericcli.Must(switchConnectedMachinesCmd.RegisterFlagCompletionFunc("rack", c.Completion.SwitchRackListCompletion))

	// TODO: add once size and machine completion are implemented
	// genericcli.Must(switchMachinesCmd.RegisterFlagCompletionFunc("size", c.Completion.SizeListCompletion))
	// genericcli.Must(switchMachinesCmd.RegisterFlagCompletionFunc("machine-id", c.Completion.MachineListCompletion))

	switchConsoleCmd := &cobra.Command{
		Use:   "console <switchID>",
		Short: "connect to the switch console",
		Long:  "this requires a network connectivity to the ip address of the console server this switch is connected to.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sw.switchConsole(args)
		},
		ValidArgsFunction: c.Completion.SwitchListCompletion,
	}

	switchDetailCmd := &cobra.Command{
		Use:   "detail",
		Short: "switch details",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sw.switchDetail()
		},
		ValidArgsFunction: c.Completion.SwitchListCompletion,
	}

	switchDetailCmd.Flags().String("id", "", "ID of the switch.")
	switchDetailCmd.Flags().String("name", "", "Name of the switch.")
	switchDetailCmd.Flags().String("os-vendor", "", "OS vendor of this switch.")
	switchDetailCmd.Flags().String("os-version", "", "OS version of this switch.")
	switchDetailCmd.Flags().String("partition", "", "Partition of this switch.")
	switchDetailCmd.Flags().String("rack", "", "Rack of this switch.")

	genericcli.Must(switchDetailCmd.RegisterFlagCompletionFunc("id", c.Completion.SwitchListCompletion))
	genericcli.Must(switchDetailCmd.RegisterFlagCompletionFunc("name", c.Completion.SwitchNameListCompletion))
	genericcli.Must(switchDetailCmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionListCompletion))
	genericcli.Must(switchDetailCmd.RegisterFlagCompletionFunc("rack", c.Completion.SwitchRackListCompletion))

	switchMigrateCmd := &cobra.Command{
		Use:               "migrate <oldSwitchID> <newSwitchID>",
		Short:             "migrate machine connections and other configuration from one switch to another",
		ValidArgsFunction: c.Completion.SwitchListCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sw.switchMigrate(args)
		},
	}

	switchPortCmd := &cobra.Command{
		Use:   "port",
		Short: "sets the given switch port state up or down",
	}
	switchPortCmd.PersistentFlags().String("port", "", "the port to be changed.")
	// TODO: implement Completion.SwitchListPorts
	// genericcli.Must(switchPortCmd.RegisterFlagCompletionFunc("port", c.Completion.SwitchListPorts))

	switchPortUpCmd := &cobra.Command{
		Use:   "up <switch ID>",
		Short: "sets the given switch port state up",
		Long:  "sets the port status to UP so the connected machine will be able to connect to the switch.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sw.port(args, apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_UP)
		},
		ValidArgsFunction: c.Completion.SwitchListCompletion,
	}

	switchPortDownCmd := &cobra.Command{
		Use:   "down <switch ID>",
		Short: "sets the given switch port state down",
		Long:  "sets the port status to DOWN so the connected machine will not be able to connect to the switch.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sw.port(args, apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_DOWN)
		},
		ValidArgsFunction: c.Completion.SwitchListCompletion,
	}

	switchPortCmd.AddCommand(switchPortUpCmd, switchPortDownCmd)

	switchReplaceCmd := &cobra.Command{
		Use:   "replace <switchID>",
		Short: "put a leaf switch into replace mode in preparation for physical replacement. For a description of the steps involved see the long help.",
		Long: `Put a leaf switch into replace mode in preparation for physical replacement

Operational steps to replace a switch:

- Put the switch that needs to be replaced in replace mode with this command
- Replace the switch MAC address in the metal-stack deployment configuration
- Make sure that interfaces on the new switch do not get connected to the PXE-bridge immediately by setting the interfaces list of the respective leaf switch to [] in the metal-stack deployment configuration
- Deploy the management servers so that the dhcp servers will serve the right address and DHCP options to the new switch
- Replace the switch physically. Be careful to ensure that the cabling mirrors the remaining leaf exactly because the new switch information will be cloned from the remaining switch! Also make sure to have console access to the switch so you can start and monitor the install process
- If the switch is not in onie install mode but already has an operating system installed, put it into install mode with "sudo onie-select -i -f -v" and reboot it. Now the switch should be provisioned with a management IP from a management server, install itself with the right software image and receive license and ssh keys through ZTP. You can check whether that process has completed successfully with the command "sudo ztp -s". The ZTP state should be disabled and the result should be success.
- Deploy the switch plane and metal-core through metal-stack deployment CI job
- The switch will now register with its metal-api, and the metal-core service will receive the cloned interface and routing information. You can verify successful switch replacement by checking the interface and BGP configuration, and checking the switch status with "metalctlv2 switch ls -o wide"; it should now be operational again`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sw.switchReplace(args)
		},
		ValidArgsFunction: c.Completion.SwitchListCompletion,
	}

	switchSSHCmd := &cobra.Command{
		Use:   "ssh <switchID>",
		Short: "connect to the switch via ssh",
		Long:  "this requires a network connectivity to the management ip address of the switch.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sw.switchSSH(args)
		},
		ValidArgsFunction: c.Completion.SwitchListCompletion,
	}

	return genericcli.NewCmds(cmdsConfig, switchConnectedMachinesCmd, switchConsoleCmd, switchDetailCmd, switchMigrateCmd, switchPortCmd, switchReplaceCmd, switchSSHCmd)
}

func (c *switchCmd) Get(id string) (*apiv2.Switch, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	res, err := c.c.Client.Adminv2().Switch().Get(ctx, &adminv2.SwitchServiceGetRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return res.Switch, nil
}

func (c *switchCmd) List() ([]*apiv2.Switch, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	vendor, err := enum.GetEnum[apiv2.SwitchOSVendor](viper.GetString("os-vendor"))
	if err != nil {
		return nil, err
	}

	res, err := c.c.Client.Adminv2().Switch().List(ctx, &adminv2.SwitchServiceListRequest{
		Query: &apiv2.SwitchQuery{
			Id:        pointer.Pointer(viper.GetString("id")),
			Partition: pointer.Pointer(viper.GetString("partition")),
			Rack:      pointer.Pointer(viper.GetString("rack")),
			Os: &apiv2.SwitchOSQuery{
				Vendor:  &vendor,
				Version: pointer.Pointer(viper.GetString("os-version")),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return res.Switches, nil
}

func (c *switchCmd) Create(rq any) (*apiv2.Switch, error) {
	panic("unimplemented")
}

func (c *switchCmd) Delete(id string) (*apiv2.Switch, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	res, err := c.c.Client.Adminv2().Switch().Delete(ctx, &adminv2.SwitchServiceDeleteRequest{
		Id:    id,
		Force: viper.GetBool("force"),
	})
	if err != nil {
		return nil, err
	}

	return res.Switch, nil
}

func (c *switchCmd) Update(rq *adminv2.SwitchServiceUpdateRequest) (*apiv2.Switch, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	res, err := c.c.Client.Adminv2().Switch().Update(ctx, rq)
	if err != nil {
		return nil, err
	}

	return res.Switch, nil
}

func (c *switchCmd) Convert(sw *apiv2.Switch) (string, any, *adminv2.SwitchServiceUpdateRequest, error) {
	return sw.Id, nil, &adminv2.SwitchServiceUpdateRequest{
		Id:             sw.Id,
		Description:    &sw.Description,
		ReplaceMode:    &sw.ReplaceMode,
		ManagementIp:   &sw.ManagementIp,
		ManagementUser: sw.ManagementUser,
		ConsoleCommand: sw.ConsoleCommand,
		Nics:           sw.Nics,
		Os:             sw.Os,
	}, nil
}

func (c *switchCmd) switchConnectedMachines() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	switches, err := c.List()
	if err != nil {
		return err
	}

	err = sorters.SwitchSorter().SortBy(switches)
	if err != nil {
		return err
	}

	var (
		id        *string
		partition *string
		rack      *string
		size      *string
	)

	if viper.IsSet("machine-id") {
		id = pointer.Pointer(viper.GetString("machine-id"))
	}
	if viper.IsSet("partition") {
		partition = pointer.Pointer(viper.GetString("partition"))
	}
	if viper.IsSet("rack") {
		rack = pointer.Pointer(viper.GetString("rack"))
	}
	if viper.IsSet("size") {
		size = pointer.Pointer(viper.GetString("size"))
	}

	resp, err := c.c.Client.Adminv2().Machine().List(ctx, &adminv2.MachineServiceListRequest{
		Query: &apiv2.MachineQuery{
			Uuid:      id,
			Partition: partition,
			Size:      size,
			Rack:      rack,
		},
		Partition: partition,
	})
	if err != nil {
		return err
	}

	machines := map[string]*apiv2.Machine{}
	for _, m := range resp.Machines {
		machines[m.Uuid] = m
	}

	return c.c.ListPrinter.Print(&tableprinters.SwitchesWithMachines{
		Switches: switches,
		Machines: machines,
	})
}

func (c *switchCmd) switchConsole(args []string) error {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	resp, err := c.Get(id)
	if err != nil {
		return err
	}

	if resp.ConsoleCommand == nil {
		return fmt.Errorf(`
	unable to connect to console because no console_command was specified for this switch
	You can add a working console_command to every switch with metalctlv2 switch edit
	A sample would look like:

	telnet console-server 7008`)
	}

	parts := strings.Fields(*resp.ConsoleCommand)

	// nolint: gosec
	cmd := exec.Command(parts[0])

	if len(parts) > 1 {
		// nolint: gosec
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

func (c *switchCmd) switchDetail() error {
	switches, err := c.List()
	if err != nil {
		return err
	}

	var switchDetails []tableprinters.SwitchDetail
	for _, sw := range switches {
		switchDetails = append(switchDetails, tableprinters.SwitchDetail{
			Switch: sw,
		})
	}

	return c.c.ListPrinter.Print(switchDetails)
}

func (c *switchCmd) switchMigrate(args []string) error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	if count := len(args); count != 2 {
		return fmt.Errorf("invalid number of arguments were provided; 2 are required, %d were passed", count)
	}

	resp, err := c.c.Client.Adminv2().Switch().Migrate(ctx, &adminv2.SwitchServiceMigrateRequest{
		OldSwitch: args[0],
		NewSwitch: args[1],
	})
	if err != nil {
		return err
	}

	return c.c.DescribePrinter.Print(resp)
}

func (c *switchCmd) port(args []string, status apiv2.SwitchPortStatus) error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	portid := viper.GetString("port")
	if portid == "" {
		return fmt.Errorf("missing port")
	}

	resp, err := c.c.Client.Adminv2().Switch().Port(ctx, &adminv2.SwitchServicePortRequest{
		Id:      id,
		NicName: portid,
		Status:  status,
	})
	if err != nil {
		return err
	}

	return c.dumpPortState(resp.Switch, portid)
}

func (c *switchCmd) switchReplace(args []string) error {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	sw, err := c.Get(id)
	if err != nil {
		return err
	}

	resp, err := c.Update(&adminv2.SwitchServiceUpdateRequest{
		Id: id,
		UpdateMeta: &apiv2.UpdateMeta{
			UpdatedAt:       timestamppb.Now(),
			LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
		},
		Description: &sw.Description,
		ReplaceMode: apiv2.SwitchReplaceMode_SWITCH_REPLACE_MODE_REPLACE.Enum(),
		Os:          sw.Os,
	})
	if err != nil {
		return err
	}

	return c.c.DescribePrinter.Print(resp)
}

func (c *switchCmd) switchSSH(args []string) error {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	resp, err := c.Get(id)
	if err != nil {
		return err
	}
	if resp.ManagementIp == "" || resp.ManagementUser == nil {
		return fmt.Errorf("unable to connect to switch by ssh because no ip and user was stored for this switch, please restart metal-core on this switch")
	}

	// nolint: gosec
	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", pointer.SafeDeref(resp.ManagementUser), resp.ManagementIp))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

func (c *switchCmd) dumpPortState(sw *apiv2.Switch, portid string) error {
	var state currentSwitchPortStateDump

	for _, con := range sw.MachineConnections {
		if pointer.SafeDeref(pointer.SafeDeref(con).Nic).Name == portid {
			state.Actual = con
			break
		}
	}
	for _, desired := range sw.Nics {
		if pointer.SafeDeref(desired).Name == portid {
			state.Desired = desired
			break
		}
	}

	if state.Actual.Nic == nil {
		return fmt.Errorf("no machine connected to port %s on switch %s", portid, sw.Id)
	}

	return c.c.DescribePrinter.Print(state)
}

type currentSwitchPortStateDump struct {
	Actual  *apiv2.MachineConnection `yaml:"actual"`
	Desired *apiv2.SwitchNic         `yaml:"desired"`
}
