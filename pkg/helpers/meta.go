package helpers

import apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"

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
