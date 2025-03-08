package tableprinters

import (
	"github.com/fatih/color"
	"github.com/metal-stack/cli/cmd/config"
)

func (t *TablePrinter) ContextTable(data *config.Contexts, wide bool) ([]string, [][]string, error) {
	var (
		header = []string{"", "Name", "Default Project"}
		rows   [][]string
	)

	for _, c := range data.Contexts {
		active := ""
		if c.Name == data.CurrentContext {
			active = color.GreenString("✔")
		}
		rows = append(rows, []string{active, c.Name, c.DefaultProject})
	}

	return header, rows, nil
}
