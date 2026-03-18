package cmd

import (
	"testing"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/tag"
)

var (
	ip1 = func() *apiv2.IP {
		return &apiv2.IP{
			Uuid:        "2e0144a2-09ef-42b7-b629-4263295db6e8",
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

// func Test_IPCmd_MultiResult(t *testing.T) {
// 	tests := []*Test[[]*apiv2.IP]{
// 		{
// 			Name: "list",
// 			Cmd: func(want []*apiv2.IP) []string {
// 				return []string{"ip", "list", "--project", "a"}
// 			},
// 			// ClientMocks: &apitests.ClientMockFns{
// 			// 	Apiv1Mocks: &apitests.Apiv1MockFns{
// 			// 		IP: func(m *mock.Mock) {
// 			// 			m.On("List", mock.Anything, connect.NewRequest(&apiv2.IPServiceListRequest{
// 			// 				Project: "a",
// 			// 			})).Return(&connect.Response[apiv2.IPServiceListResponse]{
// 			// 				Msg: &apiv2.IPServiceListResponse{
// 			// 					Ips: []*apiv2.IP{
// 			// 						ip2(),
// 			// 						ip1(),
// 			// 					},
// 			// 				},
// 			// 			}, nil)
// 			// 		},
// 			// 	},
// 			// },
// 			Want: []*apiv2.IP{
// 				ip1(),
// 				ip2(),
// 			},
// 			WantTable: new(`
// IP       PROJECT  ID                                    TYPE       NAME  ATTACHED SERVICE
// 1.1.1.1  a        2e0144a2-09ef-42b7-b629-4263295db6e8  static     a     ingress-nginx
// 4.3.2.1  b        9cef40ec-29c6-4dfa-aee8-47ee1f49223d  ephemeral  b
// `),
// 			WantWideTable: new(`
// IP       PROJECT  ID                                    TYPE       NAME  DESCRIPTION    LABELS
// 1.1.1.1  a        2e0144a2-09ef-42b7-b629-4263295db6e8  static     a     a description  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
// 4.3.2.1  b        9cef40ec-29c6-4dfa-aee8-47ee1f49223d  ephemeral  b     b description  a=b
// `),
// 			Template: new("{{ .ip }} {{ .project }}"),
// 			WantTemplate: new(`
// 1.1.1.1 a
// 4.3.2.1 b
// 			`),
// 			WantMarkdown: new(`
// | IP      | PROJECT | ID                                   | TYPE      | NAME | ATTACHED SERVICE |
// |---------|---------|--------------------------------------|-----------|------|------------------|
// | 1.1.1.1 | a       | 2e0144a2-09ef-42b7-b629-4263295db6e8 | static    | a    | ingress-nginx    |
// | 4.3.2.1 | b       | 9cef40ec-29c6-4dfa-aee8-47ee1f49223d | ephemeral | b    |                  |
// `),
// 		},
// 		// {
// 		// 	Name: "apply",
// 		// 	Cmd: func(want []*apiv2.IP) []string {
// 		// 		return appendFromFileCommonArgs("ip", "apply")
// 		// 	},
// 		// 	FsMocks: func(fs afero.Fs, want []*apiv2.IP) {
// 		// 		require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshalToMultiYAML(t, want), 0755))
// 		// 	},
// 		// 	// ClientMocks: &apitests.ClientMockFns{
// 		// 	// 	Apiv1Mocks: &apitests.Apiv1MockFns{
// 		// 	// 		IP: func(m *mock.Mock) {
// 		// 	// 			m.On("Allocate", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.IpResponseToCreate(ip1())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv2.IPServiceAllocateResponse{
// 		// 	// 				Ip: ip1(),
// 		// 	// 			}), nil)
// 		// 	// 			m.On("Allocate", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.IpResponseToCreate(ip2())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv2.IPServiceAllocateResponse{
// 		// 	// 				Ip: ip2(),
// 		// 	// 			}), nil)
// 		// 	// 			// FIXME: API does not return a conflict when already exists, so the update functionality does not work!
// 		// 	// 		},
// 		// 	// 	},
// 		// 	// },
// 		// 	Want: []*apiv2.IP{
// 		// 		ip1(),
// 		// 		ip2(),
// 		// 	},
// 		// },
// 		// {
// 		// 	Name: "update from file",
// 		// 	Cmd: func(want []*apiv2.IP) []string {
// 		// 		return appendFromFileCommonArgs("ip", "update")
// 		// 	},
// 		// 	FsMocks: func(fs afero.Fs, want []*apiv2.IP) {
// 		// 		require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshalToMultiYAML(t, want), 0755))
// 		// 	},
// 		// 	// ClientMocks: &apitests.ClientMockFns{
// 		// 	// 	Apiv1Mocks: &apitests.Apiv1MockFns{
// 		// 	// 		IP: func(m *mock.Mock) {
// 		// 	// 			m.On("Update", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv2.IPServiceUpdateRequest{
// 		// 	// 				Project: ip1().Project,
// 		// 	// 				Ip:      ip1(),
// 		// 	// 			}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv2.IPServiceUpdateResponse{
// 		// 	// 				Ip: ip1(),
// 		// 	// 			}), nil)
// 		// 	// 		},
// 		// 	// 	},
// 		// 	// },
// 		// 	Want: []*apiv2.IP{
// 		// 		ip1(),
// 		// 	},
// 		// },
// 		// {
// 		// 	Name: "create from file",
// 		// 	Cmd: func(want []*apiv2.IP) []string {
// 		// 		return appendFromFileCommonArgs("ip", "create")
// 		// 	},
// 		// 	FsMocks: func(fs afero.Fs, want []*apiv2.IP) {
// 		// 		require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshalToMultiYAML(t, want), 0755))
// 		// 	},
// 		// 	ClientMocks: &apitests.ClientMockFns{
// 		// 		Apiv1Mocks: &apitests.Apiv1MockFns{
// 		// 			IP: func(m *mock.Mock) {
// 		// 				m.On("Allocate", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.IpResponseToCreate(ip1())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv2.IPServiceAllocateResponse{
// 		// 					Ip: ip1(),
// 		// 				}), nil)
// 		// 			},
// 		// 		},
// 		// 	},
// 		// 	Want: []*apiv2.IP{
// 		// 		ip1(),
// 		// 	},
// 		// },
// 		// {
// 		// 	Name: "delete from file",
// 		// 	Cmd: func(want []*apiv2.IP) []string {
// 		// 		return appendFromFileCommonArgs("ip", "delete")
// 		// 	},
// 		// 	FsMocks: func(fs afero.Fs, want []*apiv2.IP) {
// 		// 		require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshalToMultiYAML(t, want), 0755))
// 		// 	},
// 		// 	ClientMocks: &apitests.ClientMockFns{
// 		// 		Apiv1Mocks: &apitests.Apiv1MockFns{
// 		// 			IP: func(m *mock.Mock) {
// 		// 				m.On("Delete", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv2.IPServiceDeleteRequest{
// 		// 					Uuid:    ip1().Uuid,
// 		// 					Project: ip1().Project,
// 		// 				}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv2.IPServiceDeleteResponse{
// 		// 					Ip: ip1(),
// 		// 				}), nil)
// 		// 			},
// 		// 		},
// 		// 	},
// 		// 	Want: []*apiv2.IP{
// 		// 		ip1(),
// 		// 	},
// 		// },
// 	}
// 	for _, tt := range tests {
// 		tt.TestCmd(t)
// 	}
// }

func Test_IPCmd_SingleResult(t *testing.T) {
	ip1 := ip1()

	tests := []*Test[apiv2.IPServiceGetRequest, apiv2.IPServiceGetResponse]{
		{
			Name: "describe",
			Cmd: func() []string {
				return []string{"ip", "describe", "--project", ip1.Project, ip1.Ip}
			},
			WantRequest: apiv2.IPServiceGetRequest{
				Ip:      ip1.Ip,
				Project: ip1.Project,
			},
			WantResponse: apiv2.IPServiceGetResponse{
				Ip: ip1,
			},
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
		// {
		// 	Name: "delete",
		// 	Cmd: func(want *apiv2.IP) []string {
		// 		return []string{"ip", "rm", "--project", want.Project, want.Uuid}
		// 	},
		// 	ClientMocks: &apitests.ClientMockFns{
		// 		Apiv1Mocks: &apitests.Apiv1MockFns{
		// 			IP: func(m *mock.Mock) {
		// 				m.On("Delete", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv2.IPServiceDeleteRequest{
		// 					Project: ip1().Project,
		// 					Uuid:    ip1().Uuid,
		// 				}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv2.IPServiceDeleteResponse{
		// 					Ip: ip1(),
		// 				}), nil)
		// 			},
		// 		},
		// 	},
		// 	Want: ip1(),
		// },
		// {
		// 	Name: "create",
		// 	Cmd: func(want *apiv2.IP) []string {
		// 		args := []string{"ip", "create", "--project", want.Project, "--description", want.Description, "--name", want.Name, "--tags", strings.Join(want.Tags, ",")}
		// 		if want.Type == apiv2.IPType_IP_TYPE_STATIC {
		// 			args = append(args, "--static")
		// 		}
		// 		AssertExhaustiveArgs(t, args, commonExcludedFileArgs()...)
		// 		return args
		// 	},
		// 	ClientMocks: &apitests.ClientMockFns{
		// 		Apiv1Mocks: &apitests.Apiv1MockFns{
		// 			IP: func(m *mock.Mock) {
		// 				m.On("Allocate", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.IpResponseToCreate(ip1())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv2.IPServiceAllocateResponse{
		// 					Ip: ip1(),
		// 				}), nil)
		// 			},
		// 		},
		// 	},
		// 	Want: ip1(),
		// },
		// {
		// 	Name: "update",
		// 	Cmd: func(want *apiv2.IP) []string {
		// 		args := []string{"ip", "update", want.Uuid, "--project", want.Project, "--description", want.Description, "--name", want.Name, "--tags", strings.Join(want.Tags, ",")}
		// 		if want.Type == apiv2.IPType_IP_TYPE_STATIC {
		// 			args = append(args, "--static")
		// 		}
		// 		AssertExhaustiveArgs(t, args, commonExcludedFileArgs()...)
		// 		return args
		// 	},
		// 	ClientMocks: &apitests.ClientMockFns{
		// 		Apiv1Mocks: &apitests.Apiv1MockFns{
		// 			IP: func(m *mock.Mock) {
		// 				m.On("Get", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv2.IPServiceGetRequest{
		// 					Uuid:    ip1().Uuid,
		// 					Project: ip1().Project,
		// 				}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv2.IPServiceGetResponse{
		// 					Ip: ip1(),
		// 				}), nil)

		// 				m.On("Update", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.IpResponseToUpdate(ip1())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv2.IPServiceUpdateResponse{
		// 					Ip: ip1(),
		// 				}), nil)
		// 			},
		// 		},
		// 	},
		// 	Want: ip1(),
		// },
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
