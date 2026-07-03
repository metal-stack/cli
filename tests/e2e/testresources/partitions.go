package testresources

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

var (
	Partition1 = func() *apiv2.Partition {
		return &apiv2.Partition{
			Id:          "partition-1",
			Description: "partition 1",
			MgmtServiceAddresses: []string{
				"192.168.1.1:1234",
			},
			BootConfiguration: &apiv2.PartitionBootConfiguration{
				Commandline: "commandline",
				ImageUrl:    "imageurl",
				KernelUrl:   "kernelurl",
			},
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"a": "b",
					},
				},
			},
		}
	}
	Partition2 = func() *apiv2.Partition {
		return &apiv2.Partition{
			Id:          "partition-2",
			Description: "partition 2",
			MgmtServiceAddresses: []string{
				"192.168.1.2:1234",
			},
			BootConfiguration: &apiv2.PartitionBootConfiguration{
				Commandline: "commandline",
				ImageUrl:    "imageurl",
				KernelUrl:   "kernelurl",
			},
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: nil,
				},
			},
		}
	}
)
