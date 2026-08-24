package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func ComponentSorter() *multisort.Sorter[*apiv2.Component] {
	return multisort.New(multisort.FieldMap[*apiv2.Component]{
		"type": func(a, b *apiv2.Component, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Type, b.Type, descending)
		},
		"identifier": func(a, b *apiv2.Component, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Identifier, b.Identifier, descending)
		},
		"started": func(a, b *apiv2.Component, descending bool) multisort.CompareResult {
			return multisort.Compare(a.StartedAt.AsTime().String(), b.StartedAt.AsTime().String(), descending)
		},
	}, multisort.Keys{{ID: "type", Descending: true}, {ID: "identifier", Descending: true}})
}
