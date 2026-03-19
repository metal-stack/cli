package api_e2e

import (
	"testing"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/tests/e2e"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

var (
	image1 = func() *apiv2.Image {
		return &apiv2.Image{
			Id:             "ubuntu-24.04",
			Name:           pointer.Pointer("Ubuntu 24.04"),
			Description:    pointer.Pointer("Ubuntu 24.04 LTS"),
			Features:       []apiv2.ImageFeature{apiv2.ImageFeature_IMAGE_FEATURE_MACHINE},
			Classification: apiv2.ImageClassification_IMAGE_CLASSIFICATION_SUPPORTED,
		}
	}
	image2 = func() *apiv2.Image {
		return &apiv2.Image{
			Id:             "firewall-3.0",
			Name:           pointer.Pointer("Firewall 3.0"),
			Description:    pointer.Pointer("Metal Firewall"),
			Features:       []apiv2.ImageFeature{apiv2.ImageFeature_IMAGE_FEATURE_FIREWALL},
			Classification: apiv2.ImageClassification_IMAGE_CLASSIFICATION_PREVIEW,
		}
	}
)

func Test_ImageCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.ImageServiceListResponse, *apiv2.Image]{
		{
			Name: "list",
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestClientConfig[apiv2.ImageServiceListRequest, apiv2.ImageServiceListResponse]{
				WantRequest: apiv2.ImageServiceListRequest{
					Query: &apiv2.ImageQuery{},
				},
				WantResponse: apiv2.ImageServiceListResponse{
					Images: []*apiv2.Image{
						image1(),
						image2(),
					},
				},
			}),
			CmdArgs: []string{"image", "list"},
			WantTable: new(`
ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS
ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
firewall-3.0  Firewall 3.0  Metal Firewall    firewall              preview
`),
			WantWideTable: new(`
ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS
ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
firewall-3.0  Firewall 3.0  Metal Firewall    firewall              preview
`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
ubuntu-24.04 Ubuntu 24.04
firewall-3.0 Firewall 3.0
			`),
			WantMarkdown: new(`
| ID           | NAME         | DESCRIPTION      | FEATURES | EXPIRATION | STATUS    |
|--------------|--------------|------------------|----------|------------|-----------|
| ubuntu-24.04 | Ubuntu 24.04 | Ubuntu 24.04 LTS | machine  |            | supported |
| firewall-3.0 | Firewall 3.0 | Metal Firewall   | firewall |            | preview   |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
