package cmd

import (
	"testing"
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
			Expires:     timestamppb.New(time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)),
			IssuedAt:    timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
			Permissions: nil,
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
)

func Test_TokenCmd_SingleResult(t *testing.T) {
	tk := token1()

	tests := []*Test[apiv2.TokenServiceGetRequest, apiv2.TokenServiceGetResponse]{
		{
			Name: "describe",
			Cmd: func() []string {
				return []string{"token", "describe", tk.Uuid}
			},
			WantRequest: apiv2.TokenServiceGetRequest{
				Uuid: tk.Uuid,
			},
			WantResponse: apiv2.TokenServiceGetResponse{
				Token: tk,
			},
			WantObject: tk,
			Template:   new("{{ .uuid }} {{ .description }}"),
			WantTemplate: new(`
a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 ci token
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
