package common

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func IpStaticToType(b bool) apiv2.IPType {
	if b {
		return apiv2.IPType_IP_TYPE_STATIC
	}
	return apiv2.IPType_IP_TYPE_EPHEMERAL
}
