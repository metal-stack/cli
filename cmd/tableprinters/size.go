package tableprinters

import (
	"fmt"

	"github.com/dustin/go-humanize"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func (t *TablePrinter) SizeTable(data []*apiv2.Size, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"ID", "Name", "Description", "CPU Range", "Memory Range", "Storage Range", "GPU Range"}
	)

	for _, size := range data {
		var (
			cpu     string
			memory  string
			storage string
			gpu     string
		)

		for _, c := range size.Constraints {
			switch c.Type {
			case apiv2.SizeConstraintType_SIZE_CONSTRAINT_TYPE_CORES:
				cpu = fmt.Sprintf("%d - %d", c.Min, c.Max)
			case apiv2.SizeConstraintType_SIZE_CONSTRAINT_TYPE_MEMORY:
				memory = fmt.Sprintf("%s - %s", humanize.Bytes(uint64(c.Min)), humanize.Bytes(uint64(c.Max))) //nolint:gosec
			case apiv2.SizeConstraintType_SIZE_CONSTRAINT_TYPE_STORAGE:
				storage = fmt.Sprintf("%s - %s", humanize.Bytes(uint64(c.Min)), humanize.Bytes(uint64(c.Max))) //nolint:gosec
			case apiv2.SizeConstraintType_SIZE_CONSTRAINT_TYPE_GPU:
				gpu = fmt.Sprintf("%s: %d - %d", pointer.SafeDeref(c.Identifier), c.Min, c.Max)
			}

		}

		rows = append(rows, []string{size.Id, pointer.SafeDeref(size.Name), pointer.SafeDeref(size.Description), cpu, memory, storage, gpu})
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
