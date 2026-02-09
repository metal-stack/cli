package tableprinters

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"

	"github.com/metal-stack/metal-lib/pkg/pointer"
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
				actualString, err := enum.GetStringValue(actual)
				if err != nil {
					return nil, nil, err
				}
				desiredString, err := enum.GetStringValue(*desired)
				if err != nil {
					return nil, nil, err
				}
				lastError = fmt.Sprintf("%q is %s but should be %s", c.Nic.Name, *actualString, *desiredString)
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

		var (
			os        string
			osIcon    string
			metalCore string
		)

		if s.Os != nil {
			switch s.Os.Vendor {
			case apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_CUMULUS:
				osIcon = "🐢"
			case apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_SONIC:
				osIcon = "🦔"
			default:
				osIcon = s.Os.Vendor.String()
			}

			if s.Os.Vendor != apiv2.SwitchOSVendor_SWITCH_OS_VENDOR_UNSPECIFIED {
				osString, err := enum.GetStringValue(s.Os.Vendor)
				if err != nil {
					return nil, nil, err
				}
				os = *osString
			}

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

func (t *TablePrinter) SwitchWithConnectedMachinesTable(res []*apiv2.SwitchWithMachines, wide bool) ([]string, [][]string, error) {
	var (
		header []string
		rows   [][]string
	)

	header = []string{"ID", "NIC Name", "Identifier", "Partition", "Rack", "Size", "Product Serial", "Chassis Serial"}
	if wide {
		header = []string{"ID", "", "NIC Name", "Identifier", "Partition", "Rack", "Size", "Hostname", "Product Serial", "Chassis Serial"}
	}

	for _, sm := range res {
		var (
			partition = sm.Partition
			rack      = sm.Rack
		)

		if wide {
			rows = append(rows, []string{sm.Id, "", "", "", partition, rack})
		} else {
			rows = append(rows, []string{sm.Id, "", "", partition, rack})
		}

		sort.Slice(sm.Connections, switchInterfaceNameLessFunc(sm.Connections))
		for i, con := range sm.Connections {
			if con.Nic == nil || con.Machine == nil || con.Fru == nil {
				continue
			}

			prefix := "├"
			if i == len(sm.Connections)-1 {
				prefix = "└"
			}
			prefix += "─╴"

			var (
				nic                  = con.Nic
				machine              = con.Machine
				fru                  = con.Fru
				nicName              = nic.Name
				nicIdentifier        = nic.Identifier
				fruProductSerial     = pointer.SafeDeref(fru.ProductSerial)
				fruChassisPartSerial = pointer.SafeDeref(fru.ChassisPartSerial)

				machineSize        string
				allocationHostname string
			)

			if nic.State != nil {
				state, err := enum.GetStringValue(nic.State.Actual)
				if err != nil {
					return nil, nil, err
				}
				nicName = fmt.Sprintf("%s (%s)", nicName, *state)
			}

			if nic.BgpPortState != nil && wide {
				if nic.BgpPortState.BgpState == apiv2.BGPState_BGP_STATE_ESTABLISHED {
					up := humanizeDuration(time.Since(nic.BgpPortState.BgpTimerUpEstablished.AsTime()))
					nicName = fmt.Sprintf("%s (BGP:Established(%s))", nicName, up)
				} else {
					state, err := enum.GetStringValue(nic.BgpPortState.BgpState)
					if err != nil {
						return nil, nil, err
					}
					nicName = fmt.Sprintf("%s (BGP:%s)", nicName, *state)
				}
			}

			if machine.Size != nil {
				machineSize = machine.Size.Id
			}

			if machine.Allocation != nil {
				allocationHostname = machine.Allocation.Hostname
			}

			if wide {
				rows = append(rows, []string{fmt.Sprintf("%s%s", prefix, machine.Uuid), t.getMachineStatusEmojis(machine), nicName, nicIdentifier, partition, rack, machineSize, allocationHostname, fruProductSerial, fruChassisPartSerial})
			} else {
				rows = append(rows, []string{fmt.Sprintf("%s%s", prefix, machine.Uuid), nicName, nicIdentifier, partition, rack, machineSize, fruProductSerial, fruChassisPartSerial})
			}
		}
	}

	t.t.DisableAutoWrap(true)
	return header, rows, nil
}

type SwitchDetail struct {
	*apiv2.Switch
}

func (t *TablePrinter) SwitchDetailTable(switches []SwitchDetail) ([]string, [][]string, error) {
	var (
		header = []string{"Partition", "Rack", "Switch", "Port", "Machine", "VNI-Filter", "CIDR-Filter"}
		rows   [][]string
	)

	for _, sw := range switches {
		filterByNic := map[string]*apiv2.BGPFilter{}
		for _, nic := range sw.Nics {
			if nic == nil {
				continue
			}

			if nic.BgpFilter != nil {
				filterByNic[nic.Name] = nic.BgpFilter
			}
		}

		for _, conn := range sw.MachineConnections {
			if conn == nil || conn.Nic == nil {
				continue
			}

			filter := filterByNic[conn.Nic.Name]
			row := append([]string{sw.Partition, pointer.SafeDeref(sw.Rack), sw.Id, conn.Nic.Name, conn.MachineId}, filterColumns(filter, 0)...)
			rows = append(rows, row)

			if filter == nil {
				continue
			}

			max := math.Max(float64(len(filter.Cidrs)), float64(len(filter.Vnis)))
			for i := 1; i < int(max); i++ {
				row = append([]string{"", "", "", "", ""}, filterColumns(filter, i)...)
				rows = append(rows, row)
			}
		}
	}

	return header, rows, nil
}

func (t *TablePrinter) getMachineStatusEmojis(m *apiv2.Machine) string {
	if m == nil {
		return ""
	}

	var (
		emojis []string
	)

	if status := m.Status; status != nil {
		switch status.Liveliness {
		case apiv2.MachineLiveliness_MACHINE_LIVELINESS_ALIVE:
			// noop
		case apiv2.MachineLiveliness_MACHINE_LIVELINESS_DEAD:
			emojis = append(emojis, skull)
		default:
			emojis = append(emojis, question)
		}

		if status.Condition != nil {
			switch status.Condition.State {
			case apiv2.MachineState_MACHINE_STATE_LOCKED:
				emojis = append(emojis, lock)
			case apiv2.MachineState_MACHINE_STATE_RESERVED:
				emojis = append(emojis, bark)
			default:
				// noop
			}
		}
	}

	if events := m.RecentProvisioningEvents; events != nil {
		switch events.State {
		case apiv2.MachineProvisioningEventState_MACHINE_PROVISIONING_EVENT_STATE_FAILED_RECLAIM:
			emojis = append(emojis, ambulance)
		case apiv2.MachineProvisioningEventState_MACHINE_PROVISIONING_EVENT_STATE_CRASHLOOP:
			emojis = append(emojis, loop)
		default:
			// noop

		}

		if time.Since(events.LastErrorEvent.Time.AsTime()) < t.lastEventErrorThreshold {
			emojis = append(emojis, exclamation)
		}
	}

	if m.Allocation != nil && m.Allocation.Vpn != nil && m.Allocation.Vpn.Connected {
		emojis = append(emojis, vpn)
	}

	return strings.Join(emojis, nbr)
}

func filterColumns(filter *apiv2.BGPFilter, i int) []string {
	var (
		vni  string
		cidr string
	)

	if filter == nil {
		return nil
	}

	if len(filter.Vnis) > i {
		vni = filter.Vnis[i]
	}
	if len(filter.Cidrs) > i {
		cidr = filter.Cidrs[i]
	}

	return []string{vni, cidr}
}

func switchInterfaceNameLessFunc(conns []*apiv2.SwitchNicWithMachine) func(i, j int) bool {
	numberRegex := regexp.MustCompile("([0-9]+)")

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
