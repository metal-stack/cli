package tableprinters

import (
	"fmt"
	"strings"
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (t *TablePrinter) VPNTable(data []*apiv2.VPNNode, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "Name", "Project", "IPs", "Last Seed"}

	for _, node := range data {
		lastSeen := node.LastSeen.AsTime().Format(time.DateTime + " MST")
		ips := strings.Join(node.IpAddresses, ",")

		row := []string{
			fmt.Sprintf("%d", node.Id),
			node.Name,
			node.Project,
			ips,
			lastSeen,
		}

		rows = append(rows, row)
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
