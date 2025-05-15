package completion

import (
	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) IpListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.IPServiceListRequest{
		Project: c.Project,
	}
	resp, err := c.Client.Apiv2().IP().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, s := range resp.Msg.Ips {
		names = append(names, s.Uuid+"\t"+s.Ip+"\t"+s.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) IpAddressFamilyCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{apiv2.IPAddressFamily_IP_ADDRESS_FAMILY_V4.String(), apiv2.IPAddressFamily_IP_ADDRESS_FAMILY_V6.String()}, cobra.ShellCompDirectiveNoFileComp
}
