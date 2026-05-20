package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
)

var (
	machine1 = func() *apiv2.Machine {
		return &apiv2.Machine{
			Uuid: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"test": "value",
					},
				},
			},
			Partition: &apiv2.Partition{
				Id: "test-partition",
			},
			Rack:       "rack-1",
			Size:       &apiv2.Size{Id: "test-size"},
			Hardware:   &apiv2.MachineHardware{},
			Allocation: testMachineAllocation(),
			Status: &apiv2.MachineStatus{
				Condition: &apiv2.MachineCondition{
					State:       apiv2.MachineState_MACHINE_STATE_AVAILABLE,
					Description: "available",
				},
				Liveliness: apiv2.MachineLiveliness_MACHINE_LIVELINESS_ALIVE,
			},
			RecentProvisioningEvents: &apiv2.MachineRecentProvisioningEvents{},
		}
	}
	machine2 = func() *apiv2.Machine {
		return &apiv2.Machine{
			Uuid: "b2c3d4e5-f6a7-8901-bcde-f12345678901",
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"test": "another",
					},
				},
			},
			Partition: &apiv2.Partition{
				Id: "test-partition",
			},
			Rack:       "rack-2",
			Size:       &apiv2.Size{Id: "test-size"},
			Hardware:   &apiv2.MachineHardware{},
			Allocation: testMachineAllocation2(),
			Status: &apiv2.MachineStatus{
				Condition: &apiv2.MachineCondition{
					State:       apiv2.MachineState_MACHINE_STATE_AVAILABLE,
					Description: "available",
				},
				Liveliness: apiv2.MachineLiveliness_MACHINE_LIVELINESS_DEAD,
			},
			RecentProvisioningEvents: &apiv2.MachineRecentProvisioningEvents{},
		}
	}
)

func testMachineAllocation() *apiv2.MachineAllocation {
	return &apiv2.MachineAllocation{
		Uuid:           "alloc-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		Name:           "machine-1",
		Project:        "test-project",
		Hostname:       "machine-1.test.local",
		Image:          &apiv2.Image{Name: new("ubuntu-24.04")},
		Meta:           &apiv2.Meta{CreatedAt: nil},
		AllocationType: apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_MACHINE,
	}
}

func testMachineAllocation2() *apiv2.MachineAllocation {
	return &apiv2.MachineAllocation{
		Uuid:           "alloc-b2c3d4e5-f6a7-8901-bcde-f12345678901",
		Name:           "machine-2",
		Project:        "test-project",
		Hostname:       "machine-2.test.local",
		Image:          &apiv2.Image{Name: new("ubuntu-24.04")},
		Meta:           &apiv2.Meta{CreatedAt: nil},
		AllocationType: apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_MACHINE,
	}
}

func Test_AdminMachineCmd_Get(t *testing.T) {
	tests := []*e2e.Test[adminv2.MachineServiceGetResponse, *apiv2.Machine]{
		{
			Name:    "get",
			CmdArgs: []string{"admin", "machine", "get", machine1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceGetRequest{
							Uuid: machine1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceGetResponse{
								Machine: machine1(),
							})
						},
					},
				},
			}),
			WantObject:      machine1(),
			WantProtoObject: machine1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminMachineCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.MachineServiceListResponse, apiv2.Machine]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "machine", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.MachineServiceListRequest{
							Query: &apiv2.MachineQuery{
								Labels: &apiv2.Labels{
									Labels: map[string]string{},
								},
								Allocation: &apiv2.MachineAllocationQuery{
									Labels: &apiv2.Labels{
										Labels: map[string]string{},
									},
								},
								Network: &apiv2.MachineNetworkQuery{},
								Nic:     &apiv2.MachineNicQuery{},
								Disk: &apiv2.MachineDiskQuery{
									Names: []string{},
								},
								Bmc:      &apiv2.MachineBMCQuery{},
								Fru:      &apiv2.MachineFRUQuery{},
								Hardware: &apiv2.MachineHardwareQuery{},
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.MachineServiceListResponse{
								Machines: []*apiv2.Machine{
									machine1(),
									machine2(),
								},
							})
						},
					},
				},
			}),
			Template: new("{{ .uuid }} {{ .allocation.name }}"),
			WantTemplate: new(`
a1b2c3d4-e5f6-7890-abcd-ef1234567890 machine-1
b2c3d4e5-f6a7-8901-bcde-f12345678901 machine-2
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
