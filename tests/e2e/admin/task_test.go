package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
)

func Test_AdminTaskCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.TaskServiceListResponse, adminv2.TaskInfo]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "task", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.TaskServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TaskServiceListResponse{
								Tasks: []*adminv2.TaskInfo{
									testresources.Task1(),
									testresources.Task2(),
								},
							})
						},
					},
				},
			}),
			Template: new("{{ .id }} {{ .type }}"),
			WantTemplate: new(`
00dc6a98-bd80-787a-9725-ffb692d83261 image-provision
00dc6ab4-34c0-73ce-af77-3e07748d0e0d firewall-update
			`),
			WantTable: new(`
            ID                                    QUEUE    WHEN     TYPE             STATE
            00dc6a98-bd80-787a-9725-ffb692d83261  default  1h ago   image-provision  active
            00dc6ab4-34c0-73ce-af77-3e07748d0e0d  default  30m ago  firewall-update  pending
			`),
			WantWideTable: new(`
            ID                                    QUEUE    WHEN     TYPE             STATE    ISSUED AT                      PAYLOAD                    RESULT
            00dc6a98-bd80-787a-9725-ffb692d83261  default  1h ago   image-provision  active   2000-01-01 00:00:00 +0100 CET  {"machine_id":"machine1"}
            00dc6ab4-34c0-73ce-af77-3e07748d0e0d  default  30m ago  firewall-update  pending  2000-01-01 00:30:00 +0100 CET  {"firewall_id":"fw1"}
			`),
			WantMarkdown: new(`
            | ID                                   | QUEUE   | WHEN    | TYPE            | STATE   |
            |--------------------------------------|---------|---------|-----------------|---------|
            | 00dc6a98-bd80-787a-9725-ffb692d83261 | default | 1h ago  | image-provision | active  |
            | 00dc6ab4-34c0-73ce-af77-3e07748d0e0d | default | 30m ago | firewall-update | pending |
			`),
		},
		{
			Name:    "list with queue filter",
			CmdArgs: []string{"admin", "task", "list", "--queue", "high-priority"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.TaskServiceListRequest{
							Queue: func() *string { s := "high-priority"; return &s }(),
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TaskServiceListResponse{
								Tasks: []*adminv2.TaskInfo{
									testresources.Task3(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    QUEUE          WHEN     TYPE             STATE
            00dc6ab5-1f20-7426-a397-8644fb78324e  high-priority  29m ago  machine-reimage  completed
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminTaskCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv2.TaskServiceGetResponse, *adminv2.TaskInfo]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "task", "describe", testresources.Task1().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.TaskServiceGetRequest{
							TaskId: testresources.Task1().Id,
							Queue:  "default",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TaskServiceGetResponse{
								Task: testresources.Task1(),
							})
						},
					},
				},
			}),
			WantObject:      testresources.Task1(),
			WantProtoObject: testresources.Task1(),
			Template:        new("{{ .id }} {{ .type }}"),
			WantTemplate: new(`
00dc6a98-bd80-787a-9725-ffb692d83261 image-provision
			`),
		},
		{
			Name:    "describe with queue",
			CmdArgs: []string{"admin", "task", "describe", testresources.Task3().Id, "--queue", "high-priority"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.TaskServiceGetRequest{
							TaskId: testresources.Task3().Id,
							Queue:  "high-priority",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TaskServiceGetResponse{
								Task: testresources.Task3(),
							})
						},
					},
				},
			}),
			WantObject:      testresources.Task3(),
			WantProtoObject: testresources.Task3(),
			Template:        new("{{ .id }} {{ .queue }} {{ .type }}"),
			WantTemplate: new(`
00dc6ab5-1f20-7426-a397-8644fb78324e high-priority machine-reimage
			`),
			WantTable: new(`
            ID                                    QUEUE          WHEN     TYPE             STATE
            00dc6ab5-1f20-7426-a397-8644fb78324e  high-priority  29m ago  machine-reimage  completed
			`),
			WantWideTable: new(`
            ID                                    QUEUE          WHEN     TYPE             STATE      ISSUED AT                      PAYLOAD                    RESULT
            00dc6ab5-1f20-7426-a397-8644fb78324e  high-priority  29m ago  machine-reimage  completed  2000-01-01 00:31:00 +0100 CET  {"machine_id":"machine2"}  success
			`),
			WantMarkdown: new(`
            | ID                                   | QUEUE         | WHEN    | TYPE            | STATE     |
            |--------------------------------------|---------------|---------|-----------------|-----------|
            | 00dc6ab5-1f20-7426-a397-8644fb78324e | high-priority | 29m ago | machine-reimage | completed |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminTaskQueuesCmd(t *testing.T) {
	tests := []*e2e.Test[adminv2.TaskServiceQueuesResponse, any]{
		{
			Name:    "queues",
			CmdArgs: []string{"admin", "task", "queues"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.TaskServiceQueuesRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TaskServiceQueuesResponse{
								Queues: []string{"default", "high-priority", "low-priority"},
							})
						},
					},
				},
			}),
			WantTable: new(`
			QUEUE
			default
			high-priority
			low-priority
			`),
			WantWideTable: new(`
			QUEUE
			default
			high-priority
			low-priority
			`),
			WantMarkdown: new(`
			| QUEUE         |
            |---------------|
            | default       |
            | high-priority |
            | low-priority  |
			`),
			Template: new(`{{ range .queues }}{{ . }}{{ "\n" }}{{ end }}`),
			WantTemplate: new(`
default
high-priority
low-priority
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
