package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
)

func Test_AdminNetworkCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.NetworkServiceListResponse, []*apiv2.Network]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "network", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.NetworkServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.NetworkServiceListResponse{
								Networks: []*apiv2.Network{
									testresources.Network1(),
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.Network{testresources.Network1()},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminNetworkCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.NetworkServiceDeleteResponse, *apiv2.Network]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "network", "delete", testresources.Network1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.NetworkServiceDeleteRequest{
							Id: testresources.Network1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.NetworkServiceDeleteResponse{
								Network: testresources.Network1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.Network1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
