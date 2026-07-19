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

func Test_AdminSizeReservationCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeReservationServiceListResponse, []*apiv2.SizeReservation]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "size-reservation", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SizeReservationServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SizeReservationServiceListResponse{
								SizeReservations: []*apiv2.SizeReservation{
									testresources.SizeReservation1(),
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.SizeReservation{testresources.SizeReservation1()},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSizeReservationCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeReservationServiceDeleteResponse, *apiv2.SizeReservation]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "size-reservation", "delete", testresources.SizeReservation1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SizeReservationServiceDeleteRequest{
							Id: testresources.SizeReservation1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SizeReservationServiceDeleteResponse{
								SizeReservation: testresources.SizeReservation1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.SizeReservation1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
