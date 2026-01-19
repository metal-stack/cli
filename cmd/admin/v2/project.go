package v2

import (
	"fmt"

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

type project struct {
	c *config.Config
}

func newProjectCmd(c *config.Config) *cobra.Command {
	w := &project{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv2.ProjectServiceCreateRequest, *apiv2.ProjectServiceUpdateRequest, *apiv2.Project]{
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
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *project) Get(id string) (*apiv2.Project, error) {
	panic("unimplemented")
}

func (c *project) List() ([]*apiv2.Project, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.ProjectServiceListRequest{
		Tenant: pointer.PointerOrNil(viper.GetString("tenant")),
	}

	resp, err := c.c.Client.Adminv2().Project().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return resp.GetProjects(), nil
}

func (c *project) Create(rq *apiv2.ProjectServiceCreateRequest) (*apiv2.Project, error) {
	panic("unimplemented")
}

func (c *project) Delete(id string) (*apiv2.Project, error) {
	panic("unimplemented")
}

func (c *project) Convert(r *apiv2.Project) (string, *apiv2.ProjectServiceCreateRequest, *apiv2.ProjectServiceUpdateRequest, error) {
	panic("unimplemented")
}

func (c *project) Update(rq *apiv2.ProjectServiceUpdateRequest) (*apiv2.Project, error) {
	panic("unimplemented")
}
