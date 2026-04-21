package api_e2e

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_AuditCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.ImageServiceListResponse, *apiv2.Image]{
		{
			Name:    "list",
			CmdArgs: []string{"audit", "list", "--tenant", testresources.Trace1().Tenant},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.AuditServiceListRequest{
							Login: testresources.Trace1().Tenant,
							Query: &apiv2.AuditQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.AuditServiceListResponse{
								Traces: []*apiv2.AuditTrace{
									testresources.Trace1(),
									testresources.Trace2(),
								},
							})
						},
					},
				},
			},
			),
			WantTable: new(`
            TIME                 REQUEST ID                            USER  PROJECT                               METHOD          PHASE                 CODE           
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.Create/  AUDIT_PHASE_REQUEST                  
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/    AUDIT_PHASE_RESPONSE  AlreadyExists
			`),
			WantWideTable: new(`
            TIME                 REQUEST ID                            USER  PROJECT                               METHOD          PHASE                 SOURCE IP  CODE           BODY          
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.Create/  AUDIT_PHASE_REQUEST   1.2.3.4                   request body  
            2000-01-01 00:00:00  d1ff7267-2fbb-4a63-a7c1-44f1a83381a7  me    0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/    AUDIT_PHASE_RESPONSE  1.2.3.4    AlreadyExists  result body
			`),
			Template: new("{{ .uuid }} {{ .user }} {{ .phase }}"),
			WantTemplate: new(`d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 me 1
d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 me 2`),
			WantMarkdown: new(`
            | TIME                | REQUEST ID                           | USER | PROJECT                              | METHOD         | PHASE                | CODE          |
            |---------------------|--------------------------------------|------|--------------------------------------|----------------|----------------------|---------------|
            | 2000-01-01 00:00:00 | d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 | me   | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | /apiv2.Create/ | AUDIT_PHASE_REQUEST  |               |
            | 2000-01-01 00:00:00 | d1ff7267-2fbb-4a63-a7c1-44f1a83381a7 | me   | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | /apiv2.List/   | AUDIT_PHASE_RESPONSE | AlreadyExists |
			`),
		},
		{
			Name:    "list",
			CmdArgs: []string{"audit", "list", "--tenant", testresources.Trace3().Tenant, "--from", timestamppb.New(time.Date(2000, 12, 24, 0, 0, 0, 0, time.UTC)).AsTime().Format("2006-01-02 15:04:05"), "--to", timestamppb.New(time.Date(2001, 1, 2, 0, 0, 0, 0, time.UTC)).AsTime().Format("2006-01-02 15:04:05"), "--user", testresources.Trace3().User},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.AuditServiceListRequest{
							Login: testresources.Trace3().Tenant,
							Query: &apiv2.AuditQuery{
								User: new(testresources.Trace3().User),
								From: timestamppb.New(time.Date(2000, 12, 24, 0, 0, 0, 0, time.UTC)),
								To:   timestamppb.New(time.Date(2001, 1, 2, 0, 0, 0, 0, time.UTC)),
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.AuditServiceListResponse{
								Traces: []*apiv2.AuditTrace{
									testresources.Trace3(),
								},
							})
						},
					},
				},
			},
			),
			WantTable: new(`
            TIME                 REQUEST ID                            USER   PROJECT                               METHOD        PHASE                CODE  
            2001-01-01 00:00:00  5091c4e9-e8db-483c-ab6b-fe14f82570a7  Larry  0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/  AUDIT_PHASE_REQUEST
			`),
			WantWideTable: new(`
            TIME                 REQUEST ID                            USER   PROJECT                               METHOD        PHASE                SOURCE IP  CODE  BODY         
            2001-01-01 00:00:00  5091c4e9-e8db-483c-ab6b-fe14f82570a7  Larry  0d81bca7-73f6-4da3-8397-4a8c52a0c583  /apiv2.List/  AUDIT_PHASE_REQUEST  1.2.3.4          result body
			`),
			Template: new("{{ .uuid }} {{ .user }} {{ .phase }}"),
			WantTemplate: new(`
			5091c4e9-e8db-483c-ab6b-fe14f82570a7 Larry 1
			`),
			WantMarkdown: new(`
            | TIME                | REQUEST ID                           | USER  | PROJECT                              | METHOD       | PHASE               | CODE |
            |---------------------|--------------------------------------|-------|--------------------------------------|--------------|---------------------|------|
            | 2001-01-01 00:00:00 | 5091c4e9-e8db-483c-ab6b-fe14f82570a7 | Larry | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | /apiv2.List/ | AUDIT_PHASE_REQUEST |      |
			`),
		},
		{
			Name:    "list",
			CmdArgs: []string{"audit", "list", "--tenant", "notExisting"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.AuditServiceListRequest{
							Login: "notExisting",
							Query: &apiv2.AuditQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.AuditServiceListResponse{
								Traces: []*apiv2.AuditTrace{},
							})
						},
					},
				},
			},
			),
			WantTable: new(`
        	TIME  REQUEST ID  USER  PROJECT  METHOD  PHASE  CODE
			`),
			WantWideTable: new(`
        	TIME  REQUEST ID  USER  PROJECT  METHOD  PHASE  SOURCE IP  CODE  BODY
			`),
			Template:     new("{{ .uuid }} {{ .user }} {{ .phase }}"),
			WantTemplate: new(``),
			WantMarkdown: new(`
        	| TIME | REQUEST ID | USER | PROJECT | METHOD | PHASE | CODE |
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
			CmdArgs: []string{"audit", "describe", "--tenant", testresources.Trace1().Tenant, testresources.Trace1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.AuditServiceGetRequest{
							Login: testresources.Trace1().Tenant,
							Uuid:  testresources.Trace1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.AuditServiceGetResponse{
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
