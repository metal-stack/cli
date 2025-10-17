package cmd

import (
	"log/slog"
	"os"

	client "github.com/metal-stack/api/go/client"
	"github.com/metal-stack/metal-lib/pkg/genericcli"

	adminv2 "github.com/metal-stack/cli/cmd/admin/v1"
	apiv2 "github.com/metal-stack/cli/cmd/api/v1"

	"github.com/metal-stack/cli/cmd/completion"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

func Execute() {
	cfg := &config.Config{
		Fs:         afero.NewOsFs(),
		Out:        os.Stdout,
		PromptOut:  os.Stdout,
		In:         os.Stdin,
		Completion: &completion.Completion{},
	}

	cmd := newRootCmd(cfg)

	err := cmd.Execute()
	if err != nil {
		if viper.GetBool("debug") {
			panic(err)
		}

		os.Exit(1)
	}
}

func newRootCmd(c *config.Config) *cobra.Command {
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
	rootCmd.PersistentFlags().StringP("output-format", "o", "table", "output format (table|wide|markdown|json|yaml|template|jsonraw|yamlraw), wide is a table with more columns, jsonraw and yamlraw do not translate proto enums into string types but leave the original int32 values intact.")

	genericcli.Must(rootCmd.RegisterFlagCompletionFunc("output-format", cobra.FixedCompletions([]string{"table", "wide", "markdown", "json", "yaml", "template"}, cobra.ShellCompDirectiveNoFileComp)))

	rootCmd.PersistentFlags().StringP("template", "", "", `output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.`)
	rootCmd.PersistentFlags().Bool("force-color", false, "force colored output even without tty")
	rootCmd.PersistentFlags().Bool("debug", false, "debug output")
	rootCmd.PersistentFlags().Duration("timeout", 0, "request timeout used for api requests")

	rootCmd.PersistentFlags().String("api-url", "https://api.metal-stack.io", "the url to the metal-stack.io api")
	rootCmd.PersistentFlags().String("api-token", "", "the token used for api requests")

	genericcli.Must(viper.BindPFlags(rootCmd.PersistentFlags()))

	markdownCmd := &cobra.Command{
		Use:   "markdown",
		Short: "create markdown documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(rootCmd, "./docs")
		},
		DisableAutoGenTag: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			recursiveAutoGenDisable(rootCmd)
		},
	}

	rootCmd.AddCommand(newContextCmd(c), markdownCmd, newLoginCmd(c), newLogoutCmd(c))
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
	c.Completion.Client = mc
	c.Completion.Ctx = context.Background()
	c.Completion.Project = c.GetProject()

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
