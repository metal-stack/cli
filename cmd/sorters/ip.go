package sorters

import (
	"net/netip"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func IPSorter() *multisort.Sorter[*apiv2.IP] {
	return multisort.New(multisort.FieldMap[*apiv2.IP]{
		"ip": func(a, b *apiv2.IP, descending bool) multisort.CompareResult {
			aIP, _ := netip.ParseAddr(a.Ip)
			bIP, _ := netip.ParseAddr(b.Ip)
			return multisort.WithCompareFunc(func() int {
				return aIP.Compare(bIP)
			}, descending)
		},
		"name": func(a, b *apiv2.IP, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"project": func(a, b *apiv2.IP, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Project, b.Project, descending)
		},
		"type": func(a, b *apiv2.IP, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Type, b.Type, descending)
		},
		"network": func(a, b *apiv2.IP, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Network, b.Network, descending)
		},
		"uuid": func(a, b *apiv2.IP, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
	}, multisort.Keys{{ID: "project"}, {ID: "ip"}})
}
