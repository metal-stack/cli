package api_e2e

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/tests/e2e"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	tenant1 = func() *apiv2.Tenant {
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
)

func Test_TenantCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceGetResponse, *apiv2.Tenant]{
		{
			Name:    "describe",
			CmdArgs: []string{"tenant", "describe", tenant1().Login},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.TenantServiceGetRequest{
							Login: tenant1().Login,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceGetResponse{
								Tenant: tenant1(),
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
			WantObject:      tenant1(),
			WantProtoObject: tenant1(),
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
