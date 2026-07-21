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

type (
	tenant struct {
		c *config.Config
	}
	tenantMember struct {
		c *config.Config
	}
)

func newTenantCmd(c *config.Config) *cobra.Command {
	w := &tenant{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.TenantServiceCreateRequest, any, *apiv2.Tenant]{
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
			cmd.Flags().String("email", "", "lists only tenant with the given email address")
			cmd.Flags().StringSlice("labels", nil, "lists only tenant with the given labels")
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the tenant to create")
			cmd.Flags().String("description", "", "the description of the tenant to create")
			cmd.Flags().String("email", "", "the email of the tenant to create")
			cmd.Flags().String("avatar-url", "", "the avatar url of the tenant to create")
		},
		CreateRequestFromCLI: func() (*adminv2.TenantServiceCreateRequest, error) {
			return &adminv2.TenantServiceCreateRequest{
				Name:        viper.GetString("name"),
				Description: pointer.PointerOrNil(viper.GetString("description")),
				Email:       pointer.PointerOrNil(viper.GetString("email")),
				AvatarUrl:   pointer.PointerOrNil(viper.GetString("avatar-url")),
			}, nil
		},
		OnlyCmds:    genericcli.OnlyCmds(genericcli.ListCmd, genericcli.CreateCmd),
		ValidArgsFn: w.c.Completion.AdminTenantListCompletion,
	}

	return genericcli.NewCmds(cmdsConfig, newAddMemberCmd(c), newTenantMembersCmd(c))
}

func (c *tenant) Get(id string) (*apiv2.Tenant, error) {
	panic("unimplemented")
}

func (c *tenant) List() ([]*apiv2.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.TenantServiceListRequest{
		Query: &apiv2.TenantQuery{
			Name:  pointer.PointerOrNil(viper.GetString("name")),
			Login: pointer.PointerOrNil(viper.GetString("tenant")),
		},
	}

	if labelSlice := viper.GetStringSlice("labels"); len(labelSlice) > 0 {
		labels, err := genericcli.LabelsToMap(labelSlice)
		if err != nil {
			return nil, err
		}

		req.Query.Labels = &apiv2.Labels{
			Labels: labels,
		}
	}

	resp, err := c.c.Client.Adminv2().Tenant().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	return resp.GetTenants(), nil
}

func (c *tenant) Create(rq *adminv2.TenantServiceCreateRequest) (*apiv2.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Tenant().Create(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return resp.Tenant, nil
}

func (c *tenant) Delete(id string) (*apiv2.Tenant, error) {
	panic("unimplemented")
}

func (c *tenant) Convert(r *apiv2.Tenant) (string, *adminv2.TenantServiceCreateRequest, any, error) {
	panic("unimplemented")
}

func (c *tenant) Update(rq any) (*apiv2.Tenant, error) {
	panic("unimplemented")
}

func newAddMemberCmd(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-member",
		Short: "Add a new member to a tenant",
		Long:  `Add a new member to an existing tenant by specifying the tenant ID, member's ID, and role.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewRequestContext()
			defer cancel()

			var (
				tenantId   = viper.GetString("tenant-id")
				memberId   = viper.GetString("member-id")
				memberRole = viper.GetString("role")
			)

			if tenantId == "" || memberId == "" || memberRole == "" {
				return fmt.Errorf("tenant ID, member ID, and role must all be specified")
			}

			_, err := c.Client.Adminv2().Tenant().AddMember(ctx, &adminv2.TenantServiceAddMemberRequest{
				Role:   apiv2.TenantRole(apiv2.TenantRole_value[memberRole]),
				Tenant: tenantId,
				Member: memberId,
			})
			if err != nil {
				return fmt.Errorf("failed to add member to tenant: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().String("tenant-id", "", "ID of the tenant where the member is added")
	cmd.Flags().String("member-id", "", "ID of the member to be added")
	cmd.Flags().String("role", "", "Role of the member within the tenant")
	genericcli.Must(cmd.MarkFlagRequired("tenant-id"))
	genericcli.Must(cmd.MarkFlagRequired("member-id"))
	genericcli.Must(cmd.MarkFlagRequired("role"))

	genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant-id", c.Completion.AdminTenantListCompletion))
	genericcli.Must(cmd.RegisterFlagCompletionFunc("member-id", c.Completion.AdminTenantListCompletion))
	genericcli.Must(cmd.RegisterFlagCompletionFunc("role", c.Completion.TenantRoleCompletion))

	return cmd
}

func newTenantMembersCmd(c *config.Config) *cobra.Command {
	wm := &tenantMember{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.TenantServiceAddMemberRequest, any, *adminv2.TenantServiceAddMemberResponse]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(wm).WithFS(c.Fs),
		Singular:        "member",
		Plural:          "members",
		Description:     "manage tenant-members",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ApplyCmd),
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *tenantMember) Get(id string) (*adminv2.TenantServiceAddMemberResponse, error) {
	panic("unimplemented")
}

func (c *tenantMember) List() ([]*adminv2.TenantServiceAddMemberResponse, error) {
	// ctx, cancel := c.c.NewRequestContext()
	// defer cancel()

	// req := &adminv2.TenantServiceListRequest{
	// 	Query: &apiv2.TenantQuery{
	// 		Name:  pointer.PointerOrNil(viper.GetString("name")),
	// 		Login: pointer.PointerOrNil(viper.GetString("tenant")),
	// 	},
	// }

	// resp, err := c.c.Client.Adminv2().Tenant().List(ctx, req)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to list tenants: %w", err)
	// }

	panic("unimplemented")
}

func (c *tenantMember) Create(rq *adminv2.TenantServiceAddMemberRequest) (*adminv2.TenantServiceAddMemberResponse, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Tenant().AddMember(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, fmt.Errorf("failed to create tenant member: %w", err)
	}

	return resp, nil
}

func (c *tenantMember) Delete(id string) (*adminv2.TenantServiceAddMemberResponse, error) {
	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	_, err = c.c.Client.Adminv2().Tenant().RemoveMember(ctx, &adminv2.TenantServiceRemoveMemberRequest{
		Tenant: tenant,
		Member: id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to remove tenant member: %w", err)
	}

	return &adminv2.TenantServiceAddMemberResponse{}, nil
}

func (c *tenantMember) Convert(r *adminv2.TenantServiceAddMemberResponse) (string, *adminv2.TenantServiceAddMemberRequest, any, error) {
	// FIXME: from the response object we are unable to derive the add requests
	return "", &adminv2.TenantServiceAddMemberRequest{}, nil, nil
}

func (c *tenantMember) Update(rq any) (*adminv2.TenantServiceAddMemberResponse, error) {
	// TODO: api does not provide an update endpoint :(
	return nil, nil
}
