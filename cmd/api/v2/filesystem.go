package v2

import (
	"fmt"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type filesystem struct {
	c *config.Config
}

func newFilesystemCmd(c *config.Config) *cobra.Command {
	w := &filesystem{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.FilesystemLayout]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "filesystem",
		Plural:          "filesystems",
		Description:     "a filesystem layout defines how disks, raid arrays and volumes are set up on a machine during installation",
		Sorter:          sorters.FilesystemSorter(),
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DescribeCmd),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.FilesystemLayout,
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("id", "", "id of the filesystem layout which should be listed")
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *filesystem) Get(id string) (*apiv2.FilesystemLayout, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Filesystem().Get(ctx, &apiv2.FilesystemServiceGetRequest{
		Id: id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get filesystem layout: %w", err)
	}

	return resp.FilesystemLayout, nil
}

func (c *filesystem) List() ([]*apiv2.FilesystemLayout, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Filesystem().List(ctx, &apiv2.FilesystemServiceListRequest{
		Id: pointer.PointerOrNil(viper.GetString("id")),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list filesystem layouts: %w", err)
	}

	return resp.FilesystemLayouts, nil
}

func (c *filesystem) Create(rq any) (*apiv2.FilesystemLayout, error) {
	panic("unimplemented")
}

func (c *filesystem) Delete(id string) (*apiv2.FilesystemLayout, error) {
	panic("unimplemented")
}

func (c *filesystem) Convert(r *apiv2.FilesystemLayout) (string, any, any, error) {
	panic("unimplemented")
}

func (c *filesystem) Update(rq any) (*apiv2.FilesystemLayout, error) {
	panic("unimplemented")
}
