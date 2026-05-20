package testresources

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Tenant1 = func() *apiv2.Tenant {
		return &apiv2.Tenant{
			Login:       "metal-stack",
			Name:        "Metal Stack",
			Email:       "info@metal-stack.io",
			Description: "a tenant",
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		}
	}
	Tenant2 = func() *apiv2.Tenant {
		return &apiv2.Tenant{
			Login:       "acme-corp",
			Name:        "ACME Corp",
			Email:       "admin@acme.io",
			Description: "another tenant",
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		}
	}

	Tenant1Invite = func() *apiv2.TenantInvite {
		return &apiv2.TenantInvite{
			Secret:           "secret",
			TargetTenant:     "metal-stack",
			TargetTenantName: "Metal Stack",
			Role:             apiv2.TenantRole_TENANT_ROLE_VIEWER,
			Joined:           false,
			Tenant:           "metal-stack",
			TenantName:       "Metal Stack",
			ExpiresAt:        timestamppb.New(e2e.TimeBubbleStartTime().Add(48 * time.Hour)),
		}
	}
	Tenant2Invite = func() *apiv2.TenantInvite {
		return &apiv2.TenantInvite{
			Secret:           "secret",
			TargetTenant:     "acme-corp",
			TargetTenantName: "ACME Corp",
			Role:             apiv2.TenantRole_TENANT_ROLE_EDITOR,
			Joined:           false,
			Tenant:           "acme-corp",
			TenantName:       "ACME Corp",
			ExpiresAt:        timestamppb.New(e2e.TimeBubbleStartTime().Add(48 * time.Hour)),
		}
	}

	Tenant1Members = func() *apiv2.TenantMember {
		return &apiv2.TenantMember{
			Id:        "16d6e8ba-f574-494f-8d5e-74f6cb2d8db0",
			Role:      apiv2.TenantRole_TENANT_ROLE_OWNER,
			CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			Projects:  []string{Project1().Uuid, Project2().Uuid},
		}
	}
	Tenant2Members = func() *apiv2.TenantMember {
		return &apiv2.TenantMember{
			Id:        "40c0da4b-9eb9-4371-91aa-1ae62193fa54",
			Role:      apiv2.TenantRole_TENANT_ROLE_EDITOR,
			CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			Projects:  []string{Project1().Uuid},
		}
	}
)
