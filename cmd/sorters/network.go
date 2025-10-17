package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func NetworkSorter() *multisort.Sorter[*apiv2.Network] {
	return multisort.New(multisort.FieldMap[*apiv2.Network]{
		"id": func(a, b *apiv2.Network, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"name": func(a, b *apiv2.Network, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Name), pointer.SafeDeref(b.Name), descending)
		},
		"description": func(a, b *apiv2.Network, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Description), pointer.SafeDeref(b.Description), descending)
		},
		"partition": func(a, b *apiv2.Network, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Partition), pointer.SafeDeref(b.Partition), descending)
		},
		"project": func(a, b *apiv2.Network, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Project), pointer.SafeDeref(b.Project), descending)
		},
	}, multisort.Keys{{ID: "partition"}, {ID: "id"}})
}
