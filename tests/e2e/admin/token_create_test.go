package admin_e2e

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/durationpb"
)

func Test_AdminTokenCreateCmd(t *testing.T) {
	user := "user-1"
	tests := []*e2e.Test[adminv2.TokenServiceCreateResponse, string]{
		{
			Name:    "create",
			CmdArgs: []string{"admin", "token-create", "--user", "user-1", "--description", "admin token", "--expires", "24h"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.TokenServiceCreateRequest{
							User: &user,
							TokenCreateRequest: &apiv2.TokenServiceCreateRequest{
								Description: "admin token",
								Expires:     durationpb.New(24 * time.Hour),
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TokenServiceCreateResponse{
								Token: &apiv2.Token{
									Uuid:        "token-1",
									Description: "admin token",
								},
								Secret: "secret-value",
							})
						},
					},
				},
			}),
			WantDefault: new(`Make sure to copy your personal access token now as you will not be able to see this again.

secret-value

description: admin token
uuid: token-1
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
