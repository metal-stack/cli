package v1

import (
	"github.com/metal-stack/cli/cmd/config"
	"github.com/spf13/cobra"
)

func AddCmds(cmd *cobra.Command, c *config.Config) {
	cmd.AddCommand(newVersionCmd(c))
	cmd.AddCommand(newHealthCmd(c))
	cmd.AddCommand(newTokenCmd(c))
	cmd.AddCommand(newIPCmd(c))
	cmd.AddCommand(newProjectCmd(c))
	cmd.AddCommand(newTenantCmd(c))
	cmd.AddCommand(newMethodsCmd(c))
	cmd.AddCommand(newUserCmd(c))
}
