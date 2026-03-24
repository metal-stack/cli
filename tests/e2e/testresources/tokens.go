package testresources

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Token1 = func() *apiv2.Token {
		return &apiv2.Token{
			Uuid:        "a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90",
			User:        "admin@metal-stack.io",
			Description: "ci token",
			TokenType:   apiv2.TokenType_TOKEN_TYPE_API,
			Expires:     timestamppb.New(e2e.TimeBubbleStartTime().Add(24 * time.Hour)),
			IssuedAt:    timestamppb.New(e2e.TimeBubbleStartTime()),
			Permissions: nil,
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		}
	}
	Token2 = func() *apiv2.Token {
		return &apiv2.Token{
			Uuid:        "b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01",
			User:        "dev@metal-stack.io",
			Description: "dev token",
			TokenType:   apiv2.TokenType_TOKEN_TYPE_API,
			Expires:     timestamppb.New(e2e.TimeBubbleStartTime().Add(48 * time.Hour)),
			IssuedAt:    timestamppb.New(e2e.TimeBubbleStartTime()),
			Permissions: nil,
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		}
	}
)
