package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	cliv2 "github.com/metal-stack/api/go/metalstack/cli/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func ProjectSorter() *multisort.Sorter[*apiv2.Project] {
	return multisort.New(multisort.FieldMap[*apiv2.Project]{
		"id": func(a, b *apiv2.Project, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
		"name": func(a, b *apiv2.Project, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"tenant": func(a, b *apiv2.Project, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Tenant, b.Tenant, descending)
		},
	}, multisort.Keys{{ID: "tenant"}, {ID: "name"}, {ID: "id"}})
}

func ProjectInviteSorter() *multisort.Sorter[*apiv2.ProjectInvite] {
	return multisort.New(multisort.FieldMap[*apiv2.ProjectInvite]{
		"project": func(a, b *apiv2.ProjectInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Project, b.Project, descending)
		},
		"secret": func(a, b *apiv2.ProjectInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Secret, b.Secret, descending)
		},
		"role": func(a, b *apiv2.ProjectInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Role, b.Role, descending)
		},
		"expiration": func(a, b *apiv2.ProjectInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.ExpiresAt.AsTime().UnixMilli(), b.ExpiresAt.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "project"}, {ID: "role"}, {ID: "expiration"}})
}

func ProjectMemberSorter() *multisort.Sorter[*cliv2.ProjectMember] {
	return multisort.New(multisort.FieldMap[*cliv2.ProjectMember]{
		"id": func(a, b *cliv2.ProjectMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.ProjectMember.Id, b.ProjectMember.Id, descending)
		},
		"role": func(a, b *cliv2.ProjectMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.ProjectMember.Role, b.ProjectMember.Role, descending)
		},
		"created": func(a, b *cliv2.ProjectMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.ProjectMember.CreatedAt.AsTime().UnixMilli(), b.ProjectMember.CreatedAt.AsTime().UnixMilli(), descending)
		},
		"inherited": func(a, b *cliv2.ProjectMember, descending bool) multisort.CompareResult {
			boolToInt := func(in bool) int {
				if in {
					return 1
				}
				return 0
			}
			return multisort.Compare(boolToInt(a.ProjectMember.InheritedMembership), boolToInt(b.ProjectMember.InheritedMembership), descending)
		},
	}, multisort.Keys{{ID: "inherited", Descending: false}, {ID: "role"}, {ID: "id"}})
}
