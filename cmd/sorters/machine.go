package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func MachineSorter() *multisort.Sorter[*apiv2.Machine] {
	return multisort.New(multisort.FieldMap[*apiv2.Machine]{
		"partition": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Partition.Id, b.Partition.Id, descending)
		},
		"size": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Size.Id, b.Size.Id, descending)
		},
		"uuid": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
	}, multisort.Keys{{ID: "uuid"}, {ID: "size"}, {ID: "partition"}})
}
