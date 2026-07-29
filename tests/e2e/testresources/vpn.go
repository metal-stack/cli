package testresources

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	VpnNode1 = func() *apiv2.VPNNode {
		return &apiv2.VPNNode{
			Id:          1,
			Name:        "node-1",
			Project:     Project1().Uuid,
			IpAddresses: []string{"1.2.3.4", "1.2.3.5"},
			LastSeen:    timestamppb.New(e2e.TimeBubbleStartTime()),
			Online:      false,
		}
	}
	VpnNode2 = func() *apiv2.VPNNode {
		return &apiv2.VPNNode{
			Id:          2,
			Name:        "node-2",
			Project:     Project2().Uuid,
			IpAddresses: []string{"2.3.4.5"},
			LastSeen:    timestamppb.New(e2e.TimeBubbleStartTime().Add(1 * time.Minute)),
			Online:      true,
		}
	}
)
