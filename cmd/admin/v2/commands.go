package v2

import (
	"github.com/metal-stack/cli/cmd/config"
	"github.com/spf13/cobra"
)

func AddCmds(cmd *cobra.Command, c *config.Config) {
	adminCmd := &cobra.Command{
		Use:          "admin",
		Short:        "admin commands",
		Long:         "these commands utilize the admin api, which can only be accessed by metal-stack operators.",
		SilenceUsage: true,
		Hidden:       true,
	}

	adminCmd.AddCommand(newAdminFilesystemLayoutCmd(c))
	adminCmd.AddCommand(newAdminImageUsageCmd(c))
	adminCmd.AddCommand(newAdminMachineCmd(c))
	adminCmd.AddCommand(newAdminNetworkCmd(c))
	adminCmd.AddCommand(newAdminPartitionCmd(c))
	adminCmd.AddCommand(newAdminSizeImageConstraintCmd(c))
	adminCmd.AddCommand(newAdminSizeReservationCmd(c))
	adminCmd.AddCommand(newAdminTokenCreateCmd(c))
	adminCmd.AddCommand(newAdminVPNCmd(c))
	adminCmd.AddCommand(newAuditCmd(c))
	adminCmd.AddCommand(newComponentCmd(c))
	adminCmd.AddCommand(newImageCmd(c))
	adminCmd.AddCommand(newProjectCmd(c))
	adminCmd.AddCommand(newSizeCmd(c))
	adminCmd.AddCommand(newSwitchCmd(c))
	adminCmd.AddCommand(newTaskCmd(c))
	adminCmd.AddCommand(newTenantCmd(c))
	adminCmd.AddCommand(newTokenCmd(c))

	cmd.AddCommand(adminCmd)
}
