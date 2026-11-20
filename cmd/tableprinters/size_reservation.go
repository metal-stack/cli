package tableprinters

import (
	"fmt"
	"sort"
	"strings"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
)

func (t *TablePrinter) SizeReservationTable(data []*apiv2.SizeReservation, wide bool) ([]string, [][]string, error) {
	var (
		header = []string{"ID", "Name", "Description", "Project", "Size", "Partitions", "Amount"}
		rows   [][]string
	)

	if wide {
		header = append(header, "Labels")
	}

	for _, r := range data {
		row := []string{
			r.Id,
			r.Name,
			r.Description,
			r.Project,
			r.Size,
			strings.Join(r.Partitions, "\n"),
			fmt.Sprintf("%d", r.Amount),
		}

		if wide {
			labels := genericcli.MapToLabels(r.Meta.GetLabels().GetLabels())
			sort.Strings(labels)
			row = append(row, strings.Join(labels, "\n"))
		}

		rows = append(rows, row)
	}

	return header, rows, nil
}
