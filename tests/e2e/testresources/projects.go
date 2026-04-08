package testresources

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Project1 = func() *apiv2.Project {
		return &apiv2.Project{
			Uuid:        "0d81bca7-73f6-4da3-8397-4a8c52a0c583",
			Name:        "project-a",
			Description: "first project",
			Tenant:      "metal-stack",
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		}
	}
	Project2 = func() *apiv2.Project {
		return &apiv2.Project{
			Uuid:        "f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c",
			Name:        "project-b",
			Description: "second project",
			Tenant:      "metal-stack",
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		}
	}

	Project1Invite = func() *apiv2.ProjectInvite {
		return &apiv2.ProjectInvite{
			Secret:      "secret",
			Project:     "0d81bca7-73f6-4da3-8397-4a8c52a0c583",
			Role:        apiv2.ProjectRole_PROJECT_ROLE_EDITOR,
			Joined:      false,
			ProjectName: "project-a",
			TenantName:  "metal-stack",
			ExpiresAt:   timestamppb.New(e2e.TimeBubbleStartTime().Add(48 * time.Hour)),
		}
	}
	Project2Invite = func() *apiv2.ProjectInvite {
		return &apiv2.ProjectInvite{
			Secret:      "secret",
			Project:     "f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c",
			Role:        apiv2.ProjectRole_PROJECT_ROLE_EDITOR,
			Joined:      false,
			ProjectName: "project-b",
			TenantName:  "metal-stack",
			ExpiresAt:   timestamppb.New(e2e.TimeBubbleStartTime().Add(48 * time.Hour)),
		}
	}

	Project1Members = func() *apiv2.ProjectMember {
		return &apiv2.ProjectMember{
			Id:                  "16d6e8ba-f574-494f-8d5e-74f6cb2d8db0",
			Role:                apiv2.ProjectRole_PROJECT_ROLE_OWNER,
			InheritedMembership: false,
			CreatedAt:           timestamppb.New(e2e.TimeBubbleStartTime()),
		}
	}
	Project2Members = func() *apiv2.ProjectMember {
		return &apiv2.ProjectMember{
			Id:                  "40c0da4b-9eb9-4371-91aa-1ae62193fa54",
			Role:                apiv2.ProjectRole_PROJECT_ROLE_EDITOR,
			InheritedMembership: true,
			CreatedAt:           timestamppb.New(e2e.TimeBubbleStartTime()),
		}
	}
)
