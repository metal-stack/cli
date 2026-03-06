package tableprinters

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/google/go-cmp/cmp"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_switchInterfaceNameLessFunc(t *testing.T) {
	tests := []struct {
		name  string
		conns []*apiv2.SwitchNicWithMachine
		want  []*apiv2.SwitchNicWithMachine
	}{
		{
			name: "sorts interface names for cumulus-like interface names",
			conns: []*apiv2.SwitchNicWithMachine{
				{Nic: &apiv2.SwitchNic{Name: "swp10"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s4"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s3"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s1"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s2"}},
				{Nic: &apiv2.SwitchNic{Name: "swp9"}},
			},
			want: []*apiv2.SwitchNicWithMachine{
				{Nic: &apiv2.SwitchNic{Name: "swp1s1"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s2"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s3"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s4"}},
				{Nic: &apiv2.SwitchNic{Name: "swp9"}},
				{Nic: &apiv2.SwitchNic{Name: "swp10"}},
			},
		},
		{
			name: "sorts interface names for sonic-like interface names",
			conns: []*apiv2.SwitchNicWithMachine{
				{Nic: &apiv2.SwitchNic{Name: "Ethernet3"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet49"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet10"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet2"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet1"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet11"}},
			},
			want: []*apiv2.SwitchNicWithMachine{
				{Nic: &apiv2.SwitchNic{Name: "Ethernet1"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet2"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet3"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet10"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet11"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet49"}},
			},
		},
		{
			name: "sorts interface names edge cases",
			conns: []*apiv2.SwitchNicWithMachine{
				{Nic: &apiv2.SwitchNic{Name: "123"}},
				{Nic: &apiv2.SwitchNic{Name: ""}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet1"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s4w5"}},
				{Nic: &apiv2.SwitchNic{Name: "foo"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s3w3"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet100"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s4w6"}},
				{Nic: &apiv2.SwitchNic{Name: ""}},
			},
			want: []*apiv2.SwitchNicWithMachine{
				{Nic: &apiv2.SwitchNic{Name: "swp1s3w3"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s4w5"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s4w6"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet1"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet100"}},
				{Nic: &apiv2.SwitchNic{Name: ""}},
				{Nic: &apiv2.SwitchNic{Name: ""}},
				{Nic: &apiv2.SwitchNic{Name: "123"}},
				{Nic: &apiv2.SwitchNic{Name: "foo"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Slice(tt.conns, switchInterfaceNameLessFunc(tt.conns))

			if diff := cmp.Diff(tt.conns, tt.want, protocmp.Transform()); diff != "" {
				t.Errorf("diff (+got -want):\n %s", diff)
			}
		})
	}
}

func Test_filterColumns(t *testing.T) {
	tests := []struct {
		name   string
		filter *apiv2.BGPFilter
		i      int
		want   []string
	}{
		{
			name:   "filter is nil",
			filter: nil,
			i:      1,
			want:   nil,
		},
		{
			name: "i exceeds vni and cidr length",
			filter: &apiv2.BGPFilter{
				Cidrs: []string{"1.1.1.1/32"},
				Vnis:  []string{"120"},
			},
			i:    1,
			want: []string{"", ""},
		},
		{
			name: "i exceeds vni but not cidr length",
			filter: &apiv2.BGPFilter{
				Cidrs: []string{"1.1.1.1/32", "2.2.2.2/32"},
				Vnis:  []string{"120"},
			},
			i:    1,
			want: []string{"", "2.2.2.2/32"},
		},
		{
			name: "i exceeds cidr but not vni length",
			filter: &apiv2.BGPFilter{
				Cidrs: []string{"1.1.1.1/32", "2.2.2.2/32"},
				Vnis:  []string{"120", "32", "400"},
			},
			i:    2,
			want: []string{"400", ""},
		},
		{
			name: "both vnis and cidr within range of i",
			filter: &apiv2.BGPFilter{
				Cidrs: []string{"1.1.1.1/32", "2.2.2.2/32", "3.3.3.3/32"},
				Vnis:  []string{"120", "32", "400"},
			},
			i:    2,
			want: []string{"400", "3.3.3.3/32"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterColumns(tt.filter, tt.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTablePrinter_SwitchTable(t *testing.T) {
	now := timestamppb.Now()
	tests := []struct {
		name       string
		switches   []*apiv2.Switch
		wide       bool
		wantHeader []string
		wantRows   [][]string
	}{
		{
			name:       "switches empty",
			switches:   []*apiv2.Switch{},
			wide:       false,
			wantHeader: []string{"ID", "Partition", "Rack", "OS", "Status", "Last Sync"},
			wantRows:   nil,
		},
		{
			name: "some switches",
			switches: []*apiv2.Switch{
				{
					Id:          "r01leaf01",
					Rack:        pointer.Pointer("rack01"),
					Partition:   "partition-a",
					ReplaceMode: apiv2.SwitchReplaceMode_SWITCH_REPLACE_MODE_OPERATIONAL,
					Os: &apiv2.SwitchOS{
						Vendor: apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_SONIC,
					},
					MachineConnections: []*apiv2.MachineConnection{
						{
							Nic: &apiv2.SwitchNic{
								Name: "Ethernet0",
								State: &apiv2.NicState{
									Desired: apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_DOWN.Enum(),
									Actual:  apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_UP,
								},
							},
						},
					},
					LastSync: &apiv2.SwitchSync{
						Time: now,
					},
					LastSyncError: &apiv2.SwitchSync{
						Time:  timestamppb.New(now.AsTime().Add(-7 * 24 * time.Hour)),
						Error: pointer.Pointer("sync took too long"),
					},
				},
				{
					Id:          "r01leaf02",
					Rack:        pointer.Pointer("rack01"),
					Partition:   "partition-a",
					ReplaceMode: apiv2.SwitchReplaceMode_SWITCH_REPLACE_MODE_REPLACE,
					Os: &apiv2.SwitchOS{
						Vendor: apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_CUMULUS,
					},
					LastSync: &apiv2.SwitchSync{
						Time: now,
					},
					LastSyncError: &apiv2.SwitchSync{
						Time:  timestamppb.New(now.AsTime().Add(time.Hour - 7*24*time.Hour)),
						Error: pointer.Pointer("sync took too long"),
					},
				},
				{
					Id:        "r02leaf01",
					Rack:      pointer.Pointer("rack02"),
					Partition: "partition-a",
					Os:        &apiv2.SwitchOS{},
					LastSync: &apiv2.SwitchSync{
						Time: timestamppb.New(now.AsTime().Add(-time.Hour)),
					},
					LastSyncError: &apiv2.SwitchSync{
						Time: now,
					},
				},
				{
					Id:        "r02leaf02",
					Rack:      pointer.Pointer("rack02"),
					Partition: "partition-a",
					LastSync: &apiv2.SwitchSync{
						Time: timestamppb.New(now.AsTime().Add(-10 * time.Minute)),
					},
				},
				{
					Id:        "r03leaf01",
					Rack:      pointer.Pointer("rack03"),
					Partition: "partition-a",
					LastSync: &apiv2.SwitchSync{
						Time:     now,
						Duration: durationpb.New(20 * time.Second),
					},
				},
				{
					Id:        "r03leaf02",
					Rack:      pointer.Pointer("rack03"),
					Partition: "partition-a",
					LastSync:  &apiv2.SwitchSync{},
					MachineConnections: []*apiv2.MachineConnection{
						{
							MachineId: "m1",
							Nic: &apiv2.SwitchNic{
								Name: "Ethernet1",
								State: &apiv2.NicState{
									Actual: apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_DOWN,
								},
							},
						},
					},
				},
			},
			wide:       false,
			wantHeader: []string{"ID", "Partition", "Rack", "OS", "Status", "Last Sync"},
			wantRows: [][]string{
				// FIXME: color of the dots is ignored; how to test for correct colors?
				{"r01leaf01", "partition-a", "rack01", "🦔", color.GreenString(dot), "0s ago"},                                                      // status green but error because one port is not in its desired state
				{"r01leaf02", "partition-a", "rack01", "🐢", nbr + color.RedString(dot), "0s ago"},                                                  // status red because in replace mode
				{"r02leaf01", "partition-a", "rack02", apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_UNSPECIFIED.String(), color.RedString(dot), "1h ago"}, // status red because last error came later than last sync
				{"r02leaf02", "partition-a", "rack02", "", color.RedString(dot), "10m ago"},                                                        // status red because last sync is too long ago
				{"r03leaf01", "partition-a", "rack03", "", color.YellowString(dot), "0s ago"},                                                      // status yellow because last sync duration was too long
				{"r03leaf02", "partition-a", "rack03", "", color.YellowString(dot), ""},                                                            // status yellow because not all connceted ports are up
			},
		},
		{
			name: "some switches wide",
			switches: []*apiv2.Switch{
				{
					Id:           "r01leaf01",
					Rack:         pointer.Pointer("rack01"),
					Partition:    "partition-a",
					ReplaceMode:  apiv2.SwitchReplaceMode_SWITCH_REPLACE_MODE_OPERATIONAL,
					ManagementIp: "1.1.1.1",
					Os: &apiv2.SwitchOS{
						Vendor:           apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_SONIC,
						MetalCoreVersion: "v0.15.0",
					},
					MachineConnections: []*apiv2.MachineConnection{
						{
							Nic: &apiv2.SwitchNic{
								Name: "Ethernet0",
								State: &apiv2.NicState{
									Desired: apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_DOWN.Enum(),
									Actual:  apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_UP,
								},
							},
						},
					},
					LastSync: &apiv2.SwitchSync{
						Time:     now,
						Duration: durationpb.New(time.Second),
					},
					LastSyncError: &apiv2.SwitchSync{
						Time:  timestamppb.New(now.AsTime().Add(-7 * 24 * time.Hour)),
						Error: pointer.Pointer("sync took too long"),
					},
				},
				{
					Id:           "r01leaf02",
					Rack:         pointer.Pointer("rack01"),
					Partition:    "partition-a",
					ReplaceMode:  apiv2.SwitchReplaceMode_SWITCH_REPLACE_MODE_REPLACE,
					ManagementIp: "2.2.2.2",
					Os: &apiv2.SwitchOS{
						Vendor:           apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_CUMULUS,
						MetalCoreVersion: "v0.13.0",
					},
					LastSync: &apiv2.SwitchSync{
						Time: now,
					},
					LastSyncError: &apiv2.SwitchSync{
						Time:  timestamppb.New(now.AsTime().Add(time.Hour - 7*24*time.Hour)),
						Error: pointer.Pointer("sync took too long"),
					},
				},
				{
					Id:           "r02leaf01",
					Rack:         pointer.Pointer("rack02"),
					Partition:    "partition-a",
					ManagementIp: "3.3.3.3",
					Os:           &apiv2.SwitchOS{},
					LastSync: &apiv2.SwitchSync{
						Time: timestamppb.New(now.AsTime().Add(-time.Hour)),
					},
					LastSyncError: &apiv2.SwitchSync{
						Time:  now,
						Error: pointer.Pointer("error"),
					},
				},
			},
			wide:       true,
			wantHeader: []string{"ID", "Partition", "Rack", "OS", "Metalcore", "IP", "Mode", "Last Sync", "Sync Duration", "Last Error"},
			wantRows: [][]string{
				{"r01leaf01", "partition-a", "rack01", "SONiC", "v0.15.0", "1.1.1.1", "operational", "0s ago", "1s", "\"Ethernet0\" is up but should be down"},
				{"r01leaf02", "partition-a", "rack01", "Cumulus", "v0.13.0", "2.2.2.2", "replace", "0s ago", "", "6d 23h ago: sync took too long"},
				{"r02leaf01", "partition-a", "rack02", "", "", "3.3.3.3", "operational", "1h ago", "", "0s ago: error"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := New()
			p := printers.NewTablePrinter(&printers.TablePrinterConfig{
				ToHeaderAndRows: tp.ToHeaderAndRows,
				Wide:            tt.wide,
			})
			tp.SetPrinter(p)

			gotHeader, gotRows, err := tp.SwitchTable(tt.switches, tt.wide)
			if err != nil {
				t.Errorf("TablePrinter.SwitchTable() error = %v", err)
				return
			}
			if diff := cmp.Diff(tt.wantHeader, gotHeader); diff != "" {
				t.Errorf("TablePrinter.SwitchTable() diff header = %s", diff)
			}
			if diff := cmp.Diff(tt.wantRows, gotRows); diff != "" {
				t.Errorf("TablePrinter.SwitchTable() diff rows = %s", diff)
			}
		})
	}
}

func TestTablePrinter_SwitchDetailTable(t *testing.T) {
	tests := []struct {
		name       string
		switches   []SwitchDetail
		wantHeader []string
		wantRows   [][]string
	}{
		{
			name:       "empty switches",
			switches:   []SwitchDetail{},
			wantHeader: []string{"Partition", "Rack", "Switch", "Port", "Machine", "VNI-Filter", "CIDR-Filter"},
			wantRows:   nil,
		},
		{
			name: "some switches",
			switches: []SwitchDetail{
				{
					Switch: &apiv2.Switch{
						Id:        "leaf01",
						Rack:      pointer.Pointer("rack01"),
						Partition: "partition-a",
						Nics: []*apiv2.SwitchNic{
							{
								Name: "Ethernet0",
								BgpFilter: &apiv2.BGPFilter{
									Cidrs: []string{"1.1.1.0/24", "2.2.2.0/24"},
									Vnis:  []string{"104"},
								},
							},
							{
								Name: "Ethernet1",
							},
						},
						MachineConnections: []*apiv2.MachineConnection{
							{
								MachineId: "m1",
								Nic: &apiv2.SwitchNic{
									Name: "Ethernet0",
									BgpFilter: &apiv2.BGPFilter{
										Cidrs: []string{"1.1.1.0/24", "2.2.2.0/24"},
										Vnis:  []string{"104"},
									},
								},
							},
							{
								MachineId: "m2",
								Nic: &apiv2.SwitchNic{
									Name: "Ethernet1",
								},
							},
						},
					},
				},
				{
					Switch: &apiv2.Switch{
						Id:        "leaf02",
						Rack:      pointer.Pointer("rack01"),
						Partition: "partition-a",
						Nics: []*apiv2.SwitchNic{
							{
								Name: "Ethernet0",
								BgpFilter: &apiv2.BGPFilter{
									Cidrs: []string{"1.1.1.0/24", "2.2.2.0/24"},
									Vnis:  []string{"150"},
								},
							},
							{
								Name: "Ethernet1",
							},
						},
						MachineConnections: []*apiv2.MachineConnection{
							{
								MachineId: "m1",
								Nic: &apiv2.SwitchNic{
									Name: "Ethernet0",
									BgpFilter: &apiv2.BGPFilter{
										Cidrs: []string{"1.1.1.0/24", "2.2.2.0/24"},
										Vnis:  []string{"150"},
									},
								},
							},
							{
								MachineId: "m2",
								Nic: &apiv2.SwitchNic{
									Name: "Ethernet1",
								},
							},
						},
					},
				},
			},
			wantHeader: []string{"Partition", "Rack", "Switch", "Port", "Machine", "VNI-Filter", "CIDR-Filter"},
			wantRows: [][]string{
				{"partition-a", "rack01", "leaf01", "Ethernet0", "m1", "104", "1.1.1.0/24"},
				{"", "", "", "", "", "", "2.2.2.0/24"},
				{"partition-a", "rack01", "leaf01", "Ethernet1", "m2"},
				{"partition-a", "rack01", "leaf02", "Ethernet0", "m1", "150", "1.1.1.0/24"},
				{"", "", "", "", "", "", "2.2.2.0/24"},
				{"partition-a", "rack01", "leaf02", "Ethernet1", "m2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := New()
			p := printers.NewTablePrinter(&printers.TablePrinterConfig{
				ToHeaderAndRows: tp.ToHeaderAndRows,
			})
			tp.SetPrinter(p)

			gotHeader, gotRows, err := tp.SwitchDetailTable(tt.switches)
			if err != nil {
				t.Errorf("TablePrinter.SwitchDetailTable() error = %v", err)
				return
			}
			if diff := cmp.Diff(tt.wantHeader, gotHeader); diff != "" {
				t.Errorf("TablePrinter.SwitchDetailTable() diff header = %s", diff)
			}
			if diff := cmp.Diff(tt.wantRows, gotRows); diff != "" {
				t.Errorf("TablePrinter.SwitchDetailTable() diff rows = %s", diff)
			}
		})
	}
}

func TestTablePrinter_SwitchWithConnectedMachinesTable(t *testing.T) {
	now := timestamppb.Now()

	tests := []struct {
		name       string
		res        []*apiv2.SwitchWithMachines
		wide       bool
		wantHeader []string
		wantRows   [][]string
	}{
		{
			name:       "empty response",
			res:        []*apiv2.SwitchWithMachines{},
			wide:       false,
			wantHeader: []string{"ID", "NIC Name", "Identifier", "Partition", "Rack", "Size", "Product Serial", "Chassis Serial"},
			wantRows:   nil,
		},
		{
			name: "switches with machines",
			res: []*apiv2.SwitchWithMachines{
				{
					Id:        "r01leaf01",
					Partition: "partition-a",
					Rack:      "rack01",
					Connections: []*apiv2.SwitchNicWithMachine{
						{
							Nic: &apiv2.SwitchNic{
								Name:       "Ethernet10",
								Identifier: "Eth3/3",
								BgpPortState: &apiv2.SwitchBGPPortState{
									BgpState:              apiv2.BGPState_BGP_STATE_ESTABLISHED,
									BgpTimerUpEstablished: timestamppb.New(now.AsTime().Add(-5 * 24 * time.Hour)),
								},
							},
							Machine: &apiv2.Machine{
								Uuid: "m2",
								Size: &apiv2.Size{
									Id: "medium",
								},
							},
							Fru: &apiv2.MachineFRU{
								ChassisPartSerial: pointer.Pointer("c234"),
								ProductSerial:     pointer.Pointer("p234"),
							},
						},
						{
							Nic: &apiv2.SwitchNic{
								Name:       "Ethernet2",
								Identifier: "Eth1/3",
								State: &apiv2.NicState{
									Actual: apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_DOWN,
								},
							},
							Machine: &apiv2.Machine{
								Uuid: "m1",
								Size: &apiv2.Size{
									Id: "large",
								},
							},
							Fru: &apiv2.MachineFRU{
								ChassisPartSerial: pointer.Pointer("c123"),
								ProductSerial:     pointer.Pointer("p123"),
							},
						},
					},
				},
				{
					Id:        "r02leaf02",
					Partition: "partition-b",
					Rack:      "rack02",
					Connections: []*apiv2.SwitchNicWithMachine{
						{
							Nic: &apiv2.SwitchNic{
								Name:       "Ethernet5",
								Identifier: "Eth2/2",
								BgpPortState: &apiv2.SwitchBGPPortState{
									BgpState: apiv2.BGPState_BGP_STATE_ACTIVE,
								},
							},
							Machine: &apiv2.Machine{
								Uuid: "m3",
								Size: &apiv2.Size{
									Id: "small",
								},
							},
							Fru: &apiv2.MachineFRU{
								ChassisPartSerial: pointer.Pointer("c345"),
								ProductSerial:     pointer.Pointer("p345"),
							},
						},
					},
				},
			},
			wide:       false,
			wantHeader: []string{"ID", "NIC Name", "Identifier", "Partition", "Rack", "Size", "Product Serial", "Chassis Serial"},
			wantRows: [][]string{
				{"r01leaf01", "", "", "partition-a", "rack01"},
				{"├─╴m1", "Ethernet2 (down)", "Eth1/3", "partition-a", "rack01", "large", "p123", "c123"},
				{"└─╴m2", "Ethernet10", "Eth3/3", "partition-a", "rack01", "medium", "p234", "c234"},
				{"r02leaf02", "", "", "partition-b", "rack02"},
				{"└─╴m3", "Ethernet5", "Eth2/2", "partition-b", "rack02", "small", "p345", "c345"},
			},
		},
		{
			name: "wide",
			res: []*apiv2.SwitchWithMachines{
				{
					Id:        "r01leaf01",
					Partition: "partition-a",
					Rack:      "rack01",
					Connections: []*apiv2.SwitchNicWithMachine{
						{
							Nic: &apiv2.SwitchNic{
								Name:       "Ethernet10",
								Identifier: "Eth3/3",
								BgpPortState: &apiv2.SwitchBGPPortState{
									BgpState:              apiv2.BGPState_BGP_STATE_ESTABLISHED,
									BgpTimerUpEstablished: timestamppb.New(now.AsTime().Add(-5 * 24 * time.Hour)),
								},
							},
							Machine: &apiv2.Machine{
								Uuid: "m2",
								Size: &apiv2.Size{
									Id: "medium",
								},
								Allocation: &apiv2.MachineAllocation{
									Hostname: "fw1",
									Vpn: &apiv2.MachineVPN{
										Connected: true,
									},
								},
								Status: &apiv2.MachineStatus{
									Condition: &apiv2.MachineCondition{
										State: apiv2.MachineState_MACHINE_STATE_AVAILABLE,
									},
									Liveliness: apiv2.MachineLiveliness_MACHINE_LIVELINESS_ALIVE,
								},
							},
							Fru: &apiv2.MachineFRU{
								ChassisPartSerial: pointer.Pointer("c234"),
								ProductSerial:     pointer.Pointer("p234"),
							},
						},
						{
							Nic: &apiv2.SwitchNic{
								Name:       "Ethernet2",
								Identifier: "Eth1/3",
								State: &apiv2.NicState{
									Actual: apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_DOWN,
								},
							},
							Machine: &apiv2.Machine{
								Uuid: "m1",
								Size: &apiv2.Size{
									Id: "large",
								},
								Status: &apiv2.MachineStatus{
									Condition: &apiv2.MachineCondition{
										State: apiv2.MachineState_MACHINE_STATE_LOCKED,
									},
									Liveliness: apiv2.MachineLiveliness_MACHINE_LIVELINESS_DEAD,
								},
								RecentProvisioningEvents: &apiv2.MachineRecentProvisioningEvents{
									LastErrorEvent: &apiv2.MachineProvisioningEvent{
										Time: timestamppb.New(now.AsTime().Add(-time.Hour)),
									},
									State: apiv2.MachineProvisioningEventState_MACHINE_PROVISIONING_EVENT_STATE_CRASHLOOP,
								},
							},
							Fru: &apiv2.MachineFRU{
								ChassisPartSerial: pointer.Pointer("c123"),
								ProductSerial:     pointer.Pointer("p123"),
							},
						},
					},
				},
				{
					Id:        "r02leaf02",
					Partition: "partition-b",
					Rack:      "rack02",
					Connections: []*apiv2.SwitchNicWithMachine{
						{
							Nic: &apiv2.SwitchNic{
								Name:       "Ethernet5",
								Identifier: "Eth2/2",
								State: &apiv2.NicState{
									Actual: apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_UNKNOWN,
								},
								BgpPortState: &apiv2.SwitchBGPPortState{
									BgpState: apiv2.BGPState_BGP_STATE_ACTIVE,
								},
							},
							Machine: &apiv2.Machine{
								Uuid: "m3",
								Size: &apiv2.Size{
									Id: "small",
								},
								Allocation: &apiv2.MachineAllocation{
									Hostname: "worker1",
								},
								Status: &apiv2.MachineStatus{
									Condition: &apiv2.MachineCondition{
										State: apiv2.MachineState_MACHINE_STATE_RESERVED,
									},
									Liveliness: apiv2.MachineLiveliness_MACHINE_LIVELINESS_UNKNOWN,
								},
								RecentProvisioningEvents: &apiv2.MachineRecentProvisioningEvents{
									LastErrorEvent: &apiv2.MachineProvisioningEvent{
										Time: timestamppb.New(now.AsTime().Add(-2 * time.Hour)),
									},
									State: apiv2.MachineProvisioningEventState_MACHINE_PROVISIONING_EVENT_STATE_FAILED_RECLAIM,
								},
							},
							Fru: &apiv2.MachineFRU{
								ChassisPartSerial: pointer.Pointer("c345"),
								ProductSerial:     pointer.Pointer("p345"),
							},
						},
					},
				},
			},
			wide:       true,
			wantHeader: []string{"ID", "", "NIC Name", "Identifier", "Partition", "Rack", "Size", "Hostname", "Product Serial", "Chassis Serial"},
			wantRows: [][]string{
				{"r01leaf01", "", "", "", "partition-a", "rack01"},
				{"├─╴m1", "💀 🔒 ⭕ ❗", "Ethernet2 (down)", "Eth1/3", "partition-a", "rack01", "large", "", "p123", "c123"},
				{"└─╴m2", "🛡", "Ethernet10 (BGP:Established(5d))", "Eth3/3", "partition-a", "rack01", "medium", "fw1", "p234", "c234"},
				{"r02leaf02", "", "", "", "partition-b", "rack02"},
				{"└─╴m3", "❓ 🚧 🚑", "Ethernet5 (unknown) (BGP:Active)", "Eth2/2", "partition-b", "rack02", "small", "worker1", "p345", "c345"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := New()
			p := printers.NewTablePrinter(&printers.TablePrinterConfig{
				ToHeaderAndRows: tp.ToHeaderAndRows,
			})
			tp.SetPrinter(p)
			tp.SetLastEventErrorThreshold(2 * time.Hour)

			gotHeader, gotRows, err := tp.SwitchWithConnectedMachinesTable(tt.res, tt.wide)
			if err != nil {
				t.Errorf("TablePrinter.SwitchWithConnectedMachinesTable() error = %v", err)
				return
			}
			if diff := cmp.Diff(tt.wantHeader, gotHeader); diff != "" {
				t.Errorf("TablePrinter.SwitchWithConnectedMachinesTable() header diff = %s", diff)
			}
			if diff := cmp.Diff(tt.wantRows, gotRows); diff != "" {
				t.Errorf("TablePrinter.SwitchWithConnectedMachinesTable() rows diff = %s", diff)
			}
		})
	}
}
