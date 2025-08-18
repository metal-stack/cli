package v1

import (
	"github.com/metal-stack/cli/cmd/config"
	"github.com/spf13/cobra"
)

func AddCmds(cmd *cobra.Command, c *config.Config) {
	adminCmd := &cobra.Command{
		Use:          "admin",
		Short:        "admin commands",
		Long:         "",
		SilenceUsage: true,
		Hidden:       true,
	}

	adminCmd.AddCommand(newTokenCmd(c))
	adminCmd.AddCommand(newImageCmd(c))
	adminCmd.AddCommand(newMachineCmd(c))

	cmd.AddCommand(adminCmd)
}
