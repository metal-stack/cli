package v2_test

import (
	"testing"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	apitests "github.com/metal-stack/api/go/tests"
	"github.com/metal-stack/cli/pkg/test"
	"github.com/stretchr/testify/mock"
)

// Generated with AI

var (
	testPartition1 = &apiv2.Partition{
		Id:          "1",
		Description: "partition 1",
		MgmtServiceAddresses: []string{
			"192.168.1.1:1234",
		},
		BootConfiguration: &apiv2.PartitionBootConfiguration{
			Commandline: "commandline",
			ImageUrl:    "imageurl",
			KernelUrl:   "kernelurl",
		},
		Meta: &apiv2.Meta{
			Labels: &apiv2.Labels{
				Labels: map[string]string{
					"a": "b",
				},
			},
		},
	}
	testPartition2 = &apiv2.Partition{
		Id:          "2",
		Description: "partition 2",
		MgmtServiceAddresses: []string{
			"192.168.1.2:1234",
		},
		BootConfiguration: &apiv2.PartitionBootConfiguration{
			Commandline: "commandline",
			ImageUrl:    "imageurl",
			KernelUrl:   "kernelurl",
		},
		Meta: &apiv2.Meta{
			Labels: &apiv2.Labels{
				Labels: nil,
			},
		},
	}
)

func Test_AdminPartitionCmd_List(t *testing.T) {
	tests := []*test.Test[[]*apiv2.Partition]{
		{
			Name: "list",
			Cmd: func(want []*apiv2.Partition) []string {
				return []string{"admin", "partition", "list"}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv2Mocks: &apitests.Apiv2MockFns{
					Partition: func(m *mock.Mock) {
						m.On("List", mock.Anything, mock.Anything).Return(&apiv2.PartitionServiceListResponse{
							Partitions: []*apiv2.Partition{
								testPartition1,
								testPartition2,
							},
						}, nil)
					},
				},
			},
			Want: []*apiv2.Partition{
				testPartition1,
				testPartition2,
			},
			WantTable: new(`
ID  DESCRIPTION
1   partition 1
2   partition 2
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminPartitionCmd_Describe(t *testing.T) {
	tests := []*test.Test[*apiv2.Partition]{
		{
			Name: "describe",
			Cmd: func(want *apiv2.Partition) []string {
				return []string{"admin", "partition", "describe", want.Id}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv2Mocks: &apitests.Apiv2MockFns{
					Partition: func(m *mock.Mock) {
						m.On("Get", mock.Anything, mock.Anything).Return(&apiv2.PartitionServiceGetResponse{
							Partition: testPartition1,
						}, nil)
					},
				},
			},
			Want: testPartition1,
			WantTable: new(`
ID  DESCRIPTION
1   partition 1
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminPartitionCmd_Capacity(t *testing.T) {
	tests := []*test.Test[[]*adminv2.PartitionCapacity]{
		{
			Name: "capacity",
			Cmd: func(want []*adminv2.PartitionCapacity) []string {
				return []string{"admin", "partition", "capacity"}
			},
			ClientMocks: &apitests.ClientMockFns{
				Adminv2Mocks: &apitests.Adminv2MockFns{
					Partition: func(m *mock.Mock) {
						m.On("Capacity", mock.Anything, mock.Anything).Return(&adminv2.PartitionServiceCapacityResponse{
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
						}, nil)
					},
				},
			},
			Want: []*adminv2.PartitionCapacity{
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
			WantTable: new(`
PARTITION    SIZE    ALLOCATED  FREE  UNAVAILABLE  RESERVATIONS  |  TOTAL  |  FAULTY  
partition-1  size-1  1          3     0            2 (1/3 used)  |  5      |  2       
Total                1          3     0            2             |  5      |  2
`),
		},
		{
			Name: "capacity with filters",
			Cmd: func(want []*adminv2.PartitionCapacity) []string {
				return []string{"admin", "partition", "capacity", "--id", "partition-1", "--size", "size-1", "--project", "project-123", "--sort-by", "id"}
			},
			ClientMocks: &apitests.ClientMockFns{
				Adminv2Mocks: &apitests.Adminv2MockFns{
					Partition: func(m *mock.Mock) {
						m.On("Capacity", mock.Anything, mock.Anything).Return(&adminv2.PartitionServiceCapacityResponse{
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
						}, nil)
					},
				},
			},
			Want: []*adminv2.PartitionCapacity{
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

func Test_AdminPartitionCmd_ExhaustiveArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "list",
			args: []string{"admin", "partition", "list"},
		},
		{
			name: "describe",
			args: []string{"admin", "partition", "describe", "1"},
		},
		{
			name: "capacity",
			args: []string{"admin", "partition", "capacity", "--id", "partition-1", "--size", "size-1", "--project", "project-123", "--sort-by", "id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.AssertExhaustiveArgs(t, tt.args)
		})
	}
}
