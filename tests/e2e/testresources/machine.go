package testresources

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Machine1 = func() *apiv2.Machine {
		return &apiv2.Machine{
			Uuid: "5fa2bbe1-407c-4142-92d5-e4419daf9646",
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"a": "b",
					},
				},
			},
			Partition: Partition1(),
			Rack:      *Switch1().Rack,
			Room:      *Switch1().Room,
			Size:      Size1(),
			Hardware: &apiv2.MachineHardware{
				Memory: 1024,
				Disks: []*apiv2.MachineBlockDevice{
					{
						Name: "/dev/sda",
						Size: 1,
					},
				},
				Cpus: []*apiv2.MetalCPU{
					{
						Vendor:  "Intel",
						Model:   "Xeon",
						Cores:   4,
						Threads: 4,
					},
				},
				Gpus: []*apiv2.MetalGPU{},
				Nics: []*apiv2.MachineNic{},
			},
			Allocation: nil,
			Status: &apiv2.MachineStatus{
				Condition: &apiv2.MachineCondition{
					State: apiv2.MachineState_MACHINE_STATE_AVAILABLE,
				},
				LedState: &apiv2.MachineChassisIdentifyLEDState{
					Value:       "",
					Description: "",
				},
				Liveliness:         apiv2.MachineLiveliness_MACHINE_LIVELINESS_ALIVE,
				MetalHammerVersion: "1",
			},
			RecentProvisioningEvents: &apiv2.MachineRecentProvisioningEvents{
				Events: []*apiv2.MachineProvisioningEvent{
					{
						Time:    timestamppb.New(e2e.TimeBubbleStartTime().Add(-1 * time.Minute)),
						Event:   apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_ALIVE,
						Message: "alive",
					},
				},
				LastEventTime: timestamppb.New(e2e.TimeBubbleStartTime().Add(-1 * time.Minute)),
				LastErrorEvent: &apiv2.MachineProvisioningEvent{
					Time:    timestamppb.New(e2e.TimeBubbleStartTime().Add(-1 * time.Hour)),
					Event:   apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_WAITING,
					Message: "waiting",
				},
				State: apiv2.MachineProvisioningEventState_MACHINE_PROVISIONING_EVENT_STATE_UNSPECIFIED,
			},
		}
	}
	Machine2 = func() *apiv2.Machine {
		return &apiv2.Machine{
			Uuid: "673fc473-63ca-4ea4-b9dd-b45cb2127a6fd",
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"c": "d",
					},
				},
			},
			Partition: Partition2(),
			Rack:      *Switch2().Rack,
			Room:      *Switch2().Room,
			Size:      Size1(),
			Hardware: &apiv2.MachineHardware{
				Memory: 1024,
				Disks: []*apiv2.MachineBlockDevice{
					{
						Name: "/dev/sda",
						Size: 1,
					},
				},
				Cpus: []*apiv2.MetalCPU{
					{
						Vendor:  "Intel",
						Model:   "Xeon",
						Cores:   4,
						Threads: 4,
					},
				},
				Gpus: []*apiv2.MetalGPU{},
				Nics: []*apiv2.MachineNic{},
			},
			Allocation: &apiv2.MachineAllocation{
				Uuid: "4f94e87b-b08f-4f82-b053-9b8305de60ad",
				Meta: &apiv2.Meta{
					CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime().Add(-1 * time.Minute)),
					Labels: &apiv2.Labels{
						Labels: map[string]string{
							"e": "f",
						},
					},
				},
				Name:             "machine-2",
				Description:      "machine 2",
				CreatedBy:        "foo",
				Project:          Project2().Uuid,
				Image:            Image1(),
				FilesystemLayout: &apiv2.FilesystemLayout{},
				Networks: []*apiv2.MachineNetwork{
					{
						Network: Network1().Id,
						Ips:     []string{"4.5.6.7"},
					},
					{
						Network: Network2().Id,
						Ips:     []string{"192.1.1.1"},
					},
				},
				Hostname:       "machine-2",
				SshPublicKeys:  []string{},
				Userdata:       "",
				AllocationType: apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_MACHINE,
				FirewallRules:  nil,
				DnsServers:     []*apiv2.DNSServer{},
				NtpServers:     []*apiv2.NTPServer{},
				Vpn: &apiv2.MachineVPN{
					ControlPlaneAddress: "1.2.3.4",
					AuthKey:             "abc",
					Connected:           true,
					Ips:                 []string{"3.4.5.6"},
				},
			},
			Status: &apiv2.MachineStatus{
				Condition: &apiv2.MachineCondition{
					State: apiv2.MachineState_MACHINE_STATE_AVAILABLE,
				},
				LedState: &apiv2.MachineChassisIdentifyLEDState{
					Value:       "",
					Description: "",
				},
				Liveliness:         apiv2.MachineLiveliness_MACHINE_LIVELINESS_ALIVE,
				MetalHammerVersion: "1",
			},
			RecentProvisioningEvents: &apiv2.MachineRecentProvisioningEvents{
				Events: []*apiv2.MachineProvisioningEvent{
					{
						Time:    timestamppb.New(e2e.TimeBubbleStartTime().Add(-1 * time.Minute)),
						Event:   apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_PHONED_HOME,
						Message: "phoned home",
					},
					{
						Time:    timestamppb.New(e2e.TimeBubbleStartTime().Add(-2 * time.Minute)),
						Event:   apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_ALIVE,
						Message: "alive",
					},
				},
				LastEventTime: timestamppb.New(e2e.TimeBubbleStartTime().Add(-1 * time.Minute)),
				LastErrorEvent: &apiv2.MachineProvisioningEvent{
					Time:    timestamppb.New(e2e.TimeBubbleStartTime().Add(-1 * time.Hour)),
					Event:   apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_WAITING,
					Message: "waiting",
				},
				State: apiv2.MachineProvisioningEventState_MACHINE_PROVISIONING_EVENT_STATE_UNSPECIFIED,
			},
		}
	}

	Machine1BmcDetails = &apiv2.MachineBMCDetails{
		Uuid:      Machine1().Uuid,
		Partition: Machine1().Partition.Id,
		Rack:      Machine1().Rack,
		Room:      Machine1().Room,
		Size:      Machine1().Size.Id,
		BmcReport: &apiv2.MachineBMCReport{
			Uuid: Machine1().Uuid,
			Bmc: &apiv2.MachineBMC{
				Address:    "10.0.0.1:623",
				Mac:        "02:00:00:00:00:01",
				User:       "admin",
				Password:   "password",
				Interface:  "lanplus",
				Version:    "3.1.1",
				PowerState: "on",
			},
			Bios: &apiv2.MachineBios{
				Version: "1.5.6",
				Vendor:  "American Megatrends Inc.",
				Date:    "04/01/2024",
			},
			Fru: &apiv2.MachineFRU{
				ChassisPartNumber:   new("Chassis-PN-1"),
				ChassisPartSerial:   new("Chassis-SN-1"),
				BoardMfg:            new("Supermicro"),
				BoardMfgSerial:      new("Board-SN-1"),
				BoardPartNumber:     new("Board-PN-1"),
				ProductManufacturer: new("Dell"),
				ProductPartNumber:   new("Product-PN-1"),
				ProductSerial:       new("Product-SN-1"),
			},
			PowerMetric: &apiv2.MachinePowerMetric{
				AverageConsumedWatts: 120.5,
				IntervalInMin:        5,
				MaxConsumedWatts:     180.2,
				MinConsumedWatts:     95.3,
			},
			PowerSupplies: []*apiv2.MachinePowerSupply{
				{
					Health: "OK",
					State:  "On",
				},
				{
					Health: "OK",
					State:  "On",
				},
			},
			LedState: &apiv2.MachineChassisIdentifyLEDState{
				Value:       "LED-ON",
				Description: "User initiated chassis identify",
			},
			UpdatedAt: timestamppb.New(e2e.TimeBubbleStartTime().Add(-1 * time.Minute)),
		},
	}
	Machine2BmcDetails = &apiv2.MachineBMCDetails{
		Uuid:      Machine2().Uuid,
		Partition: Machine2().Partition.Id,
		Rack:      Machine2().Rack,
		Room:      Machine2().Room,
		Size:      Machine2().Size.Id,
		BmcReport: &apiv2.MachineBMCReport{
			Uuid: Machine2().Uuid,
			Bmc: &apiv2.MachineBMC{
				Address:    "10.0.0.2:623",
				Mac:        "02:00:00:00:00:02",
				User:       "root",
				Password:   "changeme",
				Interface:  "lanplus",
				Version:    "2.4.0",
				PowerState: "off",
			},
			Bios: &apiv2.MachineBios{
				Version: "2.0.0",
				Vendor:  "Supermicro",
				Date:    "10/15/2023",
			},
			Fru: &apiv2.MachineFRU{
				ChassisPartNumber:   new("Chassis-PN-2"),
				ChassisPartSerial:   new("Chassis-SN-2"),
				BoardMfg:            new("Supermicro"),
				BoardMfgSerial:      new("Board-SN-2"),
				BoardPartNumber:     new("Board-PN-2"),
				ProductManufacturer: new("HPE"),
				ProductPartNumber:   new("Product-PN-2"),
				ProductSerial:       new("Product-SN-2"),
			},
			PowerMetric: &apiv2.MachinePowerMetric{
				AverageConsumedWatts: 0,
				IntervalInMin:        0,
				MaxConsumedWatts:     0,
				MinConsumedWatts:     0,
			},
			PowerSupplies: []*apiv2.MachinePowerSupply{
				{
					Health: "Warning",
					State:  "Absent",
				},
			},
			LedState: &apiv2.MachineChassisIdentifyLEDState{
				Value:       "LED-OFF",
				Description: "",
			},
			UpdatedAt: timestamppb.New(e2e.TimeBubbleStartTime().Add(-2 * time.Minute)),
		},
	}
)
