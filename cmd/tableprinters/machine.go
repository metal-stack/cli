package tableprinters

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (t *TablePrinter) MachineTable(data []*apiv2.Machine, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"ID", "Partition", "Size", "Hostname", "Liveliness"}
	)

	if wide {
		header = []string{"ID", "Partition", "Size", "Hostname", "Project", "Liveliness", "State", "Created"}
	}

	for _, m := range data {
		size := ""
		if m.Size != nil {
			size = m.Size.Id
		}
		hostname := ""
		if m.Allocation != nil {
			hostname = m.Allocation.Hostname
		}
		liveliness := ""
		state := ""
		if m.Status != nil {
			liveliness = m.Status.Liveliness.String()
			if m.Status.Condition != nil {
				state = m.Status.Condition.State.String()
			}
		}
		partition := ""
		if m.Partition != nil {
			partition = m.Partition.Id
		}
		project := ""
		if m.Allocation != nil {
			project = m.Allocation.Project
		}
		created := ""
		if m.Meta != nil {
			created = humanizeDuration(time.Since(m.Meta.CreatedAt.AsTime())) + " ago"
		}

		if wide {
			rows = append(rows, []string{m.Uuid, partition, size, hostname, project, liveliness, state, created})
		} else {
			rows = append(rows, []string{m.Uuid, partition, size, hostname, liveliness})
		}
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
