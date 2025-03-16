package tableprinters

import (
	"fmt"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

type TablePrinter struct {
	t *printers.TablePrinter
}

func New() *TablePrinter {
	return &TablePrinter{}
}

func (t *TablePrinter) SetPrinter(printer *printers.TablePrinter) {
	t.t = printer
}

func (t *TablePrinter) ToHeaderAndRows(data any, wide bool) ([]string, [][]string, error) {
	switch d := data.(type) {

	case *config.Contexts:
		return t.ContextTable(d, wide)

	case *apiv2.IP:
		return t.IPTable(pointer.WrapInSlice(d), wide)
	case []*apiv2.IP:
		return t.IPTable(d, wide)

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
