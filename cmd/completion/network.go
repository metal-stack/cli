package completion

import (
	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
)

func (c *Completion) NetworkListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Apiv2().Network().List(c.Ctx, connect.NewRequest(&apiv2.NetworkServiceListRequest{
		Project: c.Project,
	}))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.Msg.Networks {
		names = append(names, s.Id+"\t"+pointer.SafeDeref(s.Name))
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) NetworkTypeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var names []string
	for _, val := range apiv2.NetworkType_value {
		if e, err := enum.GetStringValue(apiv2.NetworkType(val)); err == nil {
			names = append(names, *e)
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) NetworkNatTypeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var names []string
	for _, val := range apiv2.NATType_value {
		if e, err := enum.GetStringValue(apiv2.NATType(val)); err == nil {
			names = append(names, *e)
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) NetworkAddressFamilyCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_V4.String(),
		apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_V6.String(),
		apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_DUAL_STACK.String(),
	}, cobra.ShellCompDirectiveNoFileComp
}
