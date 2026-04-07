package testresources

import (
	"time"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Task1 = func() *adminv2.TaskInfo {
		return &adminv2.TaskInfo{
			Id:            "550e8400-e29b-41d4-a716-446655440000",
			Queue:         "default",
			Type:          "image-provision",
			Payload:       []byte(`{"machine_id":"machine1"}`),
			State:         adminv2.TaskState_TASK_STATE_ACTIVE,
			MaxRetry:      3,
			Retried:       1,
			LastError:     "connection timeout",
			LastFailedAt:  timestamppb.New(e2e.TimeBubbleStartTime().Add(-5 * time.Minute)),
			Timeout:       durationpb.New(30 * time.Second),
			Deadline:      timestamppb.New(e2e.TimeBubbleStartTime().Add(5 * time.Minute)),
			NextProcessAt: timestamppb.New(e2e.TimeBubbleStartTime().Add(10 * time.Second)),
			Retention:     durationpb.New(24 * time.Hour),
			CompletedAt:   nil,
		}
	}
	Task2 = func() *adminv2.TaskInfo {
		return &adminv2.TaskInfo{
			Id:            "550e8400-e29b-41d4-a716-446655440001",
			Queue:         "default",
			Type:          "firewall-update",
			Payload:       []byte(`{"firewall_id":"fw1"}`),
			State:         adminv2.TaskState_TASK_STATE_PENDING,
			MaxRetry:      5,
			Retried:       0,
			LastError:     "",
			LastFailedAt:  nil,
			Timeout:       durationpb.New(60 * time.Second),
			Deadline:      timestamppb.New(e2e.TimeBubbleStartTime().Add(10 * time.Minute)),
			NextProcessAt: timestamppb.New(e2e.TimeBubbleStartTime().Add(2 * time.Second)),
			Retention:     durationpb.New(48 * time.Hour),
			CompletedAt:   nil,
		}
	}
	Task3 = func() *adminv2.TaskInfo {
		return &adminv2.TaskInfo{
			Id:            "550e8400-e29b-41d4-a716-446655440002",
			Queue:         "high-priority",
			Type:          "machine-reimage",
			Payload:       []byte(`{"machine_id":"machine2"}`),
			State:         adminv2.TaskState_TASK_STATE_COMPLETED,
			MaxRetry:      3,
			Retried:       0,
			LastError:     "",
			LastFailedAt:  nil,
			Timeout:       durationpb.New(5 * time.Minute),
			Deadline:      timestamppb.New(e2e.TimeBubbleStartTime().Add(30 * time.Minute)),
			NextProcessAt: nil,
			Retention:     durationpb.New(7 * 24 * time.Hour),
			CompletedAt:   timestamppb.New(e2e.TimeBubbleStartTime().Add(-10 * time.Minute)),
		}
	}
)
