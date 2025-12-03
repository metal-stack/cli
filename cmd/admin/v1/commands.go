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

	adminCmd.AddCommand(newImageCmd(c))
	adminCmd.AddCommand(newIPCmd(c))
	adminCmd.AddCommand(newNetworkCmd(c))
	adminCmd.AddCommand(newSizeCmd(c))
	adminCmd.AddCommand(newTenantCmd(c))
	adminCmd.AddCommand(newTokenCmd(c))
	adminCmd.AddCommand(newProjectCmd(c))

	cmd.AddCommand(adminCmd)
}
