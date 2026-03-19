package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/tests/e2e"
	"github.com/metal-stack/cli/pkg/tests/e2e/testresources"
)

func Test_TenantCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceGetResponse, *apiv2.Tenant]{
		{
			Name:    "describe",
			CmdArgs: []string{"tenant", "describe", testresources.Tenant1().Login},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.TenantServiceGetRequest{
							Login: testresources.Tenant1().Login,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceGetResponse{
								Tenant: testresources.Tenant1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID           NAME         EMAIL                REGISTERED  COUPONS  TERMS AND CONDITIONS
            metal-stack  Metal Stack  info@metal-stack.io  now         -
			`),
			WantWideTable: new(`
			ID           NAME         EMAIL                REGISTERED  COUPONS  TERMS AND CONDITIONS
            metal-stack  Metal Stack  info@metal-stack.io  now         -
			`),
			WantMarkdown: new(`
            | ID          | NAME        | EMAIL               | REGISTERED | COUPONS | TERMS AND CONDITIONS |
            |-------------|-------------|---------------------|------------|---------|----------------------|
            | metal-stack | Metal Stack | info@metal-stack.io | now        | -       |                      |
			`),
			WantObject:      testresources.Tenant1(),
			WantProtoObject: testresources.Tenant1(),
			Template:        new("{{ .login }} {{ .name }}"),
			WantTemplate: new(`
			metal-stack Metal Stack
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceListResponse, apiv2.Tenant]{
		{
			Name:    "list",
			CmdArgs: []string{"tenant", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.TenantServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceListResponse{
								Tenants: []*apiv2.Tenant{
									testresources.Tenant1(),
									testresources.Tenant2(),
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
