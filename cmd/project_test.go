package cmd

import (
	"testing"
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
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

func Test_ProjectCmd_SingleResult(t *testing.T) {
	p1 := project1()

	tests := []*Test[apiv2.ProjectServiceGetRequest, apiv2.ProjectServiceGetResponse]{
		{
			Name: "describe",
			Cmd: func() []string {
				return []string{"project", "describe", p1.Uuid}
			},
			WantRequest: apiv2.ProjectServiceGetRequest{
				Project: p1.Uuid,
			},
			WantResponse: apiv2.ProjectServiceGetResponse{
				Project: p1,
			},
			WantObject: p1,
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

func Test_ProjectCmd_MultiResult(t *testing.T) {
	tests := []*Test[apiv2.ProjectServiceListRequest, apiv2.ProjectServiceListResponse]{
		{
			Name: "list",
			Cmd: func() []string {
				return []string{"project", "list"}
			},
			WantRequest: apiv2.ProjectServiceListRequest{},
			WantResponse: apiv2.ProjectServiceListResponse{
				Projects: []*apiv2.Project{
					project1(),
					project2(),
				},
			},
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
