package tableprinters

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (t *TablePrinter) PartitionTable(data []*apiv2.Partition, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"ID", "Description", "Mgmt Service Address"}
	)

	for _, p := range data {
		mgmtAddress := ""
		if len(p.MgmtServiceAddresses) > 0 {
			mgmtAddress = p.MgmtServiceAddresses[0]
		}

		rows = append(rows, []string{p.Id, p.Description, mgmtAddress})
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}

func (t *TablePrinter) PartitionCapacityTable(data []*adminv2.PartitionCapacity, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"Partition", "Size", "Free", "Total", "Allocated", "Faulty", "Reserved"}
	)

	for _, pc := range data {
		for _, sc := range pc.MachineSizeCapacities {
			rows = append(rows, []string{
				pc.Partition,
				sc.Size,
				formatInt64(sc.Free),
				formatInt64(sc.Total),
				formatInt64(sc.Allocated),
				formatInt64(sc.Faulty),
				formatInt64(sc.Other),
			})
		}
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
