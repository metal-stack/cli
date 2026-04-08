package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
)

func Test_AdminAuditCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.AuditServiceListResponse, *apiv2.AuditTrace]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "audit", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.AuditServiceListRequest{
							Query: &apiv2.AuditQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.AuditServiceListResponse{
								Traces: []*apiv2.AuditTrace{
									testresources.Trace1(),
									testresources.Trace2(),
									testresources.Trace3(),
								},
							})
						},
					},
				},
			},
			),
			WantTable: new(`
            TIME                 REQUEST ID                            USER   PROJECT                               METHOD          PHASE                 CODE           
            2001-01-01 00:00:00  5091c4e9-e8db-483c-ab6b-fe14f82570a7  Larry  0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/    AUDIT_PHASE_REQUEST                  
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me     0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.Create/  AUDIT_PHASE_REQUEST                  
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me     0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/    AUDIT_PHASE_RESPONSE  AlreadyExists
			`),
			WantWideTable: new(`
            TIME                 REQUEST ID                            USER   PROJECT                               METHOD          PHASE                 SOURCE IP  CODE           BODY          
            2001-01-01 00:00:00  5091c4e9-e8db-483c-ab6b-fe14f82570a7  Larry  0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/    AUDIT_PHASE_REQUEST   1.2.3.4                   result body   
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me     0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.Create/  AUDIT_PHASE_REQUEST   1.2.3.4                   request body  
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me     0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/    AUDIT_PHASE_RESPONSE  1.2.3.4    AlreadyExists  result body
			`),
			Template: new("{{ .uuid }} {{ .user }} {{ .phase }}"),
			WantTemplate: new(`5091c4e9-e8db-483c-ab6b-fe14f82570a7 Larry 1
d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 me 1
d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 me 2`),
			WantMarkdown: new(`
            | TIME                | REQUEST ID                           | USER  | PROJECT                              | METHOD         | PHASE                | CODE          |
            |---------------------|--------------------------------------|-------|--------------------------------------|----------------|----------------------|---------------|
            | 2001-01-01 00:00:00 | 5091c4e9-e8db-483c-ab6b-fe14f82570a7 | Larry | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | /apiv2.List/   | AUDIT_PHASE_REQUEST  |               |
            | 2000-01-01 00:00:00 | d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 | me    | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | /apiv2.Create/ | AUDIT_PHASE_REQUEST  |               |
            | 2000-01-01 00:00:00 | d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 | me    | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | /apiv2.List/   | AUDIT_PHASE_RESPONSE | AlreadyExists |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AuditCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv2.AuditServiceGetResponse, *apiv2.AuditTrace]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "audit", "describe", testresources.Trace1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.AuditServiceGetRequest{
							Uuid: testresources.Trace1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.AuditServiceGetResponse{
								Trace: testresources.Trace1(),
							})
						},
					},
				},
			}),
			WantObject:      testresources.Trace1(),
			WantProtoObject: testresources.Trace1(),
			WantTable: new(`
            TIME                 REQUEST ID                            USER  PROJECT                               METHOD          PHASE                CODE
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.Create/  AUDIT_PHASE_REQUEST
			`),
			WantWideTable: new(`
            TIME                 REQUEST ID                            USER  PROJECT                               METHOD          PHASE                SOURCE IP  CODE  BODY
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.Create/  AUDIT_PHASE_REQUEST  1.2.3.4          request body
			`),
			Template: new("{{ .uuid }} {{ .user }} {{ .phase }}"),
			WantTemplate: new(`
d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 me 1
			`),
			WantMarkdown: new(`
            | TIME                | REQUEST ID                           | USER | PROJECT                              | METHOD         | PHASE               | CODE |
            |---------------------|--------------------------------------|------|--------------------------------------|----------------|---------------------|------|
            | 2000-01-01 00:00:00 | d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 | me   | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | /apiv2.Create/ | AUDIT_PHASE_REQUEST |      |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
