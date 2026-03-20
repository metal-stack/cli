package api_e2e

import (
	"fmt"
	"testing"

	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_IPCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.IPServiceListResponse, apiv2.IP]{
		{
			Name:    "list",
			CmdArgs: []string{"ip", "list", "--project", testresources.IP1().Project},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: &apiv2.IPServiceListRequest{
							Project: testresources.IP1().Project,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.IPServiceListResponse{
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
			IP       PROJECT                               ID                                    TYPE       NAME  ATTACHED SERVICE
			4.3.2.1  46bdfc45-9c8d-4268-b359-b40e3079d384  9cef40ec-29c6-4dfa-aee8-47ee1f49223d  ephemeral  b
			1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static     a
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
			| IP      | PROJECT                              | ID                                   | TYPE      | NAME | ATTACHED SERVICE |
			|---------|--------------------------------------|--------------------------------------|-----------|------|------------------|
			| 4.3.2.1 | 46bdfc45-9c8d-4268-b359-b40e3079d384 | 9cef40ec-29c6-4dfa-aee8-47ee1f49223d | ephemeral | b    |                  |
			| 1.1.1.1 | ce19a655-7933-4745-8f3e-9592b4a90488 | 2e0144a2-09ef-42b7-b629-4263295db6e8 | static    | a    |                  |
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
			CmdArgs: []string{"ip", "describe", "--project", testresources.IP1().Project, testresources.IP1().Ip},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: &apiv2.IPServiceGetRequest{
							Ip:      testresources.IP1().Ip,
							Project: testresources.IP1().Project,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.IPServiceGetResponse{
								Ip: testresources.IP1(),
							})
						},
					},
				},
			}),
			WantObject:      testresources.IP1(),
			WantProtoObject: testresources.IP1(),
			WantTable: new(`
			IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE
			1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static  a
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
			| IP      | PROJECT                              | ID                                   | TYPE   | NAME | ATTACHED SERVICE |
			|---------|--------------------------------------|--------------------------------------|--------|------|------------------|
			| 1.1.1.1 | ce19a655-7933-4745-8f3e-9592b4a90488 | 2e0144a2-09ef-42b7-b629-4263295db6e8 | static | a    |                  |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_Create(t *testing.T) {
	tests := []*e2e.Test[apiv2.IPServiceGetResponse, *apiv2.IP]{
		{
			Name:    "create",
			CmdArgs: []string{"ip", "create", "--project", testresources.IP1().Project, "--network", testresources.IP1().Network, "--static=true"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: &apiv2.IPServiceCreateRequest{
							Project: testresources.IP1().Project,
							Network: testresources.IP1().Network,
							Type:    &testresources.IP1().Type,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.IPServiceCreateResponse{
								Ip: testresources.IP1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.IP1(),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"ip", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.IP1()), 0755))
					},
					ClientCalls: []e2e.ClientCall{
						{
							WantRequest: &apiv2.IPServiceGetRequest{
								Ip:      testresources.IP1().Ip,
								Project: testresources.IP1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: testresources.IP1(),
								})
							},
						},
						{
							WantRequest: &apiv2.IPServiceCreateRequest{
								Ip:            &testresources.IP1().Ip,
								Project:       testresources.IP1().Project,
								Network:       testresources.IP1().Network,
								Name:          &testresources.IP1().Name,
								Description:   &testresources.IP1().Description,
								Labels:        testresources.IP1().Meta.Labels,
								Type:          &testresources.IP1().Type,
								AddressFamily: nil,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceCreateResponse{
									Ip: testresources.IP1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
            IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE
            1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static  a
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv2.IPServiceDeleteResponse, *apiv2.IP]{
		{
			Name:    "delete",
			CmdArgs: []string{"ip", "delete", "--project", testresources.IP1().Project, testresources.IP1().Ip},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: &apiv2.IPServiceDeleteRequest{
							Ip:      testresources.IP1().Ip,
							Project: testresources.IP1().Project,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.IPServiceDeleteResponse{
								Ip: testresources.IP1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.IP1(),
		},
		{
			Name:    "delete from file",
			CmdArgs: append([]string{"ip", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.IP1()), 0755))
					},
					ClientCalls: []e2e.ClientCall{
						{
							WantRequest: &apiv2.IPServiceGetRequest{
								Ip:      testresources.IP1().Ip,
								Project: testresources.IP1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: testresources.IP1(),
								})
							},
						},
						{
							WantRequest: &apiv2.IPServiceDeleteRequest{
								Ip:      testresources.IP1().Ip,
								Project: testresources.IP1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceDeleteResponse{
									Ip: testresources.IP1(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
		    IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE
		    1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static  a
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv2.IPServiceDeleteResponse, *apiv2.IP]{
		{
			Name:    "update",
			CmdArgs: []string{"ip", "update", "--project", testresources.IP1().Project, testresources.IP1().Ip, "--name", "foo"},
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					ClientCalls: []e2e.ClientCall{
						// TODO: the client gets the IP two times?
						{
							WantRequest: &apiv2.IPServiceGetRequest{
								Ip:      testresources.IP1().Ip,
								Project: testresources.IP1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: testresources.IP1(),
								})
							},
						},
						{
							WantRequest: &apiv2.IPServiceGetRequest{
								Ip:      testresources.IP1().Ip,
								Project: testresources.IP1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: testresources.IP1(),
								})
							},
						},
						{
							WantRequest: &apiv2.IPServiceUpdateRequest{
								Ip:      testresources.IP1().Ip,
								Project: testresources.IP1().Project,
								Name:    new("foo"),

								// TODO: these fields do not need to be sent?
								Description: &testresources.IP1().Description,
								Labels: &apiv2.UpdateLabels{
									Update: &apiv2.Labels{},
								},
								Type: &testresources.IP1().Type,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceUpdateResponse{
									Ip: testresources.IP1(),
								})
							},
						},
					},
				},
			),
			WantObject: testresources.IP1(),
		},
		{
			Name:    "update from file",
			CmdArgs: append([]string{"ip", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.IP1()), 0755))
					},
					ClientCalls: []e2e.ClientCall{
						{
							WantRequest: &apiv2.IPServiceGetRequest{
								Ip:      testresources.IP1().Ip,
								Project: testresources.IP1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: testresources.IP1(),
								})
							},
						},
						{
							WantRequest: &apiv2.IPServiceUpdateRequest{
								Ip:          testresources.IP1().Ip,
								Project:     testresources.IP1().Project,
								Description: &testresources.IP1().Description,
								Labels: &apiv2.UpdateLabels{
									Update: &apiv2.Labels{},
								},
								Name: &testresources.IP1().Name,
								Type: &testresources.IP1().Type,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceUpdateResponse{
									Ip: testresources.IP1(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
		    IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE
		    1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static  a
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv2.IPServiceDeleteResponse, *apiv2.IP]{
		{
			Name:    "apply",
			CmdArgs: append([]string{"ip", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.IP1()), 0755))
					},
					ClientCalls: []e2e.ClientCall{
						{
							WantRequest: &apiv2.IPServiceGetRequest{
								Ip:      testresources.IP1().Ip,
								Project: testresources.IP1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: testresources.IP1(),
								})
							},
						},
						{
							WantRequest: &apiv2.IPServiceCreateRequest{
								Ip:            &testresources.IP1().Ip,
								Project:       testresources.IP1().Project,
								Network:       testresources.IP1().Network,
								Name:          &testresources.IP1().Name,
								Description:   &testresources.IP1().Description,
								Labels:        testresources.IP1().Meta.Labels,
								Type:          &testresources.IP1().Type,
								AddressFamily: nil,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceCreateResponse{
									Ip: testresources.IP1(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
		    IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE
		    1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static  a
			`),
		},
		{
			Name:    "apply already exists",
			CmdArgs: append([]string{"ip", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.IP1()), 0755))
					},
					ClientCalls: []e2e.ClientCall{
						{
							WantRequest: &apiv2.IPServiceGetRequest{
								Ip:      testresources.IP1().Ip,
								Project: testresources.IP1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: testresources.IP1(),
								})
							},
						},
						{
							WantRequest: &apiv2.IPServiceCreateRequest{
								Ip:            &testresources.IP1().Ip,
								Project:       testresources.IP1().Project,
								Network:       testresources.IP1().Network,
								Name:          &testresources.IP1().Name,
								Description:   &testresources.IP1().Description,
								Labels:        testresources.IP1().Meta.Labels,
								Type:          &testresources.IP1().Type,
								AddressFamily: nil,
							},
							WantError: connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")),
						},
						{
							WantRequest: &apiv2.IPServiceGetRequest{
								Ip:      testresources.IP1().Ip,
								Project: testresources.IP1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: testresources.IP1(),
								})
							},
						},
						{
							WantRequest: &apiv2.IPServiceUpdateRequest{
								Ip:          testresources.IP1().Ip,
								Project:     testresources.IP1().Project,
								Description: &testresources.IP1().Description,
								Labels: &apiv2.UpdateLabels{
									Update: &apiv2.Labels{},
								},
								Name: &testresources.IP1().Name,
								Type: &testresources.IP1().Type,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceUpdateResponse{
									Ip: testresources.IP1(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
		    IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE
		    1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static  a
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
