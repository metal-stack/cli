package completion

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
)

func (c *Completion) FilesystemLayout(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Apiv2().Filesystem().List(cmd.Context(), &apiv2.FilesystemServiceListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, fsl := range resp.FilesystemLayouts {
		names = append(names, fsl.Id+"\t"+pointer.SafeDeref(fsl.Name))
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
