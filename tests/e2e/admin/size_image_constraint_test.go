package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
)

func Test_AdminSizeImageConstraintCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeImageConstraintServiceListResponse, []*apiv2.SizeImageConstraint]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "size-image-constraint", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SizeImageConstraintServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SizeImageConstraintServiceListResponse{
								SizeImageConstraints: []*apiv2.SizeImageConstraint{
									{
										Size:        "v1-medium-x86",
										Name:        new("constraint-1"),
										Description: new("Must use Ubuntu 22.04"),
										Meta:        &apiv2.Meta{},
									},
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.SizeImageConstraint{
				{
					Size:        "v1-medium-x86",
					Name:        new("constraint-1"),
					Description: new("Must use Ubuntu 22.04"),
					Meta:        &apiv2.Meta{},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSizeImageConstraintCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeImageConstraintServiceDeleteResponse, *apiv2.SizeImageConstraint]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "size-image-constraint", "delete", "v1-medium-x86"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SizeImageConstraintServiceDeleteRequest{
							Size: "v1-medium-x86",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SizeImageConstraintServiceDeleteResponse{
								SizeImageConstraint: &apiv2.SizeImageConstraint{
									Size: "v1-medium-x86",
								},
							})
						},
					},
				},
			}),
			WantObject: &apiv2.SizeImageConstraint{Size: "v1-medium-x86"},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
