package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
)

func Test_VersionCmd(t *testing.T) {
	tests := []*e2e.Test[apiv2.VersionServiceGetResponse, *apiv2.Version]{
		{
			Name:    "describe",
			CmdArgs: []string{"version"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.VersionServiceGetRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.VersionServiceGetResponse{
								Version: testresources.Version(),
							})
						},
					},
				},
			}),
			Template:     new("{{ .Server.build_date }} {{ .Server.git_sha1 }} {{ .Server.revision }} {{ .Server.version }}"),
			WantTemplate: new(`2026-03-21T15:35:07+00:00 477edc0b tags/v0.1.8-0-g476edc0 v0.1.8`),
			WantDefault: new(`
---
Client: version not set, please build your app with appropriate ldflags, see https://github.com/metal-stack/v
  for reference, go1.26.2-X:nodwarf5
Server:
  build_date: "2026-03-21T15:35:07+00:00"
  git_sha1: 477edc0b
  revision: tags/v0.1.8-0-g476edc0
  version: v0.1.8
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
