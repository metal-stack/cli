package tableprinters

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func (t *TablePrinter) FilesystemTable(data []*apiv2.FilesystemLayout, wide bool) ([]string, [][]string, error) {
	var (
		header = []string{"ID", "Name", "Description", "Sizes", "Images"}
		rows   [][]string
	)

	if wide {
		header = append(header, "Filesystems", "Disks", "Raid", "Volume Groups", "Logical Volumes")
	}

	for _, fsl := range data {
		row := []string{fsl.Id, pointer.SafeDeref(fsl.Name), pointer.SafeDeref(fsl.Description)}

		var (
			sizes  []string
			images []string
		)
		if fsl.Constraints != nil {
			sizes = fsl.Constraints.Sizes
			for k, v := range fsl.Constraints.Images {
				images = append(images, fmt.Sprintf("%s:%s", k, v))
			}
			sort.Strings(images)
		}

		row = append(row, strings.Join(sizes, ","), strings.Join(images, ","))

		if wide {
			row = append(row,
				strconv.Itoa(len(fsl.Filesystems)),
				strconv.Itoa(len(fsl.Disks)),
				strconv.Itoa(len(fsl.Raid)),
				strconv.Itoa(len(fsl.VolumeGroups)),
				strconv.Itoa(len(fsl.LogicalVolumes)),
			)
		}

		rows = append(rows, row)
	}

	return header, rows, nil
}
