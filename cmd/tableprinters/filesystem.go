package tableprinters

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (t *TablePrinter) FilesystemLayoutTable(data []*apiv2.FilesystemLayout, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"ID", "Name", "Description"}
	)

	for _, f := range data {
		name := ""
		if f.Name != nil {
			name = *f.Name
		}
		desc := ""
		if f.Description != nil {
			desc = *f.Description
		}

		rows = append(rows, []string{f.Id, name, desc})
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}

func (t *TablePrinter) SizeReservationTable(data []*apiv2.SizeReservation, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"ID", "Name", "Size", "Project", "Amount", "Partitions"}
	)

	for _, r := range data {
		rows = append(rows, []string{r.Id, r.Name, r.Size, r.Project, formatInt32(r.Amount), joinOrEmpty(r.Partitions)})
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}

func (t *TablePrinter) SizeImageConstraintTable(data []*apiv2.SizeImageConstraint, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"Size", "Name", "Description", "Image Constraints"}
	)

	for _, c := range data {
		name := ""
		if c.Name != nil {
			name = *c.Name
		}
		desc := ""
		if c.Description != nil {
			desc = *c.Description
		}

		constraints := []string{}
		for _, ic := range c.ImageConstraints {
			constraints = append(constraints, ic.String())
		}

		rows = append(rows, []string{c.Size, name, desc, joinOrEmpty(constraints)})
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}

func (t *TablePrinter) VPNNodeTable(data []*apiv2.VPNNode, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"ID", "Name", "Project", "Online", "Last Seen", "IPs"}
	)

	for _, n := range data {
		online := "no"
		if n.Online {
			online = "yes"
		}
		lastSeen := ""
		if n.LastSeen != nil {
			lastSeen = humanizeDuration(time.Since(n.LastSeen.AsTime())) + " ago"
		}

		rows = append(rows, []string{formatUint64(n.Id), n.Name, n.Project, online, lastSeen, joinOrEmpty(n.IpAddresses)})
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
