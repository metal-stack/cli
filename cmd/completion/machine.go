package completion

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) Machine(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Apiv2().Machine().List(cmd.Context(), &apiv2.MachineServiceListRequest{
		// FIXME this only works for end users, not admin because they have no project set
		Project: c.Proj,
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, m := range resp.Machines {
		name := m.Uuid
		if m.Allocation != nil {
			name = m.Uuid + "\t" + m.Allocation.Hostname
		}
		names = append(names, name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) AdminMachine(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Adminv2().Machine().List(cmd.Context(), &adminv2.MachineServiceListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, m := range resp.Machines {
		name := m.Uuid
		if m.Allocation != nil {
			name = m.Uuid + "\t" + m.Allocation.Hostname
		}
		names = append(names, name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) Firewall(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Adminv2().Machine().List(cmd.Context(), &adminv2.MachineServiceListRequest{
		Query: &apiv2.MachineQuery{
			Allocation: &apiv2.MachineAllocationQuery{
				AllocationType: apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_FIREWALL.Enum(),
			},
		},
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, m := range resp.Machines {
		name := m.Uuid + "\t" + m.Allocation.Hostname
		names = append(names, name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) BMCCommands(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var names []string

	for _, name := range apiv2.MachineBMCCommand_name {
		names = append(names, name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
