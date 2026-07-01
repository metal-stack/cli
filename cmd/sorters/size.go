package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func SizeSorter() *multisort.Sorter[*apiv2.Size] {
	return multisort.New(multisort.FieldMap[*apiv2.Size]{
		"id": func(a, b *apiv2.Size, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"name": func(a, b *apiv2.Size, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Name), pointer.SafeDeref(b.Name), descending)
		},
		"description": func(a, b *apiv2.Size, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Description), pointer.SafeDeref(b.Description), descending)
		},
	}, multisort.Keys{{ID: "id"}})
}
