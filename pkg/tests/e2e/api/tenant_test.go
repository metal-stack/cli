package api_e2e

import (
	"testing"
	"time"

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
				CreatedAt: timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
)

func Test_TenantCmd_Describe(t *testing.T) {
	tn := tenant1()

	tests := []*e2e.Test[apiv2.TenantServiceGetRequest, apiv2.TenantServiceGetResponse]{
		{
			Name: "describe",
			Cmd: func() []string {
				return []string{"tenant", "describe", tn.Login}
			},
			WantRequest: apiv2.TenantServiceGetRequest{
				Login: tn.Login,
			},
			WantResponse: apiv2.TenantServiceGetResponse{
				Tenant: tn,
			},
			WantObject: tn,
			Template:   new("{{ .login }} {{ .name }}"),
			WantTemplate: new(`
metal-stack Metal Stack
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
