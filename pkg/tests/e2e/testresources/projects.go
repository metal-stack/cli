package testresources

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
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
				CreatedAt: timestamppb.New(time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC)),
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
				CreatedAt: timestamppb.New(time.Date(2025, 7, 15, 14, 30, 0, 0, time.UTC)),
			},
		}
	}
)
