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

func Test_AdminTokenCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.TokenServiceListResponse, apiv2.Token]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "token", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.TokenServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TokenServiceListResponse{
								Tokens: []*apiv2.Token{
									testresources.Token1(),
									testresources.Token2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)  
            TOKEN_TYPE_API  b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01         dev@metal-stack.io    dev token    0      2      2000-01-03 00:00:00 UTC (in 2d)
			`),
			WantWideTable: new(`
            TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)  
            TOKEN_TYPE_API  b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01         dev@metal-stack.io    dev token    0      2      2000-01-03 00:00:00 UTC (in 2d)
			`),
			Template: new("{{ .uuid }} {{ .description }}"),
			WantTemplate: new(`
a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 ci token
b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01 dev token
			`),
			WantMarkdown: new(`
            | TYPE           | ID                                   | ADMIN | USER                 | DESCRIPTION | ROLES | PERMS | EXPIRES                         |
            |----------------|--------------------------------------|-------|----------------------|-------------|-------|-------|---------------------------------|
            | TOKEN_TYPE_API | a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 |       | admin@metal-stack.io | ci token    | 0     | 0     | 2000-01-02 00:00:00 UTC (in 1d) |
            | TOKEN_TYPE_API | b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01 |       | dev@metal-stack.io   | dev token   | 0     | 2     | 2000-01-03 00:00:00 UTC (in 2d) |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminTokenCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.TokenServiceRevokeResponse, *apiv2.Token]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "token", "delete", testresources.Token1().Uuid, "--user", "user-123"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.TokenServiceRevokeRequest{
							Uuid: testresources.Token1().Uuid,
							User: "user-123",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.TokenServiceRevokeResponse{})
						},
					},
				},
			}),
			WantObject: &apiv2.Token{
				Uuid: testresources.Token1().Uuid,
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
