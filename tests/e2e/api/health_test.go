package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
)

var (
	health1 = func() *apiv2.Health {
		return &apiv2.Health{
			Services: []*apiv2.HealthStatus{
				{
					Name:    apiv2.Service_SERVICE_IPAM,
					Status:  apiv2.ServiceStatus_SERVICE_STATUS_HEALTHY,
					Message: "i am healthy",
				},
			},
		}
	}
)

func Test_HealthCmd(t *testing.T) {
	tests := []*e2e.Test[apiv2.HealthServiceGetResponse, *apiv2.Health]{
		{
			Name:    "health",
			CmdArgs: []string{"health"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: &apiv2.HealthServiceGetRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.HealthServiceGetResponse{
								Health: health1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            NAME  MESSAGE
            ✔  ipam  i am healthy
			`),
			WantWideTable: new(`
			NAME  MESSAGE
            ✔  ipam  i am healthy
			`),
			WantMarkdown: new(`
            |   | NAME | MESSAGE      |
            |---|------|--------------|
            | ✔ | ipam | i am healthy |
			`),
			WantObject:      health1(),
			WantProtoObject: health1(),
			Template:        new("{{ range $s := .services }}{{ $s.message }} {{ end }}"),
			WantTemplate: new(`
            i am healthy
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
