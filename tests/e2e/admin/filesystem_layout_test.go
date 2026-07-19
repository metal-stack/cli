package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
)

func Test_AdminFilesystemLayoutCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.FilesystemServiceListResponse, []*apiv2.FilesystemLayout]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "filesystem-layout", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.FilesystemServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.FilesystemServiceListResponse{
								FilesystemLayouts: []*apiv2.FilesystemLayout{
									testresources.FilesystemLayout1(),
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.FilesystemLayout{testresources.FilesystemLayout1()},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminFilesystemLayoutCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.FilesystemServiceDeleteResponse, *apiv2.FilesystemLayout]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "filesystem-layout", "delete", testresources.FilesystemLayout1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.FilesystemServiceDeleteRequest{
							Id: testresources.FilesystemLayout1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.FilesystemServiceDeleteResponse{
								FilesystemLayout: testresources.FilesystemLayout1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.FilesystemLayout1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
