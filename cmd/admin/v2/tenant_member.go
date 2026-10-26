package v2

import (
	"fmt"

	"github.com/metal-stack/api/go/errorutil"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"

	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
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

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.TenantServiceAddMemberRequest, *apiv2.TenantMemberServiceUpdateRequest, *apiv2.TenantMember]{
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

func (c *tenantMember) Get(id string) (*apiv2.TenantMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Client.Apiv2().TenantMember().Get(ctx, &apiv2.TenantMemberServiceGetRequest{
		Login: tenant,
	})
	if err != nil {
		return nil, err
	}

	return resp.Member, nil
}

func (c *tenantMember) List() ([]*apiv2.TenantMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Client.Apiv2().TenantMember().List(ctx, &apiv2.TenantMemberServiceListRequest{
		Login: tenant,
	})
	if err != nil {
		return nil, nil
	}

	return resp.Members, nil
}

func (c *tenantMember) Create(rq *adminv2.TenantServiceAddMemberRequest) (*apiv2.TenantMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Tenant().AddMember(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, err
	}

	return resp.Member, nil
}

func (c *tenantMember) Delete(id string) (*apiv2.TenantMember, error) {
	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().TenantMember().Delete(ctx, &apiv2.TenantMemberServiceDeleteRequest{
		Login:  tenant,
		Member: id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to remove member from tenant: %w", err)
	}

	return resp.Member, nil
}

func (c *tenantMember) Convert(r *apiv2.TenantMember) (string, *adminv2.TenantServiceAddMemberRequest, *apiv2.TenantMemberServiceUpdateRequest, error) {
	return r.Member, &adminv2.TenantServiceAddMemberRequest{
			Tenant: r.Tenant,
			Member: r.Member,
			Role:   r.Role,
		},
		&apiv2.TenantMemberServiceUpdateRequest{
			Login:  r.Tenant,
			Member: r.Member,
			Role:   pointer.PointerOrNil(r.Role),
		},
		nil
}

func (c *tenantMember) Update(rq *apiv2.TenantMemberServiceUpdateRequest) (*apiv2.TenantMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().TenantMember().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update member: %w", err)
	}

	return resp.Member, nil
}

func (c *tenantMember) updateRequestFromCLI(args []string) (*apiv2.TenantMemberServiceUpdateRequest, error) {
	member, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	return &apiv2.TenantMemberServiceUpdateRequest{
		Login:  tenant,
		Member: member,
		Role:   pointer.PointerOrNil(apiv2.TenantRole(apiv2.TenantRole_value[viper.GetString("role")])),
	}, nil
}

func (c *tenantMember) createRequestFromCLI() (*adminv2.TenantServiceAddMemberRequest, error) {
	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	return &adminv2.TenantServiceAddMemberRequest{
		Tenant: tenant,
		Member: viper.GetString("id"),
		Role:   apiv2.TenantRole(apiv2.TenantRole_value[viper.GetString("role")]),
	}, nil
}
