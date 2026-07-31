package api_e2e

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_NetworkCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.NetworkServiceListResponse, apiv2.Network]{
		{
			Name: "list",
			CmdArgs: []string{"network", "list",
				"--project", *testresources.Network1().Project,
				"--addressfamily", "dual-stack",
				"--description", *testresources.Network1().Description,
				"--destination-prefixes", strings.Join(testresources.Network1().DestinationPrefixes, ","),
				"--prefixes", strings.Join(testresources.Network1().Prefixes, ","),
				"--id", testresources.Network1().Id,
				"--vrf", strconv.Itoa(int(pointer.SafeDeref(testresources.Network1().Vrf))),
				"--type", "external",
				"--labels", "a=b",
				"--name", *testresources.Network1().Name,
				"--partition", *testresources.Network1().Partition,
				"--parent-network", *testresources.Network2().ParentNetwork,
			},
			AssertExhaustiveArgs:     true,
			AssertExhaustiveExcludes: []string{"sort-by"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.NetworkServiceListRequest{
							Project: *testresources.Network1().Project,
							Query: &apiv2.NetworkQuery{
								Project:     testresources.Network1().Project,
								Id:          &testresources.Network1().Id,
								Name:        testresources.Network1().Name,
								Description: testresources.Network1().Description,
								Partition:   testresources.Network1().Partition,
								// Namespace:           new(string),
								// NatType:             ,
								Prefixes:            testresources.Network1().Prefixes,
								DestinationPrefixes: testresources.Network1().DestinationPrefixes,
								Vrf:                 testresources.Network1().Vrf,
								ParentNetwork:       testresources.Network2().ParentNetwork,
								AddressFamily:       apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_DUAL_STACK.Enum(),
								Type:                &testresources.Network1().Type,
								Labels: &apiv2.Labels{
									Labels: map[string]string{"a": "b"},
								},
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.NetworkServiceListResponse{
								Networks: []*apiv2.Network{
									testresources.Network1(),
									testresources.Network2(),
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

func Test_NetworkCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.NetworkServiceGetResponse, *apiv2.Network]{
		{
			Name:    "describe",
			CmdArgs: []string{"network", "describe", testresources.Network1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.NetworkServiceGetRequest{
							Id:      testresources.Network1().Id,
							Project: "",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.NetworkServiceGetResponse{
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

func Test_NetworkCmd_Create(t *testing.T) {
	tests := []*e2e.Test[apiv2.NetworkServiceCreateResponse, *apiv2.Network]{
		{
			Name: "create",
			CmdArgs: []string{"network", "create",
				"--project", *testresources.Network1().Project,
				"--name", *testresources.Network1().Name,
				"--partition", *testresources.Network1().Partition,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.NetworkServiceCreateRequest{
							Project:   *testresources.Network1().Project,
							Name:      testresources.Network1().Name,
							Partition: testresources.Network1().Partition,
							Length:    &apiv2.ChildPrefixLength{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.NetworkServiceCreateResponse{
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
			CmdArgs: append([]string{"network", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Network1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.NetworkServiceCreateRequest{
								Project:       *testresources.Network1().Project,
								Name:          testresources.Network1().Name,
								Description:   testresources.Network1().Description,
								Partition:     testresources.Network1().Partition,
								Labels:        testresources.Network1().Meta.Labels,
								ParentNetwork: testresources.Network1().ParentNetwork,
								AddressFamily: apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_DUAL_STACK.Enum(),
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.NetworkServiceCreateResponse{
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

func Test_NetworkCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv2.NetworkServiceDeleteResponse, *apiv2.Network]{
		{
			Name:    "delete",
			CmdArgs: []string{"network", "delete", testresources.Network1().Id, "--project", *testresources.Network1().Project},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.NetworkServiceDeleteRequest{
							Id:      testresources.Network1().Id,
							Project: *testresources.Network1().Project,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.NetworkServiceDeleteResponse{
								Network: testresources.Network1(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Network1(),
		},
		{
			Name:    "delete from file",
			CmdArgs: append([]string{"network", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Network1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.NetworkServiceDeleteRequest{
								Id:      testresources.Network1().Id,
								Project: *testresources.Network1().Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.NetworkServiceDeleteResponse{
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

func Test_NetworkCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv2.NetworkServiceUpdateResponse, *apiv2.Network]{
		{
			Name:    "update",
			CmdArgs: []string{"network", "update", testresources.Network1().Id, "--project", *testresources.Network1().Project, "--name", "foo"},
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.NetworkServiceUpdateRequest{
								Id:      testresources.Network1().Id,
								Project: *testresources.Network1().Project,
								Name:    new("foo"),
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.NetworkServiceUpdateResponse{
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
		},
		{
			Name:    "update from file",
			CmdArgs: append([]string{"network", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Network1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.NetworkServiceUpdateRequest{
								Id:          testresources.Network1().Id,
								Project:     *testresources.Network1().Project,
								Description: testresources.Network1().Description,
								Name:        testresources.Network1().Name,
								Labels: &apiv2.UpdateLabels{
									Strategy: &apiv2.UpdateLabels_Replace{
										Replace: &apiv2.Labels{
											Labels: testresources.Network1().Meta.Labels.Labels,
										},
									},
								},
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.NetworkServiceUpdateResponse{
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

func Test_NetworkCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv2.NetworkServiceCreateResponse, *apiv2.Network]{
		{
			Name:    "apply",
			CmdArgs: append([]string{"network", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Network1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.NetworkServiceCreateRequest{
								Project:       *testresources.Network1().Project,
								Name:          testresources.Network1().Name,
								Description:   testresources.Network1().Description,
								Partition:     testresources.Network1().Partition,
								Labels:        testresources.Network1().Meta.Labels,
								ParentNetwork: testresources.Network1().ParentNetwork,
								AddressFamily: apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_DUAL_STACK.Enum(),
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.NetworkServiceCreateResponse{
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
		{
			Name:    "apply already exists",
			CmdArgs: append([]string{"network", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Network1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.NetworkServiceCreateRequest{
								Project:       *testresources.Network1().Project,
								Name:          testresources.Network1().Name,
								Description:   testresources.Network1().Description,
								Partition:     testresources.Network1().Partition,
								Labels:        testresources.Network1().Meta.Labels,
								ParentNetwork: testresources.Network1().ParentNetwork,
								AddressFamily: apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_DUAL_STACK.Enum(),
							},
							WantError: connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")),
						},
						{
							WantRequest: &apiv2.NetworkServiceUpdateRequest{
								Id:          testresources.Network1().Id,
								Project:     *testresources.Network1().Project,
								Description: testresources.Network1().Description,
								Name:        testresources.Network1().Name,
								Labels: &apiv2.UpdateLabels{
									Strategy: &apiv2.UpdateLabels_Replace{
										Replace: &apiv2.Labels{
											Labels: testresources.Network1().Meta.Labels.Labels,
										},
									},
								},
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.NetworkServiceUpdateResponse{
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
