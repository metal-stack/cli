package completion

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) Tenant(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.TenantServiceListRequest{}
	resp, err := c.Client.Apiv2().Tenant().List(cmd.Context(), req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, t := range resp.Tenants {
		names = append(names, t.Login+"\t"+t.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TenantRole(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var names []string

	for value, name := range apiv2.TenantRole_name {
		if value == 0 {
			continue
		}

		names = append(names, name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TenantInvite(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	projectResp, err := c.Client.Apiv2().Project().Get(cmd.Context(), &apiv2.ProjectServiceGetRequest{
		Project: c.Proj,
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	resp, err := c.Client.Apiv2().Tenant().InvitesList(cmd.Context(), &apiv2.TenantServiceInvitesListRequest{
		Login: projectResp.Project.Tenant,
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, invite := range resp.Invites {
		names = append(names, invite.Secret+"\t"+invite.Role.String())
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TenantMember(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	projectResp, err := c.Client.Apiv2().Project().Get(cmd.Context(), &apiv2.ProjectServiceGetRequest{
		Project: c.Proj,
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	resp, err := c.Client.Apiv2().Tenant().Get(cmd.Context(), &apiv2.TenantServiceGetRequest{
		Login: projectResp.Project.Tenant,
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, member := range resp.TenantMembers {
		names = append(names, member.Id+"\t"+member.Role.String())
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) AdminTenant(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &adminv2.TenantServiceListRequest{}
	resp, err := c.Client.Adminv2().Tenant().List(cmd.Context(), req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, s := range resp.Tenants {
		names = append(names, s.Login+"\t"+s.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
