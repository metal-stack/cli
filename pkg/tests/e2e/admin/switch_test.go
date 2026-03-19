package admin_e2e

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/tests/e2e"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	switch1 = func() *apiv2.Switch {
		return &apiv2.Switch{
			Id:             "leaf01",
			Partition:      "fra-equ01",
			Rack:           new("rack-1"),
			Description:    "leaf switch 1",
			ManagementIp:   "10.0.0.1",
			ManagementUser: new("admin"),
			Os: &apiv2.SwitchOS{
				Vendor:           apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_SONIC,
				Version:          "4.2.0",
				MetalCoreVersion: "v0.9.1 (abc1234), tags/v0.9.1",
			},
			LastSync: &apiv2.SwitchSync{
				Time:     timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
				Duration: durationpb.New(100 * time.Millisecond),
			},
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
	switch2 = func() *apiv2.Switch {
		return &apiv2.Switch{
			Id:             "leaf02",
			Partition:      "fra-equ01",
			Rack:           new("rack-1"),
			Description:    "leaf switch 2",
			ManagementIp:   "10.0.0.2",
			ManagementUser: new("admin"),
			Os: &apiv2.SwitchOS{
				Vendor:           apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_SONIC,
				Version:          "4.2.0",
				MetalCoreVersion: "v0.9.1 (abc1234), tags/v0.9.1",
			},
			LastSync: &apiv2.SwitchSync{
				Time:     timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
				Duration: durationpb.New(200 * time.Millisecond),
			},
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
)

func Test_AdminSwitchCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv2.SwitchServiceGetResponse, *apiv2.Switch]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "switch", "describe", switch1().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.SwitchServiceGetRequest{
							Id: switch1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceGetResponse{
								Switch: switch1(),
							})
						},
					},
				},
			}),
			WantObject:      switch1(),
			WantProtoObject: switch1(),
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
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.SwitchServiceListRequest{
							Query: &apiv2.SwitchQuery{
								Os: &apiv2.SwitchOSQuery{},
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceListResponse{
								Switches: []*apiv2.Switch{
									switch1(),
									switch2(),
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
			CmdArgs: []string{"admin", "switch", "delete", switch1().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: adminv2.SwitchServiceDeleteRequest{
							Id: switch1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SwitchServiceDeleteResponse{
								Switch: switch1(),
							})
						},
					},
				},
			}),
			WantObject: switch1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
