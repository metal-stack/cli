package common

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func IPAddressFamilyToType(af string) *apiv2.IPAddressFamily {
	switch af {
	case "":
		return nil
	case "ipv4", "IPv4":
		return apiv2.IPAddressFamily_IP_ADDRESS_FAMILY_V4.Enum()
	case "ipv6", "IPv6":
		return apiv2.IPAddressFamily_IP_ADDRESS_FAMILY_V6.Enum()
	default:
		return apiv2.IPAddressFamily_IP_ADDRESS_FAMILY_UNSPECIFIED.Enum()
	}
}

func NetworkAddressFamilyToType(af string) *apiv2.NetworkAddressFamily {
	switch af {
	case "":
		return nil
	case "ipv4", "IPv4":
		return apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_V4.Enum()
	case "ipv6", "IPv6":
		return apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_V6.Enum()
	case "dual-stack":
		return apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_DUAL_STACK.Enum()
	default:
		return apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_UNSPECIFIED.Enum()
	}
}
