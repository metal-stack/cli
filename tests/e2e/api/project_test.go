package api_e2e

import (
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/dustin/go-humanize"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_ProjectCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceGetResponse, *apiv2.Project]{
		{
			Name:    "describe",
			CmdArgs: []string{"project", "describe", testresources.Project1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceGetRequest{
							Project: testresources.Project1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceGetResponse{
								Project: testresources.Project1(),
							})
						},
					},
				},
			}),
			WantObject:      testresources.Project1(),
			WantProtoObject: testresources.Project1(),
			WantTable: new(`
            ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE
            0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2000-01-01 00:00:00 UTC
			`),
			WantWideTable: new(`
            ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE
            0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2000-01-01 00:00:00 UTC
			`),
			Template: new("{{ .uuid }} {{ .name }}"),
			WantTemplate: new(`
			0d81bca7-73f6-4da3-8397-4a8c52a0c583 project-a
			`),
			WantMarkdown: new(`
            | ID                                   | TENANT      | NAME      | DESCRIPTION   | CREATION DATE           |
            |--------------------------------------|-------------|-----------|---------------|-------------------------|
            | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | metal-stack | project-a | first project | 2000-01-01 00:00:00 UTC |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_Create(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceCreateResponse, *apiv2.Project]{
		{
			Name:    "create",
			CmdArgs: []string{"project", "create", "--name", testresources.Project1().Name, "--description", testresources.Project1().Description, "--tenant", testresources.Project1().Tenant},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceCreateRequest{
							Login:       testresources.Project1().Tenant,
							Name:        testresources.Project1().Name,
							Description: testresources.Project1().Description,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceCreateResponse{
								Project: testresources.Project1(),
							})
						},
					},
				},
			}),
			WantObject:      testresources.Project1(),
			WantProtoObject: testresources.Project1(),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"project", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Project1()), 0755))
				},
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceCreateRequest{
							Login:       testresources.Project1().Tenant,
							Name:        testresources.Project1().Name,
							Description: testresources.Project1().Description,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceCreateResponse{
								Project: testresources.Project1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE
            0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2000-01-01 00:00:00 UTC
			`),
		},
		{
			Name:    "create many from file",
			CmdArgs: append([]string{"project", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshalToMultiYAML(t, testresources.Project1(), testresources.Project2()), 0755))
				},
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceCreateRequest{
							Login:       testresources.Project1().Tenant,
							Name:        testresources.Project1().Name,
							Description: testresources.Project1().Description,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceCreateResponse{
								Project: testresources.Project1(),
							})
						},
					},
					{
						WantRequest: &apiv2.ProjectServiceCreateRequest{
							Login:       testresources.Project2().Tenant,
							Name:        testresources.Project2().Name,
							Description: testresources.Project2().Description,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceCreateResponse{
								Project: testresources.Project2(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    TENANT       NAME       DESCRIPTION     CREATION DATE
            0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project   2000-01-01 00:00:00 UTC
            f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  metal-stack  project-b  second project  2000-01-01 00:00:00 UTC
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceDeleteResponse, *apiv2.Project]{
		{
			Name:    "delete",
			CmdArgs: []string{"project", "delete", testresources.Project1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceDeleteRequest{
							Project: testresources.Project1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceDeleteResponse{
								Project: testresources.Project1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.Project1(),
		},
		{
			Name:    "delete from file",
			CmdArgs: append([]string{"project", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Project1()), 0755))
				},
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceDeleteRequest{
							Project: testresources.Project1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceDeleteResponse{
								Project: testresources.Project1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2000-01-01 00:00:00 UTC
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceUpdateResponse, *apiv2.Project]{
		{
			Name:    "update",
			CmdArgs: []string{"project", "update", testresources.Project1().Uuid, "--name", "new-name", "--description", "new-desc"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceUpdateRequest{
							Project:     testresources.Project1().Uuid,
							Name:        new("new-name"),
							Description: new("new-desc"),
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceUpdateResponse{
								Project: testresources.Project1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.Project1(),
		},
		{
			Name:    "update from file",
			CmdArgs: append([]string{"project", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Project1()), 0755))
				},
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceUpdateRequest{
							Project:     testresources.Project1().Uuid,
							Name:        new(testresources.Project1().Name),
							Description: new(testresources.Project1().Description),
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceUpdateResponse{
								Project: testresources.Project1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2000-01-01 00:00:00 UTC
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceListResponse, apiv2.Project]{
		{
			Name:    "list",
			CmdArgs: []string{"project", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceListResponse{
								Projects: []*apiv2.Project{
									testresources.Project1(),
									testresources.Project2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION     CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project   2000-01-01 00:00:00 UTC
			f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  metal-stack  project-b  second project  2000-01-01 00:00:00 UTC
			`),
			WantWideTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION     CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project   2000-01-01 00:00:00 UTC
			f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  metal-stack  project-b  second project  2000-01-01 00:00:00 UTC
			`),
			Template: new("{{ .uuid }} {{ .name }}"),
			WantTemplate: new(`
0d81bca7-73f6-4da3-8397-4a8c52a0c583 project-a
f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c project-b
			`),
			WantMarkdown: new(`
			| ID                                   | TENANT      | NAME      | DESCRIPTION    | CREATION DATE           |
			|--------------------------------------|-------------|-----------|----------------|-------------------------|
			| 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | metal-stack | project-a | first project  | 2000-01-01 00:00:00 UTC |
			| f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | metal-stack | project-b | second project | 2000-01-01 00:00:00 UTC |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceUpdateResponse, *apiv2.Project]{
		{
			Name:    "apply",
			CmdArgs: append([]string{"project", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Project1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.ProjectServiceCreateRequest{
								Login:       testresources.Project1().Tenant,
								Name:        testresources.Project1().Name,
								Description: testresources.Project1().Description,
								Labels:      testresources.Project1().Meta.Labels,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.ProjectServiceCreateResponse{
									Project: testresources.Project1(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE            
            0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2000-01-01 00:00:00 UTC
			`),
		},
		{
			Name:    "apply already exists",
			CmdArgs: append([]string{"project", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Project1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.ProjectServiceCreateRequest{
								Login:       testresources.Project1().Tenant,
								Name:        testresources.Project1().Name,
								Description: testresources.Project1().Description,
							},
							WantError: connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")),
						},
						{
							WantRequest: &apiv2.ProjectServiceUpdateRequest{
								Project:     testresources.Project1().Uuid,
								Description: &testresources.Project1().Description,
								Name:        &testresources.Project1().Name,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.ProjectServiceUpdateResponse{
									Project: testresources.Project1(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE            
            0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2000-01-01 00:00:00 UTC
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_ListInvites(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceInvitesListResponse, apiv2.ProjectInvite]{
		{
			Name:    "list invites",
			CmdArgs: []string{"project", "invite", "list", testresources.Project1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceInvitesListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceInvitesListResponse{
								Invites: []*apiv2.ProjectInvite{
									testresources.Project1Invite(),
									testresources.Project2Invite(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            SECRET  PROJECT                               ROLE                 EXPIRES IN       
            secret  0d81bca7-73f6-4da3-8397-4a8c52a0c583  PROJECT_ROLE_EDITOR  2 days from now  
            secret  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  PROJECT_ROLE_EDITOR  2 days from now
			`),
			WantWideTable: new(`
            SECRET  PROJECT                               ROLE                 EXPIRES IN       
            secret  0d81bca7-73f6-4da3-8397-4a8c52a0c583  PROJECT_ROLE_EDITOR  2 days from now  
            secret  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  PROJECT_ROLE_EDITOR  2 days from now
			`),
			Template: new("{{ .project }} {{ .role }}"),
			WantTemplate: new(`
0d81bca7-73f6-4da3-8397-4a8c52a0c583 2
f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c 2
			`),
			WantMarkdown: new(`
            | SECRET | PROJECT                              | ROLE                | EXPIRES IN      |
            |--------|--------------------------------------|---------------------|-----------------|
            | secret | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | PROJECT_ROLE_EDITOR | 2 days from now |
            | secret | f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | PROJECT_ROLE_EDITOR | 2 days from now |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_DeleteInvite(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceInviteDeleteResponse, string]{
		{
			Name:    "delete",
			CmdArgs: []string{"project", "invite", "delete", testresources.Project1Invite().Secret, "--project", testresources.Project1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceInviteDeleteRequest{
							Project: testresources.Project1().Uuid,
							Secret:  testresources.Project1Invite().Secret,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceInviteDeleteResponse{})
						},
					},
				},
			}),
			WantMarkdown: new(""),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_CreateInvite(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceInviteRequest, string]{
		{
			Name:    "create invite",
			CmdArgs: []string{"project", "invite", "generate-join-secret", "--role", testresources.Project1Invite().Role.String(), "--project", testresources.Project1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceInviteRequest{
							Project: testresources.Project1().Uuid,
							Role:    testresources.Project1Invite().Role,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceInviteResponse{
								Invite: testresources.Project1Invite(),
							})
						},
					},
				},
			}),
			WantDefault: new(fmt.Sprintf("You can share this secret with the member to join, it expires in %s:\n\n%s (https://console.metal-stack.io/project-invite/%s)",
				humanize.RelTime(e2e.TimeBubbleStartTime(), testresources.Project1Invite().ExpiresAt.AsTime(), "from now", "ago"),
				testresources.Project1Invite().Secret,
				testresources.Project1Invite().Secret,
			)),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_Join(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceInviteAcceptResponse, string]{
		{
			Name:    "join",
			CmdArgs: []string{"project", "invite", "join", testresources.Project1Invite().Secret},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceInviteGetRequest{
							Secret: testresources.Project1Invite().Secret,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceInviteGetResponse{
								Invite: testresources.Project1Invite(),
							})
						},
					},
					{
						WantRequest: &apiv2.ProjectServiceInviteAcceptRequest{
							Secret: testresources.Project1Invite().Secret,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceInviteAcceptResponse{
								Project:     testresources.Project1().Uuid,
								ProjectName: testresources.Project1Invite().ProjectName,
							})
						},
					},
				},
			}),
			WantDefault: new(fmt.Sprintf("Do you want to join project \"%s\" as %s? [Y/n] ✔ successfully joined project \"%s\"",
				testresources.Project1Invite().ProjectName,
				testresources.Project1Invite().Role.String(),
				testresources.Project1Invite().ProjectName)),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_ListMembers(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceGetRequest, []apiv2.ProjectMember]{
		{
			Name:    "list project members",
			CmdArgs: []string{"project", "member", "list", "--project", testresources.Project1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceGetRequest{
							Project: testresources.Project1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceGetResponse{
								Project: testresources.Project1(),
								ProjectMembers: []*apiv2.ProjectMember{
									testresources.Project1Members(), testresources.Project2Members(),
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
			Template: new("{{ .id }} {{ .role }}"),
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

func Test_ProjectCmd_DeleteMember(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceRemoveMemberResponse, string]{
		{
			Name:    "delete project member",
			CmdArgs: []string{"project", "member", "delete", testresources.Project1Members().Id, "--project", testresources.Project1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceRemoveMemberRequest{
							Project: testresources.Project1().Uuid,
							Member:  testresources.Project1Members().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceRemoveMemberResponse{})
						},
					},
				},
			}),
			WantDefault: new(fmt.Sprintf("✔ successfully removed member \"%s\"", testresources.Project1Members().Id)),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_UpdateMember(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceUpdateMemberResponse, *apiv2.ProjectMember]{
		{
			Name:    "update project member",
			CmdArgs: []string{"project", "member", "update", testresources.Project1Members().Id, "--project", testresources.Project1().Uuid, "--role", testresources.Project1Members().Role.String()},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.ProjectServiceUpdateMemberRequest{
							Project: testresources.Project1().Uuid,
							Member:  testresources.Project1Members().Id,
							Role:    testresources.Project1Members().Role,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceUpdateMemberResponse{
								ProjectMember: testresources.Project1Members(),
							})
						},
					},
				},
			}),
			WantObject: testresources.Project1Members(),
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
