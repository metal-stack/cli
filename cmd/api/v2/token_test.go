package v2

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/metal-stack/api/go/errorutil"
	"github.com/metal-stack/api/go/metalstack/admin/v2/adminv2connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/api/go/metalstack/api/v2/apiv2connect"
	"github.com/metal-stack/api/go/metalstack/infra/v2/infrav2connect"
	"google.golang.org/protobuf/testing/protocmp"
)

func Test_toPermissionsByVisibility(t *testing.T) {
	tests := []struct {
		name    string
		perms   []*apiv2.MethodPermission
		want    []*apiv2.PermissionsByVisibility
		wantErr error
	}{
		{
			name: "convert",
			perms: []*apiv2.MethodPermission{
				{
					Subject: "project-a",
					Methods: []string{apiv2connect.IPServiceCreateProcedure, apiv2connect.IPServiceDeleteProcedure},
				},
				{
					Subject: "*",
					Methods: []string{apiv2connect.AuditServiceGetProcedure},
				},
				{
					Subject: "",
					Methods: []string{apiv2connect.TenantServiceListProcedure},
				},
				{
					Subject: "",
					Methods: []string{adminv2connect.AuditServiceListProcedure},
				},
				{
					Subject: "*",
					Methods: []string{infrav2connect.BootServiceRegisterProcedure},
				},
				{
					Subject: "",
					Methods: []string{apiv2connect.HealthServiceGetProcedure},
				},
			},
			want: []*apiv2.PermissionsByVisibility{
				{
					Visibility: &apiv2.PermissionsByVisibility_Project{
						Project: &apiv2.ProjectPermissions{
							Project: "project-a",
							Methods: []string{apiv2connect.IPServiceCreateProcedure, apiv2connect.IPServiceDeleteProcedure},
						},
					},
				},
				{
					Visibility: &apiv2.PermissionsByVisibility_Tenant{
						Tenant: &apiv2.TenantPermissions{
							Login:   "*",
							Methods: []string{apiv2connect.AuditServiceGetProcedure},
						},
					},
				},
				{
					Visibility: &apiv2.PermissionsByVisibility_Self{
						Self: &apiv2.SelfPermissions{
							Methods: []string{apiv2connect.TenantServiceListProcedure},
						},
					},
				},
				{
					Visibility: &apiv2.PermissionsByVisibility_Admin{
						Admin: &apiv2.AdminPermissions{
							Methods: []string{adminv2connect.AuditServiceListProcedure},
						},
					},
				},
				{
					Visibility: &apiv2.PermissionsByVisibility_Machine{
						Machine: &apiv2.MachinePermissions{
							Uuid:    "*",
							Methods: []string{infrav2connect.BootServiceRegisterProcedure},
						},
					},
				},
				{
					Visibility: &apiv2.PermissionsByVisibility_Public{
						Public: &apiv2.PublicPermissions{
							Methods: []string{apiv2connect.HealthServiceGetProcedure},
						},
					},
				},
			},
		},
		{
			name: "aggregate",
			perms: []*apiv2.MethodPermission{
				{
					Subject: "project-a",
					Methods: []string{apiv2connect.IPServiceCreateProcedure},
				},
				{
					Subject: "project-b",
					Methods: []string{apiv2connect.IPServiceDeleteProcedure},
				},
				{
					Subject: "project-a",
					Methods: []string{apiv2connect.IPServiceDeleteProcedure},
				},
			},
			want: []*apiv2.PermissionsByVisibility{
				{
					Visibility: &apiv2.PermissionsByVisibility_Project{
						Project: &apiv2.ProjectPermissions{
							Project: "project-a",
							Methods: []string{apiv2connect.IPServiceCreateProcedure, apiv2connect.IPServiceDeleteProcedure},
						},
					},
				},
				{
					Visibility: &apiv2.PermissionsByVisibility_Project{
						Project: &apiv2.ProjectPermissions{
							Project: "project-b",
							Methods: []string{apiv2connect.IPServiceDeleteProcedure},
						},
					},
				},
			},
		},
		{
			name: "unknown method",
			perms: []*apiv2.MethodPermission{
				{
					Subject: "",
					Methods: []string{"/foo"},
				},
			},
			wantErr: fmt.Errorf("method is not part of the api: /foo"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := toPermissionsByVisibility(tt.perms)
			if diff := cmp.Diff(gotErr, tt.wantErr, errorutil.ErrorStringComparer()); diff != "" {
				t.Errorf("diff = %s", diff)
				return
			}
			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Errorf("diff (+got -want):\n %s", diff)
			}
		})
	}
}
