package tableprinters

import (
	"strings"

	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/pointer"
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
		var (
			t, _            = enum.GetStringValue(ip.Type)
			attachedService = ""
		)

		var labels []string
		if ip.Meta != nil && ip.Meta.Labels != nil {
			labels = genericcli.MapToLabels(ip.Meta.Labels.Labels)
		}

		if wide {
			rows = append(rows, []string{ip.Ip, ip.Project, ip.Uuid, pointer.SafeDeref(t), ip.Name, ip.Description, strings.Join(labels, "\n")})
		} else {
			rows = append(rows, []string{ip.Ip, ip.Project, ip.Uuid, pointer.SafeDeref(t), ip.Name, attachedService})
		}
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
