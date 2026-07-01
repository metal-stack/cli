package sorters

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func TaskSorter() *multisort.Sorter[*adminv2.TaskInfo] {
	return multisort.New(multisort.FieldMap[*adminv2.TaskInfo]{
		"id": func(a, b *adminv2.TaskInfo, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"queue": func(a, b *adminv2.TaskInfo, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Queue, b.Queue, descending)
		},
		"type": func(a, b *adminv2.TaskInfo, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Type, b.Type, descending)
		},
		"state": func(a, b *adminv2.TaskInfo, descending bool) multisort.CompareResult {
			return multisort.Compare(adminv2.TaskState_name[int32(a.State)], adminv2.TaskState_name[int32(b.State)], descending)
		},
		"retried": func(a, b *adminv2.TaskInfo, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Retried, b.Retried, descending)
		},
		"completed-at": func(a, b *adminv2.TaskInfo, descending bool) multisort.CompareResult {
			return multisort.Compare(a.CompletedAt.AsTime().UnixMilli(), b.CompletedAt.AsTime().UnixMilli(), descending)
		},
		"last-failed-at": func(a, b *adminv2.TaskInfo, descending bool) multisort.CompareResult {
			return multisort.Compare(a.LastFailedAt.AsTime().UnixMilli(), b.LastFailedAt.AsTime().UnixMilli(), descending)
		},
		"deadline-at": func(a, b *adminv2.TaskInfo, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Deadline.AsTime().UnixMilli(), b.Deadline.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "id"}})
}
