package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
)

func Test_AdminComponentCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv2.ComponentServiceGetResponse, *apiv2.Component]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "component", "describe", testresources.Component1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.ComponentServiceGetRequest{
							Uuid: testresources.Component1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ComponentServiceGetResponse{
								Component: testresources.Component1(),
							})
						},
					},
				},
			}),
			WantObject:      testresources.Component1(),
			WantProtoObject: testresources.Component1(),
			WantTable: new(`
			ID                                    TYPE        IDENTIFIER    STARTED  AGE  VERSION  TOKEN                                 TOKEN EXPIRES IN
			c1a2b3d4-e5f6-7890-abcd-ef1234567890  metal-core  metal-core-1  0s       0s   v1.0.0   t1a2b3d4-e5f6-7890-abcd-ef1234567890  1d
			`),
			WantWideTable: new(`
			ID                                    TYPE        IDENTIFIER    STARTED  AGE  VERSION  TOKEN                                 TOKEN EXPIRES IN
			c1a2b3d4-e5f6-7890-abcd-ef1234567890  metal-core  metal-core-1  0s       0s   v1.0.0   t1a2b3d4-e5f6-7890-abcd-ef1234567890  1d
			`),
			Template: new("{{ .uuid }} {{ .identifier }}"),
			WantTemplate: new(`
			c1a2b3d4-e5f6-7890-abcd-ef1234567890 metal-core-1
			`),
			WantMarkdown: new(`
			| ID                                   | TYPE       | IDENTIFIER   | STARTED | AGE | VERSION | TOKEN                                | TOKEN EXPIRES IN |
			|--------------------------------------|------------|--------------|---------|-----|---------|--------------------------------------|------------------|
			| c1a2b3d4-e5f6-7890-abcd-ef1234567890 | metal-core | metal-core-1 | 0s      | 0s  | v1.0.0  | t1a2b3d4-e5f6-7890-abcd-ef1234567890 | 1d               |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminComponentCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.ComponentServiceListResponse, apiv2.Component]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "component", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.ComponentServiceListRequest{
							Query: &apiv2.ComponentQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ComponentServiceListResponse{
								Components: []*apiv2.Component{
									testresources.Component1(),
									testresources.Component2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID                                    TYPE        IDENTIFIER    STARTED  AGE  VERSION  TOKEN                                 TOKEN EXPIRES IN
			c1a2b3d4-e5f6-7890-abcd-ef1234567890  metal-core  metal-core-1  0s       0s   v1.0.0   t1a2b3d4-e5f6-7890-abcd-ef1234567890  1d
			d2b3c4e5-f6a7-8901-bcde-f12345678901  pixiecore   pixiecore-1   0s       0s   v2.0.0   t2b3c4e5-f6a7-8901-bcde-f12345678901  2d
			`),
			WantWideTable: new(`
			ID                                    TYPE        IDENTIFIER    STARTED  AGE  VERSION  TOKEN                                 TOKEN EXPIRES IN
			c1a2b3d4-e5f6-7890-abcd-ef1234567890  metal-core  metal-core-1  0s       0s   v1.0.0   t1a2b3d4-e5f6-7890-abcd-ef1234567890  1d
			d2b3c4e5-f6a7-8901-bcde-f12345678901  pixiecore   pixiecore-1   0s       0s   v2.0.0   t2b3c4e5-f6a7-8901-bcde-f12345678901  2d
			`),
			Template: new("{{ .uuid }} {{ .identifier }}"),
			WantTemplate: new(`
c1a2b3d4-e5f6-7890-abcd-ef1234567890 metal-core-1
d2b3c4e5-f6a7-8901-bcde-f12345678901 pixiecore-1
			`),
			WantMarkdown: new(`
			| ID                                   | TYPE       | IDENTIFIER   | STARTED | AGE | VERSION | TOKEN                                | TOKEN EXPIRES IN |
			|--------------------------------------|------------|--------------|---------|-----|---------|--------------------------------------|------------------|
			| c1a2b3d4-e5f6-7890-abcd-ef1234567890 | metal-core | metal-core-1 | 0s      | 0s  | v1.0.0  | t1a2b3d4-e5f6-7890-abcd-ef1234567890 | 1d               |
			| d2b3c4e5-f6a7-8901-bcde-f12345678901 | pixiecore  | pixiecore-1  | 0s      | 0s  | v2.0.0  | t2b3c4e5-f6a7-8901-bcde-f12345678901 | 2d               |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminComponentCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.ComponentServiceDeleteResponse, *apiv2.Component]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "component", "delete", testresources.Component1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.ComponentServiceDeleteRequest{
							Uuid: testresources.Component1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ComponentServiceDeleteResponse{
								Component: testresources.Component1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.Component1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
