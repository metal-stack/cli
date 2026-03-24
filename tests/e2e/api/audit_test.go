package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
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
			Project:   &testresources.Project1().Uuid,
			Method:    "/apiv2.List/",
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
			Project:    &testresources.Project1().Uuid,
			Method:     "/apiv2.List/",
			Body:       new("result body"),
			SourceIp:   "1.2.3.4",
			ResultCode: new(int32(codes.AlreadyExists)),
			Phase:      apiv2.AuditPhase_AUDIT_PHASE_RESPONSE,
		}
	}
)

func Test_AuditCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.ImageServiceListResponse, *apiv2.Image]{
		{
			Name:    "list",
			CmdArgs: []string{"audit", "list", "--tenant", "a"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.AuditServiceListRequest{
							Login: "a",
							Query: &apiv2.AuditQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.AuditServiceListResponse{
								Traces: []*apiv2.AuditTrace{
									Trace1(),
									Trace2(),
								},
							})
						},
					},
				},
			},
			),
			WantTable: new(`
            TIME                 REQUEST ID                            USER  PROJECT                               METHOD        PHASE                 CODE
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/  AUDIT_PHASE_REQUEST
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/  AUDIT_PHASE_RESPONSE  AlreadyExists
			`),
			WantWideTable: new(`
            TIME                 REQUEST ID                            USER  PROJECT                               METHOD        PHASE                 SOURCE IP  CODE           BODY
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/  AUDIT_PHASE_REQUEST   1.2.3.4                   request body
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/  AUDIT_PHASE_RESPONSE  1.2.3.4    AlreadyExists  result body
			`),
			Template: new("{{ .uuid }} {{ .user }} {{ .phase }}"),
			WantTemplate: new(`
d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 me 1
d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 me 2
			`),
			WantMarkdown: new(`
            | TIME                | REQUEST ID                           | USER | PROJECT                              | METHOD       | PHASE                | CODE          |
            |---------------------|--------------------------------------|------|--------------------------------------|--------------|----------------------|---------------|
            | 2000-01-01 00:00:00 | d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 | me   | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | /apiv2.List/ | AUDIT_PHASE_REQUEST  |               |
            | 2000-01-01 00:00:00 | d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 | me   | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | /apiv2.List/ | AUDIT_PHASE_RESPONSE | AlreadyExists |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AuditCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.AuditServiceGetResponse, *apiv2.AuditTrace]{
		{
			Name:    "describe",
			CmdArgs: []string{"audit", "describe", "--tenant", "a", Trace1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.AuditServiceGetRequest{
							Login: "a",
							Uuid:  Trace1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.AuditServiceGetResponse{
								Trace: Trace1(),
							})
						},
					},
				},
			}),
			WantObject:      Trace1(),
			WantProtoObject: Trace1(),
			WantTable: new(`
            TIME                 REQUEST ID                            USER  PROJECT                               METHOD        PHASE                CODE
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/  AUDIT_PHASE_REQUEST
			`),
			WantWideTable: new(`
            TIME                 REQUEST ID                            USER  PROJECT                               METHOD        PHASE                SOURCE IP  CODE  BODY
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/  AUDIT_PHASE_REQUEST  1.2.3.4          request body
			`),
			Template: new("{{ .uuid }} {{ .user }} {{ .phase }}"),
			WantTemplate: new(`
d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 me 1
			`),
			WantMarkdown: new(`
            | TIME                | REQUEST ID                           | USER | PROJECT                              | METHOD       | PHASE               | CODE |
            |---------------------|--------------------------------------|------|--------------------------------------|--------------|---------------------|------|
            | 2000-01-01 00:00:00 | d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 | me   | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | /apiv2.List/ | AUDIT_PHASE_REQUEST |      |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
