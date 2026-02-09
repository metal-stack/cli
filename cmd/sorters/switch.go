package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func SwitchSorter() *multisort.Sorter[*apiv2.Switch] {
	return multisort.New(multisort.FieldMap[*apiv2.Switch]{
		"id": func(a, b *apiv2.Switch, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"description": func(a, b *apiv2.Switch, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Description, b.Description, descending)
		},
		"partition": func(a, b *apiv2.Switch, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Partition, b.Partition, descending)
		},
		"rack": func(a, b *apiv2.Switch, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Rack), pointer.SafeDeref(b.Rack), descending)
		},
		"os": func(a, b *apiv2.Switch, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Os).Vendor, pointer.SafeDeref(b.Os).Vendor, descending)
		},
		"metal-core-version": func(a, b *apiv2.Switch, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Os).MetalCoreVersion, pointer.SafeDeref(b.Os).MetalCoreVersion, descending)
		},
		"management-ip": func(a, b *apiv2.Switch, descending bool) multisort.CompareResult {
			return multisort.Compare(a.ManagementIp, b.ManagementIp, descending)
		},
	}, multisort.Keys{{ID: "id"}})
}
