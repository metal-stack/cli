package sorters

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func AuditSorter() *multisort.Sorter[*apiv2.AuditTrace] {
	return multisort.New(multisort.FieldMap[*apiv2.AuditTrace]{
		"id": func(a, b *apiv2.AuditTrace, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
		"timestamp": func(a, b *apiv2.AuditTrace, descending bool) multisort.CompareResult {
			return multisort.Compare(time.Time(a.Timestamp.AsTime()).Unix(), time.Time(b.Timestamp.AsTime()).Unix(), descending)
		},
		"user": func(a, b *apiv2.AuditTrace, descending bool) multisort.CompareResult {
			return multisort.Compare(a.User, b.User, descending)
		},
		"method": func(a, b *apiv2.AuditTrace, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Method, b.Method, descending)
		},
		"project": func(a, b *apiv2.AuditTrace, descending bool) multisort.CompareResult {
			return multisort.Compare(*a.Project, *b.Project, descending)
		},
	}, multisort.Keys{{ID: "timestamp", Descending: true}})
}
