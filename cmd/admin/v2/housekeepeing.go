package v2

import (
	"context"
	"fmt"

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

	networksConsistentCmd := &cobra.Command{
		Use:   "networks-consistent",
		Short: "execute a read only check for network and ip inconsistencies",
		Long:  "identifies inconsistencies between metal networks and ips and ipam networks and ips",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.networksConsistent(cmd.Context())
		},
	}

	housekeepingCmd.AddCommand(machineLivelinessCmd)
	housekeepingCmd.AddCommand(machineVPNConnectedCmd)
	housekeepingCmd.AddCommand(machineResurrectCmd)
	housekeepingCmd.AddCommand(networksConsistentCmd)

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

func (h *housekeeping) networksConsistent(ctx context.Context) (err error) {
	result, err := h.c.Client.Adminv2().Housekeeping().NetworksConsistent(ctx, &adminv2.HousekeepingServiceNetworksConsistentRequest{})
	if err != nil {
		return err
	}

	_ = h.c.DescribePrinter.Print(result)

	if len(result.NetworksMissingInMetal) > 0 {
		fmt.Fprintf(h.c.Out, "run the following command inside of a pod of the ipam deployment to get the list of existing prefixes:\n")
		fmt.Fprintf(h.c.Out, "kubectl exec -it ipam-xxx-yyy -- /cli prefix list\n")
		fmt.Fprintln(h.c.Out)
		fmt.Fprintf(h.c.Out, "then select one of the networks listed in networksMissingInMetal and delete it in ipam with\n")
		fmt.Fprintf(h.c.Out, "kubectl exec -it ipam-xxx-yyy -- /cli prefix delete --cidr 1.2.3.0/24 [--namespace optional-namespace]\n")
		fmt.Fprintln(h.c.Out)
	}
	if len(result.IpsMissingInMetal) > 0 {
		fmt.Fprintln(h.c.Out)
		fmt.Fprintf(h.c.Out, "run the following command inside of a pod of the ipam deployment to get the list of existing ips:\n")
		fmt.Fprintf(h.c.Out, "kubectl exec -it ipam-xxx-yyy -- /cli ip list --prefix prefix-of-the-given-metal-network\n")
		fmt.Fprintln(h.c.Out)
		fmt.Fprintf(h.c.Out, "then select one of the ips listed in ipsMissingInMetal and delete it in ipam with\n")
		fmt.Fprintf(h.c.Out, "kubectl exec -it ipam-xxx-yyy -- /cli ip release --ip 1.2.3.4 --prefix 1.2.3.0/24 [--namespace optional-namespace]\n")
		fmt.Fprintln(h.c.Out)
	}

	return
}
