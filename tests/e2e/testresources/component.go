package testresources

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Component1 = func() *apiv2.Component {
		return &apiv2.Component{
			Uuid:       "c1a2b3d4-e5f6-7890-abcd-ef1234567890",
			Type:       apiv2.ComponentType_COMPONENT_TYPE_METAL_CORE,
			Identifier: "metal-core-1",
			StartedAt:  timestamppb.New(e2e.TimeBubbleStartTime()),
			ReportedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			Interval:   durationpb.New(10 * time.Second),
			Version: &apiv2.Version{
				Version: "v1.0.0",
			},
			Token: &apiv2.Token{
				Uuid:    "t1a2b3d4-e5f6-7890-abcd-ef1234567890",
				Expires: timestamppb.New(e2e.TimeBubbleStartTime().Add(24 * time.Hour)),
			},
		}
	}
	Component2 = func() *apiv2.Component {
		return &apiv2.Component{
			Uuid:       "d2b3c4e5-f6a7-8901-bcde-f12345678901",
			Type:       apiv2.ComponentType_COMPONENT_TYPE_PIXIECORE,
			Identifier: "pixiecore-1",
			StartedAt:  timestamppb.New(e2e.TimeBubbleStartTime()),
			ReportedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			Interval:   durationpb.New(10 * time.Second),
			Version: &apiv2.Version{
				Version: "v2.0.0",
			},
			Token: &apiv2.Token{
				Uuid:    "t2b3c4e5-f6a7-8901-bcde-f12345678901",
				Expires: timestamppb.New(e2e.TimeBubbleStartTime().Add(48 * time.Hour)),
			},
		}
	}
)
