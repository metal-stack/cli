package completion

import (
	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) SizeListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.SizeServiceListRequest{}
	resp, err := c.Client.Apiv2().Size().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, s := range resp.Msg.Sizes {
		names = append(names, s.Id)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
