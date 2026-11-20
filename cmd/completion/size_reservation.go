package completion

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) AdminSizeReservation(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Adminv2().SizeReservation().List(cmd.Context(), &adminv2.SizeReservationServiceListRequest{
		Query: &apiv2.SizeReservationQuery{},
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, r := range resp.SizeReservations {
		names = append(names, r.Id+"\t"+r.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
