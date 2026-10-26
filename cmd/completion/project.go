package completion

import (
	"github.com/spf13/cobra"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (c *Completion) Project(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.ProjectServiceListRequest{}
	resp, err := c.Client.Apiv2().Project().List(cmd.Context(), req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.GetProjects() {
		names = append(names, s.Uuid+"\t"+s.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ProjectRole(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var names []string

	for value, name := range apiv2.ProjectRole_name {
		if value == 0 {
			continue
		}

		names = append(names, name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ProjectInvite(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Apiv2().Project().InvitesList(cmd.Context(), &apiv2.ProjectServiceInvitesListRequest{
		Project: c.Proj,
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

func (c *Completion) ProjectMember(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Apiv2().ProjectMember().List(cmd.Context(), &apiv2.ProjectMemberServiceListRequest{
		Project: c.Proj,
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, member := range resp.Members {
		names = append(names, member.Member+"\t"+member.Role.String())
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
