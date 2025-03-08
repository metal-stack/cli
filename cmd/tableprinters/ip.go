package tableprinters

import (
	"fmt"
	"strings"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/olekukonko/tablewriter"
)

func (t *TablePrinter) IPTable(data []*apiv2.IP, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"IP", "Project", "ID", "Type", "Name", "Attached Service"}
	)

	if wide {
		header = []string{"IP", "Project", "ID", "Type", "Name", "Description", "Labels"}
	}

	for _, ip := range data {
		ip := ip

		var t string

		switch ip.Type {
		case apiv2.IPType_IP_TYPE_EPHEMERAL:
			t = "ephemeral"
		case apiv2.IPType_IP_TYPE_STATIC:
			t = "static"
		case apiv2.IPType_IP_TYPE_UNSPECIFIED:
			t = "unspecified"
		default:
			t = ip.Type.String()
		}

		attachedService := ""

		var labels []string
		if ip.Meta != nil && ip.Meta.Labels != nil {
			for k, v := range ip.Meta.Labels.Labels {
				labels = append(labels, fmt.Sprintf("%s=%s", k, v))
			}
		}

		if wide {
			rows = append(rows, []string{ip.Ip, ip.Project, ip.Uuid, t, ip.Name, ip.Description, strings.Join(labels, "\n")})
		} else {
			rows = append(rows, []string{ip.Ip, ip.Project, ip.Uuid, t, ip.Name, attachedService})
		}
	}

	t.t.MutateTable(func(table *tablewriter.Table) {
		table.SetAutoWrapText(false)
	})

	return header, rows, nil
}
