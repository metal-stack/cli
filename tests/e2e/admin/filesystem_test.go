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

func Test_AdminFileSystemCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.FilesystemServiceListResponse, *apiv2.FilesystemLayout]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "filesystem", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.FilesystemServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.FilesystemServiceListResponse{
								FilesystemLayouts: []*apiv2.FilesystemLayout{
									testresources.FilesystemLayout1(),
									testresources.FilesystemLayout2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  
			ubuntu-raid1    Ubuntu Raid1    Raid1 layout for ubuntu machines    m2.large  ubuntu:>= 24.0  
			`),
			WantWideTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          FILESYSTEMS  DISKS  RAID  VOLUME GROUPS  LOGICAL VOLUMES  
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  1            1      0     0              0                
			ubuntu-raid1    Ubuntu Raid1    Raid1 layout for ubuntu machines    m2.large  ubuntu:>= 24.0  0            0      1     1              1                
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
debian-default Debian Default
ubuntu-raid1 Ubuntu Raid1
`),
			WantMarkdown: new(`
			| ID             | NAME           | DESCRIPTION                        | SIZES    | IMAGES         |
			|----------------|----------------|------------------------------------|----------|----------------|
			| debian-default | Debian Default | Default layout for debian machines | m2.small | debian:>= 12.0 |
			| ubuntu-raid1   | Ubuntu Raid1   | Raid1 layout for ubuntu machines   | m2.large | ubuntu:>= 24.0 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminFileSystemCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.FilesystemServiceGetResponse, *apiv2.FilesystemLayout]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "filesystem", "describe", testresources.FilesystemLayout1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.FilesystemServiceGetRequest{
							Id: testresources.FilesystemLayout1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.FilesystemServiceGetResponse{
								FilesystemLayout: testresources.FilesystemLayout1(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.FilesystemLayout1(),
			WantTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  
			`),
			WantWideTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          FILESYSTEMS  DISKS  RAID  VOLUME GROUPS  LOGICAL VOLUMES  
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  1            1      0     0              0                
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
debian-default Debian Default
`),
			WantMarkdown: new(`
			| ID             | NAME           | DESCRIPTION                        | SIZES    | IMAGES         |
			|----------------|----------------|------------------------------------|----------|----------------|
			| debian-default | Debian Default | Default layout for debian machines | m2.small | debian:>= 12.0 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminFileSystemCmd_Create(t *testing.T) {
	tests := []*e2e.Test[adminv2.FilesystemServiceCreateResponse, *apiv2.FilesystemLayout]{
		{
			Name:    "create from file",
			CmdArgs: append([]string{"admin", "filesystem", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.FilesystemLayout1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.FilesystemServiceCreateRequest{
								FilesystemLayout: testresources.FilesystemLayout1(),
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.FilesystemServiceCreateResponse{
									FilesystemLayout: testresources.FilesystemLayout1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  
			`),
			WantWideTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          FILESYSTEMS  DISKS  RAID  VOLUME GROUPS  LOGICAL VOLUMES  
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  1            1      0     0              0                
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
debian-default Debian Default
`),
			WantMarkdown: new(`
			| ID             | NAME           | DESCRIPTION                        | SIZES    | IMAGES         |
			|----------------|----------------|------------------------------------|----------|----------------|
			| debian-default | Debian Default | Default layout for debian machines | m2.small | debian:>= 12.0 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminFileSystemCmd_Update(t *testing.T) {
	layout := testresources.FilesystemLayout1()

	tests := []*e2e.Test[adminv2.FilesystemServiceUpdateResponse, *apiv2.FilesystemLayout]{
		{
			Name:    "update from file",
			CmdArgs: append([]string{"admin", "filesystem", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, layout), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.FilesystemServiceUpdateRequest{
								Id:             layout.Id,
								UpdateMeta:     &apiv2.UpdateMeta{LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER},
								Name:           layout.Name,
								Description:    layout.Description,
								Filesystems:    layout.Filesystems,
								Disks:          layout.Disks,
								Raid:           layout.Raid,
								VolumeGroups:   layout.VolumeGroups,
								LogicalVolumes: layout.LogicalVolumes,
								Constraints:    layout.Constraints,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.FilesystemServiceUpdateResponse{
									FilesystemLayout: layout,
								})
							},
						},
					},
				}),
			WantTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  
			`),
			WantWideTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          FILESYSTEMS  DISKS  RAID  VOLUME GROUPS  LOGICAL VOLUMES  
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  1            1      0     0              0                
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
debian-default Debian Default
`),
			WantMarkdown: new(`
			| ID             | NAME           | DESCRIPTION                        | SIZES    | IMAGES         |
			|----------------|----------------|------------------------------------|----------|----------------|
			| debian-default | Debian Default | Default layout for debian machines | m2.small | debian:>= 12.0 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminFileSystemCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.FilesystemServiceDeleteResponse, *apiv2.FilesystemLayout]{
		{
			Name:    "delete",
			CmdArgs: []string{"admin", "filesystem", "delete", testresources.FilesystemLayout1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.FilesystemServiceDeleteRequest{
							Id: testresources.FilesystemLayout1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.FilesystemServiceDeleteResponse{
								FilesystemLayout: testresources.FilesystemLayout1(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.FilesystemLayout1(),
			WantTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  
			`),
			WantWideTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          FILESYSTEMS  DISKS  RAID  VOLUME GROUPS  LOGICAL VOLUMES  
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  1            1      0     0              0                
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
debian-default Debian Default
`),
			WantMarkdown: new(`
			| ID             | NAME           | DESCRIPTION                        | SIZES    | IMAGES         |
			|----------------|----------------|------------------------------------|----------|----------------|
			| debian-default | Debian Default | Default layout for debian machines | m2.small | debian:>= 12.0 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminFileSystemCmd_Match_SizeAndImage(t *testing.T) {
	tests := []*e2e.Test[adminv2.FilesystemServiceMatchResponse, *apiv2.FilesystemLayout]{
		{
			Name:    "match size and image",
			CmdArgs: []string{"admin", "filesystem", "match", "--size", "m2.small", "--image", "debian"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.FilesystemServiceMatchRequest{
							Match: &adminv2.FilesystemServiceMatchRequest_SizeAndImage{
								SizeAndImage: &adminv2.MatchImageAndSize{
									Size:  "m2.small",
									Image: "debian",
								},
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.FilesystemServiceMatchResponse{
								FilesystemLayout: testresources.FilesystemLayout1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  
			`),
			WantWideTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          FILESYSTEMS  DISKS  RAID  VOLUME GROUPS  LOGICAL VOLUMES  
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  1            1      0     0              0                
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
debian-default Debian Default
`),
			WantMarkdown: new(`
			| ID             | NAME           | DESCRIPTION                        | SIZES    | IMAGES         |
			|----------------|----------------|------------------------------------|----------|----------------|
			| debian-default | Debian Default | Default layout for debian machines | m2.small | debian:>= 12.0 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminFileSystemCmd_Match_MachineAndLayout(t *testing.T) {
	machineUUID := "12345678-1234-1234-1234-123456789abc"

	tests := []*e2e.Test[adminv2.FilesystemServiceMatchResponse, *apiv2.FilesystemLayout]{
		{
			Name:    "match machine and filesystem layout",
			CmdArgs: []string{"admin", "filesystem", "match", "--machine", machineUUID, "--id", testresources.FilesystemLayout1().Id},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.FilesystemServiceMatchRequest{
							Match: &adminv2.FilesystemServiceMatchRequest_MachineAndFilesystemlayout{
								MachineAndFilesystemlayout: &adminv2.MatchMachineAndFilesystemLayout{
									Machine:          machineUUID,
									FilesystemLayout: testresources.FilesystemLayout1().Id,
								},
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.FilesystemServiceMatchResponse{
								FilesystemLayout: testresources.FilesystemLayout1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  
			`),
			WantWideTable: new(`
			ID              NAME            DESCRIPTION                         SIZES     IMAGES          FILESYSTEMS  DISKS  RAID  VOLUME GROUPS  LOGICAL VOLUMES  
			debian-default  Debian Default  Default layout for debian machines  m2.small  debian:>= 12.0  1            1      0     0              0                
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
debian-default Debian Default
`),
			WantMarkdown: new(`
			| ID             | NAME           | DESCRIPTION                        | SIZES    | IMAGES         |
			|----------------|----------------|------------------------------------|----------|----------------|
			| debian-default | Debian Default | Default layout for debian machines | m2.small | debian:>= 12.0 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminFileSystemCmd_Match_InvalidFlags(t *testing.T) {
	tests := []*e2e.Test[any, any]{
		{
			Name:       "match without flags fails",
			CmdArgs:    []string{"admin", "filesystem", "match"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{}),
			WantErr:    fmt.Errorf("either both --size and --image or both --machine and --id must be set"),
		},
		{
			Name:       "match with only size fails",
			CmdArgs:    []string{"admin", "filesystem", "match", "--size", "m2.small"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{}),
			WantErr:    fmt.Errorf("either both --size and --image or both --machine and --id must be set"),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
