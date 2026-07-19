package v2

import (
	"fmt"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type filesystemLayout struct {
	c *config.Config
}

func newFilesystemLayoutCmd(c *config.Config) *cobra.Command {
	w := &filesystemLayout{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.FilesystemLayout]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "filesystem-layout",
		Plural:          "filesystem-layouts",
		Description:     "read filesystem layouts for machine disk partitioning",
		Sorter:          sorters.FilesystemLayoutSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.FilesystemLayoutListCompletion,
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *filesystemLayout) Get(id string) (*apiv2.FilesystemLayout, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.FilesystemServiceGetRequest{Id: id}

	resp, err := c.c.Client.Apiv2().Filesystem().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get filesystem layout: %w", err)
	}

	return resp.FilesystemLayout, nil
}

func (c *filesystemLayout) List() ([]*apiv2.FilesystemLayout, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.FilesystemServiceListRequest{}

	resp, err := c.c.Client.Apiv2().Filesystem().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list filesystem layouts: %w", err)
	}

	return resp.FilesystemLayouts, nil
}

func (c *filesystemLayout) Create(rq any) (*apiv2.FilesystemLayout, error) {
	panic("unimplemented")
}

func (c *filesystemLayout) Delete(id string) (*apiv2.FilesystemLayout, error) {
	panic("unimplemented")
}

func (c *filesystemLayout) Convert(r *apiv2.FilesystemLayout) (string, any, any, error) {
	panic("unimplemented")
}

func (c *filesystemLayout) Update(rq any) (*apiv2.FilesystemLayout, error) {
	panic("unimplemented")
}
