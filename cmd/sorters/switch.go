package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func SwitchSorter() *multisort.Sorter[*apiv2.Switch] {
	return multisort.New(multisort.FieldMap[*apiv2.Switch]{
		"id": func(a, b *apiv2.Switch, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"description": func(a, b *apiv2.Switch, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Description, b.Description, descending)
		},
		// TODO: also allow sorting by partition, rack, os, status, last sync time, replace mode, metal-core version, ip
	}, multisort.Keys{{ID: "id"}})
}
