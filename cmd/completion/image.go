package completion

import (
	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
)

func (c *Completion) Image(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Apiv2().Image().List(cmd.Context(), &apiv2.ImageServiceListRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, img := range resp.Images {
		names = append(names, img.Id+"\t"+*img.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ImageFeatures(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var names []string

	for i := range apiv2.ImageFeature_name {
		if i == 0 {
			continue
		}

		n, err := enum.GetStringValue(apiv2.ImageFeature(i))
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		names = append(names, *n)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ImageClassification(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var names []string

	for i := range apiv2.ImageClassification_name {
		if i == 0 {
			continue
		}

		n, err := enum.GetStringValue(apiv2.ImageClassification(i))
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		names = append(names, *n)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
