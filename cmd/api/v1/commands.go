package v1

import (
	clitypes "github.com/metal-stack/metal-lib/pkg/commands/types"
	"github.com/spf13/cobra"
)

func AddCmds(cmd *cobra.Command, c *clitypes.Config) {
	cmd.AddCommand(newVersionCmd(c))
	cmd.AddCommand(newHealthCmd(c))
	cmd.AddCommand(newTokenCmd(c))
	cmd.AddCommand(newIPCmd(c))
	cmd.AddCommand(newImageCmd(c))
	cmd.AddCommand(newProjectCmd(c))
	cmd.AddCommand(newTenantCmd(c))
	cmd.AddCommand(newMethodsCmd(c))
	cmd.AddCommand(newUserCmd(c))
}
