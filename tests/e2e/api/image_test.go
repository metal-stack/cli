package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
)

func Test_ImageCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.ImageServiceListResponse, *apiv2.Image]{
		{
			Name:    "list",
			CmdArgs: []string{"image", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ImageServiceListRequest{
							Query: &apiv2.ImageQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ImageServiceListResponse{
								Images: []*apiv2.Image{
									testresources.Image1(),
									testresources.Image2(),
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
			CmdArgs: []string{"image", "describe", testresources.Image1().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ImageServiceGetRequest{
							Id: testresources.Image1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ImageServiceGetResponse{
								Image: testresources.Image1(),
							})
						},
					},
				},
			}),
			WantObject:      testresources.Image1(),
			WantProtoObject: testresources.Image1(),
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
