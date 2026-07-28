package helpers

import (
	"strings"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func ImageFeatureFromString(feature string) *apiv2.ImageFeature {
	if feature == "" {
		return nil
	}

	switch strings.ToLower(feature) {
	case "machine":
		return apiv2.ImageFeature_IMAGE_FEATURE_MACHINE.Enum()
	case "firewall":
		return apiv2.ImageFeature_IMAGE_FEATURE_FIREWALL.Enum()
	}
	return apiv2.ImageFeature_IMAGE_FEATURE_UNSPECIFIED.Enum()
}

func ImageFeaturesFromString(features []string) []apiv2.ImageFeature {
	var result []apiv2.ImageFeature

	for _, f := range features {
		switch strings.ToLower(f) {
		case "machine":
			result = append(result, apiv2.ImageFeature_IMAGE_FEATURE_MACHINE)
		case "firewall":
			result = append(result, apiv2.ImageFeature_IMAGE_FEATURE_FIREWALL)
		}
	}

	return result
}

func ImageClassificationFromString(classification string) apiv2.ImageClassification {
	switch strings.ToLower(strings.TrimSpace(classification)) {
	case "preview":
		return apiv2.ImageClassification_IMAGE_CLASSIFICATION_PREVIEW
	case "supported":
		return apiv2.ImageClassification_IMAGE_CLASSIFICATION_SUPPORTED
	case "deprecated":
		return apiv2.ImageClassification_IMAGE_CLASSIFICATION_DEPRECATED
	}

	return apiv2.ImageClassification_IMAGE_CLASSIFICATION_UNSPECIFIED
}
