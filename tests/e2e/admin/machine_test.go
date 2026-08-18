package admin_e2e

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_MachineCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.MachineServiceListResponse, apiv2.Machine]{
		{
			Name: "list",
			CmdArgs: []string{"admin", "machine", "list",
				"--name", "name",
				"--hostname", "hostname",
				"--image", "image",
				"--partition", "partition",
				"--project", "project",
				"--size", "size",
				"--id", "uuid",
			},
			AssertExhaustiveArgs:     true,
			AssertExhaustiveExcludes: []string{"sort-by"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceListRequest{
							Query: &apiv2.MachineQuery{
								Allocation: &apiv2.MachineAllocationQuery{
									Hostname: new("hostname"),
									Name:     new("name"),
									Image:    new("image"),
									Project:  new("project"),
								},
								Partition: new("partition"),
								Size:      new("size"),
								Uuid:      new("uuid"),
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceListResponse{
								Machines: []*apiv2.Machine{
									testresources.Machine2(),
									testresources.Machine1(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                        LAST EVENT   WHEN  AGE  HOSTNAME   PROJECT                               SIZE           IMAGE         PARTITION    RACK
            5fa2bbe1-407c-4142-92d5-e4419daf9646      Alive        1m                                                          v1-medium-x86                partition-1  rack-1
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  🛡  Phoned Home  1m    1m   machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  v1-medium-x86  Ubuntu 24.04  partition-2  rack-1
            	`),
			WantWideTable: new(`
            ID                                     LAST EVENT   WHEN  AGE  DESCRIPTION  NAME       HOSTNAME   PROJECT                               IPS        SIZE           IMAGE         PARTITION    RACK    STARTED               TAGS  STATE
            5fa2bbe1-407c-4142-92d5-e4419daf9646   Alive        1m                                                                                             v1-medium-x86                partition-1  rack-1                        a=b   available
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  Phoned Home  1m    1m   machine 2    machine-2  machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  4.5.6.7    v1-medium-x86  Ubuntu 24.04  partition-2  rack-1  1999-12-31T23:59:00Z  c=d   available
                                                                                                                                                    192.1.1.1
			`),
			Template: new("{{ .allocation.uuid }} {{ .allocation.project }}"),
			WantTemplate: new(`
<no value> <no value>
4f94e87b-b08f-4f82-b053-9b8305de60ad f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c
			`),
			WantMarkdown: new(`
            | ID                                    |   | LAST EVENT  | WHEN | AGE | HOSTNAME  | PROJECT                              | SIZE          | IMAGE        | PARTITION   | RACK   |
            |---------------------------------------|---|-------------|------|-----|-----------|--------------------------------------|---------------|--------------|-------------|--------|
            | 5fa2bbe1-407c-4142-92d5-e4419daf9646  |   | Alive       | 1m   |     |           |                                      | v1-medium-x86 |              | partition-1 | rack-1 |
            | 673fc473-63ca-4ea4-b9dd-b45cb2127a6fd | 🛡 | Phoned Home | 1m   | 1m  | machine-2 | f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | v1-medium-x86 | Ubuntu 24.04 | partition-2 | rack-1 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_MachineCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.MachineServiceGetResponse, *apiv2.Machine]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "machine", "describe", testresources.Machine2().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceGetRequest{
							Uuid: testresources.Machine2().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceGetResponse{
								Machine: testresources.Machine2(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Machine2(),
			WantTable: new(`
            ID                                        LAST EVENT   WHEN  AGE  HOSTNAME   PROJECT                               SIZE           IMAGE         PARTITION    RACK
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  🛡  Phoned Home  1m    1m   machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  v1-medium-x86  Ubuntu 24.04  partition-2  rack-1
			`),
			WantWideTable: new(`
            ID                                     LAST EVENT   WHEN  AGE  DESCRIPTION  NAME       HOSTNAME   PROJECT                               IPS        SIZE           IMAGE         PARTITION    RACK    STARTED               TAGS  STATE
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  Phoned Home  1m    1m   machine 2    machine-2  machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  4.5.6.7    v1-medium-x86  Ubuntu 24.04  partition-2  rack-1  1999-12-31T23:59:00Z  c=d   available
                                                                                                                                                    192.1.1.1
			`),
			Template: new("{{ .allocation.uuid }} {{ .allocation.project }}"),
			WantTemplate: new(`
			4f94e87b-b08f-4f82-b053-9b8305de60ad f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c
			`),
			WantMarkdown: new(`
            | ID                                    |   | LAST EVENT  | WHEN | AGE | HOSTNAME  | PROJECT                              | SIZE          | IMAGE        | PARTITION   | RACK   |
            |---------------------------------------|---|-------------|------|-----|-----------|--------------------------------------|---------------|--------------|-------------|--------|
            | 673fc473-63ca-4ea4-b9dd-b45cb2127a6fd | 🛡 | Phoned Home | 1m   | 1m  | machine-2 | f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | v1-medium-x86 | Ubuntu 24.04 | partition-2 | rack-1 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_MachineCmd_Create(t *testing.T) {
	tests := []*e2e.Test[apiv2.MachineServiceGetResponse, *apiv2.Machine]{
		{
			Name: "create",
			CmdArgs: []string{"admin", "machine", "create",
				"--id", testresources.Machine2().Uuid,
				"--project", testresources.Machine2().Allocation.Project,
				"--networks", func() string {
					var names []string

					for _, nw := range testresources.Machine2().Allocation.Networks {
						nwWithIps := nw.Network
						if len(nw.Ips) > 0 {
							nwWithIps += ":" + strings.Join(nw.Ips, ";")
						}
						names = append(names, nwWithIps)
					}

					return strings.Join(names, ",")
				}(),
				"--allocation-type", testresources.Machine2().Allocation.AllocationType.String(),
				"--description", testresources.Machine2().Allocation.Description,
				"--dns-servers", "1.1.1.1",
				"--ntp-servers", "2.2.2.2,3.3.3.3",
				"--filesystem-layout", "fsl1",
				"--hostname", testresources.Machine2().Allocation.Hostname,
				"--image", testresources.Image1().Id,
				"--name", testresources.Machine2().Allocation.Name,
				"--partition", testresources.Machine2().Partition.Id,
				"--size", testresources.Machine2().Size.Id,
				"--ssh-public-key", "@.ssh/id_rsa.pub",
				"--labels", "a=b",
				"--userdata", "@ignition.json",
				"--placement-labels", "cluster-id=cluster-uuid",
			},
			AssertExhaustiveArgs:     true,
			AssertExhaustiveExcludes: e2e.CommonExcludedFileArgs(),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					genericcli.Must(fs.WriteFile(".ssh/id_rsa.pub", []byte("12345"), os.ModeAppend))
					genericcli.Must(fs.WriteFile("ignition.json", []byte("{}"), os.ModeAppend))
				},
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.MachineServiceCreateRequest{
							Uuid:             &testresources.Machine2().Uuid,
							Project:          testresources.Machine2().Allocation.Project,
							Name:             testresources.Machine2().Allocation.Name,
							Description:      &testresources.Machine2().Allocation.Description,
							Hostname:         &testresources.Machine2().Allocation.Hostname,
							Partition:        &testresources.Machine2().Partition.Id,
							Size:             &testresources.Machine2().Size.Id,
							Image:            testresources.Machine2().Allocation.Image.Id,
							FilesystemLayout: new("fsl1"),
							SshPublicKeys:    []string{"12345"},
							Userdata:         new(base64.StdEncoding.EncodeToString([]byte("{}"))),
							Labels: &apiv2.Labels{
								Labels: map[string]string{
									"a": "b",
								},
							},
							Networks: func() []*apiv2.MachineAllocationNetwork {
								var nws []*apiv2.MachineAllocationNetwork

								for _, nw := range testresources.Machine2().Allocation.Networks {
									nws = append(nws, &apiv2.MachineAllocationNetwork{
										Network: nw.Network,
										Ips:     nw.Ips,
									})
								}

								return nws
							}(),
							PlacementLabels: &apiv2.Labels{
								Labels: map[string]string{
									"cluster-id": "cluster-uuid",
								},
							},
							DnsServers:     []*apiv2.DNSServer{{Ip: "1.1.1.1"}},
							NtpServers:     []*apiv2.NTPServer{{Address: "2.2.2.2"}, {Address: "3.3.3.3"}},
							AllocationType: apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_MACHINE,
							FirewallSpec:   nil,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.MachineServiceCreateResponse{
								Machine: testresources.Machine2(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Machine2(),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"admin", "machine", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine2()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.MachineServiceCreateRequest{
								Uuid:        &testresources.Machine2().Uuid,
								Project:     testresources.Machine2().Allocation.Project,
								Name:        testresources.Machine2().Allocation.Name,
								Description: &testresources.Machine2().Allocation.Description,
								Hostname:    &testresources.Machine2().Allocation.Hostname,
								Image:       testresources.Machine2().Allocation.Image.Id,
								Userdata:    nil,
								Labels:      testresources.Machine2().Meta.Labels,
								Networks: func() []*apiv2.MachineAllocationNetwork {
									var nws []*apiv2.MachineAllocationNetwork

									for _, nw := range testresources.Machine2().Allocation.Networks {
										nws = append(nws, &apiv2.MachineAllocationNetwork{
											Network: nw.Network,
											Ips:     nw.Ips,
										})
									}

									return nws
								}(),
								AllocationType: apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_MACHINE,
								FirewallSpec:   nil,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.MachineServiceCreateResponse{
									Machine: testresources.Machine2(),
								})
							},
						},
					},
				}),
			WantTable: new(`
		    ID                                        LAST EVENT   WHEN  AGE  HOSTNAME   PROJECT                               SIZE           IMAGE         PARTITION    RACK
		    673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  🛡  Phoned Home  1m    1m   machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  v1-medium-x86  Ubuntu 24.04  partition-2  rack-1
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_MachineCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv2.MachineServiceDeleteResponse, *apiv2.Machine]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "machine", "delete", testresources.Machine2().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceDeleteRequest{
							Uuid: testresources.Machine2().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceDeleteResponse{
								Machine: testresources.Machine2(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Machine2(),
		},
		{
			Name:    "delete from file",
			CmdArgs: append([]string{"admin", "machine", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine2()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.MachineServiceDeleteRequest{
								Uuid: testresources.Machine2().Uuid,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.MachineServiceDeleteResponse{
									Machine: testresources.Machine2(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            ID                                        LAST EVENT   WHEN  AGE  HOSTNAME   PROJECT                               SIZE           IMAGE         PARTITION    RACK
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  🛡  Phoned Home  1m    1m   machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  v1-medium-x86  Ubuntu 24.04  partition-2  rack-1
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_MachineCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv2.MachineServiceUpdateResponse, *apiv2.Machine]{
		{
			Name: "update",
			CmdArgs: []string{"admin", "machine", "update",
				"--project", testresources.Project2().Uuid,
				testresources.Machine2().Uuid,
				"--description", "42",
				"--labels", "1=2",
				"--ssh-public-key", "@.ssh/id_rsa.pub",
			},
			AssertExhaustiveArgs:     true,
			AssertExhaustiveExcludes: append(e2e.CommonExcludedFileArgs(), "add-labels", "remove-labels"),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						genericcli.Must(fs.WriteFile(".ssh/id_rsa.pub", []byte("12345"), os.ModeAppend))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.MachineServiceUpdateRequest{
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
								Uuid:        testresources.Machine2().Uuid,
								Project:     testresources.Project2().Uuid,
								Description: new("42"),
								Labels: &apiv2.UpdateLabels{
									Strategy: &apiv2.UpdateLabels_Replace{
										Replace: &apiv2.Labels{
											Labels: map[string]string{
												"1": "2",
											},
										},
									},
								},
								SshPublicKeys: []string{"12345"},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.MachineServiceUpdateResponse{
									Machine: testresources.Machine2(),
								})
							},
						},
					},
				},
			),
			WantProtoObject: testresources.Machine2(),
		},
		{
			Name:    "update from file",
			CmdArgs: append([]string{"admin", "machine", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine2()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.MachineServiceUpdateRequest{
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
								Uuid:        testresources.Machine2().Uuid,
								Project:     testresources.Machine2().Allocation.Project,
								Description: &testresources.Machine2().Allocation.Description,
								Labels: &apiv2.UpdateLabels{
									Strategy: &apiv2.UpdateLabels_Replace{
										Replace: testresources.Machine2().Meta.Labels,
									},
								},
								SshPublicKeys: []string{},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.MachineServiceUpdateResponse{
									Machine: testresources.Machine2(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            ID                                        LAST EVENT   WHEN  AGE  HOSTNAME   PROJECT                               SIZE           IMAGE         PARTITION    RACK
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  🛡  Phoned Home  1m    1m   machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  v1-medium-x86  Ubuntu 24.04  partition-2  rack-1
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_MachineCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv2.MachineServiceUpdateResponse, *apiv2.Machine]{
		{
			Name:    "apply",
			CmdArgs: append([]string{"admin", "machine", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine2()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.MachineServiceCreateRequest{
								Uuid:        &testresources.Machine2().Uuid,
								Project:     testresources.Machine2().Allocation.Project,
								Name:        testresources.Machine2().Allocation.Name,
								Description: &testresources.Machine2().Allocation.Description,
								Hostname:    &testresources.Machine2().Allocation.Hostname,
								Image:       testresources.Machine2().Allocation.Image.Id,
								Userdata:    nil,
								Labels:      testresources.Machine2().Meta.Labels,
								Networks: func() []*apiv2.MachineAllocationNetwork {
									var nws []*apiv2.MachineAllocationNetwork

									for _, nw := range testresources.Machine2().Allocation.Networks {
										nws = append(nws, &apiv2.MachineAllocationNetwork{
											Network: nw.Network,
											Ips:     nw.Ips,
										})
									}

									return nws
								}(),
								AllocationType: apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_MACHINE,
								FirewallSpec:   nil,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.MachineServiceCreateResponse{
									Machine: testresources.Machine2(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            ID                                        LAST EVENT   WHEN  AGE  HOSTNAME   PROJECT                               SIZE           IMAGE         PARTITION    RACK
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  🛡  Phoned Home  1m    1m   machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  v1-medium-x86  Ubuntu 24.04  partition-2  rack-1
			`),
		},
		{
			Name:    "apply already exists",
			CmdArgs: append([]string{"admin", "machine", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine2()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.MachineServiceCreateRequest{
								Uuid:        &testresources.Machine2().Uuid,
								Project:     testresources.Machine2().Allocation.Project,
								Name:        testresources.Machine2().Allocation.Name,
								Description: &testresources.Machine2().Allocation.Description,
								Hostname:    &testresources.Machine2().Allocation.Hostname,
								Image:       testresources.Machine2().Allocation.Image.Id,
								Userdata:    nil,
								Labels:      testresources.Machine2().Meta.Labels,
								Networks: func() []*apiv2.MachineAllocationNetwork {
									var nws []*apiv2.MachineAllocationNetwork

									for _, nw := range testresources.Machine2().Allocation.Networks {
										nws = append(nws, &apiv2.MachineAllocationNetwork{
											Network: nw.Network,
											Ips:     nw.Ips,
										})
									}

									return nws
								}(),
								AllocationType: apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_MACHINE,
								FirewallSpec:   nil,
							},
							WantError: connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")),
						},
						{
							WantRequest: &apiv2.MachineServiceUpdateRequest{
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
								Uuid:        testresources.Machine2().Uuid,
								Project:     testresources.Machine2().Allocation.Project,
								Description: &testresources.Machine2().Allocation.Description,
								Labels: &apiv2.UpdateLabels{
									Strategy: &apiv2.UpdateLabels_Replace{
										Replace: testresources.Machine2().Meta.Labels,
									},
								},
								SshPublicKeys: []string{},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.MachineServiceUpdateResponse{
									Machine: testresources.Machine2(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            ID                                        LAST EVENT   WHEN  AGE  HOSTNAME   PROJECT                               SIZE           IMAGE         PARTITION    RACK
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  🛡  Phoned Home  1m    1m   machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  v1-medium-x86  Ubuntu 24.04  partition-2  rack-1
            `),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_MachineCmd_BmcList(t *testing.T) {
	tests := []*e2e.Test[adminv2.MachineServiceListBMCResponse, apiv2.MachineBMCDetails]{
		{
			Name: "list",
			CmdArgs: []string{"admin", "machine", "bmc", "list",
				"--name", "name",
				"--hostname", "hostname",
				"--image", "image",
				"--partition", "partition",
				"--project", "project",
				"--size", "size",
				"--id", "uuid",
			},
			AssertExhaustiveArgs:     true,
			AssertExhaustiveExcludes: []string{"sort-by"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceListBMCRequest{
							Query: &apiv2.MachineQuery{
								Allocation: &apiv2.MachineAllocationQuery{
									Hostname: new("hostname"),
									Name:     new("name"),
									Image:    new("image"),
									Project:  new("project"),
								},
								Partition: new("partition"),
								Size:      new("size"),
								Uuid:      new("uuid"),
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceListBMCResponse{
								BmcDetails: []*apiv2.MachineBMCDetails{
									testresources.Machine2BmcDetails,
									testresources.Machine1BmcDetails,
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                        POWER   IP            MAC                BOARD PART NUMBER  BIOS   BMC    SIZE           PARTITION    RACK    UPDATED
            5fa2bbe1-407c-4142-92d5-e4419daf9646   🟒  ⏾ 120W  10.0.0.1:623  02:00:00:00:00:01  Board-PN-1         1.5.6  3.1.1  v1-medium-x86  partition-1  rack-1  1m ago
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd     ⏾ 0W    10.0.0.2:623  02:00:00:00:00:02  Board-PN-2         2.0.0  2.4.0  v1-medium-x86  partition-2  rack-1  2m ago
            `),
			WantWideTable: new(`
            ID                                     LED      POWER                               IP            MAC                BOARD PART NUMBER  CHASSIS SERIAL  PRODUCT SERIAL  BIOS VERSION  BMC VERSION  SIZE           PARTITION    RACK    UPDATED
            5fa2bbe1-407c-4142-92d5-e4419daf9646   LED-ON   on On On 120W                       10.0.0.1:623  02:00:00:00:00:01  Board-PN-1         Chassis-SN-1    Product-SN-1    1.5.6         3.1.1        v1-medium-x86  partition-1  rack-1  1m ago
            673fc473-63ca-4ea4-b9dd-b45cb2127a6fd  LED-OFF  off Power Supply Warning Absent 0W  10.0.0.2:623  02:00:00:00:00:02  Board-PN-2         Chassis-SN-2    Product-SN-2    2.0.0         2.4.0        v1-medium-x86  partition-2  rack-1  2m ago
			`),
			Template: new("{{ .uuid }} {{ .bmc_report.bmc.address }}"),
			WantTemplate: new(`
5fa2bbe1-407c-4142-92d5-e4419daf9646 10.0.0.1:623
673fc473-63ca-4ea4-b9dd-b45cb2127a6fd 10.0.0.2:623
			`),
			WantMarkdown: new(`
            | ID                                    |   | POWER  | IP           | MAC               | BOARD PART NUMBER | BIOS  | BMC   | SIZE          | PARTITION   | RACK   | UPDATED |
            |---------------------------------------|---|--------|--------------|-------------------|-------------------|-------|-------|---------------|-------------|--------|---------|
            | 5fa2bbe1-407c-4142-92d5-e4419daf9646  | 🟒 | ⏾ 120W | 10.0.0.1:623 | 02:00:00:00:00:01 | Board-PN-1        | 1.5.6 | 3.1.1 | v1-medium-x86 | partition-1 | rack-1 | 1m ago  |
            | 673fc473-63ca-4ea4-b9dd-b45cb2127a6fd |   | ⏾ 0W   | 10.0.0.2:623 | 02:00:00:00:00:02 | Board-PN-2        | 2.0.0 | 2.4.0 | v1-medium-x86 | partition-2 | rack-1 | 2m ago  |
            `),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_MachineCmd_BmcDescribe(t *testing.T) {
	tests := []*e2e.Test[adminv2.MachineServiceListBMCResponse, apiv2.MachineBMCDetails]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "machine", "bmc", "describe", testresources.Machine1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceGetBMCRequest{
							Uuid: testresources.Machine1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceGetBMCResponse{
								BmcDetails: testresources.Machine1BmcDetails,
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                       POWER   IP            MAC                BOARD PART NUMBER  BIOS   BMC    SIZE           PARTITION    RACK    UPDATED
            5fa2bbe1-407c-4142-92d5-e4419daf9646  🟒  ⏾ 120W  10.0.0.1:623  02:00:00:00:00:01  Board-PN-1         1.5.6  3.1.1  v1-medium-x86  partition-1  rack-1  1m ago
            `),
			WantWideTable: new(`
            ID                                    LED     POWER          IP            MAC                BOARD PART NUMBER  CHASSIS SERIAL  PRODUCT SERIAL  BIOS VERSION  BMC VERSION  SIZE           PARTITION    RACK    UPDATED
            5fa2bbe1-407c-4142-92d5-e4419daf9646  LED-ON  on On On 120W  10.0.0.1:623  02:00:00:00:00:01  Board-PN-1         Chassis-SN-1    Product-SN-1    1.5.6         3.1.1        v1-medium-x86  partition-1  rack-1  1m ago
			`),
			Template: new("{{ .uuid }} {{ .bmc_report.bmc.address }}"),
			WantTemplate: new(`
            5fa2bbe1-407c-4142-92d5-e4419daf9646 10.0.0.1:623
			`),
			WantMarkdown: new(`
            | ID                                   |   | POWER  | IP           | MAC               | BOARD PART NUMBER | BIOS  | BMC   | SIZE          | PARTITION   | RACK   | UPDATED |
            |--------------------------------------|---|--------|--------------|-------------------|-------------------|-------|-------|---------------|-------------|--------|---------|
            | 5fa2bbe1-407c-4142-92d5-e4419daf9646 | 🟒 | ⏾ 120W | 10.0.0.1:623 | 02:00:00:00:00:01 | Board-PN-1        | 1.5.6 | 3.1.1 | v1-medium-x86 | partition-1 | rack-1 | 1m ago  |
            `),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_MachineCmd_BmcCommand(t *testing.T) {
	tests := []*e2e.Test[adminv2.MachineServiceListBMCResponse, apiv2.MachineBMCDetails]{
		{
			Name:    "bmc command",
			CmdArgs: []string{"admin", "machine", "bmc", "command", testresources.Machine1().Uuid, "--command", "MACHINE_BMC_COMMAND_ON"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceBMCCommandRequest{
							Uuid:    testresources.Machine1().Uuid,
							Command: apiv2.MachineBMCCommand_MACHINE_BMC_COMMAND_ON,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceBMCCommandResponse{})
						},
					},
				},
			}),
			WantDefault: new(``),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
