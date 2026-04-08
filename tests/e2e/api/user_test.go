package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
)

func Test_UserCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.UserServiceGetResponse, *apiv2.User]{
		{
			Name:    "describe",
			CmdArgs: []string{"user", "describe"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.UserServiceGetRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.UserServiceGetResponse{
								User: testresources.User(),
							})
						},
					},
				},
			}),
			WantObject:      testresources.User(),
			WantProtoObject: testresources.User(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
