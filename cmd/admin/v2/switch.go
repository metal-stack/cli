package v1

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/multisort"
	"github.com/spf13/cobra"
)

type switchCmd struct {
	c *config.Config
}

func newSwitchCmd(c *config.Config) *cobra.Command {
	sw := &switchCmd{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, *adminv2.SwitchServiceUpdateRequest, *apiv2.Switch]{
		GenericCLI: genericcli.NewGenericCLI(sw).WithFS(c.Fs),
		OnlyCmds: genericcli.OnlyCmds(
			genericcli.DescribeCmd,
			genericcli.ListCmd,
			genericcli.UpdateCmd,
			genericcli.DeleteCmd,
			genericcli.EditCmd,
		),
		BinaryName:      config.BinaryName,
		Singular:        "switch",
		Plural:          "switches",
		Description:     "view and manage network switches",
		Aliases:         []string{"sw"},
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		Sorter:          &multisort.Sorter[*apiv2.Switch]{},
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("id", "", "ID of the switch.")
			cmd.Flags().String("name", "", "Name of the switch.")
			cmd.Flags().String("os-vendor", "", "OS vendor of this switch.")
			cmd.Flags().String("os-version", "", "OS version of this switch.")
			cmd.Flags().String("partition", "", "Partition of this switch.")
			cmd.Flags().String("rack", "", "Rack of this switch.")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("id", c.Completion.SwitchListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("name", c.Completion.SwitchNameListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("rack", c.Completion.SwitchRackListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("os-vendor", c.Completion.SwitchOSVendorListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("os-version", c.Completion.SwitchOSVersionListCompletion))
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().Bool("force", false, "forcefully delete the switch accepting the risk that it still has machines connected to it")
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *switchCmd) Get(id string) (*apiv2.Switch, error) {
	panic("unimplemented")
}

func (c *switchCmd) List() ([]*apiv2.Switch, error) {
	panic("unimplemented")
}

func (c *switchCmd) Create(rq any) (*apiv2.Switch, error) {
	panic("unimplemented")
}

func (c *switchCmd) Delete(id string) (*apiv2.Switch, error) {
	panic("unimplemented")
}

func (c *switchCmd) Convert(sw *apiv2.Switch) (string, any, *adminv2.SwitchServiceUpdateRequest, error) {
	panic("unimplemented")
}

func (c *switchCmd) Update(rq *adminv2.SwitchServiceUpdateRequest) (*apiv2.Switch, error) {
	panic("unimplemented")
}
