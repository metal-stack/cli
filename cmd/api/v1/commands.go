package v1

import (
	"github.com/metal-stack/cli/cmd/config"
	"github.com/spf13/cobra"
)

func AddCmds(cmd *cobra.Command, c *config.Config) {
	cmd.AddCommand(newHealthCmd(c))
	cmd.AddCommand(newImageCmd(c))
	cmd.AddCommand(newIPCmd(c))
	cmd.AddCommand(newMethodsCmd(c))
	cmd.AddCommand(newNetworkCmd(c))
	cmd.AddCommand(newProjectCmd(c))
	cmd.AddCommand(newSizeCmd(c))
	cmd.AddCommand(newTenantCmd(c))
	cmd.AddCommand(newTokenCmd(c))
	cmd.AddCommand(newUserCmd(c))
	cmd.AddCommand(newVersionCmd(c))
}
