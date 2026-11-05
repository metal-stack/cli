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
	}

	machineLivelinessCmd := &cobra.Command{
		Use:   "machine-liveliness",
		Short: "execute machine liveliness housekeeping",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.machineLiveliness(cmd.Context())
		},
	}

	machineVPNConnectedCmd := &cobra.Command{
		Use:   "machine-vpn-connected",
		Short: "execute machine vpn connected housekeeping",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.machineVPNConnected(cmd.Context())
		},
	}

	machineResurrectCmd := &cobra.Command{
		Use:   "machine-resurrect",
		Short: "execute machine resurrect housekeeping",
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

func (h *housekeeping) Get(id string) (any, error) {
	panic("unimplemented")
}

func (h *housekeeping) List() (any, error) {
	panic("unimplemented")
}

func (h *housekeeping) Create(rq any) (any, error) {
	panic("unimplemented")
}

func (h *housekeeping) Delete(id string) (any, error) {
	panic("unimplemented")
}

func (h *housekeeping) Convert(r any) (string, any, any, error) {
	panic("unimplemented")
}

func (h *housekeeping) Update(rq any) (any, error) {
	panic("unimplemented")
}
