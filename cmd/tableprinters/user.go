package tableprinters

import (
	"strings"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

type Named interface {
	GetName() string
}

func (t *TablePrinter) UserTable(data []*apiv2.User, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"Login", "Name", "Email", "Default-Tenant"}
	)

	if wide {
		header = []string{"Login", "Name", "Email", "Default-Tenant", "Tenants", "Projects"}
	}

	for _, user := range data {
		login := user.Login
		name := user.Name
		email := user.Email
		defaultTenant := user.DefaultTenant.Name

		if wide {
			rows = append(rows, []string{login, name, email, defaultTenant, namesString(user.Tenants), namesString(user.Projects)})

		} else {
			rows = append(rows, []string{login, name, email, defaultTenant})
		}
	}

	return header, rows, nil
}

func namesString[T Named](arr []T) string {
	names := make([]string, 0, len(arr))

	for _, t := range arr {
		name := t.GetName()
		if name != "" {
			names = append(names, name)
		}
	}

	return strings.Join(names, ", ")
}
