package api_e2e

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
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
)

// func Test_TokenCmd_Describe(t *testing.T) {
// 	tk := token1()

// 	tests := []*e2e.Test[apiv2.TokenServiceGetResponse, *apiv2.Token]{
// 		{
// 			Name:    "describe",
// 			CmdArgs: []string{"token", "describe", tk.Uuid},
// 			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestClientConfig[apiv2.TokenServiceGetRequest, apiv2.TokenServiceGetResponse]{
// 				WantRequest: apiv2.TokenServiceGetRequest{
// 					Uuid: tk.Uuid,
// 				},
// 				WantResponse: apiv2.TokenServiceGetResponse{
// 					Token: tk,
// 				},
// 			}),
// 			WantTable: new(`
//             TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES
//             TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
// 			`),
// 			WantWideTable: new(`
//             TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES
//             TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
// 			`),
// 			WantMarkdown: new(`
//             | TYPE           | ID                                   | ADMIN | USER                 | DESCRIPTION | ROLES | PERMS | EXPIRES                         |
//             |----------------|--------------------------------------|-------|----------------------|-------------|-------|-------|---------------------------------|
//             | TOKEN_TYPE_API | a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 |       | admin@metal-stack.io | ci token    | 0     | 0     | 2000-01-02 00:00:00 UTC (in 1d) |
// 			`),
// 			WantObject:      tk,
// 			WantProtoObject: tk,
// 			Template:        new("{{ .uuid }} {{ .description }}"),
// 			WantTemplate: new(`
// 			a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 ci token
// 			`),
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt.TestCmd(t)
// 	}
// }
