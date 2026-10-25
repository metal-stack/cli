package testresources

import apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"

var (
	FilesystemLayout1 = func() *apiv2.FilesystemLayout {
		return &apiv2.FilesystemLayout{
			Id:          "debian-default",
			Name:        new("Debian Default"),
			Description: new("Default layout for debian machines"),
			Meta:        &apiv2.Meta{},
			Filesystems: []*apiv2.Filesystem{
				{
					Device: "/dev/sda1",
					Format: apiv2.Format_FORMAT_EXT4,
					Path:   new("/"),
				},
			},
			Disks: []*apiv2.Disk{
				{
					Device: "/dev/sda",
					Partitions: []*apiv2.DiskPartition{
						{
							Number:  1,
							Size:    1024,
							GptType: new(apiv2.GPTType_GPT_TYPE_LINUX),
						},
					},
				},
			},
			Constraints: &apiv2.FilesystemLayoutConstraints{
				Sizes:  []string{"m2.small"},
				Images: map[string]string{"debian": ">= 12.0"},
			},
		}
	}
	FilesystemLayout2 = func() *apiv2.FilesystemLayout {
		return &apiv2.FilesystemLayout{
			Id:          "ubuntu-raid1",
			Name:        new("Ubuntu Raid1"),
			Description: new("Raid1 layout for ubuntu machines"),
			Meta:        &apiv2.Meta{},
			Raid: []*apiv2.Raid{
				{
					ArrayName: "md0",
					Devices:   []string{"/dev/sda", "/dev/sdb"},
					Level:     apiv2.RaidLevel_RAID_LEVEL_1,
				},
			},
			VolumeGroups: []*apiv2.VolumeGroup{
				{
					Name:    "vg0",
					Devices: []string{"/dev/md0"},
				},
			},
			LogicalVolumes: []*apiv2.LogicalVolume{
				{
					Name:        "lv-root",
					VolumeGroup: "vg0",
					Size:        2048,
					LvmType:     apiv2.LVMType_LVM_TYPE_LINEAR,
				},
			},
			Constraints: &apiv2.FilesystemLayoutConstraints{
				Sizes:  []string{"m2.large"},
				Images: map[string]string{"ubuntu": ">= 24.0"},
			},
		}
	}
)
