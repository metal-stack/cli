package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
)

func Test_NetworkCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.NetworkServiceListResponse, []*apiv2.Network]{
		{
			Name:    "list",
			CmdArgs: []string{"network", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.NetworkServiceListRequest{
							Project: "",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.NetworkServiceListResponse{
								Networks: []*apiv2.Network{
									testresources.Network1(),
									testresources.Network2(),
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.Network{
				testresources.Network1(),
				testresources.Network2(),
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_NetworkCmd_Describe(t *testing.T) {
	n1 := testresources.Network1()
	tests := []*e2e.Test[apiv2.NetworkServiceGetResponse, *apiv2.Network]{
		{
			Name:    "describe",
			CmdArgs: []string{"network", "describe", n1.Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.NetworkServiceGetRequest{
							Id:      n1.Id,
							Project: "",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.NetworkServiceGetResponse{
								Network: n1,
							})
						},
					},
				},
			}),
			WantObject: n1,
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
