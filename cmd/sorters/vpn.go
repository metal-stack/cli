package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func VpnNodeSorter() *multisort.Sorter[*apiv2.VPNNode] {
	return multisort.New(multisort.FieldMap[*apiv2.VPNNode]{
		"id": func(a, b *apiv2.VPNNode, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"name": func(a, b *apiv2.VPNNode, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"project": func(a, b *apiv2.VPNNode, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Project, b.Project, descending)
		},
	}, multisort.Keys{{ID: "id"}})
}
