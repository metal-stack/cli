package api_e2e

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/tests/e2e"
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
