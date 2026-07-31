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
