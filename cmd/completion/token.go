package completion

import (
	"strings"

	"github.com/spf13/cobra"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (c *Completion) TokenListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.TokenServiceListRequest{}
	resp, err := c.Client.Apiv2().Token().List(c.Ctx, req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.Tokens {
		names = append(names, s.Uuid+"\t"+s.Description)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TokenProjectRolesCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	methods, err := c.Client.Apiv2().Method().TokenScopedList(c.Ctx, &apiv2.MethodServiceTokenScopedListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var roles []string

	for project, role := range methods.ProjectRoles {
		if role == apiv2.ProjectRole_PROJECT_ROLE_UNSPECIFIED {
			continue
		}
		roles = append(roles, project+"="+role.String())
	}

	return roles, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TokenTenantRolesCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	methods, err := c.Client.Apiv2().Method().TokenScopedList(c.Ctx, &apiv2.MethodServiceTokenScopedListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var roles []string

	for tenant, role := range methods.TenantRoles {
		if role == apiv2.TenantRole_TENANT_ROLE_UNSPECIFIED {
			continue
		}
		roles = append(roles, tenant+"="+role.String())
	}

	return roles, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TokenMachineRolesCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	methods, err := c.Client.Apiv2().Method().TokenScopedList(c.Ctx, &apiv2.MethodServiceTokenScopedListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var roles []string

	for machine, role := range methods.MachineRoles {
		if role == apiv2.MachineRole_MACHINE_ROLE_UNSPECIFIED {
			continue
		}
		roles = append(roles, machine+"="+role.String())
	}

	return roles, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TokenAdminRoleCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var roles []string

	for i, role := range apiv2.AdminRole_name {
		if i == int32(apiv2.AdminRole_ADMIN_ROLE_UNSPECIFIED) {
			continue
		}
		roles = append(roles, role)
	}

	return roles, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TokenInfraRoleCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var roles []string

	for i, role := range apiv2.InfraRole_name {
		if i == int32(apiv2.InfraRole_INFRA_ROLE_UNSPECIFIED) {
			continue
		}
		roles = append(roles, role)
	}

	return roles, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TokenPermissionsCompletionfunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	methods, err := c.Client.Apiv2().Method().List(c.Ctx, &apiv2.MethodServiceListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	subject := ""
	if s, _, ok := strings.Cut(toComplete, "="); ok { // This does not work ?
		subject = s
	}

	if subject == "" {
		var perms []string
		perms = append(perms, methods.Methods...)
		return perms, cobra.ShellCompDirectiveNoFileComp
	}

	var perms []string
	perms = append(perms, methods.Methods...)

	return perms, cobra.ShellCompDirectiveDefault
}
