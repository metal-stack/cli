package helpers

import (
	"net/netip"

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

func AddressFamilyFromPrefixes(prefixes ...string) (*apiv2.NetworkAddressFamily, error) {
	var (
		addressFamily *apiv2.NetworkAddressFamily
		isIPv4        bool
		isIPv6        bool
	)

	for _, pfx := range prefixes {
		p, err := netip.ParsePrefix(pfx)
		if err != nil {
			return nil, err
		}

		if p.Addr().Is4() {
			isIPv4 = true
		}
		if p.Addr().Is6() {
			isIPv6 = true
		}
	}

	switch {
	case isIPv4 && isIPv6:
		addressFamily = apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_DUAL_STACK.Enum()
	case isIPv4:
		addressFamily = apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_V4.Enum()
	case isIPv6:
		addressFamily = apiv2.NetworkAddressFamily_NETWORK_ADDRESS_FAMILY_V6.Enum()
	default:
		// noop
	}

	return addressFamily, nil
}
