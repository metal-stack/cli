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

func Test_AdminMachineCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.MachineServiceListResponse, []*apiv2.Machine]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "machine", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceListRequest{
							Query: &apiv2.MachineQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceListResponse{
								Machines: []*apiv2.Machine{
									testresources.Machine1(),
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.Machine{testresources.Machine1()},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminMachineCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.MachineServiceDeleteResponse, *apiv2.Machine]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "machine", "delete", testresources.Machine1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceDeleteRequest{
							Uuid: testresources.Machine1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceDeleteResponse{
								Machine: testresources.Machine1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.Machine1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminMachineCmd_SetState(t *testing.T) {
	tests := []*e2e.Test[adminv2.MachineServiceSetStateResponse, any]{
		{
			Name:    "set state",
			CmdArgs: []string{"admin", "machine", "set-state", testresources.Machine1().Uuid, "--state", "MACHINE_STATE_LOCKED", "--description", "maintenance"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceSetStateRequest{
							Uuid:        testresources.Machine1().Uuid,
							State:       apiv2.MachineState_MACHINE_STATE_LOCKED,
							Description: "maintenance",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceSetStateResponse{})
						},
					},
				},
			}),
			WantDefault: new(""),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminMachineCmd_ConsolePassword(t *testing.T) {
	tests := []*e2e.Test[adminv2.MachineServiceConsolePasswordResponse, *adminv2.MachineServiceConsolePasswordResponse]{
		{
			Name:    "console password",
			CmdArgs: []string{"admin", "machine", "console-password", testresources.Machine1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceConsolePasswordRequest{
							Uuid: testresources.Machine1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceConsolePasswordResponse{
								Uuid:     testresources.Machine1().Uuid,
								Password: "secret123",
							})
						},
					},
				},
			}),
			WantObject: &adminv2.MachineServiceConsolePasswordResponse{
				Uuid:     testresources.Machine1().Uuid,
				Password: "secret123",
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
