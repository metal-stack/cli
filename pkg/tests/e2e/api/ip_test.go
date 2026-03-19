package api_e2e

import (
	"fmt"
	"testing"

	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/tests/e2e"
	"github.com/metal-stack/metal-lib/pkg/tag"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var (
	ip1 = func() *apiv2.IP {
		return &apiv2.IP{
			Uuid:        "2e0144a2-09ef-42b7-b629-4263295db6e8",
			Network:     "internet",
			Ip:          "1.1.1.1",
			Name:        "a",
			Description: "a description",
			Project:     "ce19a655-7933-4745-8f3e-9592b4a90488",
			Type:        apiv2.IPType_IP_TYPE_STATIC,
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						tag.ClusterServiceFQN: "<cluster>/default/ingress-nginx",
					},
				},
			},
		}
	}
	ip2 = func() *apiv2.IP {
		return &apiv2.IP{
			Uuid:        "9cef40ec-29c6-4dfa-aee8-47ee1f49223d",
			Network:     "internet",
			Ip:          "4.3.2.1",
			Name:        "b",
			Description: "b description",
			Project:     "46bdfc45-9c8d-4268-b359-b40e3079d384",
			Type:        apiv2.IPType_IP_TYPE_EPHEMERAL,
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"a": "b",
					},
				},
			},
		}
	}
)

func Test_IPCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.IPServiceListResponse, apiv2.IP]{
		{
			Name:    "list",
			CmdArgs: []string{"ip", "list", "--project", ip1().Project},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.IPServiceListRequest{
							Project: ip1().Project,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.IPServiceListResponse{
								Ips: []*apiv2.IP{
									ip1(),
									ip2(),
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
			CmdArgs: []string{"ip", "describe", "--project", ip1().Project, ip1().Ip},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.IPServiceGetRequest{
							Ip:      ip1().Ip,
							Project: ip1().Project,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.IPServiceGetResponse{
								Ip: ip1(),
							})
						},
					},
				},
			}),
			WantObject:      ip1(),
			WantProtoObject: ip1(),
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
			CmdArgs: []string{"ip", "create", "--project", ip1().Project, "--network", ip1().Network, "--static=true"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.IPServiceCreateRequest{
							Project: ip1().Project,
							Network: ip1().Network,
							Type:    &ip1().Type,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.IPServiceCreateResponse{
								Ip: ip1(),
							})
						},
					},
				},
			}),
			WantObject: ip1(),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"ip", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, ip1()), 0755))
					},
					ClientCalls: []e2e.ClientCall{
						{
							WantRequest: apiv2.IPServiceGetRequest{
								Ip:      ip1().Ip,
								Project: ip1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: ip1(),
								})
							},
						},
						{
							WantRequest: apiv2.IPServiceCreateRequest{
								Ip:            &ip1().Ip,
								Project:       ip1().Project,
								Network:       ip1().Network,
								Name:          &ip1().Name,
								Description:   &ip1().Description,
								Labels:        ip1().Meta.Labels,
								Type:          &ip1().Type,
								AddressFamily: nil,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceCreateResponse{
									Ip: ip1(),
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
			CmdArgs: []string{"ip", "delete", "--project", ip1().Project, ip1().Ip},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.IPServiceDeleteRequest{
							Ip:      ip1().Ip,
							Project: ip1().Project,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.IPServiceDeleteResponse{
								Ip: ip1(),
							})
						},
					},
				},
			}),
			WantObject: ip1(),
		},
		{
			Name:    "delete from file",
			CmdArgs: append([]string{"ip", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, ip1()), 0755))
					},
					ClientCalls: []e2e.ClientCall{
						{
							WantRequest: apiv2.IPServiceGetRequest{
								Ip:      ip1().Ip,
								Project: ip1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: ip1(),
								})
							},
						},
						{
							WantRequest: apiv2.IPServiceDeleteRequest{
								Ip:      ip1().Ip,
								Project: ip1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceDeleteResponse{
									Ip: ip1(),
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
			CmdArgs: []string{"ip", "update", "--project", ip1().Project, ip1().Ip, "--name", "foo"},
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					ClientCalls: []e2e.ClientCall{
						// TODO: the client gets the IP two times?
						{
							WantRequest: apiv2.IPServiceGetRequest{
								Ip:      ip1().Ip,
								Project: ip1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: ip1(),
								})
							},
						},
						{
							WantRequest: apiv2.IPServiceGetRequest{
								Ip:      ip1().Ip,
								Project: ip1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: ip1(),
								})
							},
						},
						{
							WantRequest: apiv2.IPServiceUpdateRequest{
								Ip:      ip1().Ip,
								Project: ip1().Project,
								Name:    new("foo"),

								// TODO: these fields do not need to be sent?
								Description: &ip1().Description,
								Labels: &apiv2.UpdateLabels{
									Update: &apiv2.Labels{},
								},
								Type: &ip1().Type,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceUpdateResponse{
									Ip: ip1(),
								})
							},
						},
					},
				},
			),
			WantObject: ip1(),
		},
		{
			Name:    "update from file",
			CmdArgs: append([]string{"ip", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, ip1()), 0755))
					},
					ClientCalls: []e2e.ClientCall{
						{
							WantRequest: apiv2.IPServiceGetRequest{
								Ip:      ip1().Ip,
								Project: ip1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: ip1(),
								})
							},
						},
						{
							WantRequest: apiv2.IPServiceUpdateRequest{
								Ip:          ip1().Ip,
								Project:     ip1().Project,
								Description: &ip1().Description,
								Labels: &apiv2.UpdateLabels{
									Update: &apiv2.Labels{},
								},
								Name: &ip1().Name,
								Type: &ip1().Type,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceUpdateResponse{
									Ip: ip1(),
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
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, ip1()), 0755))
					},
					ClientCalls: []e2e.ClientCall{
						{
							WantRequest: apiv2.IPServiceGetRequest{
								Ip:      ip1().Ip,
								Project: ip1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: ip1(),
								})
							},
						},
						{
							WantRequest: apiv2.IPServiceCreateRequest{
								Ip:            &ip1().Ip,
								Project:       ip1().Project,
								Network:       ip1().Network,
								Name:          &ip1().Name,
								Description:   &ip1().Description,
								Labels:        ip1().Meta.Labels,
								Type:          &ip1().Type,
								AddressFamily: nil,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceCreateResponse{
									Ip: ip1(),
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
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, ip1()), 0755))
					},
					ClientCalls: []e2e.ClientCall{
						{
							WantRequest: apiv2.IPServiceGetRequest{
								Ip:      ip1().Ip,
								Project: ip1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: ip1(),
								})
							},
						},
						{
							WantRequest: apiv2.IPServiceCreateRequest{
								Ip:            &ip1().Ip,
								Project:       ip1().Project,
								Network:       ip1().Network,
								Name:          &ip1().Name,
								Description:   &ip1().Description,
								Labels:        ip1().Meta.Labels,
								Type:          &ip1().Type,
								AddressFamily: nil,
							},
							WantError: connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")),
						},
						{
							WantRequest: apiv2.IPServiceGetRequest{
								Ip:      ip1().Ip,
								Project: ip1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceGetResponse{
									Ip: ip1(),
								})
							},
						},
						{
							WantRequest: apiv2.IPServiceUpdateRequest{
								Ip:          ip1().Ip,
								Project:     ip1().Project,
								Description: &ip1().Description,
								Labels: &apiv2.UpdateLabels{
									Update: &apiv2.Labels{},
								},
								Name: &ip1().Name,
								Type: &ip1().Type,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.IPServiceUpdateResponse{
									Ip: ip1(),
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
