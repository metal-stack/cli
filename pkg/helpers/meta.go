package helpers

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
)

// UpdateMetaFromMeta returns an update for a given meta.
// When the updated at field is set, it uses client locking, otherwise server locking.
func UpdateMetaFromMeta(meta *apiv2.Meta) *apiv2.UpdateMeta {
	if meta == nil || meta.UpdatedAt == nil || meta.UpdatedAt.AsTime().IsZero() {
		return &apiv2.UpdateMeta{
			LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
		}
	}

	return &apiv2.UpdateMeta{
		LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_CLIENT,
		UpdatedAt:       meta.UpdatedAt,
	}
}

// LabelsFromSlice converts a given label slice (e.g. ["a=b"]) into apiv2 labels.
// When an empty slice is provided be aware that this function returns nil without an error.
func LabelsFromSlice(labelSlice []string) (*apiv2.Labels, error) {
	if len(labelSlice) == 0 {
		return nil, nil
	}

	labels, err := genericcli.LabelsToMap(labelSlice)
	if err != nil {
		return nil, err
	}

	return &apiv2.Labels{
		Labels: labels,
	}, nil

}
