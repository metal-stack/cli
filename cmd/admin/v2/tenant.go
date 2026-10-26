package v2

import (
	"fmt"

	"github.com/metal-stack/api/go/errorutil"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/cli/pkg/helpers"
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

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.TenantServiceCreateRequest, *apiv2.TenantServiceUpdateRequest, *apiv2.Tenant]{
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
			cmd.Flags().String("login", "", "the login id of the tenant to create")
			cmd.Flags().String("name", "", "the name of the tenant to create")
			cmd.Flags().String("description", "", "the description of the tenant to create")
			cmd.Flags().String("email", "", "the email of the tenant to create")
			cmd.Flags().String("avatar-url", "", "the avatar url of the tenant to create")
		},
		CreateRequestFromCLI: func() (*adminv2.TenantServiceCreateRequest, error) {
			return &adminv2.TenantServiceCreateRequest{
				Login: pointer.PointerOrNil(viper.GetString("login")),
				Tenant: &apiv2.TenantServiceCreateRequest{
					Name:        viper.GetString("name"),
					Description: pointer.PointerOrNil(viper.GetString("description")),
					Email:       pointer.PointerOrNil(viper.GetString("email")),
					AvatarUrl:   pointer.PointerOrNil(viper.GetString("avatar-url")),
				},
			}, nil
		},
		ValidArgsFn: w.c.Completion.AdminTenant,
	}

	return genericcli.NewCmds(cmdsConfig, newTenantMemberCmd(c))
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

	req := &adminv2.TenantServiceListRequest{
		Query: &apiv2.TenantQuery{
			Name:  pointer.PointerOrNil(viper.GetString("name")),
			Login: pointer.PointerOrNil(viper.GetString("tenant")),
		},
	}

	if labelSlice := viper.GetStringSlice("labels"); len(labelSlice) > 0 {
		var err error

		req.Query.Labels, err = helpers.LabelsFromSlice(labelSlice)
		if err != nil {
			return nil, err
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

func (c *tenant) Convert(r *apiv2.Tenant) (string, *adminv2.TenantServiceCreateRequest, *apiv2.TenantServiceUpdateRequest, error) {
	return r.Login, &adminv2.TenantServiceCreateRequest{
			Login: pointer.PointerOrNil(r.Login),
			Tenant: &apiv2.TenantServiceCreateRequest{
				Name:        r.Name,
				Description: pointer.PointerOrNil(r.Description),
				Email:       pointer.PointerOrNil(r.Email),
				AvatarUrl:   pointer.PointerOrNil(r.AvatarUrl),
				Labels:      pointer.SafeDeref(r.Meta).Labels,
			},
		},
		&apiv2.TenantServiceUpdateRequest{
			Login:       r.Login,
			Name:        pointer.PointerOrNil(r.Name),
			Email:       pointer.PointerOrNil(r.Email),
			Description: pointer.PointerOrNil(r.Description),
			AvatarUrl:   pointer.PointerOrNil(r.AvatarUrl),
			UpdateMeta:  helpers.UpdateMetaFromMeta(r.Meta),
			Labels:      helpers.UpdateLabelsFromMeta(r.Meta),
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
