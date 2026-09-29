package v2

import (
	"fmt"

	"github.com/metal-stack/api/go/errorutil"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	cliv2 "github.com/metal-stack/api/go/metalstack/cli/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type tenantMember struct {
	c *config.Config
}

func newTenantMemberCmd(c *config.Config) *cobra.Command {
	w := &tenantMember{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv2.TenantServiceAddMemberRequest, *apiv2.TenantServiceUpdateMemberRequest, *cliv2.TenantMember]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "member",
		Plural:          "members",
		Description:     "manage api tenant members",
		Sorter:          sorters.TenantMemberSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("tenant", "", "the tenant to describe the tenant members of, defaults to tenant of the default project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant", c.Completion.Tenant))
		},
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("tenant", "", "the tenant to list the tenant members of, defaults to tenant of the default project")
			cmd.Flags().String("role", "", "the role of the member")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant", c.Completion.Tenant))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("role", c.Completion.TenantRole))
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("tenant", "", "the tenant for which to create a tenant member, defaults to tenant of the default project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant", c.Completion.Tenant))
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("tenant", "", "the tenant for which to update a tenant member, defaults to tenant of the default project")
			cmd.Flags().String("role", "", "the role of the member")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant", c.Completion.Tenant))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("role", c.Completion.TenantRole))

		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("tenant", "", "the tenant from which to delete a tenant member, defaults to tenant of the default project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant", c.Completion.Tenant))
		},
		EditCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("tenant", "", "the tenant from which to edit a tenant member, defaults to tenant of the default project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant", c.Completion.Tenant))
		},
		CreateRequestFromCLI: w.createRequestFromCLI,
		UpdateRequestFromCLI: w.updateRequestFromCLI,
		ValidArgsFn:          w.c.Completion.TenantMember,
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *tenantMember) Get(id string) (*cliv2.TenantMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Client.Apiv2().Tenant().Get(ctx, &apiv2.TenantServiceGetRequest{
		Login: tenant,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	for _, member := range resp.GetTenantMembers() {
		if member.Id == id {
			return &cliv2.TenantMember{
				Tenant:       tenant,
				TenantMember: member,
			}, nil
		}
	}

	return nil, fmt.Errorf("unable to find tenant members %q in tenant %q", id, tenant)
}

func (c *tenantMember) List() ([]*cliv2.TenantMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Client.Apiv2().Tenant().Get(ctx, &apiv2.TenantServiceGetRequest{
		Login: tenant,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	var res []*cliv2.TenantMember

	for _, member := range resp.GetTenantMembers() {
		res = append(res, &cliv2.TenantMember{
			Tenant:       tenant,
			TenantMember: member,
		})
	}

	return res, nil
}

func (c *tenantMember) Create(rq *apiv2.TenantServiceAddMemberRequest) (*cliv2.TenantMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Tenant().AddMember(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, err
	}

	return &cliv2.TenantMember{
		Tenant:       rq.Login,
		TenantMember: resp.TenantMember,
	}, nil
}

func (c *tenantMember) Delete(id string) (*cliv2.TenantMember, error) {
	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Tenant().RemoveMember(ctx, &apiv2.TenantServiceRemoveMemberRequest{
		Login:  tenant,
		Member: id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to remove member from tenant: %w", err)
	}

	return &cliv2.TenantMember{
		Tenant:       tenant,
		TenantMember: resp.TenantMember,
	}, nil
}

func (c *tenantMember) Convert(r *cliv2.TenantMember) (string, *apiv2.TenantServiceAddMemberRequest, *apiv2.TenantServiceUpdateMemberRequest, error) {
	return r.TenantMember.Id, &apiv2.TenantServiceAddMemberRequest{
			Login:  r.Tenant,
			Member: r.TenantMember.Id,
			Role:   r.TenantMember.Role,
		},
		&apiv2.TenantServiceUpdateMemberRequest{
			Login:  r.Tenant,
			Member: r.TenantMember.Id,
			Role:   r.TenantMember.Role,
		},
		nil
}

func (c *tenantMember) Update(rq *apiv2.TenantServiceUpdateMemberRequest) (*cliv2.TenantMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Tenant().UpdateMember(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update member: %w", err)
	}

	return &cliv2.TenantMember{
		Tenant:       rq.Login,
		TenantMember: resp.TenantMember,
	}, nil
}

func (c *tenantMember) updateRequestFromCLI(args []string) (*apiv2.TenantServiceUpdateMemberRequest, error) {
	member, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	return &apiv2.TenantServiceUpdateMemberRequest{
		Login:  tenant,
		Member: member,
		Role:   apiv2.TenantRole(apiv2.TenantRole_value[viper.GetString("role")]),
	}, nil
}

func (c *tenantMember) createRequestFromCLI() (*apiv2.TenantServiceAddMemberRequest, error) {
	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	return &apiv2.TenantServiceAddMemberRequest{
		Login:  tenant,
		Member: viper.GetString("id"),
		Role:   apiv2.TenantRole(apiv2.TenantRole_value[viper.GetString("role")]),
	}, nil
}
