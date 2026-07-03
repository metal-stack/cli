package completion

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) PartitionListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.PartitionServiceListRequest{}
	resp, err := c.Client.Apiv2().Partition().List(c.Ctx, req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, p := range resp.Partitions {
		names = append(names, p.Id+"\t"+p.Description)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
