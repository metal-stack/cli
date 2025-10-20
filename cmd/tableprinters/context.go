package tableprinters

import (
	"github.com/fatih/color"
	clitypes "github.com/metal-stack/metal-lib/pkg/commands/types"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/viper"
)

func (t *TablePrinter) ContextTable(data *clitypes.Contexts, wide bool) ([]string, [][]string, error) {
	var (
		header = []string{"", "Name", "Provider", "Default Project"}
		rows   [][]string
	)

	if wide {
		header = append(header, "API URL")
	}

	for _, c := range data.Contexts {
		active := ""
		if c.Name == data.CurrentContext {
			active = color.GreenString("✔")
		}

		row := []string{active, c.Name, c.Provider, c.DefaultProject}
		if wide {
			url := pointer.SafeDeref(c.ApiURL)
			if url == "" {
				url = viper.GetString("api-url")
			}

			row = append(row, url)
		}

		rows = append(rows, row)
	}

	return header, rows, nil
}
