package sorters

import (
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
	}, multisort.Keys{{ID: "id"}})
}

func NetworkSorter() *multisort.Sorter[*apiv2.Network] {
	return multisort.New(multisort.FieldMap[*apiv2.Network]{
		"id": func(a, b *apiv2.Network, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"name": func(a, b *apiv2.Network, descending bool) multisort.CompareResult {
			return multisort.Compare(*a.Name, *b.Name, descending)
		},
		"type": func(a, b *apiv2.Network, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Type, b.Type, descending)
		},
	}, multisort.Keys{{ID: "id"}})
}

func MachineSorter() *multisort.Sorter[*apiv2.Machine] {
	return multisort.New(multisort.FieldMap[*apiv2.Machine]{
		"id": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
		"partition": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			ap := ""
			if a.Partition != nil {
				ap = a.Partition.Id
			}
			bp := ""
			if b.Partition != nil {
				bp = b.Partition.Id
			}
			return multisort.Compare(ap, bp, descending)
		},
		"liveliness": func(a, b *apiv2.Machine, descending bool) multisort.CompareResult {
			al := int32(0)
			bl := int32(0)
			if a.Status != nil {
				al = int32(a.Status.Liveliness)
			}
			if b.Status != nil {
				bl = int32(b.Status.Liveliness)
			}
			return multisort.Compare(al, bl, descending)
		},
	}, multisort.Keys{{ID: "id"}})
}

func FilesystemLayoutSorter() *multisort.Sorter[*apiv2.FilesystemLayout] {
	return multisort.New(multisort.FieldMap[*apiv2.FilesystemLayout]{
		"id": func(a, b *apiv2.FilesystemLayout, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
	}, multisort.Keys{{ID: "id"}})
}

func SizeReservationSorter() *multisort.Sorter[*apiv2.SizeReservation] {
	return multisort.New(multisort.FieldMap[*apiv2.SizeReservation]{
		"id": func(a, b *apiv2.SizeReservation, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"name": func(a, b *apiv2.SizeReservation, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"size": func(a, b *apiv2.SizeReservation, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Size, b.Size, descending)
		},
	}, multisort.Keys{{ID: "id"}})
}

func SizeImageConstraintSorter() *multisort.Sorter[*apiv2.SizeImageConstraint] {
	return multisort.New(multisort.FieldMap[*apiv2.SizeImageConstraint]{
		"size": func(a, b *apiv2.SizeImageConstraint, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Size, b.Size, descending)
		},
	}, multisort.Keys{{ID: "size"}})
}

func VpnNodeSorter() *multisort.Sorter[*apiv2.VPNNode] {
	return multisort.New(multisort.FieldMap[*apiv2.VPNNode]{
		"name": func(a, b *apiv2.VPNNode, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"project": func(a, b *apiv2.VPNNode, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Project, b.Project, descending)
		},
		"online": func(a, b *apiv2.VPNNode, descending bool) multisort.CompareResult {
			ai := int32(0)
			if a.Online {
				ai = 1
			}
			bi := int32(0)
			if b.Online {
				bi = 1
			}
			return multisort.Compare(ai, bi, descending)
		},
	}, multisort.Keys{{ID: "name"}})
}
