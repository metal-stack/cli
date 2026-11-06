package api_e2e

import (
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_TenantMemberCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantMemberServiceListResponse, []apiv2.TenantMember]{
		{
			Name:    "list tenant members",
			CmdArgs: []string{"tenant", "member", "list", "--tenant", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantMemberServiceListRequest{
							Login: testresources.Tenant1().Login,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantMemberServiceListResponse{
								Members: []*apiv2.TenantMember{
									testresources.Tenant1Member(),
									testresources.Tenant2Member(),
								},
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    ROLE                SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  TENANT_ROLE_OWNER   now
            40c0da4b-9eb9-4371-91aa-1ae62193fa54  TENANT_ROLE_EDITOR  now
			`),
			WantWideTable: new(`
            ID                                    ROLE                SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  TENANT_ROLE_OWNER   now
            40c0da4b-9eb9-4371-91aa-1ae62193fa54  TENANT_ROLE_EDITOR  now
			`),
			Template: new("{{ .member }} {{ .role }}"),
			WantTemplate: new(`
16d6e8ba-f574-494f-8d5e-74f6cb2d8db0 1
40c0da4b-9eb9-4371-91aa-1ae62193fa54 2
			`),
			WantMarkdown: new(`
            | ID                                   | ROLE               | SINCE |
            |--------------------------------------|--------------------|-------|
            | 16d6e8ba-f574-494f-8d5e-74f6cb2d8db0 | TENANT_ROLE_OWNER  | now   |
            | 40c0da4b-9eb9-4371-91aa-1ae62193fa54 | TENANT_ROLE_EDITOR | now   |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantMemberCmd_Get(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantMemberServiceDeleteResponse, *apiv2.TenantMember]{
		{
			Name:    "get tenant member",
			CmdArgs: []string{"tenant", "member", "get", testresources.Tenant1Member().Member, "--tenant", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantMemberServiceGetRequest{
							Login:  testresources.Tenant1().Login,
							Member: testresources.Tenant1Member().Member,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantMemberServiceGetResponse{
								Member: testresources.Tenant1Member(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    ROLE               SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  TENANT_ROLE_OWNER  now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantMemberCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantMemberServiceDeleteResponse, *apiv2.TenantMember]{
		{
			Name:    "delete tenant member",
			CmdArgs: []string{"tenant", "member", "remove", testresources.Tenant1Member().Member, "--tenant", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantMemberServiceDeleteRequest{
							Login:  testresources.Tenant1().Login,
							Member: testresources.Tenant1Member().Member,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantMemberServiceDeleteResponse{
								Member: testresources.Tenant1Member(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    ROLE               SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  TENANT_ROLE_OWNER  now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantMemberCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantMemberServiceUpdateResponse, *apiv2.TenantMember]{
		{
			Name:    "update tenant member",
			CmdArgs: []string{"tenant", "member", "update", testresources.Tenant1Member().Member, "--tenant", testresources.Tenant1().Login, "--role", testresources.Tenant1Member().Role.String()},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantMemberServiceUpdateRequest{
							Login:  testresources.Tenant1().Login,
							Member: testresources.Tenant1Member().Member,
							Role:   pointer.PointerOrNil(testresources.Tenant1Member().Role),
							UpdateMeta: &apiv2.UpdateMeta{
								LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
							},
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantMemberServiceUpdateResponse{
								Member: testresources.Tenant1Member(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    ROLE               SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  TENANT_ROLE_OWNER  now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantMemberCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantMemberServiceUpdateResponse, *apiv2.TenantMember]{
		{
			Name:    "apply",
			CmdArgs: append([]string{"tenant", "member", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, &apiv2.TenantMember{
							Tenant: testresources.Tenant1().Login,
							Member: testresources.Tenant1Member().Member,
							Role:   testresources.Tenant1Member().Role,
						}), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.TenantMemberServiceCreateRequest{
								Login:  testresources.Tenant1().Login,
								Member: testresources.Tenant1Member().Member,
								Role:   testresources.Tenant1Member().Role,
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.TenantMemberServiceCreateResponse{
									Member: testresources.Tenant1Member(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            ID                                    ROLE               SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  TENANT_ROLE_OWNER  now
			`),
		},
		{
			Name:    "apply already exists",
			CmdArgs: append([]string{"tenant", "member", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, &apiv2.TenantMember{
							Tenant: testresources.Tenant1().Login,
							Member: testresources.Tenant1Member().Member,
							Role:   testresources.Tenant1Member().Role,
						}), 0755))
					},
					ClientCalls: []client.ClientCall{
						{
							WantRequest: &apiv2.TenantMemberServiceCreateRequest{
								Login:  testresources.Tenant1().Login,
								Member: testresources.Tenant1Member().Member,
								Role:   testresources.Tenant1Member().Role,
							},
							WantError: connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")),
						},
						{
							WantRequest: &apiv2.TenantMemberServiceUpdateRequest{
								Login:  testresources.Tenant1().Login,
								Member: testresources.Tenant1Member().Member,
								Role:   &testresources.Tenant1Member().Role,
								UpdateMeta: &apiv2.UpdateMeta{
									LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
								},
							},
							WantResponse: func() connect.AnyResponse {
								return connect.NewResponse(&apiv2.TenantMemberServiceUpdateResponse{
									Member: testresources.Tenant1Member(),
								})
							},
						},
					},
				},
			),
			WantTable: new(`
            ID                                    ROLE               SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  TENANT_ROLE_OWNER  now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantMemberCmd_Leave(t *testing.T) {
	tests := []*e2e.Test[apiv2.TenantMemberServiceDeleteResponse, *apiv2.TenantMember]{
		{
			Name:    "leave tenant",
			CmdArgs: []string{"tenant", "member", "leave", "--tenant", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &apiv2.TenantMemberServiceLeaveRequest{
							Login: testresources.Tenant1().Login,
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.TenantMemberServiceLeaveResponse{
								Member: testresources.Tenant1Member(),
							})
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    ROLE               SINCE
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  TENANT_ROLE_OWNER  now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
