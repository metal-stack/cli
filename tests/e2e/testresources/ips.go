package testresources

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/tag"
)

var (
	IP1 = func() *apiv2.IP {
		return &apiv2.IP{
			Uuid:        "2e0144a2-09ef-42b7-b629-4263295db6e8",
			Network:     "internet",
			Ip:          "1.1.1.1",
			Name:        "a",
			Description: "a description",
			Project:     "ce19a655-7933-4745-8f3e-9592b4a90488",
			Type:        apiv2.IPType_IP_TYPE_STATIC,
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						tag.ClusterServiceFQN: "<cluster>/default/ingress-nginx",
					},
				},
			},
		}
	}
	IP2 = func() *apiv2.IP {
		return &apiv2.IP{
			Uuid:        "9cef40ec-29c6-4dfa-aee8-47ee1f49223d",
			Network:     "internet",
			Ip:          "4.3.2.1",
			Name:        "b",
			Description: "b description",
			Project:     "46bdfc45-9c8d-4268-b359-b40e3079d384",
			Type:        apiv2.IPType_IP_TYPE_EPHEMERAL,
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"a": "b",
					},
				},
			},
		}
	}
)
