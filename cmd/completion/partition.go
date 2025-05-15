package completion

import (
	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) PartitionListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Apiv2().Partition().List(c.Ctx, connect.NewRequest(&apiv2.PartitionServiceListRequest{}))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.Msg.Partitions {
		names = append(names, s.Id+"\t"+s.Description)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
