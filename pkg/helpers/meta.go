package helpers

import (
	"errors"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/viper"
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

// UpdateLabelsFromMeta returns update labels with replace strategy from a given meta.
func UpdateLabelsFromMeta(meta *apiv2.Meta) *apiv2.UpdateLabels {
	if meta == nil || meta.Labels == nil {
		return nil
	}

	return &apiv2.UpdateLabels{
		Strategy: &apiv2.UpdateLabels_Replace{
			Replace: &apiv2.Labels{
				Labels: pointer.SafeDeref(pointer.SafeDeref(meta).Labels).Labels,
			},
		},
	}
}

// UpdateLabelsFromCLI returns update labels from CLI args. Make sure you use --add-labels, --remove-labels and --labels in the command.
// Note that this function returns nil without an error in case none of the flags are used.
func UpdateLabelsFromCLI() (*apiv2.UpdateLabels, error) {
	var (
		addLabelSlice    = viper.GetStringSlice("add-labels")
		removeLabelSlice = viper.GetStringSlice("remove-labels")
	)

	if len(addLabelSlice) > 0 || len(removeLabelSlice) > 0 {
		if len(viper.GetStringSlice("labels")) > 0 {
			return nil, errors.New("--add-labels and --update-labels cannot be used in combination with --labels")
		}

		addLabels, err := LabelsFromSlice(addLabelSlice)
		if err != nil {
			return nil, err
		}

		return &apiv2.UpdateLabels{
			Strategy: &apiv2.UpdateLabels_Patch{
				Patch: &apiv2.LabelsPatch{
					Update: addLabels,
					Remove: removeLabelSlice,
				},
			},
		}, nil
	}

	if labelSlice := viper.GetStringSlice("labels"); len(labelSlice) > 0 {
		labels, err := LabelsFromSlice(labelSlice)
		if err != nil {
			return nil, err
		}

		return &apiv2.UpdateLabels{
			Strategy: &apiv2.UpdateLabels_Replace{
				Replace: labels,
			},
		}, nil
	}

	return nil, nil
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
