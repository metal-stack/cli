package sorters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
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

func TenantMemberSorter() *multisort.Sorter[*apiv2.TenantMember] {
	return multisort.New(multisort.FieldMap[*apiv2.TenantMember]{
		"id": func(a, b *apiv2.TenantMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"role": func(a, b *apiv2.TenantMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Role, b.Role, descending)
		},
		"created": func(a, b *apiv2.TenantMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.CreatedAt.AsTime().UnixMilli(), b.CreatedAt.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "role"}, {ID: "id"}})
}
