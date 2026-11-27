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
			CmdArgs: []string{"size", "reservation", "list", "--project", testresources.SizeReservation1().Project},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.SizeReservationServiceListRequest{
							Project: testresources.SizeReservation1().Project,
							Query:   &apiv2.SizeReservationQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.SizeReservationServiceListResponse{
								SizeReservations: []*apiv2.SizeReservation{
									testresources.SizeReservation1(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    NAME           DESCRIPTION        PROJECT                               SIZE           PARTITIONS   AMOUNT
            a1b2c3d4-e5f6-7890-abcd-ef1234567890  reservation-a  first reservation  0d81bca7-73f6-4da3-8397-4a8c52a0c583  v1-medium-x86  partition-1  3
			`),
			WantWideTable: new(`
            ID                                    NAME           DESCRIPTION        PROJECT                               SIZE           PARTITIONS   AMOUNT  LABELS
            a1b2c3d4-e5f6-7890-abcd-ef1234567890  reservation-a  first reservation  0d81bca7-73f6-4da3-8397-4a8c52a0c583  v1-medium-x86  partition-1  3       a=b
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
a1b2c3d4-e5f6-7890-abcd-ef1234567890 reservation-a
			`),
			WantMarkdown: new(`
            | ID                                   | NAME          | DESCRIPTION       | PROJECT                              | SIZE          | PARTITIONS  | AMOUNT |
            |--------------------------------------|---------------|-------------------|--------------------------------------|---------------|-------------|--------|
            | a1b2c3d4-e5f6-7890-abcd-ef1234567890 | reservation-a | first reservation | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | v1-medium-x86 | partition-1 | 3      |
			`),
		},
		{
			Name: "list with filters",
			CmdArgs: []string{"size", "reservation", "list",
				"--project", testresources.SizeReservation1().Project,
				"--id", testresources.SizeReservation1().Id,
				"--name", testresources.SizeReservation1().Name,
				"--description", testresources.SizeReservation1().Description,
				"--size", testresources.SizeReservation1().Size,
				"--partition", testresources.SizeReservation1().Partitions[0],
				"--labels", "a=b",
			},
			AssertExhaustiveArgs:     true,
			AssertExhaustiveExcludes: []string{"sort-by"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.SizeReservationServiceListRequest{
							Project: testresources.SizeReservation1().Project,
							Query: &apiv2.SizeReservationQuery{
								Id:          new(testresources.SizeReservation1().Id),
								Name:        new(testresources.SizeReservation1().Name),
								Description: new(testresources.SizeReservation1().Description),
								Size:        new(testresources.SizeReservation1().Size),
								Partition:   new(testresources.SizeReservation1().Partitions[0]),
								Labels: &apiv2.Labels{
									Labels: map[string]string{"a": "b"},
								},
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.SizeReservationServiceListResponse{
								SizeReservations: []*apiv2.SizeReservation{
									testresources.SizeReservation1(),
								},
							})
						},
					},
				},
			}),
			WantTemplate: new(`
a1b2c3d4-e5f6-7890-abcd-ef1234567890 reservation-a
			`),
			Template: new("{{ .id }} {{ .name }}"),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_SizeReservationCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.SizeReservationServiceGetResponse, *apiv2.SizeReservation]{
		{
			Name:    "describe",
			CmdArgs: []string{"size", "reservation", "describe", testresources.SizeReservation1().Id, "--project", testresources.SizeReservation1().Project},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.SizeReservationServiceGetRequest{
							Id:      testresources.SizeReservation1().Id,
							Project: testresources.SizeReservation1().Project,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.SizeReservationServiceGetResponse{
								SizeReservation: testresources.SizeReservation1(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.SizeReservation1(),
			WantTable: new(`
            ID                                    NAME           DESCRIPTION        PROJECT                               SIZE           PARTITIONS   AMOUNT
            a1b2c3d4-e5f6-7890-abcd-ef1234567890  reservation-a  first reservation  0d81bca7-73f6-4da3-8397-4a8c52a0c583  v1-medium-x86  partition-1  3
			`),
			WantMarkdown: new(`
            | ID                                   | NAME          | DESCRIPTION       | PROJECT                              | SIZE          | PARTITIONS  | AMOUNT |
            |--------------------------------------|---------------|-------------------|--------------------------------------|---------------|-------------|--------|
            | a1b2c3d4-e5f6-7890-abcd-ef1234567890 | reservation-a | first reservation | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | v1-medium-x86 | partition-1 | 3      |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
