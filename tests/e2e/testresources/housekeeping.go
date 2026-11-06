package testresources

import (
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

var (
	MachineWithLiveliness1 = &adminv2.MachineWithLiveliness{
		Uuid:       Machine1().Uuid,
		Liveliness: apiv2.MachineLiveliness_MACHINE_LIVELINESS_ALIVE,
	}
	MachineWithLiveliness2 = &adminv2.MachineWithLiveliness{
		Uuid:       Machine2().Uuid,
		Liveliness: apiv2.MachineLiveliness_MACHINE_LIVELINESS_DEAD,
	}
)
