package v1

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type tenant struct {
	c *config.Config
}

func newTenantCmd(c *config.Config) *cobra.Command {
	w := &tenant{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv2.TenantServiceCreateRequest, *apiv2.TenantServiceUpdateRequest, *apiv2.Tenant]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "tenant",
		Plural:          "tenants",
		Description:     "manage api tenants",
		Sorter:          sorters.TenantSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "lists only tenants with the given name")
			cmd.Flags().String("id", "", "lists only tenant with the given tenant id")
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the tenant to create")
			cmd.Flags().String("description", "", "the description of the tenant to create")
			cmd.Flags().String("email", "", "the email of the tenant to create")
			cmd.Flags().String("phone", "", "the phone number of the tenant to create")
			cmd.Flags().String("avatar-url", "", "the avatar url of the tenant to create")
		},
		CreateRequestFromCLI: w.createRequestFromCLI,
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the tenant to update")
			cmd.Flags().String("description", "", "the description of the tenant to update")
		},
		UpdateRequestFromCLI: w.updateRequestFromCLI,
		ValidArgsFn:          w.c.Completion.TenantListCompletion,
	}

	inviteCmd := &cobra.Command{
		Use:   "invite",
		Short: "manage tenant invites",
	}

	generateInviteCmd := &cobra.Command{
		Use:   "generate-join-secret",
		Short: "generate an invite secret to share with the new member",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.generateInvite()
		},
	}

	generateInviteCmd.Flags().String("tenant", "", "the tenant for which to generate the invite")
	generateInviteCmd.Flags().String("role", apiv2.TenantRole_TENANT_ROLE_VIEWER.String(), "the role that the new member will assume when joining through the invite secret")

	genericcli.Must(generateInviteCmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))
	genericcli.Must(generateInviteCmd.RegisterFlagCompletionFunc("role", c.Completion.TenantRoleCompletion))

	deleteInviteCmd := &cobra.Command{
		Use:     "delete <secret>",
		Aliases: []string{"destroy", "rm", "remove"},
		Short:   "deletes a pending invite",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.deleteInvite(args)
		},
		ValidArgsFunction: c.Completion.TenantInviteListCompletion,
	}

	deleteInviteCmd.Flags().String("tenant", "", "the tenant in which to delete the invite")

	genericcli.Must(deleteInviteCmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))

	listInvitesCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "lists the currently pending invites",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.listInvites()
		},
	}

	listInvitesCmd.Flags().String("tenant", "", "the tenant for which to list the invites")

	genericcli.AddSortFlag(listInvitesCmd, sorters.TenantInviteSorter())

	genericcli.Must(listInvitesCmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))

	joinTenantCmd := &cobra.Command{
		Use:   "join <secret>",
		Short: "join a tenant of someone who shared an invite secret with you",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.join(args)
		},
	}

	memberCmd := &cobra.Command{
		Use:     "member",
		Aliases: []string{"members"},
		Short:   "manage tenant members",
	}

	listMembersCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "lists members of a tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.listMembers()
		},
	}

	listMembersCmd.Flags().String("tenant", "", "the tenant in which to remove the member")

	genericcli.AddSortFlag(listMembersCmd, sorters.TenantMemberSorter())

	removeMemberCmd := &cobra.Command{
		Use:     "remove <member>",
		Short:   "remove member from a tenant",
		Aliases: []string{"destroy", "rm", "remove"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.removeMember(args)
		},
		ValidArgsFunction: c.Completion.TenantMemberListCompletion,
	}

	removeMemberCmd.Flags().String("tenant", "", "the tenant in which to remove the member")

	genericcli.Must(removeMemberCmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))

	updateMemberCmd := &cobra.Command{
		Use:   "update <member>",
		Short: "update member from a tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.updateMember(args)
		},
		ValidArgsFunction: c.Completion.TenantMemberListCompletion,
	}

	updateMemberCmd.Flags().String("tenant", "", "the tenant in which to remove the member")
	updateMemberCmd.Flags().String("role", "", "the role of the member")

	genericcli.Must(updateMemberCmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))
	genericcli.Must(updateMemberCmd.RegisterFlagCompletionFunc("role", c.Completion.TenantRoleCompletion))

	memberCmd.AddCommand(removeMemberCmd, updateMemberCmd, listMembersCmd)

	inviteCmd.AddCommand(generateInviteCmd, deleteInviteCmd, listInvitesCmd, joinTenantCmd)

	return genericcli.NewCmds(cmdsConfig, joinTenantCmd, inviteCmd, memberCmd)
}

func (c *tenant) Get(id string) (*apiv2.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.TenantServiceGetRequest{
		Login: id,
	}

	resp, err := c.c.Client.Apiv2().Tenant().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return resp.GetTenant(), nil
}

func (c *tenant) List() ([]*apiv2.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.TenantServiceListRequest{
		Name: pointer.PointerOrNil(viper.GetString("name")),
		Id:   pointer.PointerOrNil(viper.GetString("tenant")),
	}
	resp, err := c.c.Client.Apiv2().Tenant().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	return resp.GetTenants(), nil
}

func (c *tenant) Create(rq *apiv2.TenantServiceCreateRequest) (*apiv2.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Tenant().Create(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return resp.Tenant, nil
}

func (c *tenant) Delete(id string) (*apiv2.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Tenant().Delete(ctx, &apiv2.TenantServiceDeleteRequest{
		Login: id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete tenant: %w", err)
	}

	return resp.Tenant, nil
}

func (c *tenant) Convert(r *apiv2.Tenant) (string, *apiv2.TenantServiceCreateRequest, *apiv2.TenantServiceUpdateRequest, error) {

	return r.Login, &apiv2.TenantServiceCreateRequest{
			Name:        r.Name,
			Description: &r.Description,
			Email:       &r.Email,
			AvatarUrl:   &r.AvatarUrl,
		},
		&apiv2.TenantServiceUpdateRequest{
			Login:     r.Login,
			Name:      pointer.PointerOrNil(r.Name),
			Email:     pointer.PointerOrNil(r.Email),
			AvatarUrl: pointer.PointerOrNil(r.AvatarUrl),
		},
		nil
}

func (c *tenant) Update(rq *apiv2.TenantServiceUpdateRequest) (*apiv2.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Tenant().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return resp.Tenant, nil
}

func (c *tenant) createRequestFromCLI() (*apiv2.TenantServiceCreateRequest, error) {
	return &apiv2.TenantServiceCreateRequest{
		Name:        viper.GetString("name"),
		Description: pointer.PointerOrNil(viper.GetString("description")),
		Email:       pointer.PointerOrNil(viper.GetString("email")),
		AvatarUrl:   pointer.PointerOrNil(viper.GetString("avatar-url")),
	}, nil
}

func (c *tenant) updateRequestFromCLI(args []string) (*apiv2.TenantServiceUpdateRequest, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *tenant) join(args []string) error {
	secret, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Tenant().InviteGet(ctx, &apiv2.TenantServiceInviteGetRequest{
		Secret: secret,
	})
	if err != nil {
		return fmt.Errorf("failed to get tenant invite: %w", err)
	}

	err = genericcli.PromptCustom(&genericcli.PromptConfig{
		ShowAnswers: true,
		Message: fmt.Sprintf(
			"Do you want to join tenant \"%s\" as %s?",
			color.GreenString(resp.GetInvite().GetTargetTenantName()),
			resp.GetInvite().GetRole().String(),
		),
		In:  c.c.In,
		Out: c.c.Out,
	})
	if err != nil {
		return err
	}

	ctx2, cancel2 := c.c.NewRequestContext()
	defer cancel2()

	acceptResp, err := c.c.Client.Apiv2().Tenant().InviteAccept(ctx2, &apiv2.TenantServiceInviteAcceptRequest{
		Secret: secret,
	})
	if err != nil {
		return fmt.Errorf("failed to join tenant: %w", err)
	}

	_, _ = fmt.Fprintf(c.c.Out, "%s successfully joined tenant \"%s\"\n", color.GreenString("✔"), color.GreenString(acceptResp.TenantName))

	return nil
}

func (c *tenant) generateInvite() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	resp, err := c.c.Client.Apiv2().Tenant().Invite(ctx, &apiv2.TenantServiceInviteRequest{
		Login: tenant,
		Role:  apiv2.TenantRole(apiv2.TenantRole_value[viper.GetString("role")]),
	})
	if err != nil {
		return fmt.Errorf("failed to generate an invite: %w", err)
	}

	_, _ = fmt.Fprintf(c.c.Out, "You can share this secret with the member to join, it expires in %s:\n\n", humanize.Time(resp.Invite.ExpiresAt.AsTime()))
	_, _ = fmt.Fprintf(c.c.Out, "%s (https://console.metal-stack.io/organization-invite/%s)\n", resp.Invite.Secret, resp.Invite.Secret)

	return nil
}

func (c *tenant) listInvites() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	resp, err := c.c.Client.Apiv2().Tenant().InvitesList(ctx, &apiv2.TenantServiceInvitesListRequest{
		Login: tenant,
	})
	if err != nil {
		return fmt.Errorf("failed to list invites: %w", err)
	}

	err = sorters.TenantInviteSorter().SortBy(resp.Invites)
	if err != nil {
		return err
	}

	return c.c.ListPrinter.Print(resp.Invites)
}

func (c *tenant) deleteInvite(args []string) error {
	secret, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	_, err = c.c.Client.Apiv2().Tenant().InviteDelete(ctx, &apiv2.TenantServiceInviteDeleteRequest{
		Login:  tenant,
		Secret: secret,
	})
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}

func (c *tenant) removeMember(args []string) error {
	member, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	_, err = c.c.Client.Apiv2().Tenant().RemoveMember(ctx, &apiv2.TenantServiceRemoveMemberRequest{
		Login:  tenant,
		Member: member,
	})
	if err != nil {
		return fmt.Errorf("failed to remove member from tenant: %w", err)
	}

	_, _ = fmt.Fprintf(c.c.Out, "%s successfully removed member %q\n", color.GreenString("✔"), member)

	return nil
}

func (c *tenant) updateMember(args []string) error {
	member, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Tenant().UpdateMember(ctx, &apiv2.TenantServiceUpdateMemberRequest{
		Login:  tenant,
		Member: member,
		Role:   apiv2.TenantRole(apiv2.TenantRole_value[viper.GetString("role")]),
	})
	if err != nil {
		return fmt.Errorf("failed to update member: %w", err)
	}

	return c.c.DescribePrinter.Print(resp.GetTenantMember())
}

func (c *tenant) listMembers() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	resp, err := c.c.Client.Apiv2().Tenant().Get(ctx, &apiv2.TenantServiceGetRequest{
		Login: tenant,
	})
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	members := resp.GetTenantMembers()

	if err := sorters.TenantMemberSorter().SortBy(members); err != nil {
		return err
	}

	return c.c.ListPrinter.Print(members)
}
