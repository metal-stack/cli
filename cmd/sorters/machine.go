package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func MachineSorter() *multisort.Sorter[*apiv2.Machine] {
	return multisort.New(multisort.FieldMap[*apiv2.Machine]{
		"partition": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Partition).Id, pointer.SafeDeref(b.Partition).Id, descending)
		},
		"size": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Size).Id, pointer.SafeDeref(b.Size).Id, descending)
		},
		"uuid": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
		"image": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(pointer.SafeDeref(a.Allocation).Image).Id, pointer.SafeDeref(pointer.SafeDeref(b.Allocation).Image).Id, descending)
		},
		"rack": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Rack, b.Rack, descending)
		},
		"project": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Allocation).Project, pointer.SafeDeref(b.Allocation).Project, descending)
		},
		"age": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			return multisort.Compare(
				new(pointer.SafeDeref(pointer.SafeDeref(pointer.SafeDeref(a.Allocation).Meta).CreatedAt)).AsTime().Unix(),
				new(pointer.SafeDeref(pointer.SafeDeref(pointer.SafeDeref(b.Allocation).Meta).CreatedAt)).AsTime().Unix(),
				descending,
			)
		},
	}, multisort.Keys{{ID: "partition"}, {ID: "size"}, {ID: "project"}, {ID: "uuid"}})
}
