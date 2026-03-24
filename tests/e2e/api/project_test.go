package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_ProjectCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceGetResponse, *apiv2.Project]{
		{
			Name:    "describe",
			CmdArgs: []string{"project", "describe", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
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
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
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
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
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
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
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
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
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
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
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
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
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
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
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
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
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
