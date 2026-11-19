package admin_e2e

import (
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_AdminSizeReservationCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeReservationServiceListResponse, []*apiv2.SizeReservation]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "size-reservation", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SizeReservationServiceListRequest{
							Query: &apiv2.SizeReservationQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SizeReservationServiceListResponse{
								SizeReservations: []*apiv2.SizeReservation{
									testresources.SizeReservation2(),
									testresources.SizeReservation1(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    NAME           DESCRIPTION         PROJECT                               SIZE           PARTITIONS   AMOUNT
            a1b2c3d4-e5f6-7890-abcd-ef1234567890  reservation-a  first reservation   0d81bca7-73f6-4da3-8397-4a8c52a0c583  v1-medium-x86  partition-1  3
            b2c3d4e5-f6a7-8901-bcde-f12345678901  reservation-b  second reservation  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  g1-medium-x86  partition-2  1
			`),
			WantWideTable: new(`
            ID                                    NAME           DESCRIPTION         PROJECT                               SIZE           PARTITIONS   AMOUNT  LABELS
            a1b2c3d4-e5f6-7890-abcd-ef1234567890  reservation-a  first reservation   0d81bca7-73f6-4da3-8397-4a8c52a0c583  v1-medium-x86  partition-1  3       a=b
            b2c3d4e5-f6a7-8901-bcde-f12345678901  reservation-b  second reservation  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  g1-medium-x86  partition-2  1       c=d
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
a1b2c3d4-e5f6-7890-abcd-ef1234567890 reservation-a
b2c3d4e5-f6a7-8901-bcde-f12345678901 reservation-b
			`),
			WantMarkdown: new(`
            | ID                                   | NAME          | DESCRIPTION        | PROJECT                              | SIZE          | PARTITIONS  | AMOUNT |
            |--------------------------------------|---------------|--------------------|--------------------------------------|---------------|-------------|--------|
            | a1b2c3d4-e5f6-7890-abcd-ef1234567890 | reservation-a | first reservation  | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | v1-medium-x86 | partition-1 | 3      |
            | b2c3d4e5-f6a7-8901-bcde-f12345678901 | reservation-b | second reservation | f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | g1-medium-x86 | partition-2 | 1      |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSizeReservationCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeReservationServiceListResponse, *apiv2.SizeReservation]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "size-reservation", "describe", testresources.SizeReservation1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SizeReservationServiceListRequest{
							Query: &apiv2.SizeReservationQuery{
								Id: new(testresources.SizeReservation1().Id),
							},
						},
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

func Test_AdminSizeReservationCmd_Create(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeReservationServiceCreateResponse, *apiv2.SizeReservation]{
		{
			Name:    "create from file",
			CmdArgs: append([]string{"admin", "size-reservation", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.SizeReservation1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.SizeReservationServiceCreateRequest{
								SizeReservation: testresources.SizeReservation1(),
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.SizeReservationServiceCreateResponse{
									SizeReservation: testresources.SizeReservation1(),
								})
							},
						},
					},
				}),
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
		{
			Name:    "create from flags",
			CmdArgs: []string{"admin", "size-reservation", "create", "--name", "reservation-a", "--description", "first reservation", "--project", testresources.Project1().Uuid, "--size", testresources.Size1().Id, "--partitions", testresources.Partition1().Id, "--amount", "3", "--labels", "a=b"},
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.SizeReservationServiceCreateRequest{
								SizeReservation: &apiv2.SizeReservation{
									Name:        "reservation-a",
									Description: "first reservation",
									Project:     testresources.Project1().Uuid,
									Size:        testresources.Size1().Id,
									Partitions:  []string{testresources.Partition1().Id},
									Amount:      3,
									Meta: &apiv2.Meta{
										Labels: testresources.SizeReservation1().Meta.Labels,
									},
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.SizeReservationServiceCreateResponse{
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
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSizeReservationCmd_Update(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeReservationServiceUpdateResponse, *apiv2.SizeReservation]{
		{
			Name:    "update from file",
			CmdArgs: append([]string{"admin", "size-reservation", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.SizeReservation1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.SizeReservationServiceUpdateRequest{
								Id:          testresources.SizeReservation1().Id,
								Name:        new(testresources.SizeReservation1().Name),
								Description: new(testresources.SizeReservation1().Description),
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
								Partitions: testresources.SizeReservation1().Partitions,
								Amount:     new(testresources.SizeReservation1().Amount),
								Labels: &apiv2.UpdateLabels{
									Strategy: &apiv2.UpdateLabels_Replace{
										Replace: testresources.SizeReservation1().Meta.Labels,
									},
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.SizeReservationServiceUpdateResponse{
									SizeReservation: testresources.SizeReservation1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
            ID                                    NAME           DESCRIPTION        PROJECT                               SIZE           PARTITIONS   AMOUNT
            a1b2c3d4-e5f6-7890-abcd-ef1234567890  reservation-a  first reservation  0d81bca7-73f6-4da3-8397-4a8c52a0c583  v1-medium-x86  partition-1  3
			`),
		},
		{
			Name:    "update from flags",
			CmdArgs: []string{"admin", "size-reservation", "update", testresources.SizeReservation1().Id, "--description", "updated reservation", "--amount", "5", "--partitions", testresources.Partition1().Id, "--add-labels", "x=y"},
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.SizeReservationServiceUpdateRequest{
								Id: testresources.SizeReservation1().Id,
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
								Description: new("updated reservation"),
								Amount:      new(int32(5)),
								Partitions:  []string{testresources.Partition1().Id},
								Labels: &apiv2.UpdateLabels{
									Strategy: &apiv2.UpdateLabels_Patch{
										Patch: &apiv2.LabelsPatch{
											Update: &apiv2.Labels{
												Labels: map[string]string{
													"x": "y",
												},
											},
										},
									},
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.SizeReservationServiceUpdateResponse{
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
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSizeReservationCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeReservationServiceUpdateResponse, *apiv2.SizeReservation]{
		{
			Name:    "apply",
			CmdArgs: append([]string{"admin", "size-reservation", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.SizeReservation1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.SizeReservationServiceCreateRequest{
								SizeReservation: testresources.SizeReservation1(),
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.SizeReservationServiceCreateResponse{
									SizeReservation: testresources.SizeReservation1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
            ID                                    NAME           DESCRIPTION        PROJECT                               SIZE           PARTITIONS   AMOUNT
            a1b2c3d4-e5f6-7890-abcd-ef1234567890  reservation-a  first reservation  0d81bca7-73f6-4da3-8397-4a8c52a0c583  v1-medium-x86  partition-1  3
			`),
		},
		{
			Name:    "apply already exists",
			CmdArgs: append([]string{"admin", "size-reservation", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.SizeReservation1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.SizeReservationServiceCreateRequest{
								SizeReservation: testresources.SizeReservation1(),
							},
							WantError: connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")),
						},
						{
							WantRequest: &adminv2.SizeReservationServiceUpdateRequest{
								Id:          testresources.SizeReservation1().Id,
								Name:        new(testresources.SizeReservation1().Name),
								Description: new(testresources.SizeReservation1().Description),
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
								Partitions: testresources.SizeReservation1().Partitions,
								Amount:     new(testresources.SizeReservation1().Amount),
								Labels: &apiv2.UpdateLabels{
									Strategy: &apiv2.UpdateLabels_Replace{
										Replace: testresources.SizeReservation1().Meta.Labels,
									},
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.SizeReservationServiceUpdateResponse{
									SizeReservation: testresources.SizeReservation1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
            ID                                    NAME           DESCRIPTION        PROJECT                               SIZE           PARTITIONS   AMOUNT
            a1b2c3d4-e5f6-7890-abcd-ef1234567890  reservation-a  first reservation  0d81bca7-73f6-4da3-8397-4a8c52a0c583  v1-medium-x86  partition-1  3
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
