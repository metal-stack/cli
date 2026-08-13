package completion

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) Machine(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &adminv2.MachineServiceListRequest{}
	resp, err := c.Client.Adminv2().Machine().List(cmd.Context(), req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, m := range resp.Machines {
		var hostname string
		if m.Allocation != nil {
			hostname = m.Allocation.Hostname
			names = append(names, m.Uuid+"\t"+hostname)
		} else {
			names = append(names, m.Uuid)
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
