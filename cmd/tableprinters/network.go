package tableprinters

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

type network struct {
	parent   *apiv2.Network
	children []*apiv2.Network
}

type networks []*network

func (t *TablePrinter) NetworkTable(data []*apiv2.Network, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Name", "Type", "Project", "Partition", "Nat", "Prefixes", "Prefix Usage", "IP Usage"}
	if wide {
		header = []string{"ID", "Description", "Name", "Type", "Project", "Partition", "Nat", "Prefixes", "Annotations"}
	}

	nn := &networks{}
	for _, n := range data {
		if n.ParentNetwork == nil {
			*nn = append(*nn, &network{parent: n})
		}
	}
	for _, n := range data {
		if n.ParentNetwork != nil {
			if !nn.appendChild(*n.ParentNetwork, n) {
				*nn = append(*nn, &network{parent: n})
			}
		}
	}
	for _, n := range *nn {
		rows = append(rows, addNetwork("", n.parent, wide))
		for i, c := range n.children {
			prefix := "├"
			if i == len(n.children)-1 {
				prefix = "└"
			}
			prefix += "─╴"
			rows = append(rows, addNetwork(prefix, c, wide))
		}
	}

	return header, rows, nil
}

func addNetwork(prefix string, n *apiv2.Network, wide bool) []string {
	id := fmt.Sprintf("%s%s", prefix, n.Id)

	prefixes := strings.Join(n.Prefixes, ",")
	shortIPUsage := nbr
	shortPrefixUsage := nbr
	ipv4Use := 0.0
	ipv4PrefixUse := 0.0
	ipv6Use := 0.0
	ipv6PrefixUse := 0.0

	if n.Consumption != nil {
		consumption := n.Consumption
		if consumption.Ipv4 != nil {
			ipv4Consumption := consumption.Ipv4
			ipv4Use = float64(ipv4Consumption.UsedIps) / float64(ipv4Consumption.AvailableIps)

			if ipv4Consumption.AvailablePrefixes > 0 {
				ipv4PrefixUse = float64(ipv4Consumption.UsedPrefixes) / float64(ipv4Consumption.AvailablePrefixes)
			}
		}
		if consumption.Ipv6 != nil {
			ipv6Consumption := consumption.Ipv6
			ipv6Use = float64(ipv6Consumption.UsedIps) / float64(ipv6Consumption.AvailableIps)

			if ipv6Consumption.AvailablePrefixes > 0 {
				ipv6PrefixUse = float64(ipv6Consumption.UsedPrefixes) / float64(ipv6Consumption.AvailablePrefixes)
			}
		}

		if ipv4Use >= 0.9 || ipv6Use >= 0.9 {
			shortIPUsage = color.RedString(threequarterpie)
		} else if ipv4Use >= 0.7 || ipv6Use >= 0.7 {
			shortIPUsage = color.YellowString(halfpie)
		} else {
			shortIPUsage = color.GreenString(dot)
		}

		if ipv4PrefixUse >= 0.9 || ipv6PrefixUse >= 0.9 {
			shortPrefixUsage = color.RedString(threequarterpie)
		} else if ipv4PrefixUse >= 0.7 || ipv6PrefixUse >= 0.7 {
			shortPrefixUsage = color.YellowString(halfpie)
		} else {
			shortPrefixUsage = color.GreenString(dot)
		}
	}

	var (
		description = pointer.SafeDeref(n.Description)
		name        = pointer.SafeDeref(n.Name)
		project     = pointer.SafeDeref(n.Project)
		partition   = pointer.SafeDeref(n.Partition)
		natType     = pointer.SafeDeref(n.NatType).String()
	)

	if t, err := enum.GetStringValue(pointer.SafeDeref(n.NatType)); err == nil {
		natType = *t
	} else {
		fmt.Println(err)
	}

	max := getMaxLineCount(description, name, project, partition, natType, prefixes, shortIPUsage)
	for range max - 1 {
		id += "\n│"
	}

	var as []string
	if n.Meta.Labels != nil {
		for k, v := range n.Meta.Labels.Labels {
			as = append(as, k+"="+v)
		}
	}

	annotations := strings.Join(as, "\n")

	var networkType string
	nt, err := enum.GetStringValue(n.Type)
	if err != nil {
		networkType = "unknown"
	} else {
		networkType = *nt
	}

	if wide {
		return []string{id, description, name, networkType, project, partition, natType, prefixes, annotations}
	} else {
		return []string{id, name, networkType, project, partition, natType, prefixes, shortPrefixUsage, shortIPUsage}
	}
}

func (nn *networks) appendChild(parentID string, child *apiv2.Network) bool {
	for _, n := range *nn {
		if n.parent.Id == parentID {
			n.children = append(n.children, child)
			return true
		}
	}
	return false
}
