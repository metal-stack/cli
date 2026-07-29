package testresources

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	tag "github.com/metal-stack/metal-lib/pkg/tag"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Network1 = func() *apiv2.Network {
		return &apiv2.Network{
			Id:          "6988ebb0-9531-4f9b-a893-d7868258e2ef",
			Name:        new("internet"),
			Description: new("internet network"),
			Partition:   new(Partition1().Id),
			Project:     new(Project1().Uuid),
			Type:        apiv2.NetworkType_NETWORK_TYPE_EXTERNAL,
			NatType:     apiv2.NATType_NAT_TYPE_NONE,
			Prefixes:    []string{"10.0.0.0/16", "2001:db8::/32"},
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						tag.ClusterServiceFQN: "<cluster>/default/ingress-nginx",
					},
				},
			},
		}
	}
	Network2 = func() *apiv2.Network {
		return &apiv2.Network{
			Id:          "d83ffb0a-7aa6-4a66-8e03-0b5ee8b718a0",
			Name:        new("private"),
			Description: new("private network"),
			Partition:   new(Partition1().Id),
			Project:     new(Project2().Uuid),
			Type:        apiv2.NetworkType_NETWORK_TYPE_CHILD,
			NatType:     apiv2.NATType_NAT_TYPE_IPV4_MASQUERADE,
			Prefixes:    []string{"192.168.1.0/24"},
			Vrf:         new(uint32(100)),
			Meta: &apiv2.Meta{
				CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"a": "b",
					},
				},
			},
		}
	}
)
