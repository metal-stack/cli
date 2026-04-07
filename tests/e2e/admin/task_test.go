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
550e8400-e29b-41d4-a716-446655440000 image-provision
550e8400-e29b-41d4-a716-446655440001 firewall-update
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
			Template: new("{{ .id }} {{ .queue }} {{ .type }}"),
			WantTemplate: new(`
550e8400-e29b-41d4-a716-446655440002 high-priority machine-reimage
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
550e8400-e29b-41d4-a716-446655440000 image-provision
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
550e8400-e29b-41d4-a716-446655440002 high-priority machine-reimage
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
default
high-priority
low-priority
			`),
			WantWideTable: new(`
default
high-priority
low-priority
			`),
			Template: new("{{ . }}"),
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
