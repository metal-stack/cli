package tableprinters

import (
	"github.com/metal-stack/api/go/enum"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
)

func (t *TablePrinter) HousekeepingMachineLivelinessTable(data []*adminv2.MachineWithLiveliness, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"Machine", "Liveliness"}
	)

	for _, ms := range data {
		var lvl = "unknown"
		lvlString, err := enum.GetStringValue(ms.Liveliness)
		if err != nil {
			lvl = "unknown"
		}
		lvl = *lvlString

		rows = append(rows, []string{ms.Uuid, lvl})
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
