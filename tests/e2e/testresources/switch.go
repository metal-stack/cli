package testresources

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Switch1 = func() *apiv2.Switch {
		return &apiv2.Switch{
			Id:             "leaf01",
			Partition:      "fra-equ01",
			Rack:           new("rack-1"),
			Description:    "leaf switch 1",
			ManagementIp:   "10.0.0.1",
			ManagementUser: new("admin"),
			MachineConnections: []*apiv2.MachineConnection{
				&apiv2.MachineConnection{
					MachineId: "id1",
					Nic:       Nic1(),
				},
			},
			Nics: []*apiv2.SwitchNic{Nic1(), Nic2()},
			Os: &apiv2.SwitchOS{
				Vendor:           apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_SONIC,
				Version:          "4.2.0",
				MetalCoreVersion: "v0.9.1 (abc1234), tags/v0.9.1",
			},
			LastSync: &apiv2.SwitchSync{
				Time:     timestamppb.New(e2e.TimeBubbleStartTime()),
				Duration: durationpb.New(100 * time.Millisecond),
			},
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		}
	}
	Switch2 = func() *apiv2.Switch {
		return &apiv2.Switch{
			Id:             "leaf02",
			Partition:      "fra-equ01",
			Rack:           new("rack-1"),
			Description:    "leaf switch 2",
			ManagementIp:   "10.0.0.2",
			ManagementUser: new("admin"),
			Os: &apiv2.SwitchOS{
				Vendor:           apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_SONIC,
				Version:          "4.2.0",
				MetalCoreVersion: "v0.9.1 (abc1234), tags/v0.9.1",
			},
			LastSync: &apiv2.SwitchSync{
				Time:     timestamppb.New(e2e.TimeBubbleStartTime()),
				Duration: durationpb.New(200 * time.Millisecond),
			},
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		}
	}

	Nic1 = func() *apiv2.SwitchNic {
		return &apiv2.SwitchNic{
			Name:       "Ethernet0",
			Identifier: "oid:0x1000000000001",
			Mac:        "52:54:00:ab:cd:01",
			Vrf:        new("default"),
			State: &apiv2.NicState{
				Desired: new(apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_UP),
				Actual:  apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_UP,
			},
			BgpFilter: &apiv2.BGPFilter{
				Cidrs: []string{
					"10.0.0.0/24",
					"192.168.100.0/24",
				},
				Vnis: []string{
					"10001",
					"10002",
				},
			},
			BgpPortState: &apiv2.SwitchBGPPortState{
				Neighbor:              "10.0.0.2",
				PeerGroup:             "TOR-LEAFS",
				VrfName:               "default",
				BgpState:              apiv2.BGPState_BGP_STATE_ESTABLISHED,
				BgpTimerUpEstablished: timestamppb.New(time.Now().Add(-2 * time.Hour)),
				SentPrefixCounter:     120,
				AcceptedPrefixCounter: 118,
			},
		}
	}
	Nic2 = func() *apiv2.SwitchNic {
		return &apiv2.SwitchNic{
			Name:       "Ethernet4",
			Identifier: "oid:0x1000000000002",
			Mac:        "52:54:00:ab:cd:02",
			State: &apiv2.NicState{
				Desired: new(apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_UP),
				Actual:  apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_DOWN,
			},
			BgpFilter: &apiv2.BGPFilter{
				Cidrs: []string{"10.1.0.0/24"},
				Vnis:  []string{"20001"},
			},
			BgpPortState: &apiv2.SwitchBGPPortState{
				Neighbor:              "10.1.0.2",
				PeerGroup:             "TOR-LEAFS",
				VrfName:               "default",
				BgpState:              apiv2.BGPState_BGP_STATE_IDLE,
				BgpTimerUpEstablished: nil,
				SentPrefixCounter:     0,
				AcceptedPrefixCounter: 0,
			},
		}
	}
	SwitchWithMachines1 = func() *apiv2.SwitchWithMachines {
		return &apiv2.SwitchWithMachines{
			Id:        Switch1().Id,
			Partition: Switch1().Partition,
			Rack:      *Switch1().Rack,
			Connections: []*apiv2.SwitchNicWithMachine{
				{
					Nic: Switch1().Nics[0],
					Machine: &apiv2.Machine{
						Uuid: "id1",
						Partition: &apiv2.Partition{
							Id: Switch1().Partition,
						},
						Rack: *Switch1().Rack,
						Size: &apiv2.Size{
							Id: "m1-small",
						},
					},
					Fru: &apiv2.MachineFRU{
						ProductSerial:     new("ps-1"),
						ChassisPartSerial: new("cs-1"),
					},
				},
			},
		}
	}
)
