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
		header = []string{"ID", "Last Event", "When", "Age", "Description", "Name", "Hostname", "Project", "IPs", "Size", "Image", "Partition", "Rack", "Started", "Tags", "Lock/Reserve"}
	}
	for _, machine := range data {
		machineID := machine.Uuid

		if machine.Status != nil && machine.Status.LedState != nil && machine.Status.LedState.Value == "LED-ON" {
			blue := color.New(color.FgBlue).SprintFunc()
			machineID = blue(machineID)
		}

		alloc := pointer.SafeDeref(machine.Allocation)
		sizeID := pointer.SafeDeref(machine.Size).Id
		partitionID := pointer.SafeDeref(machine.Partition).Id
		project := alloc.Project
		name := alloc.Name
		desc := alloc.Description
		hostname := alloc.Hostname
		image := pointer.SafeDeref(pointer.SafeDeref(alloc.Image).Name)

		rack := machine.Rack

		truncatedHostname := genericcli.TruncateEnd(hostname, 30)

		var nwIPs []string
		for _, nw := range alloc.Networks {
			nwIPs = append(nwIPs, nw.Ips...)
		}
		ips := strings.Join(nwIPs, "\n")

		started := ""
		age := ""

		if alloc.Meta != nil && alloc.Meta.CreatedAt != nil && !alloc.Meta.CreatedAt.AsTime().IsZero() {
			started = alloc.Meta.CreatedAt.AsTime().Format(time.RFC3339)
			age = humanizeDuration(time.Since(alloc.Meta.CreatedAt.AsTime()))
		}
		tags := ""
		if machine.Meta.Labels != nil && len(machine.Meta.Labels.Labels) > 0 {
			var labels []string
			for k, v := range machine.Meta.Labels.Labels {
				labels = append(labels, k+"="+v)
			}
			tags = strings.Join(labels, ",")
		}

		reserved := ""
		if machine.Status.Condition != nil {
			stateString, err := enum.GetStringValue(machine.Status.Condition.State)
			if err != nil {
				return nil, nil, err
			}
			reserved = fmt.Sprintf("%s:%s", *stateString, machine.Status.Condition.Description)
		}

		lastEvent := ""
		when := ""
		if len(machine.RecentProvisioningEvents.Events) > 0 {
			since := time.Since(machine.RecentProvisioningEvents.LastEventTime.AsTime())
			when = humanizeDuration(since)
			lastEventString, err := enum.GetStringValue(machine.RecentProvisioningEvents.Events[0].Event)
			if err != nil {
				return nil, nil, err
			}
			lastEvent = *lastEventString
		}

		emojis, _ := t.getMachineStatusEmojis(machine.Status.Liveliness, machine.RecentProvisioningEvents, machine.Status.Condition.State, alloc.Vpn)

		if wide {
			rows = append(rows, []string{machineID, lastEvent, when, age, desc, name, hostname, project, ips, sizeID, image, partitionID, rack, started, tags, reserved})
		} else {
			rows = append(rows, []string{machineID, emojis, lastEvent, when, age, truncatedHostname, project, sizeID, image, partitionID, rack})
		}
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}

func (t *TablePrinter) getMachineStatusEmojis(liveliness apiv2.MachineLiveliness, events *apiv2.MachineRecentProvisioningEvents, state apiv2.MachineState, vpn *apiv2.MachineVPN) (string, string) {
	var (
		emojis           []string
		wide             []string
		livelinessString *string
		err              error
	)
	livelinessString, err = enum.GetStringValue(liveliness)
	if err != nil {
		livelinessString = new("unknown")
	}

	switch l := liveliness; l {
	case apiv2.MachineLiveliness_MACHINE_LIVELINESS_ALIVE:
		// noop
	case apiv2.MachineLiveliness_MACHINE_LIVELINESS_DEAD:
		emojis = append(emojis, helpers.Skull)
		wide = append(wide, *livelinessString)
	case apiv2.MachineLiveliness_MACHINE_LIVELINESS_UNKNOWN:
		emojis = append(emojis, helpers.Question)
		wide = append(wide, *livelinessString)
	default:
		emojis = append(emojis, helpers.Question)
		wide = append(wide, *livelinessString)
	}

	switch state {
	case apiv2.MachineState_MACHINE_STATE_AVAILABLE:
		// noop
	case apiv2.MachineState_MACHINE_STATE_LOCKED:
		emojis = append(emojis, helpers.Lock)
		wide = append(wide, "Locked")
	case apiv2.MachineState_MACHINE_STATE_RESERVED:
		emojis = append(emojis, helpers.Bark)
		wide = append(wide, "Reserved")
	}

	if events != nil {
		switch events.State {
		case apiv2.MachineProvisioningEventState_MACHINE_PROVISIONING_EVENT_STATE_FAILED_RECLAIM:
			emojis = append(emojis, helpers.Ambulance)
			wide = append(wide, "FailedReclaim")
		case apiv2.MachineProvisioningEventState_MACHINE_PROVISIONING_EVENT_STATE_CRASHLOOP:
			emojis = append(emojis, helpers.Loop)
			wide = append(wide, "CrashLoop")

		}

		if events.LastErrorEvent != nil && time.Since(events.LastErrorEvent.Time.AsTime()) < t.lastEventErrorThreshold {
			emojis = append(emojis, helpers.Exclamation)
			wide = append(wide, "LastEventErrors")
		}

	}

	if vpn != nil && vpn.Connected {
		emojis = append(emojis, helpers.VPN)
		wide = append(wide, "VPN")
	}

	return strings.Join(emojis, nbr), strings.Join(wide, ", ")
}
