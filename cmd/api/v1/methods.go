package v1

import (
	"fmt"
	"sort"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMethodsCmd(c *config.Config) *cobra.Command {
	methodCmd := &cobra.Command{
		Use:   "api-methods",
		Short: "show available api-methods of the metal-stack.io api",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewRequestContext()
			defer cancel()

			if viper.GetBool("scoped") {

				req := &apiv2.MethodServiceTokenScopedListRequest{}

				resp, err := c.Client.Apiv2().Method().TokenScopedList(ctx, req)
				if err != nil {
					return fmt.Errorf("failed to list methods: %w", err)
				}

				return printers.NewProtoYAMLPrinter().WithOut(c.Out).Print(resp)
			}

			var (
				methods []string
				req     = &apiv2.MethodServiceListRequest{}
			)

			resp, err := c.Client.Apiv2().Method().List(ctx, req)
			if err != nil {
				return fmt.Errorf("failed to list methods: %w", err)
			}

			methods = resp.GetMethods()

			sort.Strings(methods)

			for _, method := range methods {
				_, _ = fmt.Fprintln(c.Out, method)
			}

			return nil
		},
	}

	methodCmd.Flags().Bool("scoped", false, "show accessible api-methods depending on the api access token")

	return methodCmd
}
