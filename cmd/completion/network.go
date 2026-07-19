package completion

import (
	"github.com/spf13/cobra"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (c *Completion) PartitionListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.PartitionServiceListRequest{}
	resp, err := c.Client.Apiv2().Partition().List(c.Ctx, req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.GetPartitions() {
		names = append(names, s.Id+"\t"+s.Description)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) NetworkListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.NetworkServiceListRequest{}
	resp, err := c.Client.Apiv2().Network().List(c.Ctx, req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.GetNetworks() {
		name := ""
		if s.Name != nil {
			name = *s.Name
		}
		names = append(names, s.Id+"\t"+name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) MachineListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.MachineServiceListRequest{}
	resp, err := c.Client.Apiv2().Machine().List(c.Ctx, req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, m := range resp.GetMachines() {
		partition := ""
		if m.Partition != nil {
			partition = m.Partition.Id
		}
		names = append(names, m.Uuid+"\t"+partition)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) FilesystemLayoutListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.FilesystemServiceListRequest{}
	resp, err := c.Client.Apiv2().Filesystem().List(c.Ctx, req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.GetFilesystemLayouts() {
		name := ""
		if s.Name != nil {
			name = *s.Name
		}
		names = append(names, s.Id+"\t"+name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) SizeReservationListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv2.SizeReservationServiceListRequest{}
	resp, err := c.Client.Apiv2().SizeReservation().List(c.Ctx, req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.GetSizeReservations() {
		names = append(names, s.Id+"\t"+s.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) AdminSizeReservationListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &adminv2.SizeReservationServiceListRequest{}
	resp, err := c.Client.Adminv2().SizeReservation().List(c.Ctx, req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.GetSizeReservations() {
		names = append(names, s.Id+"\t"+s.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) SizeImageConstraintListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &adminv2.SizeImageConstraintServiceListRequest{}
	resp, err := c.Client.Adminv2().SizeImageConstraint().List(c.Ctx, req)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.GetSizeImageConstraints() {
		names = append(names, s.Size+"\t"+s.String())
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
