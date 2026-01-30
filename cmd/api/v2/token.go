package v2

import (
	"fmt"
	"strings"
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/durationpb"
)

type token struct {
	c *config.Config
}

func newTokenCmd(c *config.Config) *cobra.Command {
	w := &token{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv2.TokenServiceCreateRequest, *apiv2.TokenServiceUpdateRequest, *apiv2.Token]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "token",
		Plural:          "tokens",
		Description:     "manage api tokens",
		Sorter:          sorters.TokenSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		CreateRequestFromCLI: func() (*apiv2.TokenServiceCreateRequest, error) {
			var permissions []*apiv2.MethodPermission
			for _, r := range viper.GetStringSlice("permissions") {
				project, semicolonSeparatedMethods, ok := strings.Cut(r, "=")
				if !ok {
					return nil, fmt.Errorf("permissions must be provided in the form <project>=<methods-colon-separated>")
				}

				permissions = append(permissions, &apiv2.MethodPermission{
					Subject: project,
					Methods: strings.Split(semicolonSeparatedMethods, ":"),
				})
			}

			projectRoles := map[string]apiv2.ProjectRole{}
			for _, r := range viper.GetStringSlice("project-roles") {
				projectID, roleString, ok := strings.Cut(r, "=")
				if !ok {
					return nil, fmt.Errorf("project roles must be provided in the form <project-id>=<role>")
				}

				role, ok := apiv2.ProjectRole_value[roleString]
				if !ok {
					return nil, fmt.Errorf("unknown role: %s", roleString)
				}

				projectRoles[projectID] = apiv2.ProjectRole(role)
			}

			tenantRoles := map[string]apiv2.TenantRole{}
			for _, r := range viper.GetStringSlice("tenant-roles") {
				tenantID, roleString, ok := strings.Cut(r, "=")
				if !ok {
					return nil, fmt.Errorf("tenant roles must be provided in the form <tenant-id>=<role>")
				}

				role, ok := apiv2.TenantRole_value[roleString]
				if !ok {
					return nil, fmt.Errorf("unknown role: %s", roleString)
				}

				tenantRoles[tenantID] = apiv2.TenantRole(role)
			}

			var adminRole *apiv2.AdminRole
			if roleString := viper.GetString("admin-role"); roleString != "" {
				role, ok := apiv2.AdminRole_value[roleString]
				if !ok {
					return nil, fmt.Errorf("unknown role: %s", roleString)
				}

				adminRole = pointer.Pointer(apiv2.AdminRole(role))
			}

			return &apiv2.TokenServiceCreateRequest{
				// TODO: api should have an endpoint to list possible permissions and roles
				Description:  viper.GetString("description"),
				Permissions:  permissions,
				ProjectRoles: projectRoles,
				TenantRoles:  tenantRoles,
				AdminRole:    adminRole,
				Expires:      durationpb.New(viper.GetDuration("expires")),
			}, nil
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("description", "", "a short description for the intention to use this token for")
			cmd.Flags().StringSlice("permissions", nil, "the permissions to associate with the api token in the form <project>=<methods-colon-separated>")
			cmd.Flags().StringSlice("project-roles", nil, "the project roles to associate with the api token in the form <subject>=<role>")
			cmd.Flags().StringSlice("tenant-roles", nil, "the tenant roles to associate with the api token in the form <subject>=<role>")
			cmd.Flags().String("admin-role", "", "the admin role to associate with the api token")
			cmd.Flags().Duration("expires", 8*time.Hour, "the duration how long the api token is valid")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("permissions", c.Completion.TokenPermissionsCompletionfunc))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("project-roles", c.Completion.TokenProjectRolesCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant-roles", c.Completion.TokenTenantRolesCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("admin-role", c.Completion.TokenAdminRoleCompletion))
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Aliases = append(cmd.Aliases, "revoke")
		},
		ValidArgsFn: w.c.Completion.TokenListCompletion,
	}
	return genericcli.NewCmds(cmdsConfig)
}

func (c *token) Get(id string) (*apiv2.Token, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.TokenServiceGetRequest{
		Uuid: id,
	}

	resp, err := c.c.Client.Apiv2().Token().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return resp.GetToken(), nil
}

func (c *token) List() ([]*apiv2.Token, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.TokenServiceListRequest{}

	resp, err := c.c.Client.Apiv2().Token().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}

	return resp.GetTokens(), nil
}

func (c *token) Create(rq *apiv2.TokenServiceCreateRequest) (*apiv2.Token, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Token().Create(ctx, rq)
	if err != nil {
		return nil, err
	}

	_, _ = fmt.Fprintf(c.c.Out, "Make sure to copy your personal access token now as you will not be able to see this again.\n")
	_, _ = fmt.Fprintln(c.c.Out)
	_, _ = fmt.Fprintln(c.c.Out, resp.GetSecret())
	_, _ = fmt.Fprintln(c.c.Out)

	// TODO: allow printer in metal-lib to be silenced

	return resp.GetToken(), nil
}

func (c *token) Delete(id string) (*apiv2.Token, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.TokenServiceRevokeRequest{
		Uuid: id,
	}

	_, err := c.c.Client.Apiv2().Token().Revoke(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke token: %w", err)
	}

	return &apiv2.Token{
		Uuid: id,
	}, nil
}

func (c *token) Update(rq *apiv2.TokenServiceUpdateRequest) (*apiv2.Token, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Token().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update token: %w", err)
	}

	return resp.GetToken(), nil
}

func (c *token) Convert(r *apiv2.Token) (string, *apiv2.TokenServiceCreateRequest, *apiv2.TokenServiceUpdateRequest, error) {
	return r.Uuid, &apiv2.TokenServiceCreateRequest{
			Description:  r.GetDescription(),
			Permissions:  r.GetPermissions(),
			ProjectRoles: r.GetProjectRoles(),
			TenantRoles:  r.GetTenantRoles(),
			Expires:      durationpb.New(time.Until(r.GetExpires().AsTime())),
		}, &apiv2.TokenServiceUpdateRequest{
			Uuid:         r.Uuid,
			Description:  pointer.PointerOrNil(r.Description),
			Permissions:  r.Permissions,
			ProjectRoles: r.ProjectRoles,
			TenantRoles:  r.TenantRoles,
			AdminRole:    r.AdminRole,
		}, nil
}
