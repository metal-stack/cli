package admin_e2e

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/tests/e2e"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	component1 = func() *apiv2.Component {
		return &apiv2.Component{
			Uuid:       "c1a2b3d4-e5f6-7890-abcd-ef1234567890",
			Type:       apiv2.ComponentType_COMPONENT_TYPE_METAL_CORE,
			Identifier: "metal-core-1",
			StartedAt:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			ReportedAt: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			Interval:   durationpb.New(10 * time.Second),
			Version: &apiv2.Version{
				Version: "v1.0.0",
			},
			Token: &apiv2.Token{
				Uuid:    "t1a2b3d4-e5f6-7890-abcd-ef1234567890",
				Expires: timestamppb.New(time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
	component2 = func() *apiv2.Component {
		return &apiv2.Component{
			Uuid:       "d2b3c4e5-f6a7-8901-bcde-f12345678901",
			Type:       apiv2.ComponentType_COMPONENT_TYPE_PIXIECORE,
			Identifier: "pixiecore-1",
			StartedAt:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			ReportedAt: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			Interval:   durationpb.New(10 * time.Second),
			Version: &apiv2.Version{
				Version: "v2.0.0",
			},
			Token: &apiv2.Token{
				Uuid:    "t2b3c4e5-f6a7-8901-bcde-f12345678901",
				Expires: timestamppb.New(time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
)

func Test_AdminComponentCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv2.ComponentServiceGetResponse, *apiv2.Component]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "component", "describe", component1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.ComponentServiceGetRequest{
							Uuid: component1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ComponentServiceGetResponse{
								Component: component1(),
							})
						},
					},
				},
			}),
			WantObject:      component1(),
			WantProtoObject: component1(),
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
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.ComponentServiceListRequest{
							Query: &apiv2.ComponentQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ComponentServiceListResponse{
								Components: []*apiv2.Component{
									component1(),
									component2(),
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
			CmdArgs: []string{"admin", "component", "delete", component1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.ComponentServiceDeleteRequest{
							Uuid: component1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ComponentServiceDeleteResponse{
								Component: component1(),
							})
						},
					},
				},
			}),
			WantObject: component1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
