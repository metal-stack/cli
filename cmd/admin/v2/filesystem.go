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

type filesystem struct {
	c *config.Config
}

func newFilesystemCmd(c *config.Config) *cobra.Command {
	w := &filesystem{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.FilesystemServiceCreateRequest, *adminv2.FilesystemServiceUpdateRequest, *apiv2.FilesystemLayout]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "filesystem",
		Plural:          "filesystems",
		Description:     "manage filesystem layouts which define how disks, raid arrays and volumes are set up on machines during installation",
		Sorter:          sorters.FilesystemSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.FilesystemLayout,
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("id", "", "id of the filesystem layout to filter for")
		},
	}

	matchCmd := &cobra.Command{
		Use:   "match",
		Short: "matches a filesystem layout to a size and image or to a machine",
		Long: `matches a filesystem layout either to a given size and image combination,
or checks whether a given filesystem layout matches a given machine.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return w.match(cmd)
		},
	}
	matchCmd.Flags().String("size", "", "machine size to match against a filesystem layout")
	matchCmd.Flags().String("image", "", "machine image to match against a filesystem layout")
	matchCmd.Flags().String("machine", "", "uuid of the machine to check against a filesystem layout")
	matchCmd.Flags().String("filesystem-layout", "", "id of the filesystem layout to check against a machine")

	genericcli.Must(matchCmd.RegisterFlagCompletionFunc("size", c.Completion.Size))
	genericcli.Must(matchCmd.RegisterFlagCompletionFunc("image", c.Completion.Image))
	genericcli.Must(matchCmd.RegisterFlagCompletionFunc("machine", c.Completion.Machine))
	genericcli.Must(matchCmd.RegisterFlagCompletionFunc("filesystem-layout", c.Completion.FilesystemLayout))

	return genericcli.NewCmds(cmdsConfig, matchCmd)
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

func (c *filesystem) Create(rq *adminv2.FilesystemServiceCreateRequest) (*apiv2.FilesystemLayout, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Filesystem().Create(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, fmt.Errorf("failed to create filesystem layout: %w", err)
	}

	return resp.FilesystemLayout, nil
}

func (c *filesystem) Update(rq *adminv2.FilesystemServiceUpdateRequest) (*apiv2.FilesystemLayout, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Filesystem().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update filesystem layout: %w", err)
	}

	return resp.FilesystemLayout, nil
}

func (c *filesystem) Delete(id string) (*apiv2.FilesystemLayout, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Filesystem().Delete(ctx, &adminv2.FilesystemServiceDeleteRequest{
		Id: id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete filesystem layout: %w", err)
	}

	return resp.FilesystemLayout, nil
}

func (c *filesystem) Convert(r *apiv2.FilesystemLayout) (string, *adminv2.FilesystemServiceCreateRequest, *adminv2.FilesystemServiceUpdateRequest, error) {
	return r.Id,
		&adminv2.FilesystemServiceCreateRequest{
			FilesystemLayout: r,
		},
		&adminv2.FilesystemServiceUpdateRequest{
			Id:             r.Id,
			UpdateMeta:     helpers.UpdateMetaFromMeta(r.Meta),
			Name:           r.Name,
			Description:    r.Description,
			Filesystems:    r.Filesystems,
			Disks:          r.Disks,
			Raid:           r.Raid,
			VolumeGroups:   r.VolumeGroups,
			LogicalVolumes: r.LogicalVolumes,
			Constraints:    r.Constraints,
		},
		nil
}

func (w *filesystem) match(_ *cobra.Command) error {
	var (
		sizeSet    = viper.IsSet("size")
		imageSet   = viper.IsSet("image")
		machineSet = viper.IsSet("machine")
		layoutSet  = viper.IsSet("filesystem-layout")
		req        *adminv2.FilesystemServiceMatchRequest
	)

	switch {
	case sizeSet && imageSet:
		req = &adminv2.FilesystemServiceMatchRequest{
			Match: &adminv2.FilesystemServiceMatchRequest_SizeAndImage{
				SizeAndImage: &adminv2.MatchImageAndSize{
					Size:  viper.GetString("size"),
					Image: viper.GetString("image"),
				},
			},
		}
	case machineSet && layoutSet:
		req = &adminv2.FilesystemServiceMatchRequest{
			Match: &adminv2.FilesystemServiceMatchRequest_MachineAndFilesystemlayout{
				MachineAndFilesystemlayout: &adminv2.MatchMachineAndFilesystemLayout{
					Machine:          viper.GetString("machine"),
					FilesystemLayout: viper.GetString("filesystem-layout"),
				},
			},
		}
	default:
		return fmt.Errorf("either both --size and --image or both --machine and --filesystem-layout must be set")
	}

	ctx, cancel := w.c.NewRequestContext()
	defer cancel()

	resp, err := w.c.Client.Adminv2().Filesystem().Match(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to match filesystem layout: %w", err)
	}

	return w.c.ListPrinter.Print(resp.FilesystemLayout)
}
