package testresources

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

var (
	Partition1 = func() *apiv2.Partition {
		return &apiv2.Partition{
			Id:          "5711axzf0-4218-p7x9-3f8f-d34912dx190q",
			Description: "partition 1",
			MgmtServiceAddresses: []string{
				"192.168.1.1:1234",
			},
			BootConfiguration: &apiv2.PartitionBootConfiguration{
				Commandline: "commandline",
				ImageUrl:    "https://images.metal-stack.io/os/image.tar.lz4",
				KernelUrl:   "https://images.metal-stack.io/os/kernel",
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
			Id:          "d34912dx1-90qe-3jf8-1kd0-5711axzf08sj",
			Description: "partition 2",
			MgmtServiceAddresses: []string{
				"192.168.1.2:1234",
			},
			BootConfiguration: &apiv2.PartitionBootConfiguration{
				Commandline: "commandline",
				ImageUrl:    "https://images.metal-stack.io/os/image.tar.lz4",
				KernelUrl:   "https://images.metal-stack.io/os/kernel",
			},
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: nil,
				},
			},
		}
	}
)
