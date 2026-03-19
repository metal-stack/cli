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
