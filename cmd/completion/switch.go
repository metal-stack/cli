package completion

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
)

func (c *Completion) SwitchListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Adminv2().Switch().List(c.Ctx, &adminv2.SwitchServiceListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var ids []string
	for _, s := range resp.Switches {
		ids = append(ids, s.Id)
	}

	return ids, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) PartitionListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Adminv2().Switch().List(c.Ctx, &adminv2.SwitchServiceListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var partitions []string
	for _, s := range resp.Switches {
		partitions = append(partitions, s.Partition)
	}

	return partitions, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) SwitchRackListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Adminv2().Switch().List(c.Ctx, &adminv2.SwitchServiceListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var racks []string
	for _, s := range resp.Switches {
		racks = append(racks, pointer.SafeDeref(s.Rack))
	}

	return racks, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) SwitchOSVendorListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Adminv2().Switch().List(c.Ctx, &adminv2.SwitchServiceListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var oss []string
	for _, s := range resp.Switches {
		oss = append(oss, pointer.SafeDeref(s.Os).Vendor.String())
	}

	return oss, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) SwitchOSVersionListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Adminv2().Switch().List(c.Ctx, &adminv2.SwitchServiceListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var osVersions []string
	for _, s := range resp.Switches {
		osVersions = append(osVersions, pointer.SafeDeref(s.Os).Version)
	}

	return osVersions, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) SwitchListPorts(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		// there is no switch selected so we cannot get the list of ports
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	resp, err := c.Client.Adminv2().Switch().Get(c.Ctx, &adminv2.SwitchServiceGetRequest{
		Id: args[0],
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var nics []string
	for _, nic := range resp.Switch.Nics {
		nics = append(nics, pointer.SafeDeref(nic).Name)
	}

	return nics, cobra.ShellCompDirectiveNoFileComp
}
