package admin_e2e

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_AdminVpnCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.VPNServiceListNodesResponse, apiv2.VPNNode]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "vpn", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.VPNServiceListNodesRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.VPNServiceListNodesResponse{
								Nodes: []*apiv2.VPNNode{
									testresources.VpnNode2(),
									testresources.VpnNode1(),
								},
							})
						},
					},
				},
			}),
			Template: new("{{ .id }} {{ .project }}"),
			WantTemplate: new(`
1 0d81bca7-73f6-4da3-8397-4a8c52a0c583
2 f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c
			`),
			WantTable: new(`
            ID  NAME    PROJECT                               IPS              LAST SEEN
            1   node-1  0d81bca7-73f6-4da3-8397-4a8c52a0c583  1.2.3.4,1.2.3.5  2000-01-01 00:00:00
            2   node-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  2.3.4.5          2000-01-01 00:01:00
			`),
			WantMarkdown: new(`
            | ID | NAME   | PROJECT                              | IPS             | LAST SEEN           |
            |----|--------|--------------------------------------|-----------------|---------------------|
            | 1  | node-1 | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | 1.2.3.4,1.2.3.5 | 2000-01-01 00:00:00 |
            | 2  | node-2 | f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | 2.3.4.5         | 2000-01-01 00:01:00 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminVpnCmd_AuthKey(t *testing.T) {
	tests := []*e2e.Test[any, any]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "vpn", "auth-key", "--project", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.VPNServiceAuthKeyRequest{
							Project:   testresources.Project1().Uuid,
							Ephemeral: true,
							Expires:   durationpb.New(1 * time.Hour),
							Reason:    "",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.VPNServiceAuthKeyResponse{
								Address:   "1.2.3.4",
								AuthKey:   "abcdefghikl",
								Ephemeral: true,
								ExpiresAt: timestamppb.New(e2e.TimeBubbleStartTime().Add(1 * time.Hour)),
								CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
							})
						},
					},
				},
			}),
			WantDefault: new(`
auth-key: abcdefghikl
vpn-endpoint: 1.2.3.4
ephemeral: true
expires in: 1h0m0s
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
