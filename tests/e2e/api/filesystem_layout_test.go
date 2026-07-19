package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
)

func Test_FilesystemLayoutCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.FilesystemServiceListResponse, []*apiv2.FilesystemLayout]{
		{
			Name:    "list",
			CmdArgs: []string{"filesystem-layout", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.FilesystemServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.FilesystemServiceListResponse{
								FilesystemLayouts: []*apiv2.FilesystemLayout{
									testresources.FilesystemLayout1(),
									testresources.FilesystemLayout2(),
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.FilesystemLayout{
				testresources.FilesystemLayout1(),
				testresources.FilesystemLayout2(),
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_FilesystemLayoutCmd_Describe(t *testing.T) {
	fsl1 := testresources.FilesystemLayout1()
	tests := []*e2e.Test[apiv2.FilesystemServiceGetResponse, *apiv2.FilesystemLayout]{
		{
			Name:    "describe",
			CmdArgs: []string{"filesystem-layout", "describe", fsl1.Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.FilesystemServiceGetRequest{
							Id: fsl1.Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.FilesystemServiceGetResponse{
								FilesystemLayout: fsl1,
							})
						},
					},
				},
			}),
			WantObject: fsl1,
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
