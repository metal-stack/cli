package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
)

func Test_MachineCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.MachineServiceListResponse, apiv2.Machine]{
		{
			Name:    "list",
			CmdArgs: []string{"machine", "list", "--project", testresources.Project2().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.MachineServiceListRequest{
							Project: testresources.Project2().Uuid,
							Query:   &apiv2.MachineQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.MachineServiceListResponse{
								Machines: []*apiv2.Machine{
									testresources.Machine2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                        LAST EVENT   WHEN  AGE  HOSTNAME   PROJECT                               SIZE           IMAGE         PARTITION    RACK
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  🛡  Phoned Home  1m    1m   machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  v1-medium-x86  Ubuntu 24.04  partition-2  rack-1
			`),
			WantWideTable: new(`
            ID                                     LAST EVENT   WHEN  AGE  DESCRIPTION  NAME       HOSTNAME   PROJECT                               IPS        SIZE           IMAGE         PARTITION    RACK    STARTED               TAGS  STATE
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  Phoned Home  1m    1m   machine 2    machine-2  machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  4.5.6.7    v1-medium-x86  Ubuntu 24.04  partition-2  rack-1  1999-12-31T23:59:00Z  c=d   available
                                                                                                                                                    192.1.1.1
			`),
			Template: new("{{ .allocation.uuid }} {{ .allocation.project }}"),
			WantTemplate: new(`
			4f94e87b-b08f-4f82-b053-9b8305de60ad f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c
			`),
			WantMarkdown: new(`
            | ID                                    |   | LAST EVENT  | WHEN | AGE | HOSTNAME  | PROJECT                              | SIZE          | IMAGE        | PARTITION   | RACK   |
            |---------------------------------------|---|-------------|------|-----|-----------|--------------------------------------|---------------|--------------|-------------|--------|
            | 673fc473-63ca-4ea4-b9dd-b45cb2127a6fd | 🛡 | Phoned Home | 1m   | 1m  | machine-2 | f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | v1-medium-x86 | Ubuntu 24.04 | partition-2 | rack-1 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_MachineCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.MachineServiceGetResponse, *apiv2.Machine]{
		{
			Name:    "describe",
			CmdArgs: []string{"machine", "describe", "--project", testresources.Project2().Uuid, testresources.Machine2().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.MachineServiceGetRequest{
							Uuid:    testresources.Machine2().Uuid,
							Project: testresources.Project2().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.MachineServiceGetResponse{
								Machine: testresources.Machine2(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Machine2(),
			WantTable: new(`
            ID                                        LAST EVENT   WHEN  AGE  HOSTNAME   PROJECT                               SIZE           IMAGE         PARTITION    RACK
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  🛡  Phoned Home  1m    1m   machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  v1-medium-x86  Ubuntu 24.04  partition-2  rack-1
			`),
			WantWideTable: new(`
            ID                                     LAST EVENT   WHEN  AGE  DESCRIPTION  NAME       HOSTNAME   PROJECT                               IPS        SIZE           IMAGE         PARTITION    RACK    STARTED               TAGS  STATE
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  Phoned Home  1m    1m   machine 2    machine-2  machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  4.5.6.7    v1-medium-x86  Ubuntu 24.04  partition-2  rack-1  1999-12-31T23:59:00Z  c=d   available
                                                                                                                                                    192.1.1.1
			`),
			Template: new("{{ .allocation.uuid }} {{ .allocation.project }}"),
			WantTemplate: new(`
			4f94e87b-b08f-4f82-b053-9b8305de60ad f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c
			`),
			WantMarkdown: new(`
            | ID                                    |   | LAST EVENT  | WHEN | AGE | HOSTNAME  | PROJECT                              | SIZE          | IMAGE        | PARTITION   | RACK   |
            |---------------------------------------|---|-------------|------|-----|-----------|--------------------------------------|---------------|--------------|-------------|--------|
            | 673fc473-63ca-4ea4-b9dd-b45cb2127a6fd | 🛡 | Phoned Home | 1m   | 1m  | machine-2 | f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | v1-medium-x86 | Ubuntu 24.04 | partition-2 | rack-1 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

// func Test_MachineCmd_Create(t *testing.T) {
// 	tests := []*e2e.Test[apiv2.MachineServiceGetResponse, *apiv2.Machine]{
// 		{
// 			Name:    "create",
// 			CmdArgs: []string{"machine", "create", "--project", testresources.Machine1().Project, "--network", testresources.Machine1().Network, "--static=true"},
// 			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
// 				ClientCalls: []client.ClientCall{
// 					{
// 						WantRequest: &apiv2.MachineServiceCreateRequest{
// 							Project: testresources.Machine1().Project,
// 							Network: testresources.Machine1().Network,
// 							Type:    &testresources.Machine1().Type,
// 						},
// 						WantResponse: func() connect.AnyResponse {
// 							return connect.NewResponse(&apiv2.MachineServiceCreateResponse{
// 								Ip: testresources.Machine1(),
// 							})
// 						},
// 					},
// 				},
// 			}),
// 			WantProtoObject: testresources.Machine1(),
// 		},
// 		{
// 			Name:    "create from file",
// 			CmdArgs: append([]string{"machine", "create"}, e2e.AppendFromFileCommonArgs()...),
// 			NewRootCmd: e2erootcmd.NewRootCmd(t,
// 				&e2erootcmd.TestConfig{
// 					FsMocks: func(fs *afero.Afero) {
// 						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine1()), 0755))
// 					},
// 					ClientCalls: []client.ClientCall{
// 						{
// 							WantRequest: &apiv2.MachineServiceCreateRequest{
// 								Ip:            &testresources.Machine1().Ip,
// 								Project:       testresources.Machine1().Project,
// 								Network:       testresources.Machine1().Network,
// 								Name:          &testresources.Machine1().Name,
// 								Descrmachinetion:   &testresources.Machine1().Descrmachinetion,
// 								Labels:        testresources.Machine1().Meta.Labels,
// 								Type:          &testresources.Machine1().Type,
// 								AddressFamily: nil,
// 							},
// 							WantResponse: func() connect.AnyResponse {
// 								return connect.NewResponse(&apiv2.MachineServiceCreateResponse{
// 									Ip: testresources.Machine1(),
// 								})
// 							},
// 						},
// 					},
// 				}),
// 			WantTable: new(`
//             Machine       PROJECT                               ID                                    NETWORK   TYPE    NAME  ATTACHED SERVICE
//             1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  internet  static  a
// 			`),
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt.TestCmd(t)
// 	}
// }

// func Test_MachineCmd_Delete(t *testing.T) {
// 	tests := []*e2e.Test[apiv2.MachineServiceDeleteResponse, *apiv2.Machine]{
// 		{
// 			Name:    "delete",
// 			CmdArgs: []string{"machine", "delete", "--project", testresources.Machine1().Project, testresources.Machine1().Ip},
// 			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
// 				ClientCalls: []client.ClientCall{
// 					{
// 						WantRequest: &apiv2.MachineServiceDeleteRequest{
// 							Ip:      testresources.Machine1().Ip,
// 							Project: testresources.Machine1().Project,
// 						},
// 						WantResponse: func() connect.AnyResponse {
// 							return connect.NewResponse(&apiv2.MachineServiceDeleteResponse{
// 								Ip: testresources.Machine1(),
// 							})
// 						},
// 					},
// 				},
// 			}),
// 			WantProtoObject: testresources.Machine1(),
// 		},
// 		{
// 			Name:    "delete from file",
// 			CmdArgs: append([]string{"machine", "delete"}, e2e.AppendFromFileCommonArgs()...),
// 			NewRootCmd: e2erootcmd.NewRootCmd(t,
// 				&e2erootcmd.TestConfig{
// 					FsMocks: func(fs *afero.Afero) {
// 						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine1()), 0755))
// 					},
// 					ClientCalls: []client.ClientCall{
// 						{
// 							WantRequest: &apiv2.MachineServiceDeleteRequest{
// 								Ip:      testresources.Machine1().Ip,
// 								Project: testresources.Machine1().Project,
// 							},
// 							WantResponse: func() connect.AnyResponse {
// 								return connect.NewResponse(&apiv2.MachineServiceDeleteResponse{
// 									Ip: testresources.Machine1(),
// 								})
// 							},
// 						},
// 					},
// 				},
// 			),
// 			WantTable: new(`
//             Machine       PROJECT                               ID                                    NETWORK   TYPE    NAME  ATTACHED SERVICE
//             1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  internet  static  a
// 			`),
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt.TestCmd(t)
// 	}
// }

// func Test_MachineCmd_Update(t *testing.T) {
// 	tests := []*e2e.Test[apiv2.MachineServiceUpdateResponse, *apiv2.Machine]{
// 		{
// 			Name:    "update",
// 			CmdArgs: []string{"machine", "update", "--project", testresources.Machine1().Project, testresources.Machine1().Ip, "--name", "foo"},
// 			NewRootCmd: e2erootcmd.NewRootCmd(t,
// 				&e2erootcmd.TestConfig{
// 					ClientCalls: []client.ClientCall{
// 						{
// 							WantRequest: &apiv2.MachineServiceUpdateRequest{
// 								Ip:      testresources.Machine1().Ip,
// 								Project: testresources.Machine1().Project,
// 								Name:    new("foo"),
// 								UpdateMeta: &apiv2.UpdateMeta{
// 									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
// 								},
// 							},
// 							WantResponse: func() connect.AnyResponse {
// 								return connect.NewResponse(&apiv2.MachineServiceUpdateResponse{
// 									Ip: testresources.Machine1(),
// 								})
// 							},
// 						},
// 					},
// 				},
// 			),
// 			WantProtoObject: testresources.Machine1(),
// 			WantTable: new(`
//             Machine       PROJECT                               ID                                    NETWORK   TYPE    NAME  ATTACHED SERVICE
//             1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  internet  static  a
// 			`),
// 			WantWideTable: new(`
// 			Machine       PROJECT                               ID                                    TYPE    NAME  DESCRMachineTION    LABELS
// 			1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static  a     a descrmachinetion  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
// 			`),
// 		},
// 		{
// 			Name:    "update from file",
// 			CmdArgs: append([]string{"machine", "update"}, e2e.AppendFromFileCommonArgs()...),
// 			NewRootCmd: e2erootcmd.NewRootCmd(t,
// 				&e2erootcmd.TestConfig{
// 					FsMocks: func(fs *afero.Afero) {
// 						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine1()), 0755))
// 					},
// 					ClientCalls: []client.ClientCall{
// 						{
// 							WantRequest: &apiv2.MachineServiceUpdateRequest{
// 								Ip:          testresources.Machine1().Ip,
// 								Project:     testresources.Machine1().Project,
// 								Descrmachinetion: &testresources.Machine1().Descrmachinetion,
// 								Labels: &apiv2.UpdateLabels{
// 									Strategy: &apiv2.UpdateLabels_Replace{
// 										Replace: &apiv2.Labels{
// 											Labels: testresources.Machine1().Meta.Labels.Labels,
// 										},
// 									},
// 								},
// 								Name: &testresources.Machine1().Name,
// 								Type: &testresources.Machine1().Type,
// 								UpdateMeta: &apiv2.UpdateMeta{
// 									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
// 								},
// 							},
// 							WantResponse: func() connect.AnyResponse {
// 								return connect.NewResponse(&apiv2.MachineServiceUpdateResponse{
// 									Ip: testresources.Machine1(),
// 								})
// 							},
// 						},
// 					},
// 				},
// 			),
// 			WantTable: new(`
//             Machine       PROJECT                               ID                                    NETWORK   TYPE    NAME  ATTACHED SERVICE
//             1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  internet  static  a
// 			`),
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt.TestCmd(t)
// 	}
// }

// func Test_MachineCmd_Apply(t *testing.T) {
// 	tests := []*e2e.Test[apiv2.MachineServiceUpdateResponse, *apiv2.Machine]{
// 		{
// 			Name:    "apply",
// 			CmdArgs: append([]string{"machine", "apply"}, e2e.AppendFromFileCommonArgs()...),
// 			NewRootCmd: e2erootcmd.NewRootCmd(t,
// 				&e2erootcmd.TestConfig{
// 					FsMocks: func(fs *afero.Afero) {
// 						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine1()), 0755))
// 					},
// 					ClientCalls: []client.ClientCall{
// 						{
// 							WantRequest: &apiv2.MachineServiceCreateRequest{
// 								Ip:            &testresources.Machine1().Ip,
// 								Project:       testresources.Machine1().Project,
// 								Network:       testresources.Machine1().Network,
// 								Name:          &testresources.Machine1().Name,
// 								Descrmachinetion:   &testresources.Machine1().Descrmachinetion,
// 								Labels:        testresources.Machine1().Meta.Labels,
// 								Type:          &testresources.Machine1().Type,
// 								AddressFamily: nil,
// 							},
// 							WantResponse: func() connect.AnyResponse {
// 								return connect.NewResponse(&apiv2.MachineServiceCreateResponse{
// 									Ip: testresources.Machine1(),
// 								})
// 							},
// 						},
// 					},
// 				},
// 			),
// 			WantTable: new(`
//             Machine       PROJECT                               ID                                    NETWORK   TYPE    NAME  ATTACHED SERVICE
//             1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  internet  static  a
// 			`),
// 		},
// 		{
// 			Name:    "apply already exists",
// 			CmdArgs: append([]string{"machine", "apply"}, e2e.AppendFromFileCommonArgs()...),
// 			NewRootCmd: e2erootcmd.NewRootCmd(t,
// 				&e2erootcmd.TestConfig{
// 					FsMocks: func(fs *afero.Afero) {
// 						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine1()), 0755))
// 					},
// 					ClientCalls: []client.ClientCall{
// 						{
// 							WantRequest: &apiv2.MachineServiceCreateRequest{
// 								Ip:            &testresources.Machine1().Ip,
// 								Project:       testresources.Machine1().Project,
// 								Network:       testresources.Machine1().Network,
// 								Name:          &testresources.Machine1().Name,
// 								Descrmachinetion:   &testresources.Machine1().Descrmachinetion,
// 								Labels:        testresources.Machine1().Meta.Labels,
// 								Type:          &testresources.Machine1().Type,
// 								AddressFamily: nil,
// 							},
// 							WantError: connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")),
// 						},
// 						{
// 							WantRequest: &apiv2.MachineServiceUpdateRequest{
// 								Ip:          testresources.Machine1().Ip,
// 								Project:     testresources.Machine1().Project,
// 								Descrmachinetion: &testresources.Machine1().Descrmachinetion,
// 								Labels: &apiv2.UpdateLabels{
// 									Strategy: &apiv2.UpdateLabels_Replace{
// 										Replace: &apiv2.Labels{
// 											Labels: testresources.Machine1().Meta.Labels.Labels,
// 										},
// 									}},
// 								Name: &testresources.Machine1().Name,
// 								Type: &testresources.Machine1().Type,
// 								UpdateMeta: &apiv2.UpdateMeta{
// 									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
// 								},
// 							},
// 							WantResponse: func() connect.AnyResponse {
// 								return connect.NewResponse(&apiv2.MachineServiceUpdateResponse{
// 									Ip: testresources.Machine1(),
// 								})
// 							},
// 						},
// 					},
// 				},
// 			),
// 			WantTable: new(`
//             Machine       PROJECT                               ID                                    NETWORK   TYPE    NAME  ATTACHED SERVICE
//             1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  internet  static  a
// 			`),
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt.TestCmd(t)
// 	}
// }
