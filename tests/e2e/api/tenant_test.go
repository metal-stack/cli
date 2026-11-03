package api_e2e

import (
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/dustin/go-humanize"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"

	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_TenantCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceGetResponse, *apiv2.Tenant]{
		{
			Name:    "describe",
			CmdArgs: []string{"tenant", "describe", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceGetRequest{
							Login: testresources.Tenant1().Login,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceGetResponse{
								Tenant: testresources.Tenant1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID           NAME         EMAIL                REGISTERED
            metal-stack  Metal Stack  info@metal-stack.io  now
			`),
			WantWideTable: new(`
            ID           NAME         EMAIL                REGISTERED
            metal-stack  Metal Stack  info@metal-stack.io  now
			`),
			WantMarkdown: new(`
            | ID          | NAME        | EMAIL               | REGISTERED |
            |-------------|-------------|---------------------|------------|
            | metal-stack | Metal Stack | info@metal-stack.io | now        |
			`),
			WantProtoObject: testresources.Tenant1(),
			Template:        new("{{ .login }} {{ .name }}"),
			WantTemplate: new(`
			metal-stack Metal Stack
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceListResponse, apiv2.Tenant]{
		{
			Name:    "list",
			CmdArgs: []string{"tenant", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceListRequest{Query: &apiv2.TenantQuery{}},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceListResponse{
								Tenants: []*apiv2.Tenant{
									testresources.Tenant1(),
									testresources.Tenant2(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID           NAME         EMAIL                REGISTERED
            metal-stack  Metal Stack  info@metal-stack.io  now
            acme-corp    ACME Corp    admin@acme.io        now
			`),
			WantWideTable: new(`
            ID           NAME         EMAIL                REGISTERED
            metal-stack  Metal Stack  info@metal-stack.io  now
            acme-corp    ACME Corp    admin@acme.io        now
			`),
			Template: new("{{ .login }} {{ .name }}"),
			WantTemplate: new(`
metal-stack Metal Stack
acme-corp ACME Corp
			`),
			WantMarkdown: new(`
            | ID          | NAME        | EMAIL               | REGISTERED |
            |-------------|-------------|---------------------|------------|
            | metal-stack | Metal Stack | info@metal-stack.io | now        |
            | acme-corp   | ACME Corp   | admin@acme.io       | now        |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_Create(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceCreateResponse, *apiv2.Tenant]{
		{
			Name:    "create",
			CmdArgs: []string{"tenant", "create", "--name", testresources.Tenant1().Name, "--description", testresources.Tenant1().Description, "--email", testresources.Tenant1().Email},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceCreateRequest{
							Description: &testresources.Tenant1().Description,
							Name:        testresources.Tenant1().Name,
							Email:       &testresources.Tenant1().Email,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceCreateResponse{
								Tenant: testresources.Tenant1(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Tenant1(),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"tenant", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Tenant1()), 0755))
				},
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceCreateRequest{
							Name:        testresources.Tenant1().Name,
							Description: &testresources.Tenant1().Description,
							Email:       &testresources.Tenant1().Email,
							AvatarUrl:   nil,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceCreateResponse{
								Tenant: testresources.Tenant1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID           NAME         EMAIL                REGISTERED
            metal-stack  Metal Stack  info@metal-stack.io  now
					`),
		},
		{
			Name:    "create many from file",
			CmdArgs: append([]string{"tenant", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshalToMultiYAML(t, testresources.Tenant1(), testresources.Tenant2()), 0755))
				},
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceCreateRequest{
							Name:        testresources.Tenant1().Name,
							Description: &testresources.Tenant1().Description,
							AvatarUrl:   nil,
							Email:       &testresources.Tenant1().Email,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceCreateResponse{
								Tenant: testresources.Tenant1(),
							})
						},
					},
					{
						WantRequest: &apiv2.TenantServiceCreateRequest{
							Name:        testresources.Tenant2().Name,
							Description: &testresources.Tenant2().Description,
							Email:       &testresources.Tenant2().Email,
							AvatarUrl:   nil,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceCreateResponse{
								Tenant: testresources.Tenant2(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID           NAME         EMAIL                REGISTERED
            metal-stack  Metal Stack  info@metal-stack.io  now
            acme-corp    ACME Corp    admin@acme.io        now
					`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceDeleteResponse, *apiv2.Tenant]{
		{
			Name:    "delete",
			CmdArgs: []string{"tenant", "delete", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceDeleteRequest{
							Login: testresources.Tenant1().Login,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceDeleteResponse{
								Tenant: testresources.Tenant1(),
							})
						},
					},
				},
			}),
			WantProtoObject: testresources.Tenant1(),
		},
		{
			Name:    "delete from file",
			CmdArgs: append([]string{"tenant", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Tenant1()), 0755))
				},
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceDeleteRequest{
							Login: testresources.Tenant1().Login,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceDeleteResponse{
								Tenant: testresources.Tenant1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID           NAME         EMAIL                REGISTERED
            metal-stack  Metal Stack  info@metal-stack.io  now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceUpdateResponse, *apiv2.Tenant]{
		{
			Name:    "update",
			CmdArgs: []string{"tenant", "update", testresources.Tenant1().Login, "--name", "new-name", "--description", "new-desc"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceUpdateRequest{
							Login:       testresources.Tenant1().Login,
							Name:        new("new-name"),
							Description: new("new-desc"),
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceUpdateResponse{
								Tenant: testresources.Tenant1(),
							})
						},
					},
				},
			}),
			WantErr: fmt.Errorf("not implemented"),
		},
		{
			Name:    "update from file",
			CmdArgs: append([]string{"tenant", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Tenant1()), 0755))
				},
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceUpdateRequest{
							Login:       testresources.Tenant1().Login,
							Name:        new(testresources.Tenant1().Name),
							Email:       new(testresources.Tenant1().Email),
							Description: new(testresources.Tenant1().Description),
							UpdateMeta: &apiv2.UpdateMeta{
								LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceUpdateResponse{
								Tenant: testresources.Tenant1(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID           NAME         EMAIL                REGISTERED
            metal-stack  Metal Stack  info@metal-stack.io  now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceUpdateResponse, *apiv2.Tenant]{
		{
			Name:    "apply",
			CmdArgs: append([]string{"tenant", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Tenant1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.TenantServiceCreateRequest{
								Email:       new(testresources.Tenant1().Email),
								Description: new(testresources.Tenant1().Description),
								Name:        testresources.Tenant1().Name,
								AvatarUrl:   nil,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.TenantServiceCreateResponse{
									Tenant: testresources.Tenant1(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            ID           NAME         EMAIL                REGISTERED
            metal-stack  Metal Stack  info@metal-stack.io  now
			`),
		},
		{
			Name:    "apply already exists",
			CmdArgs: append([]string{"tenant", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Tenant1()), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.TenantServiceCreateRequest{
								Email:       new(testresources.Tenant1().Email),
								Description: new(testresources.Tenant1().Description),
								Name:        testresources.Tenant1().Name,
								AvatarUrl:   nil,
							},
							WantError: connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")),
						},
						{
							WantRequest: &apiv2.TenantServiceUpdateRequest{
								Login:       testresources.Tenant1().Login,
								Email:       new(testresources.Tenant1().Email),
								Description: new(testresources.Tenant1().Description),
								Name:        new(testresources.Tenant1().Name),
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.TenantServiceUpdateResponse{
									Tenant: testresources.Tenant1(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            ID           NAME         EMAIL                REGISTERED
            metal-stack  Metal Stack  info@metal-stack.io  now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_ListInvites(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceInvitesListResponse, apiv2.TenantInvite]{
		{
			Name:    "list invites",
			CmdArgs: []string{"tenant", "invite", "list", "--tenant", testresources.Tenant2().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceInvitesListRequest{
							Login: testresources.Tenant2().Login,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceInvitesListResponse{
								Invites: []*apiv2.TenantInvite{
									testresources.Tenant1Invite(),
									testresources.Tenant2Invite(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            SECRET  TENANT       INVITED BY   ROLE                EXPIRES IN
            secret  acme-corp    acme-corp    TENANT_ROLE_EDITOR  2 days from now
            secret  metal-stack  metal-stack  TENANT_ROLE_VIEWER  2 days from now
			`),
			WantWideTable: new(`
            SECRET  TENANT       INVITED BY   ROLE                EXPIRES IN
            secret  acme-corp    acme-corp    TENANT_ROLE_EDITOR  2 days from now
            secret  metal-stack  metal-stack  TENANT_ROLE_VIEWER  2 days from now
			`),
			Template: new("{{ .tenant }} {{ .role }}"),
			WantTemplate: new(`
acme-corp 2
metal-stack 3
			`),
			WantMarkdown: new(`
            | SECRET | TENANT      | INVITED BY  | ROLE               | EXPIRES IN      |
            |--------|-------------|-------------|--------------------|-----------------|
            | secret | acme-corp   | acme-corp   | TENANT_ROLE_EDITOR | 2 days from now |
            | secret | metal-stack | metal-stack | TENANT_ROLE_VIEWER | 2 days from now |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_DeleteInvite(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceInviteDeleteResponse, string]{
		{
			Name:    "delete invite",
			CmdArgs: []string{"tenant", "invite", "delete", testresources.Tenant1Invite().Secret, "--tenant", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceInviteDeleteRequest{
							Login:  testresources.Tenant1().Login,
							Secret: testresources.Tenant1Invite().Secret,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceInviteDeleteResponse{})
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

func Test_TenantCmd_CreateInvite(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceInviteRequest, string]{
		{
			Name:    "create invite",
			CmdArgs: []string{"tenant", "invite", "generate-join-secret", "--role", testresources.Tenant1Invite().Role.String(), "--tenant", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceInviteRequest{
							Login: testresources.Tenant1().Login,
							Role:  testresources.Tenant1Invite().Role,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceInviteResponse{
								Invite: testresources.Tenant1Invite(),
							})
						},
					},
				},
			}),
			WantDefault: new(fmt.Sprintf("You can share this secret with the member to join, it expires in %s:\n\n%s (https://console.metal-stack.io/organization-invite/%s)",
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

func Test_TenantCmd_Join(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantServiceInviteAcceptResponse, string]{
		{
			Name:    "join",
			CmdArgs: []string{"tenant", "invite", "join", testresources.Tenant1Invite().Secret},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantServiceInviteGetRequest{
							Secret: testresources.Tenant1Invite().Secret,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceInviteGetResponse{
								Invite: testresources.Tenant1Invite(),
							})
						},
					},
					{
						WantRequest: &apiv2.TenantServiceInviteAcceptRequest{
							Secret: testresources.Tenant1Invite().Secret,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantServiceInviteAcceptResponse{
								Tenant:     testresources.Tenant1().Login,
								TenantName: testresources.Tenant1Invite().TargetTenantName,
							})
						},
					},
				},
			}),
			WantDefault: new(fmt.Sprintf("Do you want to join tenant \"%s\" as %s? [Y/n] ✔ successfully joined tenant \"%s\"",
				testresources.Tenant1Invite().TenantName,
				testresources.Tenant1Invite().Role.String(),
				testresources.Tenant1Invite().TenantName)),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
