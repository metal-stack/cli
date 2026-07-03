package admin_e2e

import (
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

func Test_AdminSizeCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeServiceDeleteResponse, *apiv2.Size]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "size", "delete", testresources.Size1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.SizeServiceDeleteRequest{
							Id: testresources.Size1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.SizeServiceDeleteResponse{
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

func Test_AdminSizeCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.SizeServiceListResponse, []*apiv2.Size]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "size", "list"},
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

func Test_AdminSizeCmd_Create(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeServiceCreateResponse, *apiv2.Size]{
		{
			Name:    "create from file",
			CmdArgs: append([]string{"admin", "size", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Size1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.SizeServiceCreateRequest{
								Size: testresources.Size1(),
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.SizeServiceCreateResponse{
									Size: testresources.Size1(),
								})
							},
						},
					},
				}),
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

func Test_AdminSizeCmd_Update(t *testing.T) {
	tests := []*e2e.Test[adminv2.SizeServiceUpdateResponse, *apiv2.Size]{
		{
			Name:    "update from file",
			CmdArgs: append([]string{"admin", "size", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Size1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.SizeServiceUpdateRequest{
								Id:          testresources.Size1().Id,
								Name:        testresources.Size1().Name,
								Description: testresources.Size1().Description,
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_CLIENT,
								},
								Labels: &apiv2.UpdateLabels{
									Strategy: &apiv2.UpdateLabels_Replace{
										Replace: testresources.Size1().Meta.Labels,
									},
								},
								Constraints: testresources.Size1().Constraints,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.SizeServiceUpdateResponse{
									Size: testresources.Size1(),
								})
							},
						},
					},
				}),
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
