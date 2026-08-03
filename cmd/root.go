package cmd

import (
	"log/slog"
	"os"

	client "github.com/metal-stack/api/go/client"
	"github.com/metal-stack/metal-lib/pkg/genericcli"

	adminv2 "github.com/metal-stack/cli/cmd/admin/v2"
	apiv2 "github.com/metal-stack/cli/cmd/api/v2"
	"github.com/metal-stack/cli/cmd/mcp"

	"github.com/metal-stack/cli/cmd/completion"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

func Execute() {
	cfg := &config.Config{
		Fs:        afero.NewOsFs(),
		Out:       os.Stdout,
		PromptOut: os.Stdout,
		In:        os.Stdin,
	}

	cmd := NewRootCmd(cfg)

	err := cmd.Execute()
	if err != nil {
		if viper.GetBool("debug") {
			panic(err)
		}

		os.Exit(1)
	}
}

func NewRootCmd(c *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          config.BinaryName,
		Aliases:      []string{"m"},
		Short:        "cli for managing entities in metal-stack",
		Long:         "",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			viper.SetFs(c.Fs)

			genericcli.Must(viper.BindPFlags(cmd.Flags()))
			genericcli.Must(viper.BindPFlags(cmd.PersistentFlags()))

			return initConfigWithViperCtx(c)
		},
	}
	rootCmd.PersistentFlags().StringP("config", "c", "", "alternative config file path, (default is ~/.metal-stack/config.yaml)")
	rootCmd.PersistentFlags().StringP("output-format", "o", "table", "output format (table|wide|markdown|json|yaml|template), wide is a table with more columns.")

	genericcli.Must(rootCmd.RegisterFlagCompletionFunc("output-format", cobra.FixedCompletions([]string{"table", "wide", "markdown", "json", "yaml", "template"}, cobra.ShellCompDirectiveNoFileComp)))

	rootCmd.PersistentFlags().StringP("template", "", "", `output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.`)
	rootCmd.PersistentFlags().Bool("force-color", false, "force colored output even without tty")
	rootCmd.PersistentFlags().Bool("debug", false, "debug output")
	rootCmd.PersistentFlags().Duration("timeout", 0, "request timeout used for api requests")

	rootCmd.PersistentFlags().String("api-url", "", "the url to the metal-stack.io api")
	rootCmd.PersistentFlags().String("api-token", "", "the token used for api requests")

	genericcli.Must(viper.BindPFlags(rootCmd.PersistentFlags()))

	markdownCmd := &cobra.Command{
		Use:   "markdown",
		Short: "create markdown documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := os.MkdirAll("./docs", 0755)
			if err != nil {
				return err
			}
			err = doc.GenMarkdownTree(rootCmd, "./docs")
			if err != nil {
				return err
			}

			err = genAdminDocs(rootCmd)
			if err != nil {
				return err
			}

			return nil
		},
		DisableAutoGenTag: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			recursiveAutoGenDisable(rootCmd)
		},
	}

	rootCmd.AddCommand(newContextCmd(c), markdownCmd, newLoginCmd(c), newLogoutCmd(c), mcp.NewMCPCmd(c))
	adminv2.AddCmds(rootCmd, c)
	apiv2.AddCmds(rootCmd, c)

	return rootCmd
}

func initConfigWithViperCtx(c *config.Config) error {
	c.Context = c.MustDefaultContext()

	listPrinter, err := newPrinterFromCLI(c.Out)
	if err != nil {
		return err
	}
	describePrinter, err := defaultToYAMLPrinter(c.Out)
	if err != nil {
		return err
	}

	c.ListPrinter = listPrinter
	c.DescribePrinter = describePrinter

	if c.Client != nil {
		return nil
	}

	mc, err := newApiClient(c.GetApiURL(), c.GetToken())
	if err != nil {
		return err
	}

	c.Client = mc
	c.Completion = completion.New(mc, c.GetProject())
	return nil
}

func newApiClient(apiURL, token string) (client.Client, error) {
	logLevel := slog.LevelInfo
	if viper.GetBool("debug") {
		logLevel = slog.LevelDebug
	}

	dialConfig := &client.DialConfig{
		BaseURL:   apiURL,
		Token:     token,
		UserAgent: "metal-stack-cli",
		Log:       slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})),
	}

	return client.New(dialConfig)
}

func recursiveAutoGenDisable(cmd *cobra.Command) {
	cmd.DisableAutoGenTag = true
	for _, child := range cmd.Commands() {
		recursiveAutoGenDisable(child)
	}
}

func genAdminDocs(rootCmd *cobra.Command) error {
	adminCmd, _, err := rootCmd.Find([]string{"admin"})
	if err != nil {
		return err
	}

	err = os.MkdirAll("./docs/admin", 0755)
	if err != nil {
		return err
	}

	hidden := adminCmd.Hidden
	adminCmd.Hidden = false
	defer func() { adminCmd.Hidden = hidden }()

	return doc.GenMarkdownTree(adminCmd, "./docs/admin")
}
