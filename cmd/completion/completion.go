package completion

import (
	"github.com/metal-stack/api/go/client"
	"github.com/spf13/cobra"
)

type Completion struct {
	Client client.Client
	Proj   string
}

func OutputFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"table", "wide", "markdown", "json", "yaml", "template"}, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) genericEnums(enums map[int32]string) ([]string, cobra.ShellCompDirective) {
	var names []string

	for id, name := range enums {
		if id == 0 {
			// skip UNSPECIFIED
			continue
		}
		names = append(names, name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
