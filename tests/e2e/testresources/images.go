package testresources

import apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"

var (
	Image1 = func() *apiv2.Image {
		return &apiv2.Image{
			Id:             "ubuntu-24.04",
			Name:           new("Ubuntu 24.04"),
			Description:    new("Ubuntu 24.04 LTS"),
			Features:       []apiv2.ImageFeature{apiv2.ImageFeature_IMAGE_FEATURE_MACHINE},
			Classification: apiv2.ImageClassification_IMAGE_CLASSIFICATION_SUPPORTED,
		}
	}
	Image2 = func() *apiv2.Image {
		return &apiv2.Image{
			Id:             "firewall-3.0",
			Name:           new("Firewall 3.0"),
			Description:    new("Metal Firewall"),
			Features:       []apiv2.ImageFeature{apiv2.ImageFeature_IMAGE_FEATURE_FIREWALL},
			Classification: apiv2.ImageClassification_IMAGE_CLASSIFICATION_PREVIEW,
		}
	}
)
