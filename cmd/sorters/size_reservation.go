package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func SizeReservationSorter() *multisort.Sorter[*apiv2.SizeReservation] {
	return multisort.New(multisort.FieldMap[*apiv2.SizeReservation]{
		"id": func(a, b *apiv2.SizeReservation, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"name": func(a, b *apiv2.SizeReservation, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"description": func(a, b *apiv2.SizeReservation, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Description, b.Description, descending)
		},
		"project": func(a, b *apiv2.SizeReservation, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Project, b.Project, descending)
		},
		"size": func(a, b *apiv2.SizeReservation, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Size, b.Size, descending)
		},
		"amount": func(a, b *apiv2.SizeReservation, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Amount, b.Amount, descending)
		},
	}, multisort.Keys{{ID: "id"}})
}
