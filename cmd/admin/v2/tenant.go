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

type tenant struct {
	c *config.Config
}

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

	return genericcli.NewCmds(cmdsConfig, newAddMemberCmd(c))
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
