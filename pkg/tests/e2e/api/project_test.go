package api_e2e

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/tests/e2e"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	project1 = func() *apiv2.Project {
		return &apiv2.Project{
			Uuid:        "0d81bca7-73f6-4da3-8397-4a8c52a0c583",
			Name:        "project-a",
			Description: "first project",
			Tenant:      "metal-stack",
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC)),
			},
		}
	}
	project2 = func() *apiv2.Project {
		return &apiv2.Project{
			Uuid:        "f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c",
			Name:        "project-b",
			Description: "second project",
			Tenant:      "metal-stack",
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(time.Date(2025, 7, 15, 14, 30, 0, 0, time.UTC)),
			},
		}
	}
)

func Test_ProjectCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.ProjectServiceGetResponse, *apiv2.Project]{
		{
			Name:    "describe",
			CmdArgs: []string{"project", "describe", project1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.ProjectServiceGetRequest{
							Project: project1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceGetResponse{
								Project: project1(),
							})
						},
					},
				},
			}),
			WantObject:      project1(),
			WantProtoObject: project1(),
			WantTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2025-06-01 10:00:00 UTC
			`),
			WantWideTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2025-06-01 10:00:00 UTC
			`),
			Template: new("{{ .uuid }} {{ .name }}"),
			WantTemplate: new(`
			0d81bca7-73f6-4da3-8397-4a8c52a0c583 project-a
			`),
			WantMarkdown: new(`
			| ID                                   | TENANT      | NAME      | DESCRIPTION   | CREATION DATE           |
			|--------------------------------------|-------------|-----------|---------------|-------------------------|
			| 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | metal-stack | project-a | first project | 2025-06-01 10:00:00 UTC |
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
			CmdArgs: []string{"project", "create", "--name", project1().Name, "--description", project1().Description, "--tenant", project1().Tenant},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.ProjectServiceCreateRequest{
							Login:       project1().Tenant,
							Name:        project1().Name,
							Description: project1().Description,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceCreateResponse{
								Project: project1(),
							})
						},
					},
				},
			}),
			WantObject:      project1(),
			WantProtoObject: project1(),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"project", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, project1()), 0755))
				},
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.ProjectServiceCreateRequest{
							Login:       project1().Tenant,
							Name:        project1().Name,
							Description: project1().Description,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceCreateResponse{
								Project: project1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2025-06-01 10:00:00 UTC
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
			CmdArgs: []string{"project", "delete", project1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.ProjectServiceDeleteRequest{
							Project: project1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceDeleteResponse{
								Project: project1(),
							})
						},
					},
				},
			}),
			WantObject: project1(),
		},
		{
			Name:    "delete from file",
			CmdArgs: append([]string{"project", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, project1()), 0755))
				},
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.ProjectServiceDeleteRequest{
							Project: project1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceDeleteResponse{
								Project: project1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2025-06-01 10:00:00 UTC
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
			CmdArgs: []string{"project", "update", project1().Uuid, "--name", "new-name", "--description", "new-desc"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.ProjectServiceUpdateRequest{
							Project:     project1().Uuid,
							Name:        new("new-name"),
							Description: new("new-desc"),
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceUpdateResponse{
								Project: project1(),
							})
						},
					},
				},
			}),
			WantObject: project1(),
		},
		{
			Name:    "update from file",
			CmdArgs: append([]string{"project", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, project1()), 0755))
				},
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.ProjectServiceUpdateRequest{
							Project:     project1().Uuid,
							Name:        new(project1().Name),
							Description: new(project1().Description),
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceUpdateResponse{
								Project: project1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION    CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project  2025-06-01 10:00:00 UTC
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
				ClientCalls: []e2e.ClientCall{
					{
						WantRequest: apiv2.ProjectServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.ProjectServiceListResponse{
								Projects: []*apiv2.Project{
									project1(),
									project2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION     CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project   2025-06-01 10:00:00 UTC
			f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  metal-stack  project-b  second project  2025-07-15 14:30:00 UTC
			`),
			WantWideTable: new(`
			ID                                    TENANT       NAME       DESCRIPTION     CREATION DATE
			0d81bca7-73f6-4da3-8397-4a8c52a0c583  metal-stack  project-a  first project   2025-06-01 10:00:00 UTC
			f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  metal-stack  project-b  second project  2025-07-15 14:30:00 UTC
			`),
			Template: new("{{ .uuid }} {{ .name }}"),
			WantTemplate: new(`
0d81bca7-73f6-4da3-8397-4a8c52a0c583 project-a
f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c project-b
			`),
			WantMarkdown: new(`
			| ID                                   | TENANT      | NAME      | DESCRIPTION    | CREATION DATE           |
			|--------------------------------------|-------------|-----------|----------------|-------------------------|
			| 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | metal-stack | project-a | first project  | 2025-06-01 10:00:00 UTC |
			| f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | metal-stack | project-b | second project | 2025-07-15 14:30:00 UTC |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
