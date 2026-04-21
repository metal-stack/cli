package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
)

func Test_AdminPartitionCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.PartitionServiceListResponse, apiv2.Partition]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "partition", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
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
			ID           DESCRIPTION
			partition-1  partition 1
			partition-2  partition 2
			`),
			WantWideTable: new(`
			ID           DESCRIPTION  LABELS
			partition-1  partition 1  a=b
			partition-2  partition 2
			`),
			Template: new("{{ .id }} {{ .description }}"),
			WantTemplate: new(`
partition-1 partition 1
partition-2 partition 2
			`),
			WantMarkdown: new(`
			| ID          | DESCRIPTION |
			|-------------|-------------|
			| partition-1 | partition 1 |
			| partition-2 | partition 2 |
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
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
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
			WantObject:      testresources.Partition1(),
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
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.PartitionServiceCapacityRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.PartitionServiceCapacityResponse{
								PartitionCapacity: []*adminv2.PartitionCapacity{
									{
										Partition: "partition-1",
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
			PARTITION    SIZE    ALLOCATED  FREE  UNAVAILABLE  RESERVATIONS  |  TOTAL  |  FAULTY
			partition-1  size-1  1          3     0            2 (1/3 used)  |  5      |  2
			Total                1          3     0            2             |  5      |  2
			`),
			WantWideTable: new(`
			PARTITION    SIZE    ALLOCATED  FREE  UNAVAILABLE  RESERVATIONS  |  TOTAL  |  FAULTY  PHONED HOME  WAITING  OTHER
			partition-1  size-1  1          3     0            2 (1/3 used)  |  5      |  2       0            0        0
			Total                1          3     0            2             |  5      |  2       0            0        0
			`),
			Template: new("{{ .partition }} {{ (index .machine_size_capacities 0).size }}"),
			WantTemplate: new(`
partition-1 size-1
			`),
			WantMarkdown: new(`
			| PARTITION   | SIZE   | ALLOCATED | FREE | UNAVAILABLE | RESERVATIONS | | | TOTAL | | | FAULTY |
			|-------------|--------|-----------|------|-------------|--------------|---|-------|---|--------|
			| partition-1 | size-1 | 1         | 3    | 0           | 2 (1/3 used) | | | 5     | | | 2      |
			| Total       |        | 1         | 3    | 0           | 2            | | | 5     | | | 2      |
			`),
		},
		{
			Name:    "capacity with filters",
			CmdArgs: []string{"admin", "partition", "capacity", "--id", "partition-1", "--size", "size-1", "--project", "project-123"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.PartitionServiceCapacityRequest{
							Id:      new("partition-1"),
							Size:    new("size-1"),
							Project: new("project-123"),
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.PartitionServiceCapacityResponse{
								PartitionCapacity: []*adminv2.PartitionCapacity{
									{
										Partition: "partition-1",
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
			PARTITION    SIZE    ALLOCATED  FREE  UNAVAILABLE  RESERVATIONS  |  TOTAL  |  FAULTY
			partition-1  size-1  1          3     0            2 (1/3 used)  |  5      |  2
			Total                1          3     0            2             |  5      |  2
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
