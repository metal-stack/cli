package v2

import (
	"fmt"
	"sort"
	"strings"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/cli/pkg/helpers"
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
	var (
		w = &partition{
			c: c,
		}

		addMutableFlags = func(cmd *cobra.Command) {
			cmd.Flags().String("description", "", "the description of the partition")
			cmd.Flags().String("image-url", "", "the url of the boot image used by metal-hammer")
			cmd.Flags().String("kernel-url", "", "the url of the kernel used by metal-hammer")
			cmd.Flags().String("commandline", "", "the kernel commandline used by metal-hammer")
			cmd.Flags().StringSlice("dns-servers", nil, "the dns servers of this partition")
			cmd.Flags().StringSlice("ntp-servers", nil, "the ntp servers of this partition")
			cmd.Flags().StringSlice("mgmt-service-addresses", nil, "the management service addresses of this partition, each in the form <ip|host>:<port>")
		}
	)

	gcli := genericcli.NewGenericCLI(w).WithFS(c.Fs)

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.PartitionServiceCreateRequest, *adminv2.PartitionServiceUpdateRequest, *apiv2.Partition]{
		BinaryName:      config.BinaryName,
		GenericCLI:      gcli,
		Singular:        "partition",
		Plural:          "partitions",
		Description:     "manage partitions",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.PartitionListCompletion,
		Sorter:          sorters.PartitionSorter(),
		OnlyCmds: genericcli.OnlyCmds(
			genericcli.DescribeCmd,
			genericcli.ListCmd,
			genericcli.CreateCmd,
			genericcli.UpdateCmd,
			genericcli.DeleteCmd,
			genericcli.EditCmd,
		),
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("id", "", "the id of the partition to create")
			addMutableFlags(cmd)
		},
		CreateRequestFromCLI: func() (*adminv2.PartitionServiceCreateRequest, error) {
			return &adminv2.PartitionServiceCreateRequest{
				Partition: &apiv2.Partition{
					Id:                   viper.GetString("id"),
					Description:          viper.GetString("description"),
					BootConfiguration:    partitionBootConfigurationFromCLI(),
					DnsServers:           dnsServersFromCLI(viper.GetStringSlice("dns-servers")),
					NtpServers:           ntpServersFromCLI(viper.GetStringSlice("ntp-servers")),
					MgmtServiceAddresses: viper.GetStringSlice("mgmt-service-addresses"),
				},
			}, nil
		},
		UpdateCmdMutateFn:    addMutableFlags,
		UpdateRequestFromCLI: w.updateRequestFromCLI,
	}

	capacityCmd := &cobra.Command{
		Use:   "capacity",
		Short: "show partition capacity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.capacity()
		},
	}

	capacityCmd.Flags().StringP("id", "", "", "filter on partition id.")
	capacityCmd.Flags().StringP("size", "", "", "filter on size id.")
	capacityCmd.Flags().StringP("project", "", "", "consider project-specific counts, e.g. size reservations.")
	capacityCmd.Flags().StringSlice("sort-by", []string{}, fmt.Sprintf("order by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: %s", strings.Join(sorters.PartitionCapacitySorter().AvailableKeys(), "|")))
	genericcli.Must(capacityCmd.RegisterFlagCompletionFunc("id", c.Completion.PartitionListCompletion))
	genericcli.Must(capacityCmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
	genericcli.Must(capacityCmd.RegisterFlagCompletionFunc("size", c.Completion.SizeListCompletion))
	genericcli.Must(capacityCmd.RegisterFlagCompletionFunc("sort-by", cobra.FixedCompletions(sorters.PartitionCapacitySorter().AvailableKeys(), cobra.ShellCompDirectiveNoFileComp)))

	return genericcli.NewCmds(cmdsConfig, capacityCmd)
}

func (c *partition) capacity() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Partition().Capacity(ctx, &adminv2.PartitionServiceCapacityRequest{
		Id:      pointer.PointerOrNil(viper.GetString("id")),
		Size:    pointer.PointerOrNil(viper.GetString("size")),
		Project: pointer.PointerOrNil(viper.GetString("project")),
	})
	if err != nil {
		return fmt.Errorf("failed to get partition capacity: %w", err)
	}

	for _, cap := range resp.PartitionCapacity {
		sort.SliceStable(cap.MachineSizeCapacities, func(i, j int) bool {
			return cap.MachineSizeCapacities[i].Size < cap.MachineSizeCapacities[j].Size
		})
	}

	err = sorters.PartitionCapacitySorter().SortBy(resp.PartitionCapacity)
	if err != nil {
		return err
	}

	return c.c.ListPrinter.Print(resp.PartitionCapacity)
}

func (c *partition) Get(id string) (*apiv2.Partition, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Partition().Get(ctx, &apiv2.PartitionServiceGetRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("failed to get partition: %w", err)
	}

	return resp.Partition, nil
}

func (c *partition) List() ([]*apiv2.Partition, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.PartitionServiceListRequest{
		Query: &apiv2.PartitionQuery{
			Id: pointer.PointerOrNil(viper.GetString("id")),
		},
	}

	if labelSlice := viper.GetStringSlice("labels"); len(labelSlice) > 0 {
		var err error

		req.Query.Labels, err = helpers.LabelsFromSlice(labelSlice)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.c.Client.Apiv2().Partition().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get partitions: %w", err)
	}

	return resp.Partitions, nil
}

func (c *partition) Create(rq *adminv2.PartitionServiceCreateRequest) (*apiv2.Partition, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Partition().Create(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to create partition: %w", err)
	}

	return resp.Partition, nil
}

func (c *partition) Delete(id string) (*apiv2.Partition, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.PartitionServiceDeleteRequest{Id: id}

	resp, err := c.c.Client.Adminv2().Partition().Delete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete partition: %w", err)
	}

	return resp.Partition, nil
}

func (c *partition) Convert(r *apiv2.Partition) (string, *adminv2.PartitionServiceCreateRequest, *adminv2.PartitionServiceUpdateRequest, error) {
	return r.Id,
		&adminv2.PartitionServiceCreateRequest{
			Partition: r,
		}, &adminv2.PartitionServiceUpdateRequest{
			Id:                   r.Id,
			Description:          new(r.Description),
			BootConfiguration:    r.BootConfiguration,
			DnsServers:           r.DnsServers,
			NtpServers:           r.NtpServers,
			MgmtServiceAddresses: r.MgmtServiceAddresses,
			UpdateMeta:           helpers.UpdateMetaFromMeta(r.Meta),
			Labels: &apiv2.UpdateLabels{
				Strategy: &apiv2.UpdateLabels_Replace{
					Replace: &apiv2.Labels{
						Labels: pointer.SafeDeref(pointer.SafeDeref(r.Meta).Labels).Labels,
					},
				},
			},
		}, nil
}

func (c *partition) Update(rq *adminv2.PartitionServiceUpdateRequest) (*apiv2.Partition, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Partition().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update partition: %w", err)
	}

	return resp.Partition, nil
}

func (c *partition) updateRequestFromCLI(args []string) (*adminv2.PartitionServiceUpdateRequest, error) {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	req := &adminv2.PartitionServiceUpdateRequest{
		Id: id,
		UpdateMeta: &apiv2.UpdateMeta{
			LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
		},
		Description:          pointer.PointerOrNil(viper.GetString("description")),
		BootConfiguration:    partitionBootConfigurationFromCLI(),
		DnsServers:           dnsServersFromCLI(viper.GetStringSlice("dns-servers")),
		NtpServers:           ntpServersFromCLI(viper.GetStringSlice("ntp-servers")),
		MgmtServiceAddresses: viper.GetStringSlice("mgmt-service-addresses"),
		Labels:               nil,
	}

	// TODO: we need a helper function that evaluates update labels (handles both strategies for all cmds)
	if labelSlice := viper.GetStringSlice("labels"); len(labelSlice) > 0 {
		labels, err := helpers.LabelsFromSlice(labelSlice)
		if err != nil {
			return nil, err
		}

		req.Labels = &apiv2.UpdateLabels{
			Strategy: &apiv2.UpdateLabels_Replace{
				Replace: labels,
			},
		}
	}

	return req, nil
}

func partitionBootConfigurationFromCLI() *apiv2.PartitionBootConfiguration {
	if !viper.IsSet("image-url") && !viper.IsSet("kernel-url") && !viper.IsSet("commandline") {
		return nil
	}

	return &apiv2.PartitionBootConfiguration{
		ImageUrl:    viper.GetString("image-url"),
		KernelUrl:   viper.GetString("kernel-url"),
		Commandline: viper.GetString("commandline"),
	}
}

func dnsServersFromCLI(ips []string) []*apiv2.DNSServer {
	var servers []*apiv2.DNSServer
	for _, ip := range ips {
		servers = append(servers, &apiv2.DNSServer{Ip: ip})
	}
	return servers
}

func ntpServersFromCLI(addresses []string) []*apiv2.NTPServer {
	var servers []*apiv2.NTPServer
	for _, address := range addresses {
		servers = append(servers, &apiv2.NTPServer{Address: address})
	}
	return servers
}
