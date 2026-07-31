package v2

import (
	"fmt"
	"strings"
	"time"

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
	"google.golang.org/protobuf/types/known/durationpb"
)

type token struct {
	c *config.Config
}

func newTokenCmd(c *config.Config) *cobra.Command {
	w := &token{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.TokenServiceCreateRequest, *apiv2.TokenServiceUpdateRequest, *apiv2.Token]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "token",
		Plural:          "tokens",
		Description:     "manage api tokens for accessing the metal-stack.io api",
		Sorter:          sorters.TokenSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("user", "", "the uuid of the user to list the tokens for")
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("user", "", "the uuid of the user who owns the token")

			cmd.Aliases = append(cmd.Aliases, "revoke")
		},
		CreateRequestFromCLI: func() (*adminv2.TokenServiceCreateRequest, error) {
			var methodPermissions []*apiv2.MethodPermission
			for _, m := range viper.GetStringSlice("permissions") {
				subject, colonSeparatedMethods, ok := strings.Cut(m, "=")
				if !ok {
					colonSeparatedMethods = subject
				}

				methodPermissions = append(methodPermissions, &apiv2.MethodPermission{
					Subject: subject,
					Methods: strings.Split(colonSeparatedMethods, ":"),
				})
			}

			perms, err := helpers.ToPermissionsByVisibility(methodPermissions)
			if err != nil {
				return nil, err
			}

			projectRoles, err := helpers.ToProjectRoles(viper.GetStringSlice("project-roles"))
			if err != nil {
				return nil, err
			}

			tenantRoles, err := helpers.ToTenantRoles(viper.GetStringSlice("tenant-roles"))
			if err != nil {
				return nil, err
			}

			machineRoles, err := helpers.ToMachineRoles(viper.GetStringSlice("machine-roles"))
			if err != nil {
				return nil, err
			}

			adminRole, err := helpers.ToAdminRole(viper.GetString("admin-role"))
			if err != nil {
				return nil, err
			}

			infraRole, err := helpers.ToInfraRole(viper.GetString("infra-role"))
			if err != nil {
				return nil, err
			}

			var user *string
			if userString := viper.GetString("user"); userString != "" {
				user = new(userString)
			}

			return &adminv2.TokenServiceCreateRequest{
				User: user,
				TokenCreateRequest: &apiv2.TokenServiceCreateRequest{
					Description:  viper.GetString("description"),
					Permissions:  perms,
					ProjectRoles: projectRoles,
					TenantRoles:  tenantRoles,
					MachineRoles: machineRoles,
					AdminRole:    adminRole,
					InfraRole:    infraRole,
					Expires:      durationpb.New(viper.GetDuration("expires")),
				},
			}, nil
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("user", "", "user for which this token should be created")
			cmd.Flags().String("description", "", "a short description for the intention to use this token for")
			cmd.Flags().StringSlice("permissions", nil, "the permissions to associate with the api token in the form [<subject>=]<methods-colon-separated>")
			cmd.Flags().StringSlice("project-roles", nil, "the project roles to associate with the api token in the form <subject>=<role>")
			cmd.Flags().StringSlice("tenant-roles", nil, "the tenant roles to associate with the api token in the form <subject>=<role>")
			cmd.Flags().StringSlice("machine-roles", nil, "the machine roles to associate with the api token in the form <subject>=<role>")
			cmd.Flags().String("admin-role", "", "the admin role to associate with the api token")
			cmd.Flags().String("infra-role", "", "the infra role to associate with the api token")
			cmd.Flags().Duration("expires", 8*time.Hour, "the duration how long the api token is valid")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("user", c.Completion.Tenant))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("permissions", c.Completion.TokenPermissions))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("project-roles", c.Completion.TokenProjectRoles))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant-roles", c.Completion.TokenTenantRoles))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("admin-role", c.Completion.TokenAdminRole))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("infra-role", c.Completion.TokenInfraRole))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("machine-roles", c.Completion.TokenMachineRoles))
		},

		ValidArgsFn: w.c.Completion.Token,
	}
	return genericcli.NewCmds(cmdsConfig)
}

func (t *token) Get(id string) (*apiv2.Token, error) {
	ctx, cancel := t.c.NewRequestContext()
	defer cancel()

	// the admin token api does not have a token get and the one from API scopes it to self
	req := &adminv2.TokenServiceListRequest{
		Query: &apiv2.TokenQuery{
			Uuid: &id,
		},
	}

	resp, err := t.c.Client.Adminv2().Token().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	switch len(resp.Tokens) {
	case 0:
		return nil, errorutil.NotFound("token not found")
	case 1:
		return resp.Tokens[0], nil
	default:
		return nil, errorutil.Internal("token uuid exists multiple times")
	}
}

func (t *token) List() ([]*apiv2.Token, error) {
	ctx, cancel := t.c.NewRequestContext()
	defer cancel()

	req := &adminv2.TokenServiceListRequest{}

	if viper.IsSet("user") {
		req.Query = &apiv2.TokenQuery{
			User: new(viper.GetString("user")),
		}
	}

	resp, err := t.c.Client.Adminv2().Token().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}

	return resp.GetTokens(), nil
}

func (t *token) Create(rq *adminv2.TokenServiceCreateRequest) (*apiv2.Token, error) {
	ctx, cancel := t.c.NewRequestContext()
	defer cancel()

	resp, err := t.c.Client.Adminv2().Token().Create(ctx, rq)
	if err != nil {
		return nil, err
	}

	_, _ = fmt.Fprintf(t.c.Out, "Make sure to copy your personal access token now as you will not be able to see this again.\n")
	_, _ = fmt.Fprintln(t.c.Out)
	_, _ = fmt.Fprintln(t.c.Out, resp.GetSecret())
	_, _ = fmt.Fprintln(t.c.Out)

	// TODO: allow printer in metal-lib to be silenced

	return resp.GetToken(), nil
}

func (t *token) Delete(id string) (*apiv2.Token, error) {
	ctx, cancel := t.c.NewRequestContext()
	defer cancel()

	if !viper.IsSet("user") {
		return nil, fmt.Errorf("user is required to be set")
	}

	req := &adminv2.TokenServiceRevokeRequest{
		Uuid: id,
		User: viper.GetString("user"),
	}

	_, err := t.c.Client.Adminv2().Token().Revoke(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke token: %w", err)
	}

	return &apiv2.Token{
		Uuid: id,
	}, nil
}

func (t *token) Convert(r *apiv2.Token) (string, *adminv2.TokenServiceCreateRequest, *apiv2.TokenServiceUpdateRequest, error) {
	perms, err := helpers.ToPermissionsByVisibility(r.Permissions)
	if err != nil {
		return "", nil, nil, err
	}

	return r.Uuid, &adminv2.TokenServiceCreateRequest{
			User: &r.User,
			TokenCreateRequest: &apiv2.TokenServiceCreateRequest{
				Description:  r.GetDescription(),
				Permissions:  perms,
				ProjectRoles: r.GetProjectRoles(),
				TenantRoles:  r.GetTenantRoles(),
				Expires:      durationpb.New(time.Until(r.GetExpires().AsTime())),
				Labels:       pointer.SafeDeref(r.Meta).Labels,
			},
		}, &apiv2.TokenServiceUpdateRequest{
			Uuid:         r.Uuid,
			Description:  pointer.PointerOrNil(r.Description),
			Permissions:  perms,
			ProjectRoles: r.ProjectRoles,
			TenantRoles:  r.TenantRoles,
			AdminRole:    r.AdminRole,
			Labels: &apiv2.UpdateLabels{
				Strategy: &apiv2.UpdateLabels_Replace{
					Replace: &apiv2.Labels{
						Labels: pointer.SafeDeref(pointer.SafeDeref(r.Meta).Labels).Labels,
					},
				},
			},
			UpdateMeta: helpers.UpdateMetaFromMeta(r.Meta),
		}, nil
}

func (t *token) Update(rq *apiv2.TokenServiceUpdateRequest) (*apiv2.Token, error) {
	ctx, cancel := t.c.NewRequestContext()
	defer cancel()

	resp, err := t.c.Client.Apiv2().Token().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update token: %w", err)
	}

	return resp.Token, nil
}
