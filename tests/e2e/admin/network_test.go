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
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_AdminNetworkCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.NetworkServiceListResponse, apiv2.Network]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "network", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.NetworkServiceListRequest{
							Query: &apiv2.NetworkQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.NetworkServiceListResponse{
								Networks: []*apiv2.Network{
									testresources.Network2(),
									testresources.Network1(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                       NAME      TYPE      PROJECT                               PARTITION    NAT        PREFIXES                   PREFIX USAGE  IP USAGE
            6988ebb0-9531-4f9b-a893-d7868258e2ef     internet  external  0d81bca7-73f6-4da3-8397-4a8c52a0c583  partition-1  none       10.0.0.0/16,2001:db8::/32
            └─╴d83ffb0a-7aa6-4a66-8e03-0b5ee8b718a0  private   child     f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  partition-1  ipv4-masq  192.168.1.0/24
			`),
			WantWideTable: new(`
            ID                                       DESCRIPTION       NAME      TYPE      PROJECT                               PARTITION    NAT        PREFIXES                   ANNOTATIONS
            6988ebb0-9531-4f9b-a893-d7868258e2ef     internet network  internet  external  0d81bca7-73f6-4da3-8397-4a8c52a0c583  partition-1  none       10.0.0.0/16,2001:db8::/32  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
            └─╴d83ffb0a-7aa6-4a66-8e03-0b5ee8b718a0  private network   private   child     f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  partition-1  ipv4-masq  192.168.1.0/24             a=b
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
6988ebb0-9531-4f9b-a893-d7868258e2ef internet
d83ffb0a-7aa6-4a66-8e03-0b5ee8b718a0 private
			`),
			WantMarkdown: new(`
            | ID                                      | NAME     | TYPE     | PROJECT                              | PARTITION   | NAT       | PREFIXES                  | PREFIX USAGE | IP USAGE |
            |-----------------------------------------|----------|----------|--------------------------------------|-------------|-----------|---------------------------|--------------|----------|
            | 6988ebb0-9531-4f9b-a893-d7868258e2ef    | internet | external | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | partition-1 | none      | 10.0.0.0/16,2001:db8::/32 |              |          |
            | └─╴d83ffb0a-7aa6-4a66-8e03-0b5ee8b718a0 | private  | child    | f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | partition-1 | ipv4-masq | 192.168.1.0/24            |              |          |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminNetworkCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv2.NetworkServiceGetResponse, *apiv2.Network]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "network", "describe", testresources.Network1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.NetworkServiceGetRequest{
							Id: testresources.Network1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.NetworkServiceGetResponse{
								Network: testresources.Network1(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Network1(),
			WantTable: new(`
            ID                                    NAME      TYPE      PROJECT                               PARTITION    NAT   PREFIXES                   PREFIX USAGE  IP USAGE
            6988ebb0-9531-4f9b-a893-d7868258e2ef  internet  external  0d81bca7-73f6-4da3-8397-4a8c52a0c583  partition-1  none  10.0.0.0/16,2001:db8::/32
			`),
			WantWideTable: new(`
            ID                                    DESCRIPTION       NAME      TYPE      PROJECT                               PARTITION    NAT   PREFIXES                   ANNOTATIONS
            6988ebb0-9531-4f9b-a893-d7868258e2ef  internet network  internet  external  0d81bca7-73f6-4da3-8397-4a8c52a0c583  partition-1  none  10.0.0.0/16,2001:db8::/32  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
			6988ebb0-9531-4f9b-a893-d7868258e2ef internet
			`),
			WantMarkdown: new(`
            | ID                                   | NAME     | TYPE     | PROJECT                              | PARTITION   | NAT  | PREFIXES                  | PREFIX USAGE | IP USAGE |
            |--------------------------------------|----------|----------|--------------------------------------|-------------|------|---------------------------|--------------|----------|
            | 6988ebb0-9531-4f9b-a893-d7868258e2ef | internet | external | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | partition-1 | none | 10.0.0.0/16,2001:db8::/32 |              |          |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminNetworkCmd_Create(t *testing.T) {
	tests := []*e2e.Test[adminv2.NetworkServiceCreateResponse, *apiv2.Network]{
		{
			Name: "create",
			CmdArgs: []string{"admin", "network", "create",
				"--id", testresources.Network1().Id,
				"--name", *testresources.Network1().Name,
				"--type", "external",
				"--nat-type", "none",
				"--partition", *testresources.Network1().Partition,
				"--project", *testresources.Network1().Project,
				"--prefixes", "10.0.0.0/16,2001:db8::/32",
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.NetworkServiceCreateRequest{
							Id:        &testresources.Network1().Id,
							Name:      testresources.Network1().Name,
							Partition: testresources.Network1().Partition,
							Project:   testresources.Network1().Project,
							Type:      apiv2.NetworkType_NETWORK_TYPE_EXTERNAL,
							NatType:   apiv2.NATType_NAT_TYPE_NONE.Enum(),
							Prefixes:  []string{"10.0.0.0/16", "2001:db8::/32"},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.NetworkServiceCreateResponse{
								Network: testresources.Network1(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Network1(),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"admin", "network", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Network1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.NetworkServiceCreateRequest{
								Id:                         &testresources.Network1().Id,
								Name:                       testresources.Network1().Name,
								Description:                testresources.Network1().Description,
								Partition:                  testresources.Network1().Partition,
								Project:                    testresources.Network1().Project,
								Type:                       apiv2.NetworkType_NETWORK_TYPE_EXTERNAL,
								Prefixes:                   testresources.Network1().Prefixes,
								NatType:                    apiv2.NATType_NAT_TYPE_NONE.Enum(),
								ParentNetwork:              testresources.Network1().ParentNetwork,
								AdditionalAnnouncableCidrs: testresources.Network1().AdditionalAnnouncableCidrs,
								Labels:                     testresources.Network1().Meta.Labels,
								AddressFamily:              apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_DUAL_STACK.Enum(),
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.NetworkServiceCreateResponse{
									Network: testresources.Network1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
            ID                                    NAME      TYPE      PROJECT                               PARTITION    NAT   PREFIXES                   PREFIX USAGE  IP USAGE
            6988ebb0-9531-4f9b-a893-d7868258e2ef  internet  external  0d81bca7-73f6-4da3-8397-4a8c52a0c583  partition-1  none  10.0.0.0/16,2001:db8::/32
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminNetworkCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.NetworkServiceDeleteResponse, *apiv2.Network]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "network", "delete", testresources.Network1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.NetworkServiceDeleteRequest{
							Id: testresources.Network1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.NetworkServiceDeleteResponse{
								Network: testresources.Network1(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Network1(),
			WantTable: new(`
            ID                                    NAME      TYPE      PROJECT                               PARTITION    NAT   PREFIXES                   PREFIX USAGE  IP USAGE
            6988ebb0-9531-4f9b-a893-d7868258e2ef  internet  external  0d81bca7-73f6-4da3-8397-4a8c52a0c583  partition-1  none  10.0.0.0/16,2001:db8::/32
			`),
		},
		{
			Name:    "delete from file",
			CmdArgs: append([]string{"admin", "network", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Network1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.NetworkServiceDeleteRequest{
								Id: testresources.Network1().Id,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.NetworkServiceDeleteResponse{
									Network: testresources.Network1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
            ID                                    NAME      TYPE      PROJECT                               PARTITION    NAT   PREFIXES                   PREFIX USAGE  IP USAGE
            6988ebb0-9531-4f9b-a893-d7868258e2ef  internet  external  0d81bca7-73f6-4da3-8397-4a8c52a0c583  partition-1  none  10.0.0.0/16,2001:db8::/32
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminNetworkCmd_Update(t *testing.T) {
	tests := []*e2e.Test[adminv2.NetworkServiceUpdateResponse, *apiv2.Network]{
		{
			Name:    "update",
			CmdArgs: []string{"admin", "network", "update", testresources.Network1().Id, "--name", "foo"},
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.NetworkServiceUpdateRequest{
								Id:   testresources.Network1().Id,
								Name: new("foo"),
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.NetworkServiceUpdateResponse{
									Network: testresources.Network1(),
								})
							},
						},
					},
				}),
			WantProtoObject: testresources.Network1(),
			WantTable: new(`
            ID                                    NAME      TYPE      PROJECT                               PARTITION    NAT   PREFIXES                   PREFIX USAGE  IP USAGE
            6988ebb0-9531-4f9b-a893-d7868258e2ef  internet  external  0d81bca7-73f6-4da3-8397-4a8c52a0c583  partition-1  none  10.0.0.0/16,2001:db8::/32
			`),
		},
		{
			Name:    "update from file",
			CmdArgs: append([]string{"admin", "network", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Network1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.NetworkServiceUpdateRequest{
								Id:          testresources.Network1().Id,
								Description: testresources.Network1().Description,
								Name:        testresources.Network1().Name,
								Labels: &apiv2.UpdateLabels{
									Strategy: &apiv2.UpdateLabels_Replace{
										Replace: &apiv2.Labels{
											Labels: testresources.Network1().Meta.Labels.Labels,
										},
									},
								},
								NatType:  &testresources.Network1().NatType,
								Prefixes: testresources.Network1().Prefixes,
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.NetworkServiceUpdateResponse{
									Network: testresources.Network1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
            ID                                    NAME      TYPE      PROJECT                               PARTITION    NAT   PREFIXES                   PREFIX USAGE  IP USAGE
            6988ebb0-9531-4f9b-a893-d7868258e2ef  internet  external  0d81bca7-73f6-4da3-8397-4a8c52a0c583  partition-1  none  10.0.0.0/16,2001:db8::/32
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
