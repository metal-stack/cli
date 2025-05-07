package tableprinters

import (
	"sort"

	"github.com/fatih/color"
	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/pointer"
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

			name, err := enum.GetStringValue(s.Name)
			if err != nil {
				name = pointer.Pointer("service status unknown")
			}
			message := "All systems operational"
			if s.Message != "" {
				message = s.Message
			}

			rows = append(rows, []string{statusIcon(s.Status), *name, message})

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
