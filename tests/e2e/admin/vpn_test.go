package admin_e2e

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_AdminVPNCmd_AuthKey(t *testing.T) {
	tests := []*e2e.Test[adminv2.VPNServiceAuthKeyResponse, *adminv2.VPNServiceAuthKeyResponse]{
		{
			Name:    "auth key",
			CmdArgs: []string{"admin", "vpn", "auth-key", "--project", "project-1", "--ephemeral", "--expires", "1h", "--reason", "debugging"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.VPNServiceAuthKeyRequest{
							Project:   "project-1",
							Ephemeral: true,
							Expires:   durationpb.New(1 * time.Hour),
							Reason:    "debugging",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.VPNServiceAuthKeyResponse{
								Address:   "vpn.example.com:443",
								AuthKey:   "key-12345",
								Ephemeral: true,
								ExpiresAt: timestamppb.New(e2e.TimeBubbleStartTime().Add(1 * time.Hour)),
								CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
							})
						},
					},
				},
			}),
			WantObject: &adminv2.VPNServiceAuthKeyResponse{
				Address:   "vpn.example.com:443",
				AuthKey:   "key-12345",
				Ephemeral: true,
				ExpiresAt: timestamppb.New(e2e.TimeBubbleStartTime().Add(1 * time.Hour)),
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminVPNCmd_ListNodes(t *testing.T) {
	tests := []*e2e.Test[adminv2.VPNServiceListNodesResponse, []*apiv2.VPNNode]{
		{
			Name:    "list nodes",
			CmdArgs: []string{"admin", "vpn", "list-nodes"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.VPNServiceListNodesRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.VPNServiceListNodesResponse{
								Nodes: []*apiv2.VPNNode{
									{
										Id:          1,
										Name:        "node-1",
										Project:     "project-1",
										Online:      true,
										LastSeen:    timestamppb.New(e2e.TimeBubbleStartTime()),
										IpAddresses: []string{"10.0.0.1"},
									},
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.VPNNode{
				{
					Id:          1,
					Name:        "node-1",
					Project:     "project-1",
					Online:      true,
					LastSeen:    timestamppb.New(e2e.TimeBubbleStartTime()),
					IpAddresses: []string{"10.0.0.1"},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
