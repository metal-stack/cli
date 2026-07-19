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

func Test_PartitionCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.PartitionServiceListResponse, []*apiv2.Partition]{
		{
			Name:    "list",
			CmdArgs: []string{"partition", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.PartitionServiceListRequest{
							Query: &apiv2.PartitionQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.PartitionServiceListResponse{
								Partitions: []*apiv2.Partition{
									testresources.Partition1(),
									testresources.Partition2(),
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.Partition{
				testresources.Partition1(),
				testresources.Partition2(),
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_PartitionCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.PartitionServiceGetResponse, *apiv2.Partition]{
		{
			Name:    "describe",
			CmdArgs: []string{"partition", "describe", testresources.Partition1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.PartitionServiceGetRequest{
							Id: testresources.Partition1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.PartitionServiceGetResponse{
								Partition: testresources.Partition1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.Partition1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
