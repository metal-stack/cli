package tableprinters

import (
	"fmt"
	"sort"
	"strings"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
)

func (t *TablePrinter) PartitionTable(data []*apiv2.Partition, wide bool) ([]string, [][]string, error) {
	var (
		header = []string{"ID", "Description"}
		rows   [][]string
	)

	if wide {
		header = []string{"ID", "Description", "Labels"}
	}

	for _, p := range data {
		row := []string{p.Id, p.Description}

		if wide {
			labels := genericcli.MapToLabels(p.Meta.Labels.Labels)
			sort.Strings(labels)
			row = append(row, strings.Join(labels, "\n"))
		}

		rows = append(rows, row)
	}

	return header, rows, nil
}

func (t *TablePrinter) PartitionCapacityTable(data []*adminv2.PartitionCapacity, wide bool) ([]string, [][]string, error) {
	var (
		header = []string{"Partition", "Size", "Allocated", "Free", "Unavailable", "Reservations", "|", "Total", "|", "Faulty"}
		rows   [][]string

		allocatedCount       int64
		faultyCount          int64
		freeCount            int64
		otherCount           int64
		phonedHomeCount      int64
		reservationCount     int64
		totalCount           int64
		unavailableCount     int64
		usedReservationCount int64
		waitingCount         int64
	)

	if wide {
		header = append(header, "Phoned Home", "Waiting", "Other")
	}

	for _, pc := range data {
		for _, c := range pc.MachineSizeCapacities {
			id := c.Size
			var (
				allocated    = fmt.Sprintf("%d", c.Allocated)
				faulty       = fmt.Sprintf("%d", c.Faulty)
				free         = fmt.Sprintf("%d", c.Free)
				other        = fmt.Sprintf("%d", c.Other)
				phonedHome   = fmt.Sprintf("%d", c.PhonedHome)
				reservations = "0"
				total        = fmt.Sprintf("%d", c.Total)
				unavailable  = fmt.Sprintf("%d", c.Unavailable)
				waiting      = fmt.Sprintf("%d", c.Waiting)
			)

			if c.Reservations > 0 {
				reservations = fmt.Sprintf("%d (%d/%d used)", c.Reservations-c.UsedReservations, c.UsedReservations, c.Reservations)
			}

			allocatedCount += c.Allocated
			faultyCount += c.Faulty
			freeCount += c.Free
			otherCount += c.Other
			phonedHomeCount += c.PhonedHome
			reservationCount += c.Reservations
			totalCount += c.Total
			unavailableCount += c.Unavailable
			usedReservationCount += c.UsedReservations
			waitingCount += c.Waiting

			row := []string{pc.Partition, id, allocated, free, unavailable, reservations, "|", total, "|", faulty}
			if wide {
				row = append(row, phonedHome, waiting, other)
			}

			rows = append(rows, row)
		}
	}

	footerRow := ([]string{
		"Total",
		"",
		fmt.Sprintf("%d", allocatedCount),
		fmt.Sprintf("%d", freeCount),
		fmt.Sprintf("%d", unavailableCount),
		fmt.Sprintf("%d", reservationCount-usedReservationCount),
		"|",
		fmt.Sprintf("%d", totalCount),
		"|",
		fmt.Sprintf("%d", faultyCount),
	})

	if wide {
		footerRow = append(footerRow, []string{
			fmt.Sprintf("%d", phonedHomeCount),
			fmt.Sprintf("%d", waitingCount),
			fmt.Sprintf("%d", otherCount),
		}...)
	}

	// if t.markdown {
	// 	// for markdown we already have enough dividers, remove them
	// 	removeDivider := func(e string) bool {
	// 		return e == "|"
	// 	}
	// 	header = slices.DeleteFunc(header, removeDivider)
	// 	footerRow = slices.DeleteFunc(footerRow, removeDivider)
	// 	for i, row := range rows {
	// 		rows[i] = slices.DeleteFunc(row, removeDivider)
	// 	}
	// }

	rows = append(rows, footerRow)

	return header, rows, nil
}
