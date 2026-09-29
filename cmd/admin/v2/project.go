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

type project struct {
	c *config.Config
}

func newProjectCmd(c *config.Config) *cobra.Command {
	w := &project{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.ProjectServiceCreateRequest, *apiv2.ProjectServiceUpdateRequest, *apiv2.Project]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "project",
		Plural:          "projects",
		Description:     "manage api projects",
		Sorter:          sorters.ProjectSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("tenant", "", "lists only projects with the given tenant")
			cmd.Flags().StringSlice("labels", nil, "lists only projects with the given labels")
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *project) Get(id string) (*apiv2.Project, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.ProjectServiceGetRequest{
		Project: id,
	}

	resp, err := c.c.Client.Apiv2().Project().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return resp.GetProject(), nil
}

func (c *project) List() ([]*apiv2.Project, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.ProjectServiceListRequest{
		Query: &apiv2.ProjectQuery{
			Tenant: pointer.PointerOrNil(viper.GetString("tenant")),
		},
	}

	if labelSlice := viper.GetStringSlice("labels"); len(labelSlice) > 0 {
		var err error

		req.Query.Labels, err = helpers.LabelsFromSlice(labelSlice)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.c.Client.Adminv2().Project().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return resp.GetProjects(), nil
}

func (c *project) Create(rq *adminv2.ProjectServiceCreateRequest) (*apiv2.Project, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Project().Create(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, err
	}

	return resp.Project, nil
}

func (c *project) Delete(id string) (*apiv2.Project, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Project().Delete(ctx, &apiv2.ProjectServiceDeleteRequest{
		Project: id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete project: %w", err)
	}

	return resp.Project, nil
}

func (c *project) Convert(r *apiv2.Project) (string, *adminv2.ProjectServiceCreateRequest, *apiv2.ProjectServiceUpdateRequest, error) {
	return r.Uuid, &adminv2.ProjectServiceCreateRequest{
		Login:       r.Tenant,
		Name:        r.Name,
		Description: r.Description,
		AvatarUrl:   r.AvatarUrl,
		Labels:      pointer.SafeDeref(r.Meta).Labels,
		Project:     pointer.PointerOrNil(r.Uuid),
	}, &apiv2.ProjectServiceUpdateRequest{
		Project:     r.Uuid,
		Name:        pointer.PointerOrNil(r.Name),
		Description: pointer.PointerOrNil(r.Description),
		AvatarUrl:   r.AvatarUrl,
		UpdateMeta:  helpers.UpdateMetaFromMeta(r.Meta),
		Labels:      helpers.UpdateLabelsFromMeta(r.Meta),
	}, nil
}

func (c *project) Update(rq *apiv2.ProjectServiceUpdateRequest) (*apiv2.Project, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Project().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return resp.Project, nil
}
