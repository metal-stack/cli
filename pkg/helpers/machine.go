package helpers

import (
	"encoding/base64"
	"fmt"
	"net/netip"
	"os"
	osuser "os/user"
	"path/filepath"
	"strings"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/completion"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func MachineResponseToCreate(r *apiv2.Machine) (*apiv2.MachineServiceCreateRequest, error) {
	if r.Allocation == nil {
		return nil, fmt.Errorf("allocation is nil")
	}

	var (
		networks     []*apiv2.MachineAllocationNetwork
		firewallSpec *apiv2.FirewallSpec
	)

	for _, nw := range r.Allocation.Networks {
		networks = append(networks, &apiv2.MachineAllocationNetwork{
			Network: nw.Network,
			Ips:     nw.Ips,
		})
	}

	if r.Allocation.AllocationType.Enum() == apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_FIREWALL.Enum() {
		firewallSpec = &apiv2.FirewallSpec{
			FirewallRules: &apiv2.FirewallRules{
				Egress:  r.Allocation.FirewallRules.Egress,
				Ingress: r.Allocation.FirewallRules.Ingress,
			},
		}
	}

	return &apiv2.MachineServiceCreateRequest{
		Project:          r.Allocation.Project,
		Name:             r.Allocation.Name,
		Description:      &r.Allocation.Description,
		Hostname:         &r.Allocation.Hostname,
		Partition:        new(pointer.SafeDeref(r.Partition).Id),
		Size:             new(pointer.SafeDeref(r.Size).Id),
		Image:            r.Allocation.Image.Id,
		FilesystemLayout: pointer.PointerOrNil(r.Allocation.FilesystemLayout.Id),
		SshPublicKeys:    r.Allocation.SshPublicKeys,
		Userdata:         pointer.PointerOrNil(r.Allocation.Userdata),
		Labels:           pointer.SafeDeref(r.Meta).Labels,
		Networks:         networks,
		DnsServers:       r.Allocation.DnsServers,
		NtpServers:       r.Allocation.NtpServers,
		AllocationType:   r.Allocation.AllocationType,
		FirewallSpec:     firewallSpec,
		// PlacementTags:    r.Allocation.Plac, // TODO: should be stored in the allocation to see what was provided
	}, nil
}

func MachineResponseToUpdate(r *apiv2.Machine) (*apiv2.MachineServiceUpdateRequest, error) {
	if r.Allocation == nil {
		return nil, fmt.Errorf("allocation is nil")
	}

	return &apiv2.MachineServiceUpdateRequest{
		Uuid:          r.Uuid,
		UpdateMeta:    UpdateMetaFromMeta(r.Meta),
		Labels:        UpdateLabelsFromMeta(r.Meta),
		Project:       r.Allocation.Project,
		Description:   &r.Allocation.Description,
		SshPublicKeys: r.Allocation.SshPublicKeys,
	}, nil
}

func MachineCreateRequestFromCLI(c *config.Config) (*apiv2.MachineServiceCreateRequest, error) {
	var (
		keys           []string
		dnsServers     []*apiv2.DNSServer
		ntpServers     []*apiv2.NTPServer
		allocationType apiv2.MachineAllocationType
		firewallSpec   *apiv2.FirewallSpec

		sshPublicKeyArgument = viper.GetString("ssh-public-key")
		dnsServersArgument   = viper.GetStringSlice("dns-servers")
		ntpServersArgument   = viper.GetStringSlice("ntp-servers")
	)

	if strings.HasPrefix(sshPublicKeyArgument, "@") {
		var err error
		sshPublicKeyArgument, err = readFromFile(c.Fs, sshPublicKeyArgument[1:])
		if err != nil {
			return nil, err
		}
	}

	if len(sshPublicKeyArgument) == 0 {
		sshKey, err := SearchSSHKey()
		if err != nil {
			return nil, err
		}
		sshPublicKey := sshKey + ".pub"
		sshPublicKeyArgument, err = readFromFile(c.Fs, sshPublicKey)
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
		userDataArgument, err = readFromFile(c.Fs, userDataArgument[1:])
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

	labels, err := LabelsFromSlice(viper.GetStringSlice("labels"))
	if err != nil {
		return nil, err
	}

	var filesystemlayout *string
	if viper.IsSet("filesystem-layout") {
		filesystemlayout = new(viper.GetString("filesystem-layout"))
	}
	var size *string
	if viper.IsSet("size") {
		size = new(viper.GetString("size"))
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

	return &apiv2.MachineServiceCreateRequest{
		Description:      description,
		Partition:        partition,
		Hostname:         hostname,
		Image:            viper.GetString("image"),
		Name:             viper.GetString("name"),
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
	}, nil
}

func MachineUpdateRequestFromCLI(c *config.Config, args []string) (*apiv2.MachineServiceUpdateRequest, error) {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	updateLabels, err := UpdateLabelsFromCLI()
	if err != nil {
		return nil, err
	}

	sshPublicKeyArgument := viper.GetString("ssh-public-key")

	if strings.HasPrefix(sshPublicKeyArgument, "@") {
		var err error
		sshPublicKeyArgument, err = readFromFile(c.Fs, sshPublicKeyArgument[1:])
		if err != nil {
			return nil, err
		}
	}

	return &apiv2.MachineServiceUpdateRequest{
		Uuid: id,
		UpdateMeta: &apiv2.UpdateMeta{
			LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
		},
		Project:       c.GetProject(),
		Description:   pointer.PointerOrNil(viper.GetString("description")),
		Labels:        updateLabels,
		SshPublicKeys: []string{sshPublicKeyArgument},
	}, nil
}

var defaultSSHKeys = [...]string{"id_ed25519", "id_ecdsa", "id_rsa", "id_dsa"}

func SearchSSHKey() (string, error) {
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

func readFromFile(fs *afero.Afero, filePath string) (string, error) {
	filePath, err := expandFilepath(filePath)
	if err != nil {
		return "", err
	}

	content, err := fs.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to read from given file %q: %w", filePath, err)
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
			for ip := range strings.SplitSeq(ipsString, ";") {
				if ip == "" {
					continue
				}
				_, err := netip.ParseAddr(ip)
				if err != nil {
					return nil, fmt.Errorf("malformed ip %q: %w", ip, err)
				}
				man.Ips = append(man.Ips, ip)
			}
		}
		result = append(result, man)
	}

	return result, nil
}

func AddMachineCreateFlags(cmd *cobra.Command, name string, completion *completion.Completion) {
	cmd.Flags().String("description", "", "Description of the "+name+" to create. [optional]")
	cmd.Flags().String("partition", "", "partition/datacenter where the "+name+" is created. [required, except for reserved machines]")
	cmd.Flags().String("hostname", "", "Hostname of the "+name+". [required]")
	cmd.Flags().String("image", "", "OS Image to install. [required]")
	cmd.Flags().String("filesystem-layout", "", "Filesystemlayout to use during machine installation. [optional]")
	cmd.Flags().String("name", "", "Name of the "+name+". [optional]")
	cmd.Flags().StringP("project", "p", "", "Project where the "+name+" should belong to. [required]")
	cmd.Flags().String("size", "", "Size of the "+name+". [required, except for reserved machines]")
	cmd.Flags().String("allocation-type", "machine", "allocation type, can be either machine|firewall")
	cmd.Flags().StringP("ssh-public-key", "i", "",
		`SSH public key for access via ssh and console. [optional]
Can be either the public key as string, or pointing to the public key file to use e.g.: "@~/.ssh/id_rsa.pub".
If ~/.ssh/[id_ed25519.pub | id_rsa.pub | id_dsa.pub] is present it will be picked as default, matching the first one in this order.`)
	cmd.Flags().StringSlice("labels", []string{}, "labels to add to the "+name+", use it like: --labels \"a=b\" or --labels \"a=\".")
	cmd.Flags().String("userdata", "", `cloud-init.io compatible userdata. [optional]
Can be either the userdata as string, or pointing to the userdata file to use e.g.: "@/tmp/userdata.cfg".`)
	cmd.Flags().StringSlice("dns-servers", []string{}, "dns servers to add to the machine or firewall. [optional]")
	cmd.Flags().StringSlice("ntp-servers", []string{}, "ntp servers to add to the machine or firewall. [optional]")

	cmd.Flags().StringSlice("networks", []string{},
		`Adds a network. Usage: [--networks NETWORK[:ip[;ip]][,NETWORK[:ip[;ip]]...
NETWORK specifies the name or id of an existing network.
IPs can be added per network colon separated, these ips must be already allocated upfront. If no ip(s) are specified per network, one ip per network is allocated.
`)
	cmd.Flags().StringSlice("placement-tags", []string{}, "placement tags used for rack spreading")

	cmd.MarkFlagsMutuallyExclusive("file", "project")
	cmd.MarkFlagsRequiredTogether("project", "networks", "hostname", "image")
	cmd.MarkFlagsRequiredTogether("size", "partition")

	// Completion for arguments
	genericcli.Must(cmd.RegisterFlagCompletionFunc("networks", completion.Network))
	genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", completion.Partition))
	genericcli.Must(cmd.RegisterFlagCompletionFunc("size", completion.Size))
	genericcli.Must(cmd.RegisterFlagCompletionFunc("project", completion.Project))
	genericcli.Must(cmd.RegisterFlagCompletionFunc("image", completion.Image))
	// FIXME implement
	// genericcli.Must(cmd.RegisterFlagCompletionFunc("filesystem-layout", c.c.Completion.FilesystemLayoutListCompletion))
}
