package completion

import (
	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) Ip(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.IPServiceListRequest{
		Project: c.Proj,
	}
	resp, err := c.Client.Apiv2().IP().List(cmd.Context(), req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, s := range resp.Ips {
		names = append(names, s.Uuid+"\t"+s.Ip+"\t"+s.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) AddressFamily(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var afs []string
	for _, af := range []apiv2.IPAddressFamily{
		apiv2.IPAddressFamily_IP_ADDRESS_FAMILY_V4,
		apiv2.IPAddressFamily_IP_ADDRESS_FAMILY_V6} {
		stringValue, err := enum.GetStringValue(af)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		afs = append(afs, *stringValue)
	}

	return afs, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) IPType(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var afs []string
	for _, af := range []apiv2.IPType{
		apiv2.IPType_IP_TYPE_STATIC,
		apiv2.IPType_IP_TYPE_EPHEMERAL} {
		stringValue, err := enum.GetStringValue(af)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		afs = append(afs, *stringValue)
	}

	return afs, cobra.ShellCompDirectiveNoFileComp
}
