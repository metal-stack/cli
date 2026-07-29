package helpers

import (
	"strings"

	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func ImageFeatureFromString(feature string) (*apiv2.ImageFeature, error) {
	if feature == "" {
		return nil, nil
	}

	e, err := enum.GetEnum[apiv2.ImageFeature](strings.ToLower(feature))
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func ImageFeaturesFromStringSlice(features []string) ([]apiv2.ImageFeature, error) {
	var result []apiv2.ImageFeature

	for _, f := range features {
		e, err := enum.GetEnum[apiv2.ImageFeature](strings.ToLower(f))
		if err != nil {
			return nil, err
		}

		result = append(result, e)
	}

	return result, nil
}

func ImageClassificationFromString(classification string) (apiv2.ImageClassification, error) {
	return enum.GetEnum[apiv2.ImageClassification](strings.ToLower(classification))
}
