package admin_e2e

import (
	"fmt"
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

func Test_AdminPartitionCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.PartitionServiceListResponse, apiv2.Partition]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "partition", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.PartitionServiceListRequest{
							Query: &apiv2.PartitionQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.PartitionServiceListResponse{
								Partitions: []*apiv2.Partition{
									testresources.Partition1(),
									testresources.Partition2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID                                     DESCRIPTION
			5711axzf0-4218-p7x9-3f8f-d34912dx190q  partition 1
			d34912dx1-90qe-3jf8-1kd0-5711axzf08sj  partition 2
			`),
			WantWideTable: new(`
			ID                                     DESCRIPTION  LABELS
			5711axzf0-4218-p7x9-3f8f-d34912dx190q  partition 1  a=b
			d34912dx1-90qe-3jf8-1kd0-5711axzf08sj  partition 2
			`),
			Template: new("{{ .id }} {{ .description }}"),
			WantTemplate: new(`
5711axzf0-4218-p7x9-3f8f-d34912dx190q partition 1
d34912dx1-90qe-3jf8-1kd0-5711axzf08sj partition 2
			`),
			WantMarkdown: new(`
			| ID                                    | DESCRIPTION |
			|---------------------------------------|-------------|
			| 5711axzf0-4218-p7x9-3f8f-d34912dx190q | partition 1 |
			| d34912dx1-90qe-3jf8-1kd0-5711axzf08sj | partition 2 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminPartitionCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.PartitionServiceGetResponse, *apiv2.Partition]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "partition", "describe", testresources.Partition1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.PartitionServiceGetRequest{
							Id: testresources.Partition1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.PartitionServiceGetResponse{
								Partition: testresources.Partition1(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Partition1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminPartitionCmd_Capacity(t *testing.T) {
	tests := []*e2e.Test[adminv2.PartitionServiceCapacityResponse, adminv2.PartitionCapacity]{
		{
			Name:    "capacity",
			CmdArgs: []string{"admin", "partition", "capacity"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.PartitionServiceCapacityRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.PartitionServiceCapacityResponse{
								PartitionCapacity: []*adminv2.PartitionCapacity{
									{
										Partition: "5711axzf0-4218-p7x9-3f8f-d34912dx190q",
										MachineSizeCapacities: []*adminv2.MachineSizeCapacity{
											{
												Size:             "size-1",
												Free:             3,
												Allocated:        1,
												Total:            5,
												Faulty:           2,
												Reservations:     3,
												UsedReservations: 1,
											},
										},
									},
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            PARTITION                              SIZE    ALLOCATED  FREE  UNAVAILABLE  RESERVATIONS  |  TOTAL  |  FAULTY  
            5711axzf0-4218-p7x9-3f8f-d34912dx190q  size-1  1          3     0            2 (1/3 used)  |  5      |  2       
            Total                                          1          3     0            2             |  5      |  2      
			`),
			WantWideTable: new(`
            PARTITION                              SIZE    ALLOCATED  FREE  UNAVAILABLE  RESERVATIONS  |  TOTAL  |  FAULTY  PHONED HOME  WAITING  OTHER  
            5711axzf0-4218-p7x9-3f8f-d34912dx190q  size-1  1          3     0            2 (1/3 used)  |  5      |  2       0            0        0      
            Total                                          1          3     0            2             |  5      |  2       0            0        0
			`),
			Template: new("{{ .partition }} {{ (index .machine_size_capacities 0).size }}"),
			WantTemplate: new(`
5711axzf0-4218-p7x9-3f8f-d34912dx190q size-1
			`),
			WantMarkdown: new(`
			| PARTITION                             | SIZE   | ALLOCATED | FREE | UNAVAILABLE | RESERVATIONS | | | TOTAL | | | FAULTY |
			|---------------------------------------|--------|-----------|------|-------------|--------------|---|-------|---|--------|
			| 5711axzf0-4218-p7x9-3f8f-d34912dx190q | size-1 | 1         | 3    | 0           | 2 (1/3 used) | | | 5     | | | 2      |
			| Total                                 |        | 1         | 3    | 0           | 2            | | | 5     | | | 2      |
			`),
		},
		{
			Name:    "capacity with filters",
			CmdArgs: []string{"admin", "partition", "capacity", "--id", "5711axzf0-4218-p7x9-3f8f-d34912dx190q", "--size", "size-1", "--project", "0d81bca7-73f6-4da3-8397-4a8c52a0c583"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.PartitionServiceCapacityRequest{
							Id:      new("5711axzf0-4218-p7x9-3f8f-d34912dx190q"),
							Size:    new("size-1"),
							Project: new("0d81bca7-73f6-4da3-8397-4a8c52a0c583"),
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.PartitionServiceCapacityResponse{
								PartitionCapacity: []*adminv2.PartitionCapacity{
									{
										Partition: "5711axzf0-4218-p7x9-3f8f-d34912dx190q",
										MachineSizeCapacities: []*adminv2.MachineSizeCapacity{
											{
												Size:             "size-1",
												Free:             3,
												Allocated:        1,
												Total:            5,
												Faulty:           2,
												Reservations:     3,
												UsedReservations: 1,
											},
										},
									},
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            PARTITION                              SIZE    ALLOCATED  FREE  UNAVAILABLE  RESERVATIONS  |  TOTAL  |  FAULTY  
            5711axzf0-4218-p7x9-3f8f-d34912dx190q  size-1  1          3     0            2 (1/3 used)  |  5      |  2       
            Total                                          1          3     0            2             |  5      |  2
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminPartitionCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[adminv2.PartitionServiceUpdateResponse, *apiv2.Partition]{
		{
			Name:    "apply",
			CmdArgs: append([]string{"admin", "partition", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Partition1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.PartitionServiceCreateRequest{
								Partition: testresources.Partition1(),
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.PartitionServiceCreateResponse{
									Partition: testresources.Partition1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
			ID                                     DESCRIPTION
			5711axzf0-4218-p7x9-3f8f-d34912dx190q  partition 1
			`),
		},
		{
			Name:    "apply already exists",
			CmdArgs: append([]string{"admin", "partition", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Partition1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.PartitionServiceCreateRequest{
								Partition: testresources.Partition1(),
							},
							WantError: connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")),
						},
						{
							WantRequest: &adminv2.PartitionServiceUpdateRequest{
								Id:                   testresources.Partition1().Id,
								Description:          new(testresources.Partition1().Description),
								BootConfiguration:    testresources.Partition1().BootConfiguration,
								MgmtServiceAddresses: testresources.Partition1().MgmtServiceAddresses,
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
								Labels: &apiv2.UpdateLabels{
									Strategy: &apiv2.UpdateLabels_Replace{
										Replace: testresources.Partition1().Meta.Labels,
									},
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.PartitionServiceUpdateResponse{
									Partition: testresources.Partition1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
			ID                                     DESCRIPTION
			5711axzf0-4218-p7x9-3f8f-d34912dx190q  partition 1
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
