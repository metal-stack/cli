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
)

func Test_AdminPartitionCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.PartitionServiceListResponse, []*apiv2.Partition]{
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
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.Partition{testresources.Partition1()},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminPartitionCmd_Capacity(t *testing.T) {
	tests := []*e2e.Test[adminv2.PartitionServiceCapacityResponse, []*adminv2.PartitionCapacity]{
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
										Partition: "fra-equ01",
										MachineSizeCapacities: []*adminv2.MachineSizeCapacity{
											{
												Size:      "v1-medium-x86",
												Total:     100,
												Free:      50,
												Allocated: 30,
												Faulty:    5,
												Other:     15,
											},
										},
									},
								},
							})
						},
					},
				},
			}),
			WantObject: []*adminv2.PartitionCapacity{
				{
					Partition: "fra-equ01",
					MachineSizeCapacities: []*adminv2.MachineSizeCapacity{
						{
							Size:      "v1-medium-x86",
							Total:     100,
							Free:      50,
							Allocated: 30,
							Faulty:    5,
							Other:     15,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminPartitionCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.PartitionServiceDeleteResponse, *apiv2.Partition]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "partition", "delete", testresources.Partition1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.PartitionServiceDeleteRequest{
							Id: testresources.Partition1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.PartitionServiceDeleteResponse{
								Partition: testresources.Partition1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.Partition1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
