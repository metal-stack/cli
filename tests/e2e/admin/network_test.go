package admin_e2e

import (
	"strconv"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_AdminNetworkCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.NetworkServiceListResponse, apiv2.Network]{
		{
			Name: "list",
			CmdArgs: []string{"admin", "network", "list",
				"--project", *testresources.Network1().Project,
				"--addressfamily", "dual-stack",
				"--description", *testresources.Network1().Description,
				"--destination-prefixes", strings.Join(testresources.Network1().DestinationPrefixes, ","),
				"--prefixes", strings.Join(testresources.Network1().Prefixes, ","),
				"--id", testresources.Network1().Id,
				"--vrf", strconv.Itoa(int(pointer.SafeDeref(testresources.Network1().Vrf))),
				"--type", "external",
				"--nat-type", "none",
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
						WantRequest: &adminv2.NetworkServiceListRequest{
							Query: &apiv2.NetworkQuery{
								Project:     testresources.Network1().Project,
								Id:          new(testresources.Network1().Id),
								Name:        testresources.Network1().Name,
								Description: testresources.Network1().Description,
								Partition:   testresources.Network1().Partition,
								// Namespace:           new(string),
								NatType:             apiv2.NATType_NAT_TYPE_NONE.Enum(),
								Prefixes:            testresources.Network1().Prefixes,
								DestinationPrefixes: testresources.Network1().DestinationPrefixes,
								Vrf:                 testresources.Network1().Vrf,
								ParentNetwork:       testresources.Network2().ParentNetwork,
								AddressFamily:       apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_DUAL_STACK.Enum(),
								Type:                new(testresources.Network1().Type),
								Labels: &apiv2.Labels{
									Labels: map[string]string{"a": "b"},
								},
							},
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
				"--description", *testresources.Network1().Description,
				"--type", "external",
				"--nat-type", "none",
				"--parent-network", *testresources.Network2().ParentNetwork,
				"--partition", *testresources.Network1().Partition,
				"--project", *testresources.Network1().Project,
				"--prefixes", "10.0.0.0/16,2001:db8::/32",
				"--destination-prefixes", strings.Join(testresources.Network1().DestinationPrefixes, ","),
				"--additional-announcable-cidrs", "1.2.3.4/32",
				"--default-ipv4-prefix-length", "17",
				"--default-ipv6-prefix-length", "64",
				"--ipv4-prefix-length", "18",
				"--ipv6-prefix-length", "65",
				"--labels", "a=b",
				"--min-ipv4-prefix-length", "14",
				"--min-ipv6-prefix-length", "60",
				"--vrf", strconv.Itoa(int(pointer.SafeDeref(testresources.Network1().Vrf))),
			},
			AssertExhaustiveArgs:     true,
			AssertExhaustiveExcludes: append(e2e.CommonExcludedFileArgs(), "addressfamily"),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.NetworkServiceCreateRequest{
							Id:                         new(testresources.Network1().Id),
							Name:                       testresources.Network1().Name,
							Partition:                  testresources.Network1().Partition,
							Project:                    testresources.Network1().Project,
							Type:                       apiv2.NetworkType_NETWORK_TYPE_EXTERNAL,
							NatType:                    apiv2.NATType_NAT_TYPE_NONE.Enum(),
							Prefixes:                   []string{"10.0.0.0/16", "2001:db8::/32"},
							AdditionalAnnouncableCidrs: []string{"1.2.3.4/32"},
							Description:                testresources.Network1().Description,
							DestinationPrefixes:        testresources.Network1().DestinationPrefixes,
							Vrf:                        new(uint32(0)),
							ParentNetwork:              testresources.Network2().ParentNetwork,
							AddressFamily:              nil,
							Labels: &apiv2.Labels{
								Labels: map[string]string{"a": "b"},
							},
							DefaultChildPrefixLength: &apiv2.ChildPrefixLength{
								Ipv4: new(uint32(17)),
								Ipv6: new(uint32(64)),
							},
							MinChildPrefixLength: &apiv2.ChildPrefixLength{
								Ipv4: new(uint32(14)),
								Ipv6: new(uint32(60)),
							},
							Length: &apiv2.ChildPrefixLength{
								Ipv4: new(uint32(18)),
								Ipv6: new(uint32(65)),
							},
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
								Id:                         new(testresources.Network1().Id),
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
								AddressFamily:              nil,
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
								NatType:  new(testresources.Network1().NatType),
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

func Test_AdminNetwork_Issues(t *testing.T) {
	response := &adminv2.NetworkServiceIssuesResponse{
		NetworksMissingInIpam:  []*apiv2.Network{testresources.Network1()},
		NetworksMissingInMetal: []*apiv2.Network{testresources.Network2()},
		IpsMissingInIpam:       []*apiv2.IP{testresources.IP1()},
		IpsMissingInMetal:      []*apiv2.IP{testresources.IP2()},
	}

	tests := []*e2e.Test[adminv2.NetworkServiceIssuesResponse, *adminv2.NetworkServiceIssuesResponse]{
		{
			Name:    "network-issues",
			CmdArgs: []string{"admin", "network", "issues"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.NetworkServiceIssuesRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(response)
						},
					},
				},
			}),
			WantDefault: new(`
networksMissingInIpam:
    - id: 6988ebb0-9531-4f9b-a893-d7868258e2ef
      meta:
        labels:
            labels:
                cluster.metal-stack.io/id/namespace/service: <cluster>/default/ingress-nginx
        createdAt: "2000-01-01T00:00:00Z"
      name: internet
      description: internet network
      partition: partition-1
      project: 0d81bca7-73f6-4da3-8397-4a8c52a0c583
      prefixes:
        - 10.0.0.0/16
        - 2001:db8::/32
      type: NETWORK_TYPE_EXTERNAL
      natType: NAT_TYPE_NONE
networksMissingInMetal:
    - id: d83ffb0a-7aa6-4a66-8e03-0b5ee8b718a0
      meta:
        labels:
            labels:
                a: b
        createdAt: "2000-01-01T00:00:00Z"
      name: private
      description: private network
      partition: partition-1
      project: f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c
      prefixes:
        - 192.168.1.0/24
      type: NETWORK_TYPE_CHILD
      natType: NAT_TYPE_IPV4_MASQUERADE
      vrf: 100
      parentNetwork: 6988ebb0-9531-4f9b-a893-d7868258e2ef
ipsMissingInIpam:
    - uuid: 2e0144a2-09ef-42b7-b629-4263295db6e8
      meta:
        labels:
            labels:
                cluster.metal-stack.io/id/namespace/service: <cluster>/default/ingress-nginx
      ip: 1.1.1.1
      name: aa
      description: a description
      network: internet
      project: ce19a655-7933-4745-8f3e-9592b4a90488
      type: IP_TYPE_STATIC
ipsMissingInMetal:
    - uuid: 9cef40ec-29c6-4dfa-aee8-47ee1f49223d
      meta:
        labels:
            labels:
                a: b
      ip: 4.3.2.1
      name: bb
      description: b description
      network: internet
      project: 46bdfc45-9c8d-4268-b359-b40e3079d384
      type: IP_TYPE_EPHEMERAL
run the following command inside of a pod of the ipam deployment to get the list of existing prefixes:
kubectl exec -it ipam-xxx-yyy -- /cli prefix list

then select one of the networks listed in networksMissingInMetal and delete it in ipam with
kubectl exec -it ipam-xxx-yyy -- /cli prefix delete --cidr 1.2.3.0/24 [--namespace optional-namespace]


run the following command inside of a pod of the ipam deployment to get the list of existing ips:
kubectl exec -it ipam-xxx-yyy -- /cli ip list --prefix prefix-of-the-given-metal-network

then select one of the ips listed in ipsMissingInMetal and delete it in ipam with
kubectl exec -it ipam-xxx-yyy -- /cli ip release --ip 1.2.3.4 --prefix 1.2.3.0/24 [--namespace optional-namespace]

			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
