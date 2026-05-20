package testresources

import apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"

var (
	User = func() *apiv2.User {
		return &apiv2.User{
			Name:          "Larry",
			Email:         "larry@metal-stack.io",
			Login:         "larry@metal-stack.io@openid-connect",
			AvatarUrl:     "",
			DefaultTenant: Tenant1(),
			Tenants:       []*apiv2.Tenant{Tenant1(), Tenant2()},
			Projects:      []*apiv2.Project{Project1(), Project2()},
		}
	}
)
