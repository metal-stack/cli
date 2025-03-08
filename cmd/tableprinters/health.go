package tableprinters

import (
	"sort"

	"github.com/fatih/color"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (t *TablePrinter) HealthTable(data []*apiv2.Health, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"", "Name", "Message"}

		statusIcon = func(s apiv2.ServiceStatus) string {
			switch s {
			case apiv2.ServiceStatus_SERVICE_STATUS_HEALTHY:
				return color.GreenString("✔")
			case apiv2.ServiceStatus_SERVICE_STATUS_DEGRADED:
				return color.YellowString("✗")
			case apiv2.ServiceStatus_SERVICE_STATUS_UNHEALTHY:
				return color.RedString("✗")
			case apiv2.ServiceStatus_SERVICE_STATUS_UNSPECIFIED:
				return color.YellowString("?")
			default:
				return color.YellowString("?")
			}
		}
	)

	for _, h := range data {
		h := h
		for _, s := range h.Services {
			s := s

			name := ""
			switch s.Name {
			case apiv2.Service_SERVICE_MACHINES:
				name = "Machines"
			case apiv2.Service_SERVICE_IPAM:
				name = "Ipam"
			case apiv2.Service_SERVICE_RETHINK:
				name = "RethinkDB"
			case apiv2.Service_SERVICE_VPN:
				name = "VPN"
			case apiv2.Service_SERVICE_AUDIT:
				name = "Audit"
			case apiv2.Service_SERVICE_UNSPECIFIED:
				name = "Unspecified"
			default:
				name = s.Name.String()
			}

			message := "All systems operational"
			if s.Message != "" {
				message = s.Message
			}

			rows = append(rows, []string{statusIcon(s.Status), name, message})

			type partitionStatus struct {
				ID string
				*apiv2.PartitionHealth
			}

			var partitions []partitionStatus
			for id, p := range s.Partitions {
				p := p

				partitions = append(partitions, partitionStatus{
					ID:              id,
					PartitionHealth: p,
				})
			}

			sort.Slice(partitions, func(i, j int) bool {
				return partitions[i].ID < partitions[j].ID
			})

			for i, status := range partitions {
				status := status

				prefix := "├"
				if i == len(partitions)-1 {
					prefix = "└"
				}
				prefix += "─╴"

				message := "All systems operational"
				if s.Message != "" {
					message = s.Message
				}

				rows = append(rows, []string{
					statusIcon(status.Status),
					prefix + status.ID,
					message,
				})
			}
		}
	}

	return header, rows, nil
}
