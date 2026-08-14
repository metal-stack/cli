package helpers

import (
	"fmt"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func MachineResponseToCreate(r *apiv2.Machine) (*apiv2.MachineServiceCreateRequest, error) {
	if r.Allocation == nil {
		return nil, fmt.Errorf("allocation is nil")
	}

	var (
		networks     []*apiv2.MachineAllocationNetwork
		firewallSpec *apiv2.FirewallSpec
	)

	for _, nw := range r.Allocation.Networks {
		networks = append(networks, &apiv2.MachineAllocationNetwork{
			Network: nw.Network,
			Ips:     nw.Ips,
		})
	}

	if r.Allocation.AllocationType.Enum() == apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_FIREWALL.Enum() {
		firewallSpec = &apiv2.FirewallSpec{
			FirewallRules: &apiv2.FirewallRules{
				Egress:  r.Allocation.FirewallRules.Egress,
				Ingress: r.Allocation.FirewallRules.Ingress,
			},
		}
	}

	return &apiv2.MachineServiceCreateRequest{
		Project:          r.Allocation.Project,
		Name:             r.Allocation.Name,
		Description:      &r.Allocation.Description,
		Hostname:         &r.Allocation.Hostname,
		Partition:        new(pointer.SafeDeref(r.Partition).Id),
		Size:             new(pointer.SafeDeref(r.Size).Id),
		Image:            r.Allocation.Image.Id,
		FilesystemLayout: pointer.PointerOrNil(r.Allocation.FilesystemLayout.Id),
		SshPublicKeys:    r.Allocation.SshPublicKeys,
		Userdata:         pointer.PointerOrNil(r.Allocation.Userdata),
		Labels:           pointer.SafeDeref(r.Meta).Labels,
		Networks:         networks,
		DnsServers:       r.Allocation.DnsServers,
		NtpServers:       r.Allocation.NtpServers,
		AllocationType:   r.Allocation.AllocationType,
		FirewallSpec:     firewallSpec,
		// PlacementTags:    r.Allocation.Plac, // TODO: should be stored in the allocation to see what was provided
	}, nil
}

func MachineResponseToUpdate(r *apiv2.Machine) (*apiv2.MachineServiceUpdateRequest, error) {
	if r.Allocation == nil {
		return nil, fmt.Errorf("allocation is nil")
	}

	return &apiv2.MachineServiceUpdateRequest{
		Uuid:          r.Uuid,
		UpdateMeta:    UpdateMetaFromMeta(r.Meta),
		Labels:        UpdateLabelsFromMeta(r.Meta),
		Project:       r.Allocation.Project,
		Description:   &r.Allocation.Description,
		SshPublicKeys: r.Allocation.SshPublicKeys,
	}, nil
}
