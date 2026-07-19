package v2

import (
	"fmt"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type adminFilesystemLayout struct {
	c *config.Config
}

func newAdminFilesystemLayoutCmd(c *config.Config) *cobra.Command {
	w := &adminFilesystemLayout{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, *adminv2.FilesystemServiceUpdateRequest, *apiv2.FilesystemLayout]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "filesystem-layout",
		Plural:          "filesystem-layouts",
		Description:     "manage filesystem layouts",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd, genericcli.DeleteCmd),
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *adminFilesystemLayout) Get(id string) (*apiv2.FilesystemLayout, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.FilesystemServiceGetRequest{Id: id}

	resp, err := c.c.Client.Apiv2().Filesystem().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get filesystem layout: %w", err)
	}

	return resp.FilesystemLayout, nil
}

func (c *adminFilesystemLayout) List() ([]*apiv2.FilesystemLayout, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.FilesystemServiceListRequest{}

	resp, err := c.c.Client.Apiv2().Filesystem().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list filesystem layouts: %w", err)
	}

	return resp.FilesystemLayouts, nil
}

func (c *adminFilesystemLayout) Create(rq any) (*apiv2.FilesystemLayout, error) {
	panic("unimplemented")
}

func (c *adminFilesystemLayout) Delete(id string) (*apiv2.FilesystemLayout, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Filesystem().Delete(ctx, &adminv2.FilesystemServiceDeleteRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("failed to delete filesystem layout: %w", err)
	}

	return resp.FilesystemLayout, nil
}

func (c *adminFilesystemLayout) Convert(r *apiv2.FilesystemLayout) (string, any, *adminv2.FilesystemServiceUpdateRequest, error) {
	return r.Id, nil, &adminv2.FilesystemServiceUpdateRequest{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Disks:       r.Disks,
		Constraints: r.Constraints,
	}, nil
}

func (c *adminFilesystemLayout) Update(rq *adminv2.FilesystemServiceUpdateRequest) (*apiv2.FilesystemLayout, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Filesystem().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update filesystem layout: %w", err)
	}

	return resp.FilesystemLayout, nil
}
