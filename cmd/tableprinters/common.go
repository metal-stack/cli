package tableprinters

import (
	"fmt"
	"math"
	"strings"
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

type TablePrinter struct {
	t                       *printers.TablePrinter
	lastEventErrorThreshold time.Duration
}

func New() *TablePrinter {
	return &TablePrinter{}
}

func (t *TablePrinter) SetPrinter(printer *printers.TablePrinter) {
	t.t = printer
}

func (t *TablePrinter) SetLastEventErrorThreshold(threshold time.Duration) {
	t.lastEventErrorThreshold = threshold
}

func (t *TablePrinter) ToHeaderAndRows(data any, wide bool) ([]string, [][]string, error) {
	switch d := data.(type) {

	case *config.Contexts:
		return t.ContextTable(d, wide)

	case *apiv2.IP:
		return t.IPTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.IP:
		return t.IPTable(d, wide)

	case *apiv2.Machine:
		return t.MachineTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.Machine:
		return t.MachineTable(d, wide)

	case *apiv2.Image:
		return t.ImageTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.Image:
		return t.ImageTable(d, wide)

	case *apiv2.Project:
		return t.ProjectTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.Project:
		return t.ProjectTable(d, wide)
	case *apiv2.ProjectInvite:
		return t.ProjectInviteTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.ProjectInvite:
		return t.ProjectInviteTable(d, wide)
	case *apiv2.ProjectMember:
		return t.ProjectMemberTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.ProjectMember:
		return t.ProjectMemberTable(d, wide)

	case *apiv2.Token:
		return t.TokenTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.Token:
		return t.TokenTable(d, wide)

	case *apiv2.Tenant:
		return t.TenantTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.Tenant:
		return t.TenantTable(d, wide)
	case *apiv2.TenantInvite:
		return t.TenantInviteTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.TenantInvite:
		return t.TenantInviteTable(d, wide)
	case *apiv2.TenantMember:
		return t.TenantMemberTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.TenantMember:
		return t.TenantMemberTable(d, wide)

	case *apiv2.Health:
		return t.HealthTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.Health:
		return t.HealthTable(d, wide)

	default:
		return nil, nil, fmt.Errorf("unknown table printer for type: %T", d)
	}
}

func humanizeDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"d", days},
		{"h", hours},
		{"m", minutes},
		{"s", seconds},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		default:
			parts = append(parts, fmt.Sprintf("%d%s", chunk.amount, chunk.singularName))
		}
	}

	if len(parts) == 0 {
		return "0s"
	}
	if len(parts) > 2 {
		parts = parts[:2]
	}
	return strings.Join(parts, " ")
}
