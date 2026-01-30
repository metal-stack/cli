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
		conns []*apiv2.MachineConnection
		want  []*apiv2.MachineConnection
	}{
		{
			name: "sorts interface names for cumulus-like interface names",
			conns: []*apiv2.MachineConnection{
				{Nic: &apiv2.SwitchNic{Name: "swp10"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s4"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s3"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s1"}},
				{Nic: &apiv2.SwitchNic{Name: "swp1s2"}},
				{Nic: &apiv2.SwitchNic{Name: "swp9"}},
			},
			want: []*apiv2.MachineConnection{
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
			conns: []*apiv2.MachineConnection{
				{Nic: &apiv2.SwitchNic{Name: "Ethernet3"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet49"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet10"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet2"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet1"}},
				{Nic: &apiv2.SwitchNic{Name: "Ethernet11"}},
			},
			want: []*apiv2.MachineConnection{
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
			conns: []*apiv2.MachineConnection{
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
			want: []*apiv2.MachineConnection{
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
				{"r01leaf01", "partition-a", "rack01", apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_SONIC.String(), "v0.15.0", "1.1.1.1", "operational", "0s ago", "1s", "\"Ethernet0\" is UP but should be DOWN"},
				{"r01leaf02", "partition-a", "rack01", apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_CUMULUS.String(), "v0.13.0", "2.2.2.2", "replace", "0s ago", "", "6d 23h ago: sync took too long"},
				{"r02leaf01", "partition-a", "rack02", apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_UNSPECIFIED.String(), "", "3.3.3.3", "operational", "1h ago", "", "0s ago: error"},
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
	type fields struct {
		t *printers.TablePrinter
	}
	type args struct {
		switches []SwitchDetail
		wide     bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		want1   [][]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TablePrinter{
				t: tt.fields.t,
			}
			got, got1, err := tr.SwitchDetailTable(tt.args.switches, tt.args.wide)
			if (err != nil) != tt.wantErr {
				t.Errorf("TablePrinter.SwitchDetailTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TablePrinter.SwitchDetailTable() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("TablePrinter.SwitchDetailTable() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
