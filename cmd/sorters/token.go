package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func TokenSorter() *multisort.Sorter[*apiv2.Token] {
	return multisort.New(multisort.FieldMap[*apiv2.Token]{
		"id": func(a, b *apiv2.Token, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
		"user": func(a, b *apiv2.Token, descending bool) multisort.CompareResult {
			return multisort.Compare(a.UserId, b.UserId, descending)
		},
		"type": func(a, b *apiv2.Token, descending bool) multisort.CompareResult {
			return multisort.Compare(a.TokenType.String(), b.TokenType.String(), descending)
		},
		"description": func(a, b *apiv2.Token, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Description, b.Description, descending)
		},
		"expires": func(a, b *apiv2.Token, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Expires.AsTime().UnixMilli(), b.Expires.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "type"}, {ID: "user"}, {ID: "expires"}, {ID: "id"}})
}
