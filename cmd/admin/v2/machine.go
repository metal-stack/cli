package v2

import (
	"context"
	"fmt"
	"net/netip"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/metal-stack/api/go/errorutil"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/cli/pkg/helpers"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	metalssh "github.com/metal-stack/metal-lib/pkg/ssh"
	metalvpn "github.com/metal-stack/metal-lib/pkg/vpn"

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

	cmdsConfig := &genericcli.CmdsConfig[*apiv2.MachineServiceCreateRequest, *apiv2.MachineServiceUpdateRequest, *apiv2.Machine]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "machine",
		Plural:          "machines",
		Description:     "manage machines",
		Sorter:          sorters.MachineSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("id", "", "id of machine which should be listed")
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
			genericcli.Must(cmd.RegisterFlagCompletionFunc("id", c.Completion.AdminMachine))

			cmd.Long = cmd.Short + "\n" + helpers.EmojiHelpText()
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the machine")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
		},
		CreateRequestFromCLI: func() (*apiv2.MachineServiceCreateRequest, error) {
			rq, err := helpers.MachineCreateRequestFromCLI(c)
			if err != nil {
				return nil, err
			}

			if viper.IsSet("id") {
				rq.Uuid = new(viper.GetString("id"))
			}

			return rq, nil
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			helpers.AddMachineCreateFlags(cmd, "machine", c.Completion)

			cmd.Flags().String("id", "", "id of the machine to create. [optional]")
			genericcli.Must(cmd.RegisterFlagCompletionFunc("id", c.Completion.AdminMachine))

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
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Short = "Delete a machine from the database. This can only be done if the machine is offline and dead."
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project from where machines should be listed")
			cmd.Flags().String("description", "", "description of the machine")
			cmd.Flags().StringSlice("labels", nil, "labels to replace for the machine")
			cmd.Flags().StringSlice("add-labels", nil, "labels to add to the machine")
			cmd.Flags().StringSlice("remove-labels", nil, "labels to remove to the machine")
			cmd.Flags().StringP("ssh-public-key", "i", "",
				`SSH public key for access via ssh and console. [optional]
Can be either the public key as string, or pointing to the public key file to use e.g.: "@~/.ssh/id_rsa.pub".
If ~/.ssh/[id_ed25519.pub | id_rsa.pub | id_dsa.pub] is present it will be picked as default, matching the first one in this order.`)

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
		},
		UpdateRequestFromCLI: func(args []string) (*apiv2.MachineServiceUpdateRequest, error) {
			return helpers.MachineUpdateRequestFromCLI(c, args)
		},
		ValidArgsFn: c.Completion.AdminMachine,
	}

	bmcCommandCmd := &cobra.Command{
		Use:   "bmc-command",
		Short: "send a command to the bmc of a machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.bmcCommand(args)
		},
		ValidArgsFunction: c.Completion.AdminMachine,
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
		ValidArgsFunction: c.Completion.AdminMachine,
	}
	lockCmd.Flags().String("description", "", "description of why the machine was locked")
	lockCmd.Flags().Bool("remove", false, "if set to true, machine will be unlocked")

	taintCmd := &cobra.Command{
		Use:   "taint",
		Short: "taint or untaint a machine, e.g. machine will not be automatically selected on machine create, only admins can create them",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.lockOrTaint(args, apiv2.MachineState_MACHINE_STATE_TAINTED)
		},
		ValidArgsFunction: c.Completion.AdminMachine,
	}
	taintCmd.Flags().String("description", "", "description of why the machine was tainted")
	taintCmd.Flags().Bool("remove", false, "if set to true, machine will be untainted")

	consoleCmd := &cobra.Command{
		Use:   "console",
		Short: "establishes a connection to the serial console of a machine. for authentication at the metal-console it uses the token such that no machine ssh key is required for access (unlike the corresponding user API command).",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.console(cmd.Context(), args)
		},
		ValidArgsFunction: c.Completion.AdminMachine,
	}
	consoleCmd.Flags().Bool("ipmi", false, "if set to true, the serial console will be opened using ipmitool (requires ipmitool to be present)")
	consoleCmd.Flags().Int("metal-console-port", 5222, "port open on our control-plane to connect via ssh to get machine console access")

	firewallSSHCmd := &cobra.Command{
		Use:   "ssh <firewall ID>",
		Short: "SSH to a firewall",
		Long:  `SSH to a firewall via VPN.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.firewallSSH(cmd.Context(), args)
		},
		ValidArgsFunction: c.Completion.Firewall,
	}
	firewallSSHCmd.Flags().StringP("identity", "i", "~/.ssh/id_rsa", "specify identity file to SSH to the firewall like: -i path/to/id_rsa")
	firewallSSHCmd.Flags().String("reason", "", "the reason why to connect to the firewall through SSH")

	return genericcli.NewCmds(cmdsConfig, bmcCommandCmd, lockCmd, taintCmd, consoleCmd, firewallSSHCmd)
}

func (c *machine) Create(rq *apiv2.MachineServiceCreateRequest) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Machine().Create(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, err
	}

	return resp.Machine, nil
}

func (c *machine) Delete(id string) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Machine().Delete(ctx, &adminv2.MachineServiceDeleteRequest{
		Uuid: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Machine, nil
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

	if viper.IsSet("hostname") || viper.IsSet("name") || viper.IsSet("project") || viper.IsSet("image") {
		allocation = &apiv2.MachineAllocationQuery{
			Hostname: pointer.PointerOrNil(viper.GetString("hostname")),
			Name:     pointer.PointerOrNil(viper.GetString("name")),
			Project:  pointer.PointerOrNil(viper.GetString("project")),
			Image:    pointer.PointerOrNil(viper.GetString("image")),
		}
	}

	resp, err := c.c.Client.Adminv2().Machine().List(ctx, &adminv2.MachineServiceListRequest{
		Query: &apiv2.MachineQuery{
			Uuid:       pointer.PointerOrNil(viper.GetString("id")),
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

func (c *machine) Update(rq *apiv2.MachineServiceUpdateRequest) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Machine().Update(ctx, rq)
	if err != nil {
		return nil, err
	}

	return resp.Machine, nil
}

func (c *machine) Convert(r *apiv2.Machine) (string, *apiv2.MachineServiceCreateRequest, *apiv2.MachineServiceUpdateRequest, error) {
	update, err := helpers.MachineResponseToUpdate(r)
	if err != nil {
		return "", nil, nil, err
	}

	create, err := helpers.MachineResponseToCreate(r)
	if err != nil {
		return "", nil, nil, err
	}

	if r.Uuid != "" {
		create.Uuid = &r.Uuid
		create.Partition = nil
		create.Size = nil
	}

	return r.Uuid, create, update, err
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

func (c *machine) console(ctx context.Context, args []string) error {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	useIpmi := viper.GetBool("ipmi")
	if useIpmi {
		return c.impitool(ctx, id)
	}

	parsedurl, err := url.Parse(pointer.SafeDeref(c.c.Context.ApiURL))
	if err != nil {
		return err
	}

	err = sshClient(id, viper.GetString("sshidentity"), parsedurl.Host, viper.GetInt("metal-console-port"), &c.c.Context.Token, true)
	if err != nil {
		return fmt.Errorf("machine console error:%w", err)
	}

	return nil
}

func (c *machine) impitool(ctx context.Context, id string) error {
	path, err := exec.LookPath("ipmitool")
	if err != nil {
		return fmt.Errorf("unable to locate ipmitool in path")
	}

	resp, err := c.c.Client.Adminv2().Machine().GetBMC(context.Background(), &adminv2.MachineServiceGetBMCRequest{
		Uuid: id,
	})
	if err != nil {
		return err
	}

	bmc := resp.Bmc.Bmc
	intf := "lanplus"

	// -I lanplus  -H 192.168.2.19 -U ADMIN -P ADMIN sol activate
	hostAndPort := strings.Split(bmc.Address, ":")
	if len(hostAndPort) < 2 {
		hostAndPort = append(hostAndPort, "623")
	}
	usr := bmc.User
	if bmc.User == "" {
		_, _ = fmt.Fprintf(c.c.Out, "no ipmi user stored, please specify with --ipmiuser\n")
	}
	ipmiuser := viper.GetString("ipmiuser")
	if ipmiuser != "" {
		usr = ipmiuser
	}
	password := bmc.Password
	if bmc.Password == "" {
		_, _ = fmt.Fprintf(c.c.Out, "no ipmi password stored, please specify with --ipmipassword\n")
	}

	bmcpassword := viper.GetString("ipmipassword")
	if bmcpassword != "" {
		password = bmcpassword
	}

	err = os.Setenv("IPMITOOL_PASSWORD", password)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Unsetenv("IPMITOOL_PASSWORD")
	}()

	args := []string{"-I", intf, "-H", hostAndPort[0], "-p", hostAndPort[1], "-U", usr, "-E", "sol", "activate"}
	_, _ = fmt.Fprintf(c.c.Out, "connecting to console with:\n%s %s\nExit with ~.\n\n", path, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	return cmd.Run()
}

func (c *machine) firewallSSH(ctx context.Context, args []string) (err error) {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	machine, err := c.Get(id)
	if err != nil {
		return fmt.Errorf("failed to find firewall: %w", err)
	}

	if machine.Allocation == nil {
		return fmt.Errorf("firewall allocation is nil")
	}

	if machine.Allocation.AllocationType != apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_FIREWALL {
		return fmt.Errorf("ssh can only be used for connecting to firewalls")
	}

	projectID := machine.Allocation.Project
	_, _ = fmt.Fprintf(c.c.Out, "accessing firewall through vpn ")
	authKeyResp, err := c.c.Client.Adminv2().VPN().AuthKey(ctx, &adminv2.VPNServiceAuthKeyRequest{
		Project:   projectID,
		Ephemeral: true,
		Reason:    viper.GetString("reason"),
	})
	if err != nil {
		return fmt.Errorf("failed to get VPN auth key: %w", err)
	}

	var vpnopts = []metalvpn.ConnectOpt{}

	if machine.Allocation.Vpn != nil {
		for _, ip := range machine.Allocation.Vpn.Ips {
			parsed, err := netip.ParseAddr(ip)
			if err != nil {
				return err
			}
			if parsed.Is4() {
				vpnopts = append(vpnopts, metalvpn.ConnectOptWithVpnIPAddress(ip))
			}
		}
	}

	v, err := metalvpn.Connect(ctx, machine.Uuid, authKeyResp.Address, authKeyResp.AuthKey, vpnopts...)
	if err != nil {
		return err
	}
	defer func() {
		_ = v.Close()
	}()

	privateKeyFile := viper.GetString("identity")
	if strings.HasPrefix(privateKeyFile, "~/") {
		home, _ := os.UserHomeDir()
		privateKeyFile = filepath.Join(home, privateKeyFile[2:])
	}

	privateKey, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return err
	}

	opts := []metalssh.ConnectOpt{metalssh.ConnectOptOutputPrivateKey(privateKey)}

	s, err := metalssh.NewClientWithConnection("metal", v.TargetIP, v.Conn, opts...)
	if err != nil {
		return err
	}
	return s.Connect(nil)
}

// sshClient opens an interactive ssh session to the host on port with user, authenticated by the key.
func sshClient(user, keyfile, host string, port int, idToken *string, passwordAuth bool) error {
	var opts []metalssh.ConnectOpt

	if passwordAuth {
		opts = append(opts, metalssh.ConnectOptOutputPassword(*idToken))
	} else {
		if keyfile == "" {
			var err error
			keyfile, err = helpers.SearchSSHKey()
			if err != nil {
				return err
			}
		}

		privateKey, err := os.ReadFile(keyfile)
		if err != nil {
			return err
		}

		opts = append(opts, metalssh.ConnectOptOutputPrivateKey(privateKey))
	}

	s, err := metalssh.NewClient(user, host, port, opts...)
	if err != nil {
		return err
	}

	var env *metalssh.Env

	if idToken != nil {
		env = &metalssh.Env{"LC_METAL_STACK_OIDC_TOKEN": *idToken}
	}

	return s.Connect(env)
}
