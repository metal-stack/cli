package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
)

func Test_ProjectMemberCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectMemberServiceListResponse, []apiv2.ProjectMember]{
		{
			Name:    "list project members",
			CmdArgs: []string{"project", "member", "list", "--project", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectMemberServiceListRequest{
							Project: testresources.Project1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectMemberServiceListResponse{
								Members: []*apiv2.ProjectMember{
									testresources.Project1Member(),
									testresources.Project2Member(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    ROLE                 INHERITED  SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  PROJECT_ROLE_OWNER   false      now
            40c0da4b-9eb9-4371-91aa-1ae62193fa54  PROJECT_ROLE_EDITOR  true       now
			`),
			WantWideTable: new(`
            ID                                    ROLE                 INHERITED  SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  PROJECT_ROLE_OWNER   false      now
            40c0da4b-9eb9-4371-91aa-1ae62193fa54  PROJECT_ROLE_EDITOR  true       now
			`),
			Template: new("{{ .member }} {{ .role }}"),
			WantTemplate: new(`
16d6e8ba-f574-494f-8d5e-74f6cb2d8db0 1
40c0da4b-9eb9-4371-91aa-1ae62193fa54 2
			`),
			WantMarkdown: new(`
            | ID                                   | ROLE                | INHERITED | SINCE |
            |--------------------------------------|---------------------|-----------|-------|
            | 16d6e8ba-f574-494f-8d5e-74f6cb2d8db0 | PROJECT_ROLE_OWNER  | false     | now   |
            | 40c0da4b-9eb9-4371-91aa-1ae62193fa54 | PROJECT_ROLE_EDITOR | true      | now   |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectMemberCmd_Get(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectMemberServiceGetResponse, *apiv2.ProjectMember]{
		{
			Name:    "get project member",
			CmdArgs: []string{"project", "member", "describe", testresources.Project1Member().Member, "--project", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectMemberServiceGetRequest{
							Project: testresources.Project1().Uuid,
							Member:  testresources.Project1Member().Member,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectMemberServiceGetResponse{
								Member: testresources.Project1Member(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    ROLE                INHERITED  SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  PROJECT_ROLE_OWNER  false      now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectMemberCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectMemberServiceDeleteResponse, *apiv2.ProjectMember]{
		{
			Name:    "delete project member",
			CmdArgs: []string{"project", "member", "delete", testresources.Project1Member().Member, "--project", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectMemberServiceDeleteRequest{
							Project: testresources.Project1().Uuid,
							Member:  testresources.Project1Member().Member,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectMemberServiceDeleteResponse{
								Member: testresources.Project1Member(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    ROLE                INHERITED  SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  PROJECT_ROLE_OWNER  false      now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectMemberCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectMemberServiceUpdateResponse, *apiv2.ProjectMember]{
		{
			Name:    "update project member",
			CmdArgs: []string{"project", "member", "update", testresources.Project1Member().Member, "--project", testresources.Project1().Uuid, "--role", testresources.Project1Member().Role.String()},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectMemberServiceUpdateRequest{
							Project: testresources.Project1().Uuid,
							Member:  testresources.Project1Member().Member,
							Role:    new(testresources.Project1Member().Role),
							UpdateMeta: &apiv2.UpdateMeta{
								LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectMemberServiceUpdateResponse{
								Member: testresources.Project1Member(),
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID                                    ROLE                INHERITED  SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  PROJECT_ROLE_OWNER  false      now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectMemberCmd_Leave(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectMemberServiceLeaveResponse, *apiv2.ProjectMember]{
		{
			Name:    "get project member",
			CmdArgs: []string{"project", "member", "leave", "--project", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectMemberServiceLeaveRequest{
							Project: testresources.Project1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectMemberServiceLeaveResponse{
								Member: testresources.Project1Member(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    ROLE                INHERITED  SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  PROJECT_ROLE_OWNER  false      now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
