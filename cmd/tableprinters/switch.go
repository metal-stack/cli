package tableprinters

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func (t *TablePrinter) SwitchTable(data []*apiv2.Switch, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Partition", "Rack", "OS", "Status", "Last Sync"}
	if wide {
		header = []string{"ID", "Partition", "Rack", "OS", "Metalcore", "IP", "Mode", "Last Sync", "Sync Duration", "Last Error"}

		t.t.DisableAutoWrap(true)
	}

	for _, s := range data {
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

		// FIXME: nil pointer checks and refactor
		if s.LastSync != nil {
			syncTime = s.LastSync.Time.AsTime()
			syncAge := time.Since(syncTime)
			syncDur := s.LastSync.Duration.AsDuration().Round(time.Millisecond)

			if syncAge >= time.Minute*10 || syncDur >= 30*time.Second {
				shortStatus = color.RedString(dot)
			} else if syncAge >= time.Minute*1 || syncDur >= 20*time.Second {
				shortStatus = color.YellowString(dot)
			} else {
				shortStatus = color.GreenString(dot)
				if !allUp {
					shortStatus = color.YellowString(dot)
				}
			}

			syncLast = humanizeDuration(syncAge) + " ago"
			syncDurStr = fmt.Sprintf("%v", syncDur)
		}

		// FIXME: nil pointer checks and refactor
		if s.LastSyncError != nil {
			errorTime := s.LastSyncError.Time.AsTime()
			// after 7 days we do not show sync errors anymore
			if !errorTime.IsZero() && time.Since(errorTime) < 7*24*time.Hour {
				lastError = fmt.Sprintf("%s ago: %s", humanizeDuration(time.Since(errorTime)), s.LastSyncError.Error)

				if errorTime.After(s.LastSync.Time.AsTime()) {
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
