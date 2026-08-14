package completion

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) Machine(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Apiv2().Machine().List(cmd.Context(), &apiv2.MachineServiceListRequest{
		Project: c.Proj,
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, s := range resp.Machines {
		names = append(names, s.Uuid)
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
