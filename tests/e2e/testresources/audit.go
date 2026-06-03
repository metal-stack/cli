package testresources

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Trace1 = func() *apiv2.AuditTrace {
		return &apiv2.AuditTrace{
			Uuid:      "d1ff7267-2fbb-4a63-a7c1-44f1a83381a7",
			Timestamp: timestamppb.New(e2e.TimeBubbleStartTime()),
			User:      "me",
			Tenant:    "a",
			Project:   new(Project1().Uuid),
			Method:    "/apiv2.Create/",
			Body:      new("request body"),
			SourceIp:  "1.2.3.4",
			Phase:     apiv2.AuditPhase_AUDIT_PHASE_REQUEST,
		}
	}

	Trace2 = func() *apiv2.AuditTrace {
		return &apiv2.AuditTrace{
			Uuid:       "d1ff7267-2fbb-4a63-a7c1-44f1a83381a7",
			Timestamp:  timestamppb.New(e2e.TimeBubbleStartTime()),
			User:       "me",
			Tenant:     "a",
			Project:    new(Project1().Uuid),
			Method:     "/apiv2.List/",
			Body:       new("result body"),
			SourceIp:   "1.2.3.4",
			ResultCode: new(int32(codes.AlreadyExists)),
			Phase:      apiv2.AuditPhase_AUDIT_PHASE_RESPONSE,
		}
	}

	Trace3 = func() *apiv2.AuditTrace {
		return &apiv2.AuditTrace{
			Uuid:      "5091c4e9-e8db-483c-ab6b-fe14f82570a7",
			Timestamp: timestamppb.New(e2e.TimeBubbleStartTime().AddDate(1, 0, 0)),
			User:      "Larry",
			Tenant:    "b",
			Project:   new(Project1().Uuid),
			Method:    "/apiv2.List/",
			Body:      new("result body"),
			SourceIp:  "1.2.3.4",
			Phase:     apiv2.AuditPhase_AUDIT_PHASE_REQUEST,
		}
	}
)
