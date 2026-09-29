package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	cliv2 "github.com/metal-stack/api/go/metalstack/cli/v2"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func TenantSorter() *multisort.Sorter[*apiv2.Tenant] {
	return multisort.New(multisort.FieldMap[*apiv2.Tenant]{
		"id": func(a, b *apiv2.Tenant, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Login, b.Login, descending)
		},
		"name": func(a, b *apiv2.Tenant, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"since": func(a, b *apiv2.Tenant, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Meta.CreatedAt.AsTime().UnixMilli(), b.Meta.CreatedAt.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "since", Descending: true}})
}

func TenantInviteSorter() *multisort.Sorter[*apiv2.TenantInvite] {
	return multisort.New(multisort.FieldMap[*apiv2.TenantInvite]{
		"tenant": func(a, b *apiv2.TenantInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Tenant, b.Tenant, descending)
		},
		"secret": func(a, b *apiv2.TenantInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Secret, b.Secret, descending)
		},
		"role": func(a, b *apiv2.TenantInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Role, b.Role, descending)
		},
		"expiration": func(a, b *apiv2.TenantInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.ExpiresAt.AsTime().UnixMilli(), b.ExpiresAt.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "tenant"}, {ID: "role"}, {ID: "expiration"}})
}

func TenantMemberSorter() *multisort.Sorter[*cliv2.TenantMember] {
	return multisort.New(multisort.FieldMap[*cliv2.TenantMember]{
		"id": func(a, b *cliv2.TenantMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.TenantMember.Id, b.TenantMember.Id, descending)
		},
		"role": func(a, b *cliv2.TenantMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.TenantMember.Role, b.TenantMember.Role, descending)
		},
		"created": func(a, b *cliv2.TenantMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.TenantMember.CreatedAt.AsTime().UnixMilli(), b.TenantMember.CreatedAt.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "role"}, {ID: "id"}})
}
