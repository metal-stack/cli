package api_e2e

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
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
			Name:    "list",
			CmdArgs: []string{"machine", "list", "--project", testresources.Project2().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.MachineServiceListRequest{
							Project: testresources.Project2().Uuid,
							Query:   &apiv2.MachineQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.MachineServiceListResponse{
								Machines: []*apiv2.Machine{
									testresources.Machine2(),
								},
							})
						},
					},
				},
			}),
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

func Test_MachineCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.MachineServiceGetResponse, *apiv2.Machine]{
		{
			Name:    "describe",
			CmdArgs: []string{"machine", "describe", "--project", testresources.Project2().Uuid, testresources.Machine2().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.MachineServiceGetRequest{
							Uuid:    testresources.Machine2().Uuid,
							Project: testresources.Project2().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.MachineServiceGetResponse{
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
			CmdArgs: []string{"machine", "create",
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
				"--placement-tags", "cluster-id=cluster-uuid",
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
							PlacementTags:  []string{"cluster-id=cluster-uuid"},
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
			CmdArgs: append([]string{"machine", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine2()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.MachineServiceCreateRequest{
								Project:     testresources.Machine2().Allocation.Project,
								Name:        testresources.Machine2().Allocation.Name,
								Description: &testresources.Machine2().Allocation.Description,
								Hostname:    &testresources.Machine2().Allocation.Hostname,
								Partition:   &testresources.Machine2().Partition.Id,
								Size:        &testresources.Machine2().Size.Id,
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
			CmdArgs: []string{"machine", "delete", "--project", testresources.Project2().Uuid, testresources.Machine2().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.MachineServiceDeleteRequest{
							Project: testresources.Project2().Uuid,
							Uuid:    testresources.Machine2().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.MachineServiceDeleteResponse{
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
			CmdArgs: append([]string{"machine", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine2()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.MachineServiceDeleteRequest{
								Uuid:    testresources.Machine2().Uuid,
								Project: testresources.Machine2().Allocation.Project,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.MachineServiceDeleteResponse{
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
			CmdArgs: []string{"machine", "update",
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
			CmdArgs: append([]string{"machine", "update"}, e2e.AppendFromFileCommonArgs()...),
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
			CmdArgs: append([]string{"machine", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine2()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.MachineServiceCreateRequest{
								Project:     testresources.Machine2().Allocation.Project,
								Name:        testresources.Machine2().Allocation.Name,
								Description: &testresources.Machine2().Allocation.Description,
								Hostname:    &testresources.Machine2().Allocation.Hostname,
								Partition:   &testresources.Machine2().Partition.Id,
								Size:        &testresources.Machine2().Size.Id,
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
			CmdArgs: append([]string{"machine", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Machine2()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.MachineServiceCreateRequest{
								Project:     testresources.Machine2().Allocation.Project,
								Name:        testresources.Machine2().Allocation.Name,
								Description: &testresources.Machine2().Allocation.Description,
								Hostname:    &testresources.Machine2().Allocation.Hostname,
								Partition:   &testresources.Machine2().Partition.Id,
								Size:        &testresources.Machine2().Size.Id,
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
