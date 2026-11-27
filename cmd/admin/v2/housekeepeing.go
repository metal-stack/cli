package v2

import (
	"context"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/spf13/cobra"
)

type housekeeping struct {
	c *config.Config
}

func newHousekeepingCmd(c *config.Config) *cobra.Command {
	w := &housekeeping{
		c: c,
	}

	housekeepingCmd := &cobra.Command{
		Use:   "housekeeping",
		Short: "execute housekeeping commands",
		Long:  "housekeeping commands are usually as cronjobs in the control-plane",
	}

	machineLivelinessCmd := &cobra.Command{
		Use:   "machine-liveliness",
		Short: "execute machine liveliness housekeeping",
		Long:  "updates the liveleness of all machines, only the modified machines are returned",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.machineLiveliness(cmd.Context())
		},
	}

	machineVPNConnectedCmd := &cobra.Command{
		Use:   "machine-vpn-connected",
		Short: "execute machine vpn connected housekeeping",
		Long:  "updates the vpn connection details of machines which have vpn configured, only the modified machines are returned",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.machineVPNConnected(cmd.Context())
		},
	}

	machineResurrectCmd := &cobra.Command{
		Use:   "machine-resurrect",
		Short: "execute machine resurrect housekeeping",
		Long:  "tries to resurrect machines which are not allocated but not in waiting state by forcing a restart, only machines where a restart was triggered are returned",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.machineResurrect(cmd.Context())
		},
	}

	housekeepingCmd.AddCommand(machineLivelinessCmd)
	housekeepingCmd.AddCommand(machineVPNConnectedCmd)
	housekeepingCmd.AddCommand(machineResurrectCmd)

	return housekeepingCmd
}

func (h *housekeeping) machineLiveliness(ctx context.Context) (err error) {
	result, err := h.c.Client.Adminv2().Housekeeping().MachineLiveliness(ctx, &adminv2.HousekeepingServiceMachineLivelinessRequest{})
	if err != nil {
		return err
	}

	return h.c.DescribePrinter.Print(result)
}

func (h *housekeeping) machineVPNConnected(ctx context.Context) (err error) {
	result, err := h.c.Client.Adminv2().Housekeeping().MachineVPNConnected(ctx, &adminv2.HousekeepingServiceMachineVPNConnectedRequest{})
	if err != nil {
		return err
	}

	return h.c.ListPrinter.Print(result.ModifiedMachines)
}

func (h *housekeeping) machineResurrect(ctx context.Context) (err error) {
	result, err := h.c.Client.Adminv2().Housekeeping().MachineResurrect(ctx, &adminv2.HousekeepingServiceMachineResurrectRequest{})
	if err != nil {
		return err
	}

	return h.c.ListPrinter.Print(result.ModifiedMachines)
}
