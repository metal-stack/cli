package tableprinters

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"google.golang.org/protobuf/testing/protocmp"
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
