package helpers

import (
	"fmt"
	"slices"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/api/go/permissions"
)

func ToPermissionsByVisibility(perms []*apiv2.MethodPermission) ([]*apiv2.PermissionsByVisibility, error) {
	var res []*apiv2.PermissionsByVisibility

	for _, perm := range perms {
		for _, method := range perm.Methods {
			if _, ok := permissions.GetServicePermissions().Visibility.Admin[method]; ok {
				idx := slices.IndexFunc(res, func(p *apiv2.PermissionsByVisibility) bool {
					return p.GetAdmin() != nil
				})

				if idx < 0 {
					res = append(res, &apiv2.PermissionsByVisibility{
						Visibility: &apiv2.PermissionsByVisibility_Admin{
							Admin: &apiv2.AdminPermissions{
								Methods: []string{method},
							},
						},
					})

					continue
				}

				res[idx].GetAdmin().Methods = append(res[idx].GetAdmin().Methods, method)

				continue
			}

			if _, ok := permissions.GetServicePermissions().Visibility.Infra[method]; ok {
				idx := slices.IndexFunc(res, func(p *apiv2.PermissionsByVisibility) bool {
					return p.GetInfra() != nil
				})

				if idx < 0 {
					res = append(res, &apiv2.PermissionsByVisibility{
						Visibility: &apiv2.PermissionsByVisibility_Infra{
							Infra: &apiv2.InfraPermissions{
								Methods: []string{method},
							},
						},
					})

					continue
				}

				res[idx].GetInfra().Methods = append(res[idx].GetInfra().Methods, method)

				continue
			}

			if _, ok := permissions.GetServicePermissions().Visibility.Machine[method]; ok {
				idx := slices.IndexFunc(res, func(p *apiv2.PermissionsByVisibility) bool {
					return p.GetMachine() != nil && p.GetMachine().Uuid == perm.Subject
				})

				if idx < 0 {
					res = append(res, &apiv2.PermissionsByVisibility{
						Visibility: &apiv2.PermissionsByVisibility_Machine{
							Machine: &apiv2.MachinePermissions{
								Uuid:    perm.Subject,
								Methods: []string{method},
							},
						},
					})

					continue
				}

				res[idx].GetMachine().Methods = append(res[idx].GetMachine().Methods, method)

				continue
			}

			if _, ok := permissions.GetServicePermissions().Visibility.Project[method]; ok {
				idx := slices.IndexFunc(res, func(p *apiv2.PermissionsByVisibility) bool {
					return p.GetProject() != nil && p.GetProject().Project == perm.Subject
				})

				if idx < 0 {
					res = append(res, &apiv2.PermissionsByVisibility{
						Visibility: &apiv2.PermissionsByVisibility_Project{
							Project: &apiv2.ProjectPermissions{
								Project: perm.Subject,
								Methods: []string{method},
							},
						},
					})

					continue
				}

				res[idx].GetProject().Methods = append(res[idx].GetProject().Methods, method)

				continue
			}

			if _, ok := permissions.GetServicePermissions().Visibility.Public[method]; ok {
				idx := slices.IndexFunc(res, func(p *apiv2.PermissionsByVisibility) bool {
					return p.GetPublic() != nil
				})

				if idx < 0 {
					res = append(res, &apiv2.PermissionsByVisibility{
						Visibility: &apiv2.PermissionsByVisibility_Public{
							Public: &apiv2.PublicPermissions{
								Methods: []string{method},
							},
						},
					})

					continue
				}

				res[idx].GetPublic().Methods = append(res[idx].GetPublic().Methods, method)

				continue
			}

			if _, ok := permissions.GetServicePermissions().Visibility.Self[method]; ok {
				idx := slices.IndexFunc(res, func(p *apiv2.PermissionsByVisibility) bool {
					return p.GetSelf() != nil
				})

				if idx < 0 {
					res = append(res, &apiv2.PermissionsByVisibility{
						Visibility: &apiv2.PermissionsByVisibility_Self{
							Self: &apiv2.SelfPermissions{
								Methods: []string{method},
							},
						},
					})

					continue
				}

				res[idx].GetSelf().Methods = append(res[idx].GetSelf().Methods, method)

				continue
			}

			if _, ok := permissions.GetServicePermissions().Visibility.Tenant[method]; ok {
				idx := slices.IndexFunc(res, func(p *apiv2.PermissionsByVisibility) bool {
					return p.GetTenant() != nil && p.GetTenant().Login == perm.Subject
				})

				if idx < 0 {
					res = append(res, &apiv2.PermissionsByVisibility{
						Visibility: &apiv2.PermissionsByVisibility_Tenant{
							Tenant: &apiv2.TenantPermissions{
								Login:   perm.Subject,
								Methods: []string{method},
							},
						},
					})

					continue
				}

				res[idx].GetTenant().Methods = append(res[idx].GetTenant().Methods, method)

				continue
			}

			return nil, fmt.Errorf("method is not part of the api: %s", method)
		}
	}

	return res, nil
}
