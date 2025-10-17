package completion

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) MachineListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.MachineServiceListRequest{
		Project: c.Project,
	}
	resp, err := c.Client.Apiv2().Machine().List(c.Ctx, req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, s := range resp.Machines {
		names = append(names, s.Uuid)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
