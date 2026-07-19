package testresources

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Partition1 = func() *apiv2.Partition {
		return &apiv2.Partition{
			Id:          "fra-equ01",
			Description: "Frankfurt Equinix 1",
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"location": "fra",
					},
				},
			},
			MgmtServiceAddresses: []string{"10.0.0.1:8080", "10.0.0.2:8080"},
		}
	}
	Partition2 = func() *apiv2.Partition {
		return &apiv2.Partition{
			Id:          "fra-equ02",
			Description: "Frankfurt Equinix 2",
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
			MgmtServiceAddresses: []string{"10.0.1.1:8080"},
		}
	}
	Network1 = func() *apiv2.Network {
		return &apiv2.Network{
			Id:          "n-1",
			Name:        new("internal-net"),
			Description: new("Internal network"),
			Project:     new("project-1"),
			Partition:   new("fra-equ01"),
			Type:        apiv2.NetworkType_NETWORK_TYPE_CHILD,
			Prefixes:    []string{"10.0.0.0/16", "10.1.0.0/16"},
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
			NatType: apiv2.NATType_NAT_TYPE_NONE,
			Vrf:     new(uint32(100)),
		}
	}
	Network2 = func() *apiv2.Network {
		return &apiv2.Network{
			Id:          "n-2",
			Name:        new("external-net"),
			Description: new("External network"),
			Project:     new("project-2"),
			Partition:   new("fra-equ02"),
			Type:        apiv2.NetworkType_NETWORK_TYPE_EXTERNAL,
			Prefixes:    []string{"172.16.0.0/12"},
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
			NatType: apiv2.NATType_NAT_TYPE_IPV4_MASQUERADE,
		}
	}
	Machine1 = func() *apiv2.Machine {
		return &apiv2.Machine{
			Uuid:      "m-1",
			Partition: &apiv2.Partition{Id: "fra-equ01"},
			Rack:      "rack-1",
			Size:      &apiv2.Size{Id: "v1-medium-x86"},
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
			Allocation: &apiv2.MachineAllocation{
				Hostname: "node-1",
				Project:  "project-1",
			},
			Status: &apiv2.MachineStatus{
				Liveliness: apiv2.MachineLiveliness_MACHINE_LIVELINESS_ALIVE,
				Condition: &apiv2.MachineCondition{
					State: apiv2.MachineState_MACHINE_STATE_AVAILABLE,
				},
			},
		}
	}
	Machine2 = func() *apiv2.Machine {
		return &apiv2.Machine{
			Uuid:      "m-2",
			Partition: &apiv2.Partition{Id: "fra-equ02"},
			Rack:      "rack-2",
			Size:      &apiv2.Size{Id: "g1-medium-x86"},
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
			Allocation: &apiv2.MachineAllocation{
				Hostname: "node-2",
				Project:  "project-2",
			},
			Status: &apiv2.MachineStatus{
				Liveliness: apiv2.MachineLiveliness_MACHINE_LIVELINESS_ALIVE,
				Condition: &apiv2.MachineCondition{
					State: apiv2.MachineState_MACHINE_STATE_AVAILABLE,
				},
			},
		}
	}
	SizeReservation1 = func() *apiv2.SizeReservation {
		return &apiv2.SizeReservation{
			Id:          "sr-1",
			Name:        "reservation-1",
			Description: "Reservation for project-1",
			Project:     "project-1",
			Size:        "v1-medium-x86",
			Partitions:  []string{"fra-equ01"},
			Amount:      5,
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		}
	}
	SizeReservation2 = func() *apiv2.SizeReservation {
		return &apiv2.SizeReservation{
			Id:          "sr-2",
			Name:        "reservation-2",
			Description: "Reservation for project-2",
			Project:     "project-2",
			Size:        "g1-medium-x86",
			Partitions:  []string{"fra-equ02"},
			Amount:      3,
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		}
	}
	FilesystemLayout1 = func() *apiv2.FilesystemLayout {
		return &apiv2.FilesystemLayout{
			Id:          "fsl-1",
			Name:        new("default-ext4"),
			Description: new("Default ext4 layout"),
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
			Disks: []*apiv2.Disk{
				{
					Device: "/dev/sda",
					Partitions: []*apiv2.DiskPartition{
						{
							Number: 1,
							Size:   100000000000,
						},
					},
				},
			},
		}
	}
	FilesystemLayout2 = func() *apiv2.FilesystemLayout {
		return &apiv2.FilesystemLayout{
			Id:          "fsl-2",
			Name:        new("default-xfs"),
			Description: new("Default xfs layout"),
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		}
	}
)
