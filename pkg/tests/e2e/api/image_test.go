package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/tests/e2e"
)

var (
	image1 = func() *apiv2.Image {
		return &apiv2.Image{
			Id:             "ubuntu-24.04",
			Name:           new("Ubuntu 24.04"),
			Description:    new("Ubuntu 24.04 LTS"),
			Features:       []apiv2.ImageFeature{apiv2.ImageFeature_IMAGE_FEATURE_MACHINE},
			Classification: apiv2.ImageClassification_IMAGE_CLASSIFICATION_SUPPORTED,
		}
	}
	image2 = func() *apiv2.Image {
		return &apiv2.Image{
			Id:             "firewall-3.0",
			Name:           new("Firewall 3.0"),
			Description:    new("Metal Firewall"),
			Features:       []apiv2.ImageFeature{apiv2.ImageFeature_IMAGE_FEATURE_FIREWALL},
			Classification: apiv2.ImageClassification_IMAGE_CLASSIFICATION_PREVIEW,
		}
	}
)

func Test_ImageCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.ImageServiceListResponse, *apiv2.Image]{
		{
			Name:    "list",
			CmdArgs: []string{"image", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.ImageServiceListRequest{
							Query: &apiv2.ImageQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ImageServiceListResponse{
								Images: []*apiv2.Image{
									image1(),
									image2(),
								},
							})
						},
					},
				},
			},
			),
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

func Test_ImageCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.ImageServiceGetResponse, *apiv2.Image]{
		{
			Name:    "describe",
			CmdArgs: []string{"image", "describe", image1().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.ImageServiceGetRequest{
							Id: image1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ImageServiceGetResponse{
								Image: image1(),
							})
						},
					},
				},
			}),
			WantObject:      image1(),
			WantProtoObject: image1(),
			WantTable: new(`
			ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS
			ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
			`),
			WantWideTable: new(`
			ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS
			ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
			ubuntu-24.04 Ubuntu 24.04
			`),
			WantMarkdown: new(`
            | ID           | NAME         | DESCRIPTION      | FEATURES | EXPIRATION | STATUS    |
            |--------------|--------------|------------------|----------|------------|-----------|
            | ubuntu-24.04 | Ubuntu 24.04 | Ubuntu 24.04 LTS | machine  |            | supported |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
