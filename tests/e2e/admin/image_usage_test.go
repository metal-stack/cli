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

func Test_AdminImageUsageCmd(t *testing.T) {
	tests := []*e2e.Test[adminv2.ImageServiceUsageResponse, []*apiv2.ImageUsage]{
		{
			Name:    "usage",
			CmdArgs: []string{"admin", "image-usage"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.ImageServiceUsageRequest{
							Query: &apiv2.ImageQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ImageServiceUsageResponse{
								ImageUsage: []*apiv2.ImageUsage{
									{
										Image: &apiv2.Image{
											Id: "ubuntu-22.04",
										},
										UsedBy: []string{"machine-1", "machine-2"},
									},
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.ImageUsage{
				{
					Image: &apiv2.Image{
						Id: "ubuntu-22.04",
					},
					UsedBy: []string{"machine-1", "machine-2"},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
