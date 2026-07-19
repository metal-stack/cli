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

func Test_MachineCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.MachineServiceListResponse, []*apiv2.Machine]{
		{
			Name:    "list",
			CmdArgs: []string{"machine", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.MachineServiceListRequest{
							Query: &apiv2.MachineQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.MachineServiceListResponse{
								Machines: []*apiv2.Machine{
									testresources.Machine1(),
									testresources.Machine2(),
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.Machine{
				testresources.Machine1(),
				testresources.Machine2(),
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_MachineCmd_Describe(t *testing.T) {
	m1 := testresources.Machine1()
	tests := []*e2e.Test[apiv2.MachineServiceGetResponse, *apiv2.Machine]{
		{
			Name:    "describe",
			CmdArgs: []string{"machine", "describe", m1.Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.MachineServiceGetRequest{
							Uuid:    m1.Uuid,
							Project: "",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.MachineServiceGetResponse{
								Machine: m1,
							})
						},
					},
				},
			}),
			WantObject: m1,
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
