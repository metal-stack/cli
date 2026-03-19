package api_e2e

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/tests/e2e"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	token1 = func() *apiv2.Token {
		return &apiv2.Token{
			Uuid:        "a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90",
			User:        "admin@metal-stack.io",
			Description: "ci token",
			TokenType:   apiv2.TokenType_TOKEN_TYPE_API,
			Expires:     timestamppb.New(time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)),
			IssuedAt:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			Permissions: nil,
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
	token2 = func() *apiv2.Token {
		return &apiv2.Token{
			Uuid:        "b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01",
			User:        "dev@metal-stack.io",
			Description: "dev token",
			TokenType:   apiv2.TokenType_TOKEN_TYPE_API,
			Expires:     timestamppb.New(time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC)),
			IssuedAt:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			Permissions: nil,
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
)

func Test_TokenCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.TokenServiceGetResponse, *apiv2.Token]{
		{
			Name:    "describe",
			CmdArgs: []string{"token", "describe", token1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.TokenServiceGetRequest{
							Uuid: token1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TokenServiceGetResponse{
								Token: token1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES
            TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
			`),
			WantWideTable: new(`
            TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES
            TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
			`),
			WantMarkdown: new(`
            | TYPE           | ID                                   | ADMIN | USER                 | DESCRIPTION | ROLES | PERMS | EXPIRES                         |
            |----------------|--------------------------------------|-------|----------------------|-------------|-------|-------|---------------------------------|
            | TOKEN_TYPE_API | a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 |       | admin@metal-stack.io | ci token    | 0     | 0     | 2000-01-02 00:00:00 UTC (in 1d) |
			`),
			WantObject:      token1(),
			WantProtoObject: token1(),
			Template:        new("{{ .uuid }} {{ .description }}"),
			WantTemplate: new(`
			a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 ci token
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TokenCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.TokenServiceListResponse, apiv2.Token]{
		{
			Name:    "list",
			CmdArgs: []string{"token", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.TokenServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TokenServiceListResponse{
								Tokens: []*apiv2.Token{
									token1(),
									token2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
			TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES
			TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
			TOKEN_TYPE_API  b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01         dev@metal-stack.io    dev token    0      0      2000-01-03 00:00:00 UTC (in 2d)
			`),
			WantWideTable: new(`
			TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES
			TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
			TOKEN_TYPE_API  b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01         dev@metal-stack.io    dev token    0      0      2000-01-03 00:00:00 UTC (in 2d)
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
			| TOKEN_TYPE_API | b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01 |       | dev@metal-stack.io   | dev token   | 0     | 0     | 2000-01-03 00:00:00 UTC (in 2d) |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
