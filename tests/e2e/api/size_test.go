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

func Test_SizeCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.SizeServiceListResponse, []*apiv2.Size]{
		{
			Name:    "list",
			CmdArgs: []string{"size", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.SizeServiceListRequest{
							Query: &apiv2.SizeQuery{},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.SizeServiceListResponse{
								Sizes: []*apiv2.Size{
									testresources.Size2(),
									testresources.Size1(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID             NAME           DESCRIPTION                CPU RANGE  MEMORY RANGE     STORAGE RANGE    GPU RANGE
            g1-medium-x86  g1-medium-x86  A medium sized GPU server  32 - 32    275 GB - 275 GB  1.8 TB - 1.8 TB  AD102GL [RTX 6000 Ada Generation]: 1 - 1
            v1-medium-x86  v1-medium-x86  Virtual size for mini-lab  4 - 4      500 MB - 4.0 GB  1.0 GB - 100 GB
			`),
			WantWideTable: new(`
            ID             NAME           DESCRIPTION                CPU RANGE  MEMORY RANGE     STORAGE RANGE    GPU RANGE
            g1-medium-x86  g1-medium-x86  A medium sized GPU server  32 - 32    275 GB - 275 GB  1.8 TB - 1.8 TB  AD102GL [RTX 6000 Ada Generation]: 1 - 1
            v1-medium-x86  v1-medium-x86  Virtual size for mini-lab  4 - 4      500 MB - 4.0 GB  1.0 GB - 100 GB
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
g1-medium-x86 g1-medium-x86
v1-medium-x86 v1-medium-x86
			`),
			WantMarkdown: new(`
            | ID            | NAME          | DESCRIPTION               | CPU RANGE | MEMORY RANGE    | STORAGE RANGE   | GPU RANGE                                |
            |---------------|---------------|---------------------------|-----------|-----------------|-----------------|------------------------------------------|
            | g1-medium-x86 | g1-medium-x86 | A medium sized GPU server | 32 - 32   | 275 GB - 275 GB | 1.8 TB - 1.8 TB | AD102GL [RTX 6000 Ada Generation]: 1 - 1 |
            | v1-medium-x86 | v1-medium-x86 | Virtual size for mini-lab | 4 - 4     | 500 MB - 4.0 GB | 1.0 GB - 100 GB |                                          |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_SizeCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.SizeServiceGetResponse, *apiv2.Size]{
		{
			Name:    "get",
			CmdArgs: []string{"size", "describe", testresources.Size1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.SizeServiceGetRequest{
							Id: testresources.Size1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.SizeServiceGetResponse{
								Size: testresources.Size1(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Size1(),
			WantTable: new(`
            ID             NAME           DESCRIPTION                CPU RANGE  MEMORY RANGE     STORAGE RANGE    GPU RANGE
            v1-medium-x86  v1-medium-x86  Virtual size for mini-lab  4 - 4      500 MB - 4.0 GB  1.0 GB - 100 GB
			`),
			WantWideTable: new(`
            ID             NAME           DESCRIPTION                CPU RANGE  MEMORY RANGE     STORAGE RANGE    GPU RANGE
            v1-medium-x86  v1-medium-x86  Virtual size for mini-lab  4 - 4      500 MB - 4.0 GB  1.0 GB - 100 GB
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
			v1-medium-x86 v1-medium-x86
			`),
			WantMarkdown: new(`
            | ID            | NAME          | DESCRIPTION               | CPU RANGE | MEMORY RANGE    | STORAGE RANGE   | GPU RANGE |
            |---------------|---------------|---------------------------|-----------|-----------------|-----------------|-----------|
            | v1-medium-x86 | v1-medium-x86 | Virtual size for mini-lab | 4 - 4     | 500 MB - 4.0 GB | 1.0 GB - 100 GB |           |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
