package sorters

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func PartitionSorter() *multisort.Sorter[*apiv2.Partition] {
	return multisort.New(multisort.FieldMap[*apiv2.Partition]{
		"id": func(a, b *apiv2.Partition, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"description": func(a, b *apiv2.Partition, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Description, b.Description, descending)
		},
	}, multisort.Keys{{ID: "id"}, {ID: "description"}})
}

func PartitionCapacitySorter() *multisort.Sorter[*adminv2.PartitionCapacity] {
	return multisort.New(multisort.FieldMap[*adminv2.PartitionCapacity]{
		"id": func(a, b *adminv2.PartitionCapacity, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Partition, b.Partition, descending)
		},
		"size": func(a, b *adminv2.PartitionCapacity, descending bool) multisort.CompareResult {
			// FIXME implement
			return multisort.Compare(a.Partition, b.Partition, descending)
		},
	}, multisort.Keys{{ID: "id"}, {ID: "size"}})
}
