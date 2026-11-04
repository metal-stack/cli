package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
)

func Test_IPCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.IPServiceListResponse, apiv2.IP]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "ip", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.IPServiceListRequest{
							Query: &apiv2.IPQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.IPServiceListResponse{
								Ips: []*apiv2.IP{
									testresources.IP1(),
									testresources.IP2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            IP       PROJECT                               ID                                    NETWORK   TYPE       NAME  ATTACHED SERVICE
            4.3.2.1  46bdfc45-9c8d-4268-b359-b40e3079d384  9cef40ec-29c6-4dfa-aee8-47ee1f49223d  internet  ephemeral  b
            1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  internet  static     a
			`),
			WantWideTable: new(`
			IP       PROJECT                               ID                                    TYPE       NAME  DESCRIPTION    LABELS
			4.3.2.1  46bdfc45-9c8d-4268-b359-b40e3079d384  9cef40ec-29c6-4dfa-aee8-47ee1f49223d  ephemeral  b     b description  a=b
			1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static     a     a description  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
			`),
			Template: new("{{ .ip }} {{ .project }}"),
			WantTemplate: new(`
4.3.2.1 46bdfc45-9c8d-4268-b359-b40e3079d384
1.1.1.1 ce19a655-7933-4745-8f3e-9592b4a90488
			`),
			WantMarkdown: new(`
            | IP      | PROJECT                              | ID                                   | NETWORK  | TYPE      | NAME | ATTACHED SERVICE |
            |---------|--------------------------------------|--------------------------------------|----------|-----------|------|------------------|
            | 4.3.2.1 | 46bdfc45-9c8d-4268-b359-b40e3079d384 | 9cef40ec-29c6-4dfa-aee8-47ee1f49223d | internet | ephemeral | b    |                  |
            | 1.1.1.1 | ce19a655-7933-4745-8f3e-9592b4a90488 | 2e0144a2-09ef-42b7-b629-4263295db6e8 | internet | static    | a    |                  |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.IPServiceGetResponse, *apiv2.IP]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "ip", "describe", testresources.IP1().Ip},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.IPServiceListRequest{
							Query: &apiv2.IPQuery{
								Ip: &testresources.IP1().Ip,
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.IPServiceListResponse{
								Ips: []*apiv2.IP{
									testresources.IP1(),
								},
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.IP1(),
			WantTable: new(`
            IP       PROJECT                               ID                                    NETWORK   TYPE    NAME  ATTACHED SERVICE
            1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  internet  static  a
			`),
			WantWideTable: new(`
			IP       PROJECT                               ID                                    TYPE    NAME  DESCRIPTION    LABELS
			1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static  a     a description  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
			`),
			Template: new("{{ .ip }} {{ .project }}"),
			WantTemplate: new(`
			1.1.1.1 ce19a655-7933-4745-8f3e-9592b4a90488
			`),
			WantMarkdown: new(`
            | IP      | PROJECT                              | ID                                   | NETWORK  | TYPE   | NAME | ATTACHED SERVICE |
            |---------|--------------------------------------|--------------------------------------|----------|--------|------|------------------|
            | 1.1.1.1 | ce19a655-7933-4745-8f3e-9592b4a90488 | 2e0144a2-09ef-42b7-b629-4263295db6e8 | internet | static | a    |                  |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
