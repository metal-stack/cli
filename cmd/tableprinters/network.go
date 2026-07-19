package tableprinters

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (t *TablePrinter) NetworkTable(data []*apiv2.Network, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"ID", "Name", "Type", "Partition", "Project", "Prefixes"}
	)

	if wide {
		header = []string{"ID", "Name", "Type", "Partition", "Project", "Prefixes", "Dest. Prefixes", "VRF"}
	}

	for _, n := range data {
		name := ""
		if n.Name != nil {
			name = *n.Name
		}
		partition := ""
		if n.Partition != nil {
			partition = *n.Partition
		}
		project := ""
		if n.Project != nil {
			project = *n.Project
		}
		vrf := ""
		if n.Vrf != nil {
			vrf = formatUint32(*n.Vrf)
		}

		prefixes := n.Prefixes
		if len(prefixes) > 3 {
			prefixes = append(prefixes[:3], "...")
		}

		if wide {
			rows = append(rows, []string{n.Id, name, n.Type.String(), partition, project, joinOrEmpty(n.Prefixes), joinOrEmpty(n.DestinationPrefixes), vrf})
		} else {
			rows = append(rows, []string{n.Id, name, n.Type.String(), partition, project, joinOrEmpty(prefixes)})
		}
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
