package tableprinters

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/helpers"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func (t *TablePrinter) MachineTable(data []*apiv2.Machine, wide bool) ([]string, [][]string, error) {

	var (
		rows   [][]string
		header = []string{"ID", "", "Last Event", "When", "Age", "Hostname", "Project", "Size", "Image", "Partition", "Rack"}
	)

	if wide {
		header = []string{"ID", "Last Event", "When", "Age", "Description", "Name", "Hostname", "Project", "Ips", "Size", "Image", "Partition", "Rack", "Started", "Tags", "State"}
	}

	for _, machine := range data {
		machineID := machine.Uuid

		if machine.Status != nil && machine.Status.LedState != nil && machine.Status.LedState.Value == "LED-ON" {
			blue := color.New(color.FgBlue).SprintFunc()
			machineID = blue(machineID)
		}

		var (
			alloc             = pointer.SafeDeref(machine.Allocation)
			sizeID            = pointer.SafeDeref(machine.Size).Id
			partitionID       = pointer.SafeDeref(machine.Partition).Id
			project           = alloc.Project
			name              = alloc.Name
			desc              = alloc.Description
			hostname          = alloc.Hostname
			image             = pointer.SafeDeref(pointer.SafeDeref(alloc.Image).Name)
			rack              = machine.Rack
			truncatedHostname = genericcli.TruncateEnd(hostname, 30)

			nwIPs []string
		)

		for _, nw := range alloc.Networks {
			nwIPs = append(nwIPs, nw.Ips...)
		}

		var (
			ips       = strings.Join(nwIPs, "\n")
			started   = ""
			age       = ""
			tags      = ""
			reserved  = ""
			lastEvent = ""
			when      = ""
		)

		if alloc.Meta != nil && alloc.Meta.CreatedAt != nil && !alloc.Meta.CreatedAt.AsTime().IsZero() {
			started = alloc.Meta.CreatedAt.AsTime().Format(time.RFC3339)
			age = humanizeDuration(time.Since(alloc.Meta.CreatedAt.AsTime()))
		}

		if machine.Meta.Labels != nil && len(machine.Meta.Labels.Labels) > 0 {
			var labels []string
			for k, v := range machine.Meta.Labels.Labels {
				labels = append(labels, k+"="+v)
			}
			tags = strings.Join(labels, ",")
		}

		if machine.Status.Condition != nil {
			stateString, err := enum.GetStringValue(machine.Status.Condition.State)
			if err != nil {
				return nil, nil, err
			}

			reserved = *stateString
			if machine.Status.Condition.Description != "" {
				reserved += ":" + machine.Status.Condition.Description
			}
		}

		if len(machine.RecentProvisioningEvents.Events) > 0 {
			since := time.Since(machine.RecentProvisioningEvents.LastEventTime.AsTime())
			when = humanizeDuration(since)
			lastEventString, err := enum.GetStringValue(machine.RecentProvisioningEvents.Events[0].Event)
			if err != nil {
				return nil, nil, err
			}
			lastEvent = *lastEventString
		}

		emojis := t.getMachineStatusEmojis(machine)

		if wide {
			rows = append(rows, []string{machineID, lastEvent, when, age, desc, name, hostname, project, ips, sizeID, image, partitionID, rack, started, tags, reserved})
		} else {
			rows = append(rows, []string{machineID, emojis, lastEvent, when, age, truncatedHostname, project, sizeID, image, partitionID, rack})
		}
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}

func (t *TablePrinter) MachineBMCTable(data []*apiv2.MachineBMCDetails, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"ID", "", "Power", "IP", "Mac", "Board Part Number", "Bios", "BMC", "Size", "Partition", "Rack", "Updated"}
	)

	if wide {
		header = []string{"ID", "LED", "Power", "IP", "Mac", "Board Part Number", "Chassis Serial", "Product Serial", "Bios Version", "BMC Version", "Size", "Partition", "Rack", "Updated"}
	}

	for _, detail := range data {
		var (
			size      = detail.Size
			partition = detail.Partition
			rack      = detail.Rack
			machineID = detail.Uuid

			ipAddress     = ""
			mac           = ""
			bpn           = ""
			chassisSerial = ""
			productSerial = ""
			bmcVersion    = ""
			lastUpdated   = "never"
			biosVersion   = ""
			power         = ""
			powerText     = ""
			ledState      = ""

			emojis []string
		)

		if report := detail.BmcReport; report != nil {
			if led := report.LedState; led != nil {
				ledState = led.Value
			}

			if fru := report.Fru; fru != nil {
				bpn = pointer.SafeDeref(fru.BoardPartNumber)
				chassisSerial = pointer.SafeDeref(fru.ChassisPartSerial)
				productSerial = pointer.SafeDeref(fru.ProductSerial)
			}

			if bmc := report.Bmc; bmc != nil {
				ipAddress = bmc.Address
				mac = bmc.Mac
				bmcVersion = bmc.Version
			}

			power, powerText = extractPowerState(report)

			if report.UpdatedAt != nil && !report.UpdatedAt.AsTime().IsZero() {
				lastUpdated = fmt.Sprintf("%s ago", humanizeDuration(time.Since(report.UpdatedAt.AsTime())))
			}

			if bios := report.Bios; bios != nil {
				biosVersion = bios.Version
			}
		}

		if strings.EqualFold(ledState, "LED-ON") {
			blue := color.New(color.FgBlue).SprintFunc()
			emojis = append(emojis, blue(ledon))
		}

		if wide {
			rows = append(rows, []string{machineID, ledState, powerText, ipAddress, mac, bpn, chassisSerial, productSerial, biosVersion, bmcVersion, size, partition, rack, lastUpdated})
		} else {
			rows = append(rows, []string{machineID, strings.Join(emojis, nbr), power, ipAddress, mac, bpn, biosVersion, bmcVersion, size, partition, rack, lastUpdated})
		}
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}

func extractPowerState(bmc *apiv2.MachineBMCReport) (short, wide string) {
	if bmc == nil || bmc.Bmc == nil {
		return color.WhiteString(poweron), wide
	}

	state := bmc.Bmc.PowerState
	switch state {
	case "ON":
		short = color.GreenString(poweron)
	case "OFF":
		short = color.GreenString(powersleep)
	default:
		short = color.WhiteString(poweron)
	}

	wide = state
	for _, ps := range bmc.PowerSupplies {
		if ps.Health != "OK" {
			short = color.RedString(poweron)
			wide = wide + nbr + "Power Supply" + nbr + ps.Health
		}
		if ps.State != "Enabled" {
			short = color.RedString(powersleep)
			wide = wide + nbr + ps.State
		}
	}

	if bmc.PowerMetric != nil {
		short = fmt.Sprintf("%s"+nbr+"%.0fW", short, bmc.PowerMetric.AverageConsumedWatts)
		wide = fmt.Sprintf("%s %.0fW", wide, bmc.PowerMetric.AverageConsumedWatts)
	}

	return short, wide
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
			emojis = append(emojis, helpers.Skull)
		default:
			emojis = append(emojis, helpers.Question)
		}

		if status.Condition != nil {
			switch status.Condition.State {
			case apiv2.MachineState_MACHINE_STATE_LOCKED:
				emojis = append(emojis, helpers.Lock)
			case apiv2.MachineState_MACHINE_STATE_TAINTED:
				emojis = append(emojis, helpers.Bark)
			default:
				// noop
			}
		}
	}

	if events := m.RecentProvisioningEvents; events != nil {
		switch events.State {
		case apiv2.MachineProvisioningEventState_MACHINE_PROVISIONING_EVENT_STATE_FAILED_RECLAIM:
			emojis = append(emojis, helpers.Ambulance)
		case apiv2.MachineProvisioningEventState_MACHINE_PROVISIONING_EVENT_STATE_CRASHLOOP:
			emojis = append(emojis, helpers.Loop)
		default:
			// noop

		}

		if time.Since(events.LastErrorEvent.Time.AsTime()) < t.lastEventErrorThreshold {
			emojis = append(emojis, helpers.Exclamation)
		}
	}

	if m.Allocation != nil && m.Allocation.Vpn != nil && m.Allocation.Vpn.Connected {
		emojis = append(emojis, helpers.VPN)
	}

	return strings.Join(emojis, nbr)
}
