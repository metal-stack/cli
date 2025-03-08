package tableprinters

import (
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (t *TablePrinter) TenantTable(data []*apiv2.Tenant, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Name", "Email", "Registered", "Coupons", "Terms And Conditions"}
	if wide {
		header = []string{"ID", "Name", "Email", "Registered", "Coupons", "Terms And Conditions"}
	}

	for _, tenant := range data {
		id := tenant.Login
		name := tenant.Name
		email := tenant.Email
		since := humanize.Time(tenant.Meta.CreatedAt.AsTime())
		coupons := "-"
		couponsWide := coupons
		termsAndConditions := ""

		if wide {
			rows = append(rows, []string{id, name, email, since, couponsWide, termsAndConditions})
		} else {
			rows = append(rows, []string{id, name, email, since, coupons, termsAndConditions})
		}
	}

	return header, rows, nil
}

func (t *TablePrinter) TenantMemberTable(data []*apiv2.TenantMember, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "Role", "Since"}

	for _, member := range data {
		row := []string{
			member.Id,
			member.Role.String(),
			humanize.Time(member.CreatedAt.AsTime()),
		}

		rows = append(rows, row)
	}

	return header, rows, nil
}

func (t *TablePrinter) TenantInviteTable(data []*apiv2.TenantInvite, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"Secret", "Tenant", "Invited By", "Role", "Expires in"}

	for _, invite := range data {
		row := []string{
			invite.Secret,
			invite.TargetTenant,
			invite.Tenant,
			invite.Role.String(),
			humanize.Time(invite.ExpiresAt.AsTime()),
		}

		rows = append(rows, row)
	}

	t.t.MutateTable(func(table *tablewriter.Table) {
		table.SetAutoWrapText(false)
	})

	return header, rows, nil
}
