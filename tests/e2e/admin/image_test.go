package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_AdminImageCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv2.ImageServiceDeleteResponse, *apiv2.Image]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "image", "delete", testresources.Image1().Id},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.ImageServiceDeleteRequest{
							Id: testresources.Image1().Id,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ImageServiceDeleteResponse{
								Image: testresources.Image1(),
							})
						},
					},
				},
			}),
			WantObject:      testresources.Image1(),
			WantProtoObject: testresources.Image1(),
			WantTable: new(`
			ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS
			ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
			`),
			WantWideTable: new(`
			ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS
			ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
			`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
			ubuntu-24.04 Ubuntu 24.04
			`),
			WantMarkdown: new(`
            | ID           | NAME         | DESCRIPTION      | FEATURES | EXPIRATION | STATUS    |
            |--------------|--------------|------------------|----------|------------|-----------|
            | ubuntu-24.04 | Ubuntu 24.04 | Ubuntu 24.04 LTS | machine  |            | supported |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminImageCmd_Create(t *testing.T) {
	tests := []*e2e.Test[adminv2.ImageServiceCreateResponse, *apiv2.Image]{
		{
			Name: "create",
			CmdArgs: []string{"admin", "image", "create",
				"--id", testresources.Image1().Id,
				"--url", testresources.Image1().Url,
				"--features", "machine",
				"--classification", "supported",
				"--description", *testresources.Image1().Description,
				"--name", *testresources.Image1().Name},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.ImageServiceCreateRequest{
							Image: testresources.Image1(),
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ImageServiceCreateResponse{
								Image: testresources.Image1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.Image1(),
			WantTable: new(`
			ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS     
            ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
			`),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"admin", "image", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Image1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.ImageServiceCreateRequest{
								Image: testresources.Image1(),
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.ImageServiceCreateResponse{
									Image: testresources.Image1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
					ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS
					ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
					`),
			WantWideTable: new(`
					ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS
					ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
					`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
					ubuntu-24.04 Ubuntu 24.04
					`),
			WantMarkdown: new(`
		            | ID           | NAME         | DESCRIPTION      | FEATURES | EXPIRATION | STATUS    |
		            |--------------|--------------|------------------|----------|------------|-----------|
		            | ubuntu-24.04 | Ubuntu 24.04 | Ubuntu 24.04 LTS | machine  |            | supported |
					`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminImageCmd_Update(t *testing.T) {
	tests := []*e2e.Test[adminv2.ImageServiceUpdateResponse, *apiv2.Image]{
		{
			Name:    "update",
			CmdArgs: []string{"admin", "image", "update", testresources.Image1().Id, "--name", *testresources.Image1().Name},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.ImageServiceUpdateRequest{
							Id:   testresources.Image1().Id,
							Name: testresources.Image1().Name,
							UpdateMeta: &apiv2.UpdateMeta{
								LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.ImageServiceUpdateResponse{
								Image: testresources.Image1(),
							})
						},
					},
				},
			}),
			WantObject: testresources.Image1(),
			WantTable: new(`
			ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS     
            ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
			`),
		},
		{
			Name:    "update from file",
			CmdArgs: append([]string{"admin", "image", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Image1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &adminv2.ImageServiceUpdateRequest{
								Id:             testresources.Image1().Id,
								Name:           testresources.Image1().Name,
								Description:    testresources.Image1().Description,
								Features:       testresources.Image1().Features,
								Classification: testresources.Image1().Classification,
								ExpiresAt:      testresources.Image1().ExpiresAt,
								Url:            new(testresources.Image1().Url),
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_CLIENT,
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&adminv2.ImageServiceUpdateResponse{
									Image: testresources.Image1(),
								})
							},
						},
					},
				}),
			WantTable: new(`
					ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS
					ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
					`),
			WantWideTable: new(`
					ID            NAME          DESCRIPTION       FEATURES  EXPIRATION  STATUS
					ubuntu-24.04  Ubuntu 24.04  Ubuntu 24.04 LTS  machine               supported
					`),
			Template: new("{{ .id }} {{ .name }}"),
			WantTemplate: new(`
					ubuntu-24.04 Ubuntu 24.04
					`),
			WantMarkdown: new(`
		            | ID           | NAME         | DESCRIPTION      | FEATURES | EXPIRATION | STATUS    |
		            |--------------|--------------|------------------|----------|------------|-----------|
		            | ubuntu-24.04 | Ubuntu 24.04 | Ubuntu 24.04 LTS | machine  |            | supported |
					`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
