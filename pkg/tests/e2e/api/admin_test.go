package api_e2e

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
	adminTenant1 = func() *apiv2.Tenant {
		return &apiv2.Tenant{
			Login:       "metal-stack",
			Name:        "Metal Stack",
			Email:       "info@metal-stack.io",
			Description: "a tenant",
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
	adminTenant2 = func() *apiv2.Tenant {
		return &apiv2.Tenant{
			Login:       "acme-corp",
			Name:        "ACME Corp",
			Email:       "admin@acme.io",
			Description: "another tenant",
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
)

func Test_AdminTenantCmd_Create(t *testing.T) {
	tests := []*e2e.Test[adminv2.TenantServiceCreateResponse, *apiv2.Tenant]{
		{
			Name:    "create",
			CmdArgs: []string{"admin", "tenant", "create", "--name", adminTenant1().Name, "--description", adminTenant1().Description, "--email", adminTenant1().Email},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.TenantServiceCreateRequest{
							Name:        adminTenant1().Name,
							Description: new(adminTenant1().Description),
							Email:       new(adminTenant1().Email),
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TenantServiceCreateResponse{
								Tenant: adminTenant1(),
							})
						},
					},
				},
			}),
			WantObject:      adminTenant1(),
			WantProtoObject: adminTenant1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminTenantCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.TenantServiceListResponse, apiv2.Tenant]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "tenant", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.TenantServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TenantServiceListResponse{
								Tenants: []*apiv2.Tenant{
									adminTenant1(),
									adminTenant2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID           NAME         EMAIL                REGISTERED  COUPONS  TERMS AND CONDITIONS
			metal-stack  Metal Stack  info@metal-stack.io  now         -
			acme-corp    ACME Corp    admin@acme.io        now         -
			`),
			WantWideTable: new(`
			ID           NAME         EMAIL                REGISTERED  COUPONS  TERMS AND CONDITIONS
			metal-stack  Metal Stack  info@metal-stack.io  now         -
			acme-corp    ACME Corp    admin@acme.io        now         -
			`),
			Template: new("{{ .login }} {{ .name }}"),
			WantTemplate: new(`
metal-stack Metal Stack
acme-corp ACME Corp
			`),
			WantMarkdown: new(`
			| ID          | NAME        | EMAIL               | REGISTERED | COUPONS | TERMS AND CONDITIONS |
			|-------------|-------------|---------------------|------------|---------|----------------------|
			| metal-stack | Metal Stack | info@metal-stack.io | now        | -       |                      |
			| acme-corp   | ACME Corp   | admin@acme.io       | now        | -       |                      |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminTokenCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.TokenServiceListResponse, apiv2.Token]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "token", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.TokenServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TokenServiceListResponse{
								Tokens: []*apiv2.Token{
									token1(),
									token2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
			TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES
			TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
			TOKEN_TYPE_API  b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01         dev@metal-stack.io    dev token    0      0      2000-01-03 00:00:00 UTC (in 2d)
			`),
			WantWideTable: new(`
			TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES
			TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
			TOKEN_TYPE_API  b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01         dev@metal-stack.io    dev token    0      0      2000-01-03 00:00:00 UTC (in 2d)
			`),
			Template: new("{{ .uuid }} {{ .description }}"),
			WantTemplate: new(`
a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 ci token
b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01 dev token
			`),
			WantMarkdown: new(`
			| TYPE           | ID                                   | ADMIN | USER                 | DESCRIPTION | ROLES | PERMS | EXPIRES                         |
			|----------------|--------------------------------------|-------|----------------------|-------------|-------|-------|---------------------------------|
			| TOKEN_TYPE_API | a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 |       | admin@metal-stack.io | ci token    | 0     | 0     | 2000-01-02 00:00:00 UTC (in 1d) |
			| TOKEN_TYPE_API | b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01 |       | dev@metal-stack.io   | dev token   | 0     | 0     | 2000-01-03 00:00:00 UTC (in 2d) |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminTokenCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.TokenServiceRevokeResponse, *apiv2.Token]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "token", "delete", token1().Uuid, "--user", "user-123"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.TokenServiceRevokeRequest{
							Uuid: token1().Uuid,
							User: "user-123",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TokenServiceRevokeResponse{})
						},
					},
				},
			}),
			WantObject: &apiv2.Token{
				Uuid: token1().Uuid,
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminProjectCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.ProjectServiceListResponse, apiv2.Project]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "project", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.ProjectServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ProjectServiceListResponse{
								Projects: []*apiv2.Project{
									project1(),
									project2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION     CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project   2025-06-01 10:00:00 UTC
			f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  metal-stack  project-b  second project  2025-07-15 14:30:00 UTC
			`),
			WantWideTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION     CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project   2025-06-01 10:00:00 UTC
			f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  metal-stack  project-b  second project  2025-07-15 14:30:00 UTC
			`),
			Template: new("{{ .uuid }} {{ .name }}"),
			WantTemplate: new(`
0d81bca7-73f6-4da3-8397-4a8c52a0c583 project-a
f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c project-b
			`),
			WantMarkdown: new(`
			| ID                                   | TENANT      | NAME      | DESCRIPTION    | CREATION DATE           |
			|--------------------------------------|-------------|-----------|----------------|-------------------------|
			| 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | metal-stack | project-a | first project  | 2025-06-01 10:00:00 UTC |
			| f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | metal-stack | project-b | second project | 2025-07-15 14:30:00 UTC |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

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

var (
	switch1 = func() *apiv2.Switch {
		return &apiv2.Switch{
			Id:          "leaf01",
			Partition:   "fra-equ01",
			Rack:        new("rack-1"),
			Description: "leaf switch 1",
			ManagementIp:   "10.0.0.1",
			ManagementUser: new("admin"),
			Os: &apiv2.SwitchOS{
				Vendor:           apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_SONIC,
				Version:          "4.2.0",
				MetalCoreVersion: "v0.9.1 (abc1234), tags/v0.9.1",
			},
			LastSync: &apiv2.SwitchSync{
				Time:     timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
				Duration: durationpb.New(100 * time.Millisecond),
			},
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
	switch2 = func() *apiv2.Switch {
		return &apiv2.Switch{
			Id:          "leaf02",
			Partition:   "fra-equ01",
			Rack:        new("rack-1"),
			Description: "leaf switch 2",
			ManagementIp:   "10.0.0.2",
			ManagementUser: new("admin"),
			Os: &apiv2.SwitchOS{
				Vendor:           apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_SONIC,
				Version:          "4.2.0",
				MetalCoreVersion: "v0.9.1 (abc1234), tags/v0.9.1",
			},
			LastSync: &apiv2.SwitchSync{
				Time:     timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
				Duration: durationpb.New(200 * time.Millisecond),
			},
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
)

func Test_AdminSwitchCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServiceGetResponse, *apiv2.Switch]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "switch", "describe", switch1().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.SwitchServiceGetRequest{
							Id: switch1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceGetResponse{
								Switch: switch1(),
							})
						},
					},
				},
			}),
			WantObject:      switch1(),
			WantProtoObject: switch1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSwitchCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServiceListResponse, apiv2.Switch]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "switch", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.SwitchServiceListRequest{
							Query: &apiv2.SwitchQuery{
								Os: &apiv2.SwitchOSQuery{},
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceListResponse{
								Switches: []*apiv2.Switch{
									switch1(),
									switch2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID      PARTITION  RACK    OS  STATUS  LAST SYNC
			leaf01  fra-equ01  rack-1  🦔  ●
			leaf02  fra-equ01  rack-1  🦔  ●
			`),
			WantWideTable: new(`
			ID      PARTITION  RACK    OS             METALCORE         IP        MODE         LAST SYNC  SYNC DURATION  LAST ERROR
			leaf01  fra-equ01  rack-1  SONiC (4.2.0)  v0.9.1 (abc1234)  10.0.0.1  operational             100ms
			leaf02  fra-equ01  rack-1  SONiC (4.2.0)  v0.9.1 (abc1234)  10.0.0.2  operational             200ms
			`),
			Template: new("{{ .id }} {{ .partition }}"),
			WantTemplate: new(`
leaf01 fra-equ01
leaf02 fra-equ01
			`),
			WantMarkdown: new(`
			| ID     | PARTITION | RACK   | OS | STATUS | LAST SYNC |
			|--------|-----------|--------|----|--------|-----------|
			| leaf01 | fra-equ01 | rack-1 | 🦔 | ●      |           |
			| leaf02 | fra-equ01 | rack-1 | 🦔 | ●      |           |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSwitchCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServiceDeleteResponse, *apiv2.Switch]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "switch", "delete", switch1().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.SwitchServiceDeleteRequest{
							Id: switch1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceDeleteResponse{
								Switch: switch1(),
							})
						},
					},
				},
			}),
			WantObject: switch1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
