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

func Test_SizeReservationCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.SizeReservationServiceListResponse, []*apiv2.SizeReservation]{
		{
			Name:    "list",
			CmdArgs: []string{"size-reservation", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.SizeReservationServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.SizeReservationServiceListResponse{
								SizeReservations: []*apiv2.SizeReservation{
									testresources.SizeReservation1(),
									testresources.SizeReservation2(),
								},
							})
						},
					},
				},
			}),
			WantObject: []*apiv2.SizeReservation{
				testresources.SizeReservation1(),
				testresources.SizeReservation2(),
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_SizeReservationCmd_Describe(t *testing.T) {
	sr1 := testresources.SizeReservation1()
	tests := []*e2e.Test[apiv2.SizeReservationServiceGetResponse, *apiv2.SizeReservation]{
		{
			Name:    "describe",
			CmdArgs: []string{"size-reservation", "describe", sr1.Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.SizeReservationServiceGetRequest{
							Id: sr1.Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.SizeReservationServiceGetResponse{
								SizeReservation: sr1,
							})
						},
					},
				},
			}),
			WantObject: sr1,
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
