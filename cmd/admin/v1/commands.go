package v1

import (
	"github.com/spf13/cobra"

	clitypes "github.com/metal-stack/metal-lib/pkg/commands/types"
)

func AddCmds(cmd *cobra.Command, c *clitypes.Config) {
	adminCmd := &cobra.Command{
		Use:          "admin",
		Short:        "admin commands",
		Long:         "",
		SilenceUsage: true,
		Hidden:       true,
	}

	adminCmd.AddCommand(newTokenCmd(c))
	adminCmd.AddCommand(newImageCmd(c))

	cmd.AddCommand(adminCmd)
}
