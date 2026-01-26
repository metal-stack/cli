package tableprinters

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"

	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/viper"
)

func (t *TablePrinter) SwitchTable(switches []*apiv2.Switch, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Partition", "Rack", "OS", "Status", "Last Sync"}
	if wide {
		header = []string{"ID", "Partition", "Rack", "OS", "Metalcore", "IP", "Mode", "Last Sync", "Sync Duration", "Last Error"}

		t.t.DisableAutoWrap(true)
	}

	for _, s := range switches {
		var (
			id        = s.Id
			partition = s.Partition
			rack      = pointer.SafeDeref(s.Rack)

			syncTime    time.Time
			syncLast    = ""
			syncDurStr  = ""
			lastError   = ""
			shortStatus = nbr
			allUp       = true
		)

		for _, c := range s.MachineConnections {
			if c.Nic == nil {
				continue
			}

			if c.Nic.State == nil {
				allUp = false
				lastError = fmt.Sprintf("port status of %q is unknown", c.Nic.Name)
				break
			}

			desired := c.Nic.State.Desired
			actual := c.Nic.State.Actual
			allUp = allUp && actual == apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_UP

			if desired != nil && actual != *desired {
				lastError = fmt.Sprintf("%q is %s but should be %s", c.Nic.Name, c.Nic.State.Actual, desired)
				break
			}

			if !allUp {
				lastError = fmt.Sprintf("%q is %s", c.Nic.Name, c.Nic.State.Actual)
				break
			}
		}

		if s.LastSync != nil {
			var (
				syncAge time.Duration
				syncDur time.Duration
			)

			if s.LastSync.Time != nil && !s.LastSync.Time.AsTime().IsZero() {
				syncTime = s.LastSync.Time.AsTime()
				syncAge = time.Since(syncTime)
			}
			if s.LastSync.Duration != nil {
				syncDur = s.LastSync.Duration.AsDuration().Round(time.Millisecond)
			}

			switch {
			case syncAge >= 10*time.Minute, syncDur >= 30*time.Second:
				shortStatus = color.RedString(dot)
			case syncAge >= time.Minute, syncDur >= 20*time.Second, !allUp:
				shortStatus = color.YellowString(dot)
			default:
				shortStatus = color.GreenString(dot)
			}

			if syncAge > 0 {
				syncLast = humanizeDuration(syncAge) + " ago"
			}
			if syncDur > 0 {
				syncDurStr = fmt.Sprintf("%v", syncDur)
			}
		}

		if s.LastSyncError != nil {
			var (
				errorTime time.Time
				error     string
			)

			if s.LastSyncError.Time != nil {
				errorTime = s.LastSyncError.Time.AsTime()
			}
			if s.LastSyncError.Error != nil {
				error = *s.LastSyncError.Error
			}
			// after 7 days we do not show sync errors anymore
			if !errorTime.IsZero() && time.Since(errorTime) < 7*24*time.Hour {
				lastError = fmt.Sprintf("%s ago: %s", humanizeDuration(time.Since(errorTime)), error)

				if errorTime.After(syncTime) {
					shortStatus = color.RedString(dot)
				}
			}
		}

		var mode string
		switch s.ReplaceMode {
		case apiv2.SwitchReplaceMode_SWITCH_REPLACE_MODE_REPLACE:
			shortStatus = nbr + color.RedString(dot)
			mode = "replace"
		default:
			mode = "operational"
		}

		os := ""
		osIcon := ""
		metalCore := ""
		if s.Os != nil {
			switch s.Os.Vendor {
			case apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_CUMULUS:
				osIcon = "🐢"
			case apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_SONIC:
				osIcon = "🦔"
			default:
				osIcon = s.Os.Vendor.String()
			}

			os = s.Os.Vendor.String()
			if s.Os.Version != "" {
				os = fmt.Sprintf("%s (%s)", os, s.Os.Version)
			}
			// metal core version is very long: v0.9.1 (1d5e42ea), tags/v0.9.1-0-g1d5e42e, go1.20.5
			metalCore = strings.Split(s.Os.MetalCoreVersion, ",")[0]
		}

		if wide {
			rows = append(rows, []string{id, partition, rack, os, metalCore, s.ManagementIp, mode, syncLast, syncDurStr, lastError})
		} else {
			rows = append(rows, []string{id, partition, rack, osIcon, shortStatus, syncLast})
		}
	}

	return header, rows, nil
}

type SwitchesWithMachines struct {
	Switches []*apiv2.Switch           `yaml:"switches"`
	Machines map[string]*apiv2.Machine `yaml:"machines"`
}

func (t *TablePrinter) SwitchWithConnectedMachinesTable(data *SwitchesWithMachines, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "NIC Name", "Identifier", "Partition", "Rack", "Size", "Product Serial", "Chassis Serial"}
	if wide {
		header = []string{"ID", "", "NIC Name", "Identifier", "Partition", "Rack", "Size", "Hostname", "Product Serial", "Chassis Serial"}
	}

	t.t.DisableAutoWrap(true)

	for _, s := range data.Switches {
		rack := pointer.SafeDeref(s.Rack)

		if wide {
			rows = append(rows, []string{s.Id, "", "", "", s.Partition, rack})
		} else {
			rows = append(rows, []string{s.Id, "", "", s.Partition, rack})
		}

		conns := s.MachineConnections
		if viper.IsSet("size") || viper.IsSet("machine-id") {
			filteredConns := []*apiv2.MachineConnection{}

			for _, conn := range s.MachineConnections {
				m, ok := data.Machines[conn.MachineId]
				if !ok {
					continue
				}

				if viper.IsSet("machine-id") && m.Uuid == viper.GetString("machine-id") {
					filteredConns = append(filteredConns, conn)
				}

				if viper.IsSet("size") && m.Size.Id == viper.GetString("size") {
					filteredConns = append(filteredConns, conn)
				}
			}

			conns = filteredConns
		}

		sort.Slice(conns, switchInterfaceNameLessFunc(conns))

		for i, conn := range conns {
			if conn == nil {
				continue
			}

			prefix := "├"
			if i == len(conns)-1 {
				prefix = "└"
			}
			prefix += "─╴"

			nic := pointer.SafeDeref(conn.Nic)
			m, ok := data.Machines[conn.MachineId]
			if !ok {
				return nil, nil, fmt.Errorf("switch port %s is connected to a machine which does not exist: %q", nic.Name, conn.MachineId)
			}

			identifier := nic.Identifier
			if identifier == "" {
				identifier = nic.Mac
			}

			nicname := nic.Name
			nicstate := pointer.SafeDeref(nic.State).Actual
			bgpstate := pointer.SafeDeref(nic.BgpPortState)
			if nicstate != apiv2.SwitchPortStatus_SWITCH_PORT_STATUS_UP {
				nicname = fmt.Sprintf("%s (%s)", nicname, color.RedString(nicstate.String()))
			}
			if wide {
				switch bgpstate.BgpState {
				case apiv2.BGPState_BGP_STATE_ESTABLISHED:
					uptime := time.Since(time.Unix(pointer.SafeDeref(bgpstate.BgpTimerUpEstablished).Seconds, 0))
					nicname = fmt.Sprintf("%s (BGP:%s(%s))", nicname, bgpstate.BgpState, uptime)
				default:
					nicname = fmt.Sprintf("%s (BGP:%s)", nicname, bgpstate.BgpState)
				}
			}

			if wide {
				// TODO: add emojis once machine functions are implemented
				// emojis, _ := t.getMachineStatusEmojis(m.Liveliness, m.Events, m.State, pointer.SafeDeref(m.Allocation).Vpn)

				rows = append(rows, []string{
					fmt.Sprintf("%s%s", prefix, m.Uuid),
					// emojis,
					nicname,
					identifier,
					pointer.SafeDeref(m.Partition).Id,
					m.Rack,
					pointer.SafeDeref(m.Size).Id,
					pointer.SafeDeref(m.Allocation).Hostname,
					// TODO: where to get ipmi information?
					// pointer.SafeDeref(pointer.SafeDeref(m.Ipmi).Fru).ProductSerial,
					// pointer.SafeDeref(pointer.SafeDeref(m.Ipmi).Fru).ChassisPartSerial,
				})
			} else {
				rows = append(rows, []string{
					fmt.Sprintf("%s%s", prefix, m.Uuid),
					nicname,
					identifier,
					pointer.SafeDeref(m.Partition).Id,
					m.Rack,
					pointer.SafeDeref(m.Size).Id,
					// TODO: where to get ipmi information?
					// pointer.SafeDeref(pointer.SafeDeref(m.Ipmi).Fru).ProductSerial,
					// pointer.SafeDeref(pointer.SafeDeref(m.Ipmi).Fru).ChassisPartSerial,
				})
			}
		}
	}

	return header, rows, nil
}

var numberRegex = regexp.MustCompile("([0-9]+)")

func switchInterfaceNameLessFunc(conns []*apiv2.MachineConnection) func(i, j int) bool {
	return func(i, j int) bool {
		var (
			a = pointer.SafeDeref(pointer.SafeDeref(conns[i]).Nic).Name
			b = pointer.SafeDeref(pointer.SafeDeref(conns[j]).Nic).Name

			aMatch = numberRegex.FindAllStringSubmatch(a, -1)
			bMatch = numberRegex.FindAllStringSubmatch(b, -1)
		)

		for i := range aMatch {
			if i >= len(bMatch) {
				return true
			}

			interfaceNumberA, aErr := strconv.Atoi(aMatch[i][0])
			interfaceNumberB, bErr := strconv.Atoi(bMatch[i][0])

			if aErr == nil && bErr == nil {
				if interfaceNumberA < interfaceNumberB {
					return true
				}
				if interfaceNumberA != interfaceNumberB {
					return false
				}
			}
		}

		return a < b
	}
}

type SwitchDetail struct {
	*apiv2.Switch
}

func (t *TablePrinter) SwitchDetailTable(data []*SwitchDetail, wide bool) ([]string, [][]string, error) {
	var (
		header = []string{"Partition", "Rack", "Switch", "Port", "Machine", "VNI-Filter", "CIDR-Filter"}
		rows   [][]string
	)

	for _, sw := range data {
		filterBySwp := map[string]*apiv2.BGPFilter{}
		for _, nic := range sw.Nics {
			if nic == nil {
				continue
			}

			if nic.BgpFilter != nil {
				filterBySwp[nic.Name] = nic.BgpFilter
			}
		}

		for _, conn := range sw.MachineConnections {
			if conn == nil {
				continue
			}

			nicName := pointer.SafeDeref(conn.Nic).Name

			f := filterBySwp[nicName]
			row := []string{sw.Partition, pointer.SafeDeref(sw.Rack), sw.Id, nicName, conn.MachineId}
			row = append(row, filterColumns(f, 0)...)
			max := len(f.Vnis)
			if len(f.Cidrs) > max {
				max = len(f.Cidrs)
			}
			rows = append(rows, row)
			for i := 1; i < max; i++ {
				row = append([]string{"", "", "", "", ""}, filterColumns(f, i)...)
				rows = append(rows, row)
			}
		}
	}

	return header, rows, nil
}

func filterColumns(filter *apiv2.BGPFilter, i int) []string {
	if filter == nil {
		return nil
	}

	v := ""
	if len(filter.Vnis) > i {
		v = filter.Vnis[i]
	}
	c := ""
	if len(filter.Cidrs) > i {
		c = filter.Cidrs[i]
	}
	return []string{v, c}
}
