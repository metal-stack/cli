package config

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/metal-stack/api/go/client"
	"github.com/metal-stack/cli/cmd/completion"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

const (
	// BinaryName is the name of the cli in all help texts
	BinaryName = "metalctlv2"
	// ConfigDir is the directory in either the homedir or in /etc where the cli searches for a file config.yaml
	// also used as prefix for environment based configuration, e.g. METAL_STACK_CLOUD_ will be the variable prefix.
	ConfigDir  = "metal-stack"
	keyTimeout = "timeout"
)

type Config struct {
	Fs              afero.Fs
	In              io.Reader
	Out             io.Writer
	PromptOut       io.Writer
	Client          client.Client
	ListPrinter     printers.Printer
	DescribePrinter printers.Printer
	Completion      *completion.Completion
	Context         genericcli.Context
	ContextConfig   genericcli.ContextConfig
}

func (c *Config) NewRequestContext() (context.Context, context.CancelFunc) {
	timeout := c.Context.Timeout
	if timeout == nil {
		timeout = pointer.Pointer(30 * time.Second)
	}
	if viper.IsSet(genericcli.KeyTimeout) {
		timeout = pointer.Pointer(viper.GetDuration(genericcli.KeyTimeout))
	}

	return context.WithTimeout(context.Background(), *timeout)
}

func (c *Config) GetProject() string {
	if viper.IsSet("project") {
		return viper.GetString("project")
	}
	return c.Context.DefaultProject
}

func (c *Config) GetTenant() (string, error) {
	if viper.IsSet("tenant") {
		return viper.GetString("tenant"), nil
	}

	if c.GetProject() == "" {
		return "", fmt.Errorf("tenant is not set")
	}

	ctx, cancel := c.NewRequestContext()
	defer cancel()

	projectResp, err := c.Client.Apiv2().Project().Get(ctx, &apiv2.ProjectServiceGetRequest{
		Project: c.GetProject(),
	})
	if err != nil {
		return "", fmt.Errorf("unable to derive tenant from project: %w", err)
	}

	return projectResp.Project.Tenant, nil
}

func (c *Config) GetToken() string {
	if viper.IsSet(genericcli.KeyAPIToken) {
		return viper.GetString(genericcli.KeyAPIToken)
	}
	return c.Context.APIToken
}

func (c *Config) GetApiURL() string {
	if viper.IsSet(genericcli.KeyAPIURL) {
		return viper.GetString(genericcli.KeyAPIURL)
	}
	if c.Context.APIURL != nil {
		return *c.Context.APIURL
	}

	// fallback to the default specified by viper
	return viper.GetString(genericcli.KeyAPIURL)
}

func (c *Config) GetProvider() string {
	if viper.IsSet(genericcli.KeyProvider) {
		return viper.GetString(genericcli.KeyProvider)
	}
	return c.Context.Provider
}
