package api_e2e

import (
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
)

func Test_TokenCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.TokenServiceGetResponse, *apiv2.Token]{
		{
			Name:    "describe",
			CmdArgs: []string{"token", "describe", testresources.Token1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TokenServiceGetRequest{
							Uuid: testresources.Token1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TokenServiceGetResponse{
								Token: testresources.Token1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES
            TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
			`),
			WantWideTable: new(`
            TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES
            TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
			`),
			WantMarkdown: new(`
            | TYPE           | ID                                   | ADMIN | USER                 | DESCRIPTION | ROLES | PERMS | EXPIRES                         |
            |----------------|--------------------------------------|-------|----------------------|-------------|-------|-------|---------------------------------|
            | TOKEN_TYPE_API | a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 |       | admin@metal-stack.io | ci token    | 0     | 0     | 2000-01-02 00:00:00 UTC (in 1d) |
			`),
			WantObject:      testresources.Token1(),
			WantProtoObject: testresources.Token1(),
			Template:        new("{{ .uuid }} {{ .description }}"),
			WantTemplate: new(`
			a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 ci token
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TokenCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.TokenServiceListResponse, apiv2.Token]{
		{
			Name:    "list",
			CmdArgs: []string{"token", "list"},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TokenServiceListRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TokenServiceListResponse{
								Tokens: []*apiv2.Token{
									testresources.Token1(),
									testresources.Token2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)  
            TOKEN_TYPE_API  b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01         dev@metal-stack.io    dev token    0      2      2000-01-03 00:00:00 UTC (in 2d)
			`),
			WantWideTable: new(`
            TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)  
            TOKEN_TYPE_API  b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01         dev@metal-stack.io    dev token    0      2      2000-01-03 00:00:00 UTC (in 2d)
			`),
			Template: new("{{ .uuid }} {{ .description }}"),
			WantTemplate: new(`
a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 ci token
b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01 dev token
			`),
			WantMarkdown: new(`
            | TYPE           | ID                                   | ADMIN | USER                 | DESCRIPTION | ROLES | PERMS | EXPIRES                         |
            |----------------|--------------------------------------|-------|----------------------|-------------|-------|-------|---------------------------------|
            | TOKEN_TYPE_API | a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 |       | admin@metal-stack.io | ci token    | 0     | 0     | 2000-01-02 00:00:00 UTC (in 1d) |
            | TOKEN_TYPE_API | b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01 |       | dev@metal-stack.io   | dev token   | 0     | 2     | 2000-01-03 00:00:00 UTC (in 2d) |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TokenCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv2.TokenServiceRevokeResponse, *apiv2.Token]{
		{
			Name:    "delete",
			CmdArgs: []string{"token", "delete", testresources.Token1().Uuid},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TokenServiceRevokeRequest{
							Uuid: testresources.Token1().Uuid,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TokenServiceRevokeResponse{})
						},
					},
				},
			}),
			WantMarkdown: new(`
			| TYPE                   | ID                                   | ADMIN | USER | DESCRIPTION | ROLES | PERMS | EXPIRES                              |
            |------------------------|--------------------------------------|-------|------|-------------|-------|-------|--------------------------------------|
            | TOKEN_TYPE_UNSPECIFIED | a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 |       |      |             | 0     | 0     | 1970-01-01 00:00:00 UTC (in -10957d) |
			`),
		},
		{
			Name:    "delete from file",
			CmdArgs: append([]string{"token", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Token1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.TokenServiceRevokeRequest{
								Uuid: testresources.Token1().Uuid,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.TokenServiceRevokeResponse{})
							},
						},
					},
				},
			),
			WantTable: new(`
            TYPE                    ID                                    ADMIN  USER  DESCRIPTION  ROLES  PERMS  EXPIRES                               
            TOKEN_TYPE_UNSPECIFIED  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90                            0      0      1970-01-01 00:00:00 UTC (in -10957d)
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TokenCmd_Create(t *testing.T) {
	tests := []*e2e.Test[apiv2.TokenServiceGetResponse, *apiv2.Token]{
		{
			Name:    "create",
			CmdArgs: []string{"token", "create", "--description", testresources.Token1().Description, "--expires", durationpb.New(testresources.Token1().Expires.AsTime().Sub(e2e.TimeBubbleStartTime())).AsDuration().String()},
			NewRootCmd: e2e.NewRootCmd(t, &e2e.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TokenServiceCreateRequest{
							Description: testresources.Token1().Description,
							Expires:     durationpb.New(testresources.Token1().Expires.AsTime().Sub(e2e.TimeBubbleStartTime())),
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TokenServiceCreateResponse{
								Token:  testresources.Token1(),
								Secret: "token-secret",
							})
						},
					},
				},
			}),
			WantMarkdown: new(`
            Make sure to copy your personal access token now as you will not be able to see this again.
            
            token-secret
            
            | TYPE           | ID                                   | ADMIN | USER                 | DESCRIPTION | ROLES | PERMS | EXPIRES                         |
            |----------------|--------------------------------------|-------|----------------------|-------------|-------|-------|---------------------------------|
            | TOKEN_TYPE_API | a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90 |       | admin@metal-stack.io | ci token    | 0     | 0     | 2000-01-02 00:00:00 UTC (in 1d) |			
			`),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"token", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Token1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.TokenServiceCreateRequest{
								Description: testresources.Token1().Description,
								Expires:     durationpb.New(testresources.Token1().Expires.AsTime().Sub(e2e.TimeBubbleStartTime())),
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.TokenServiceCreateResponse{
									Token:  testresources.Token1(),
									Secret: "token-secret",
								})
							},
						},
					},
				}),
			WantTable: new(`
            Make sure to copy your personal access token now as you will not be able to see this again.
            
            token-secret
            
            TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_UpdateCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv2.TokenServiceUpdateResponse, *apiv2.Token]{
		{
			Name:    "update from file",
			CmdArgs: append([]string{"token", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Token2()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.TokenServiceUpdateRequest{
								Uuid:         testresources.Token2().Uuid,
								Description:  new(testresources.Token2().Description),
								Permissions:  testresources.Token2().Permissions,
								ProjectRoles: testresources.Token2().ProjectRoles,
								TenantRoles:  testresources.Token2().TenantRoles,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.TokenServiceUpdateResponse{
									Token: testresources.Token2(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            TYPE            ID                                    ADMIN  USER                DESCRIPTION  ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_API  b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01         dev@metal-stack.io  dev token    0      2      2000-01-03 00:00:00 UTC (in 2d)
				`),
			Template:     new("{{ .uuid }} {{ .permissions }}"),
			WantTemplate: new(`b4c2e7f3-5a9d-4b8e-a1c3-2d6f9e4b8a01 [map[methods:[api/method1 api/method2] subject:0d81bca7-73f6-4da3-8397-4a8c52a0c583] map[methods:[api/method3] subject:metal-stack]]`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TokenCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv2.TokenServiceUpdateResponse, *apiv2.Token]{
		{
			Name:    "apply",
			CmdArgs: append([]string{"token", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Token1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.TokenServiceCreateRequest{
								Description: testresources.Token1().Description,
								Labels:      testresources.Token1().Meta.Labels,
								Expires:     durationpb.New(testresources.Token1().Expires.AsTime().Sub(e2e.TimeBubbleStartTime())),
								Permissions: testresources.Token1().Permissions,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.TokenServiceCreateResponse{
									Token: testresources.Token1(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            Make sure to copy your personal access token now as you will not be able to see this again.
            
            
            
            TYPE            ID                                    ADMIN  USER                  DESCRIPTION  ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_API  a3b1f6d2-4e8c-4f7a-9d2e-1b5c8f3a7e90         admin@metal-stack.io  ci token     0      0      2000-01-02 00:00:00 UTC (in 1d)
			`),
		},
		{
			Name:    "apply already exists",
			CmdArgs: append([]string{"token", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2e.NewRootCmd(t,
				&e2e.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Token1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.TokenServiceCreateRequest{
								Description: testresources.Token1().Description,
								Labels:      testresources.Token1().Meta.Labels,
								Expires:     durationpb.New(testresources.Token1().Expires.AsTime().Sub(e2e.TimeBubbleStartTime())),
								Permissions: testresources.Token1().Permissions,
							},
							WantError: connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already_exists")),
						},
					},
				},
			),
			WantErr: fmt.Errorf("error creating entity: already_exists: already_exists"),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
