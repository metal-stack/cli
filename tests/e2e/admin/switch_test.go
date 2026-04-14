package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_AdminSwitchCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServiceGetResponse, *apiv2.Switch]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "switch", "describe", testresources.Switch2().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SwitchServiceGetRequest{
							Id: testresources.Switch2().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceGetResponse{
								Switch: testresources.Switch2(),
							})
						},
					},
				},
			}),
			WantObject:      testresources.Switch2(),
			WantProtoObject: testresources.Switch2(),
			WantDefault: new(`
description: leaf switch 2
id: leaf02
lastSync:
  duration: 0.200s
  time: "2000-01-01T00:00:00Z"
managementIp: 10.0.0.2
managementUser: admin
meta:
  createdAt: "2000-01-01T00:00:00Z"
os:
  metalCoreVersion: v0.9.1 (abc1234), tags/v0.9.1
  vendor: SWITCH_OS_VENDOR_SONIC
  version: 4.2.0
partition: fra-equ01
rack: rack-1
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSwitchCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServiceListResponse, apiv2.Switch]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "switch", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SwitchServiceListRequest{
							Query: &apiv2.SwitchQuery{
								Os: &apiv2.SwitchOSQuery{},
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceListResponse{
								Switches: []*apiv2.Switch{
									testresources.Switch1(),
									testresources.Switch2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID      PARTITION  RACK    OS  STATUS  LAST SYNC
			leaf01  fra-equ01  rack-1  🦔  ●
			leaf02  fra-equ01  rack-1  🦔  ●
			`),
			WantWideTable: new(`
			ID      PARTITION  RACK    OS             METALCORE         IP        MODE         LAST SYNC  SYNC DURATION  LAST ERROR
			leaf01  fra-equ01  rack-1  SONiC (4.2.0)  v0.9.1 (abc1234)  10.0.0.1  operational             100ms
			leaf02  fra-equ01  rack-1  SONiC (4.2.0)  v0.9.1 (abc1234)  10.0.0.2  operational             200ms
			`),
			Template: new("{{ .id }} {{ .partition }}"),
			WantTemplate: new(`
leaf01 fra-equ01
leaf02 fra-equ01
			`),
			WantMarkdown: new(`
			| ID     | PARTITION | RACK   | OS | STATUS | LAST SYNC |
			|--------|-----------|--------|----|--------|-----------|
			| leaf01 | fra-equ01 | rack-1 | 🦔 | ●      |           |
			| leaf02 | fra-equ01 | rack-1 | 🦔 | ●      |           |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSwitchCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServiceDeleteResponse, *apiv2.Switch]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "switch", "delete", testresources.Switch2().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SwitchServiceDeleteRequest{
							Id: testresources.Switch2().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceDeleteResponse{
								Switch: testresources.Switch2(),
							})
						},
					},
				},
			}),
			WantObject: testresources.Switch2(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSwitchCmd_Update(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServiceUpdateResponse, *apiv2.Switch]{
		{
			Name:    "update from file",
			CmdArgs: append([]string{"admin", "switch", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Switch2()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.SwitchServiceUpdateRequest{
								Id:             testresources.Switch2().Id,
								Description:    new(testresources.Switch2().Description),
								ManagementIp:   new(testresources.Switch2().ManagementIp),
								ManagementUser: testresources.Switch2().ManagementUser,
								Os:             testresources.Switch2().Os,
								ReplaceMode:    new(testresources.Switch2().ReplaceMode),
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.SwitchServiceUpdateResponse{
									Switch: testresources.Switch2(),
								})
							},
						},
					},
				}),
			WantTable: new(`
            ID      PARTITION  RACK    OS  STATUS  LAST SYNC  
            leaf02  fra-equ01  rack-1  🦔  ●
					`),
			WantWideTable: new(`
            ID      PARTITION  RACK    OS             METALCORE         IP        MODE         LAST SYNC  SYNC DURATION  LAST ERROR  
            leaf02  fra-equ01  rack-1  SONiC (4.2.0)  v0.9.1 (abc1234)  10.0.0.2  operational             200ms
					`),
			Template:     new("{{ .id }} {{ .os.metal_core_version }}"),
			WantTemplate: new(`leaf02 v0.9.1 (abc1234), tags/v0.9.1`),
			WantMarkdown: new(`
            | ID     | PARTITION | RACK   | OS | STATUS | LAST SYNC |
            |--------|-----------|--------|----|--------|-----------|
            | leaf02 | fra-equ01 | rack-1 | 🦔 | ●      |           |
					`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSwitchCmd_ConnectedMachines(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServiceConnectedMachinesResponse, *apiv2.Switch]{
		{
			Name:    "connected machines",
			CmdArgs: []string{"admin", "switch", "connected-machines"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SwitchServiceConnectedMachinesRequest{
							Query:        &apiv2.SwitchQuery{},
							MachineQuery: &apiv2.MachineQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceConnectedMachinesResponse{
								SwitchesWithMachines: []*apiv2.SwitchWithMachines{testresources.SwitchWithMachines1()},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID      NIC NAME        IDENTIFIER           PARTITION  RACK    SIZE      PRODUCT SERIAL  CHASSIS SERIAL  
            leaf01                                       fra-equ01  rack-1                                            
            └─╴id1  Ethernet0 (up)  oid:0x1000000000001  fra-equ01  rack-1  m1-small  ps-1            cs-1
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSwitchCmd_Detail(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServiceListResponse, apiv2.Switch]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "switch", "detail", testresources.Switch1().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SwitchServiceListRequest{
							Query: &apiv2.SwitchQuery{
								Os: &apiv2.SwitchOSQuery{},
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceListResponse{
								Switches: []*apiv2.Switch{
									testresources.Switch1(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            PARTITION  RACK    SWITCH  PORT       MACHINE  VNI - FILTER  CIDR - FILTER     
            fra-equ01  rack-1  leaf01  Ethernet0  id1      10001         10.0.0.0/24       
                                                           10002         192.168.100.0/24
			`),
			Template:     new("{{ .id }} {{ .partition }} {{ .rack }}"),
			WantTemplate: new(`leaf01 fra-equ01 rack-1`),
			WantMarkdown: new(`
            | PARTITION | RACK   | SWITCH | PORT      | MACHINE | VNI - FILTER | CIDR - FILTER    |
            |-----------|--------|--------|-----------|---------|--------------|------------------|
            | fra-equ01 | rack-1 | leaf01 | Ethernet0 | id1     | 10001        | 10.0.0.0/24      |
            |           |        |        |           |         | 10002        | 192.168.100.0/24 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSwitchCmd_Migrate(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServiceMigrateResponse, *apiv2.Switch]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "switch", "migrate", testresources.Switch1().Id, testresources.Switch2().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SwitchServiceMigrateRequest{
							OldSwitch: testresources.Switch1().Id,
							NewSwitch: testresources.Switch2().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceMigrateResponse{
								Switch: testresources.Switch2(),
							})
						},
					},
				},
			}),
			WantDefault: new(`
switch:
  description: leaf switch 2
  id: leaf02
  lastSync:
    duration: 0.200s
    time: "2000-01-01T00:00:00Z"
  managementIp: 10.0.0.2
  managementUser: admin
  meta:
    createdAt: "2000-01-01T00:00:00Z"
  os:
    metalCoreVersion: v0.9.1 (abc1234), tags/v0.9.1
    vendor: SWITCH_OS_VENDOR_SONIC
    version: 4.2.0
  partition: fra-equ01
  rack: rack-1
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSwitchCmd_Port(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServicePortResponse, *apiv2.Switch]{
		{
			Name:    "port up",
			CmdArgs: []string{"admin", "switch", "port", "up", testresources.Switch1().Id, "--port", testresources.Nic1().Name},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SwitchServicePortRequest{
							Id:      testresources.Switch1().Id,
							NicName: testresources.Nic1().Name,
							Status:  apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_UP,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServicePortResponse{
								Switch: testresources.Switch1(),
							})
						},
					},
				},
			}),
			Template:     new("{{ .Actual.machine_id }} {{ .Actual.nic.name }} {{ .Desired.name }}"),
			WantTemplate: new(`id1 Ethernet0 Ethernet0`),
		},
		{
			Name:    "down",
			CmdArgs: []string{"admin", "switch", "port", "down", testresources.Switch1().Id, "--port", testresources.Nic1().Name},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SwitchServicePortRequest{
							Id:      testresources.Switch1().Id,
							NicName: testresources.Nic1().Name,
							Status:  apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_DOWN,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServicePortResponse{
								Switch: testresources.Switch1(),
							})
						},
					},
				},
			}),
			Template:     new("{{ .Actual.machine_id }} {{ .Actual.nic.name }} {{ .Desired.name }}"),
			WantTemplate: new(`id1 Ethernet0 Ethernet0`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
