package v2

import (
	"encoding/base64"
	"fmt"
	"net/netip"
	"os"
	osuser "os/user"
	"path/filepath"
	"strings"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/cli/pkg/helpers"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
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
		BinaryName:           config.BinaryName,
		GenericCLI:           genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:             "machine",
		Plural:               "machines",
		Description:          "an machine of metal-stack.io",
		Sorter:               sorters.MachineSorter(),
		DescribePrinter:      func() printers.Printer { return c.DescribePrinter },
		ListPrinter:          func() printers.Printer { return c.ListPrinter },
		CreateRequestFromCLI: w.createRequestFromCLI,
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			w.addMachineCreateFlags(cmd, "machine")
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

func (c *machine) Create(rq *apiv2.MachineServiceCreateRequest) (*apiv2.Machine, error) {
	var (
		keys           []string
		dnsServers     []*apiv2.DNSServer
		ntpServers     []*apiv2.NTPServer
		allocationType apiv2.MachineAllocationType
		firewallSpec   *apiv2.FirewallSpec
		labels         *apiv2.Labels
	)

	sshPublicKeyArgument := viper.GetString("sshpublickey")
	dnsServersArgument := viper.GetStringSlice("dnsservers")
	ntpServersArgument := viper.GetStringSlice("ntpservers")

	if strings.HasPrefix(sshPublicKeyArgument, "@") {
		var err error
		sshPublicKeyArgument, err = readFromFile(sshPublicKeyArgument[1:])
		if err != nil {
			return nil, err
		}
	}

	if len(sshPublicKeyArgument) == 0 {
		sshKey, err := searchSSHKey()
		if err != nil {
			return nil, err
		}
		sshPublicKey := sshKey + ".pub"
		sshPublicKeyArgument, err = readFromFile(sshPublicKey)
		if err != nil {
			return nil, err
		}
	}

	if sshPublicKeyArgument != "" {
		keys = append(keys, sshPublicKeyArgument)
	}

	userDataArgument := viper.GetString("userdata")
	if strings.HasPrefix(userDataArgument, "@") {
		var err error
		userDataArgument, err = readFromFile(userDataArgument[1:])
		if err != nil {
			return nil, err
		}
	}
	if userDataArgument != "" {
		userDataArgument = base64.StdEncoding.EncodeToString([]byte(userDataArgument))
	}

	possibleNetworks := viper.GetStringSlice("networks")
	networks, err := parseNetworks(possibleNetworks)
	if err != nil {
		return nil, err
	}

	for _, s := range dnsServersArgument {
		dnsServers = append(dnsServers, &apiv2.DNSServer{Ip: s})
	}

	for _, s := range ntpServersArgument {
		ntpServers = append(ntpServers, &apiv2.NTPServer{Address: s})
	}

	allocationType = apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_MACHINE
	if viper.GetString("allocation-type") == "firewall" {
		allocationType = apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_FIREWALL
	}

	for k, v := range viper.GetStringMap("labels") {
		if labels == nil {
			labels = &apiv2.Labels{}
		} else {
			value, ok := v.(string)
			if ok {
				labels.Labels[k] = value
			} else {
				labels.Labels[k] = ""
			}
		}
	}

	var filesystemlayout *string
	if viper.IsSet("filesystemlayout") {
		filesystemlayout = new(viper.GetString("filesystemlayout"))
	}
	var size *string
	if viper.IsSet("size") {
		size = new(viper.GetString("size"))
	}
	var uuid *string
	if viper.IsSet("id") {
		uuid = new(viper.GetString("id"))
	}
	var partition *string
	if viper.IsSet("partition") {
		partition = new(viper.GetString("partition"))
	}
	var hostname *string
	if viper.IsSet("hostname") {
		hostname = new(viper.GetString("hostname"))
	}
	var description *string
	if viper.IsSet("description") {
		description = new(viper.GetString("description"))
	}

	mcr := &apiv2.MachineServiceCreateRequest{
		Description:      description,
		Partition:        partition,
		Hostname:         hostname,
		Image:            viper.GetString("image"),
		Name:             viper.GetString("name"),
		Uuid:             uuid,
		Project:          viper.GetString("project"),
		Size:             size,
		SshPublicKeys:    keys,
		Labels:           labels,
		Userdata:         new(userDataArgument),
		Networks:         networks,
		DnsServers:       dnsServers,
		NtpServers:       ntpServers,
		FilesystemLayout: filesystemlayout,
		PlacementTags:    viper.GetStringSlice("placement-tags"),
		AllocationType:   allocationType,
		FirewallSpec:     firewallSpec,
	}
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()
	resp, err := c.c.Client.Apiv2().Machine().Create(ctx, mcr)
	if err != nil {
		return nil, fmt.Errorf("unable to create machine:%w", err)
	}
	return resp.Machine, nil
}

func (c *machine) createRequestFromCLI() (*apiv2.MachineServiceCreateRequest, error) {
	var (
		keys           []string
		dnsServers     []*apiv2.DNSServer
		ntpServers     []*apiv2.NTPServer
		allocationType apiv2.MachineAllocationType
		firewallSpec   *apiv2.FirewallSpec
		labels         *apiv2.Labels
	)

	sshPublicKeyArgument := viper.GetString("sshpublickey")
	dnsServersArgument := viper.GetStringSlice("dnsservers")
	ntpServersArgument := viper.GetStringSlice("ntpservers")

	if strings.HasPrefix(sshPublicKeyArgument, "@") {
		var err error
		sshPublicKeyArgument, err = readFromFile(sshPublicKeyArgument[1:])
		if err != nil {
			return nil, err
		}
	}

	if len(sshPublicKeyArgument) == 0 {
		sshKey, err := searchSSHKey()
		if err != nil {
			return nil, err
		}
		sshPublicKey := sshKey + ".pub"
		sshPublicKeyArgument, err = readFromFile(sshPublicKey)
		if err != nil {
			return nil, err
		}
	}

	if sshPublicKeyArgument != "" {
		keys = append(keys, sshPublicKeyArgument)
	}

	userDataArgument := viper.GetString("userdata")
	if strings.HasPrefix(userDataArgument, "@") {
		var err error
		userDataArgument, err = readFromFile(userDataArgument[1:])
		if err != nil {
			return nil, err
		}
	}
	if userDataArgument != "" {
		userDataArgument = base64.StdEncoding.EncodeToString([]byte(userDataArgument))
	}

	possibleNetworks := viper.GetStringSlice("networks")
	networks, err := parseNetworks(possibleNetworks)
	if err != nil {
		return nil, err
	}

	for _, s := range dnsServersArgument {
		dnsServers = append(dnsServers, &apiv2.DNSServer{Ip: s})
	}

	for _, s := range ntpServersArgument {
		ntpServers = append(ntpServers, &apiv2.NTPServer{Address: s})
	}

	allocationType = apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_MACHINE
	if viper.GetString("allocation-type") == "firewall" {
		allocationType = apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_FIREWALL
	}

	for k, v := range viper.GetStringMap("labels") {
		if labels == nil {
			labels = &apiv2.Labels{}
		} else {
			value, ok := v.(string)
			if ok {
				labels.Labels[k] = value
			} else {
				labels.Labels[k] = ""
			}
		}
	}

	var filesystemlayout *string
	if viper.IsSet("filesystemlayout") {
		filesystemlayout = new(viper.GetString("filesystemlayout"))
	}
	var size *string
	if viper.IsSet("size") {
		size = new(viper.GetString("size"))
	}
	var uuid *string
	if viper.IsSet("id") {
		uuid = new(viper.GetString("id"))
	}
	var partition *string
	if viper.IsSet("partition") {
		partition = new(viper.GetString("partition"))
	}
	var hostname *string
	if viper.IsSet("hostname") {
		hostname = new(viper.GetString("hostname"))
	}
	var description *string
	if viper.IsSet("description") {
		description = new(viper.GetString("description"))
	}

	mcr := &apiv2.MachineServiceCreateRequest{
		Description:      description,
		Partition:        partition,
		Hostname:         hostname,
		Image:            viper.GetString("image"),
		Name:             viper.GetString("name"),
		Uuid:             uuid,
		Project:          viper.GetString("project"),
		Size:             size,
		SshPublicKeys:    keys,
		Labels:           labels,
		Userdata:         new(userDataArgument),
		Networks:         networks,
		DnsServers:       dnsServers,
		NtpServers:       ntpServers,
		FilesystemLayout: filesystemlayout,
		PlacementTags:    viper.GetStringSlice("placement-tags"),
		AllocationType:   allocationType,
		FirewallSpec:     firewallSpec,
	}
	return mcr, nil
}

func (c *machine) Delete(id string) (*apiv2.Machine, error) {
	panic("unimplemented")
}

func (c *machine) Get(id string) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Machine().Get(ctx, &apiv2.MachineServiceGetRequest{
		Project: c.c.GetProject(),
		Uuid:    id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Machine, nil
}

func (c *machine) List() ([]*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Machine().List(ctx, &apiv2.MachineServiceListRequest{
		Project: c.c.GetProject(),
		Query:   &apiv2.MachineQuery{
			// FIXME implement
		},
	})
	if err != nil {
		return nil, err
	}

	return resp.Machines, nil
}

func (c *machine) Update(rq *apiv2.MachineServiceUpdateRequest) (*apiv2.Machine, error) {
	panic("unimplemented")
}

func (c *machine) Convert(r *apiv2.Machine) (string, *apiv2.MachineServiceCreateRequest, *apiv2.MachineServiceUpdateRequest, error) {
	responseToUpdate, err := c.MachineResponseToUpdate(r)
	return helpers.EncodeProject(r.Uuid, r.Allocation.Project), c.MachineResponseToCreate(r), responseToUpdate, err
}

func (c *machine) MachineResponseToCreate(r *apiv2.Machine) *apiv2.MachineServiceCreateRequest {
	return &apiv2.MachineServiceCreateRequest{
		// FIXME
	}
}

func (c *machine) MachineResponseToUpdate(desired *apiv2.Machine) (*apiv2.MachineServiceUpdateRequest, error) {
	panic("unimplemented")
}

var defaultSSHKeys = [...]string{"id_ed25519", "id_ecdsa", "id_rsa", "id_dsa"}

func searchSSHKey() (string, error) {
	currentUser, err := osuser.Current()
	if err != nil {
		return "", fmt.Errorf("unable to determine current user for expanding userdata path:%w", err)
	}
	homeDir := currentUser.HomeDir
	defaultDir := filepath.Join(homeDir, "/.ssh/")
	var key string
	for _, k := range defaultSSHKeys {
		possibleKey := filepath.Join(defaultDir, k)
		_, err := os.ReadFile(possibleKey)
		if err == nil {
			fmt.Printf("using SSH identity: %s. Another identity can be specified with --sshidentity/-p\n",
				possibleKey)
			key = possibleKey
			break
		}
	}

	if key == "" {
		return "", fmt.Errorf("failure to locate a SSH identity in default location (%s), "+
			"another identity can be specified with --sshidentity/-p", defaultDir)
	}
	return key, nil
}
func readFromFile(filePath string) (string, error) {
	filePath, err := expandFilepath(filePath)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to read from given file %s error:%w", filePath, err)
	}
	return strings.TrimSpace(string(content)), nil
}

func expandFilepath(filePath string) (string, error) {
	currentUser, err := osuser.Current()
	if err != nil {
		return "", fmt.Errorf("unable to determine current user for expanding userdata path:%w", err)
	}
	homeDir := currentUser.HomeDir

	if filePath == "~" {
		filePath = homeDir
	} else if strings.HasPrefix(filePath, "~/") {
		filePath = filepath.Join(homeDir, filePath[2:])
	}

	return filePath, nil
}
func parseNetworks(possibleNetworks []string) ([]*apiv2.MachineAllocationNetwork, error) {
	var result []*apiv2.MachineAllocationNetwork
	for _, n := range possibleNetworks {
		if n == "" {
			continue
		}
		man := &apiv2.MachineAllocationNetwork{
			Network: n,
		}
		nw, ipsString, found := strings.Cut(n, ":")
		if found {
			man.Network = nw
			for ip := range strings.SplitSeq(ipsString, ",") {
				if ip == "" {
					continue
				}
				_, err := netip.ParseAddr(ip)
				if err != nil {
					return nil, fmt.Errorf("malformed ip:%s %w", ip, err)
				}
				man.Ips = append(man.Ips, ip)
			}
		}
		result = append(result, man)
	}
	return result, nil
}

func (c *machine) addMachineCreateFlags(cmd *cobra.Command, name string) {
	cmd.Flags().StringP("description", "d", "", "Description of the "+name+" to create. [optional]")
	cmd.Flags().StringP("partition", "S", "", "partition/datacenter where the "+name+" is created. [required, except for reserved machines]")
	cmd.Flags().StringP("hostname", "H", "", "Hostname of the "+name+". [required]")
	cmd.Flags().StringP("image", "i", "", "OS Image to install. [required]")
	cmd.Flags().StringP("filesystemlayout", "", "", "Filesystemlayout to use during machine installation. [optional]")
	cmd.Flags().StringP("name", "n", "", "Name of the "+name+". [optional]")
	cmd.Flags().StringP("id", "I", "", "ID of a specific "+name+" to allocate, if given, size and partition are ignored. Need to be set to reserved (--reserve) state before.")
	cmd.Flags().StringP("project", "P", "", "Project where the "+name+" should belong to. [required]")
	cmd.Flags().StringP("size", "s", "", "Size of the "+name+". [required, except for reserved machines]")
	cmd.Flags().StringP("allocation-type", "t", "machine", "allocation type, can be either machine|firewall")
	cmd.Flags().StringP("sshpublickey", "p", "",
		`SSH public key for access via ssh and console. [optional]
Can be either the public key as string, or pointing to the public key file to use e.g.: "@~/.ssh/id_rsa.pub".
If ~/.ssh/[id_ed25519.pub | id_rsa.pub | id_dsa.pub] is present it will be picked as default, matching the first one in this order.`)
	cmd.Flags().StringSlice("tags", []string{}, "tags to add to the "+name+", use it like: --tags \"tag1,tag2\" or --tags \"tag3\".")
	cmd.Flags().StringP("userdata", "", "", `cloud-init.io compatible userdata. [optional]
Can be either the userdata as string, or pointing to the userdata file to use e.g.: "@/tmp/userdata.cfg".`)
	cmd.Flags().StringSlice("dnsservers", []string{}, "dns servers to add to the machine or firewall. [optional]")
	cmd.Flags().StringSlice("ntpservers", []string{}, "ntp servers to add to the machine or firewall. [optional]")

	cmd.Flags().StringSlice("networks", []string{},
		`Adds a network. Usage: [--networks NETWORK,[ip:ip:ip][,NETWORK...]...
NETWORK specifies the name or id of an existing network.
IPs can be added per network colon separated, these ips must be already allocated upfront. If no ip(s) are specified per network, one ip per network is allocated.
`)

	cmd.MarkFlagsMutuallyExclusive("file", "project")
	cmd.MarkFlagsRequiredTogether("project", "networks", "hostname", "image")
	cmd.MarkFlagsRequiredTogether("size", "partition")

	// Completion for arguments
	genericcli.Must(cmd.RegisterFlagCompletionFunc("networks", c.c.Completion.NetworkListCompletion))
	genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.c.Completion.PartitionListCompletion))
	genericcli.Must(cmd.RegisterFlagCompletionFunc("size", c.c.Completion.SizeListCompletion))
	genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.c.Completion.ProjectListCompletion))
	genericcli.Must(cmd.RegisterFlagCompletionFunc("id", c.c.Completion.MachineListCompletion))
	// FIXME implement
	// genericcli.Must(cmd.RegisterFlagCompletionFunc("image", c.c.Completion.ImageListCompletion))
	// genericcli.Must(cmd.RegisterFlagCompletionFunc("filesystemlayout", c.c.Completion.FilesystemLayoutListCompletion))
}
