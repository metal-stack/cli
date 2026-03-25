package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
)

func Test_AdminProjectCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv2.ProjectServiceListResponse, apiv2.Project]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "project", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.ProjectServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ProjectServiceListResponse{
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
