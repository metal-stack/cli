package v2

import (
	"fmt"

	"github.com/metal-stack/api/go/enum"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type component struct {
	c *config.Config
}

func newComponentCmd(c *config.Config) *cobra.Command {
	w := &component{
		c: c,
	}
	gcli := genericcli.NewGenericCLI(w).WithFS(c.Fs)

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.Component]{
		BinaryName:      config.BinaryName,
		GenericCLI:      gcli,
		Singular:        "component",
		Plural:          "components",
		Description:     "list status of components, e.g. microservices connected to the metal-apiserver",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd, genericcli.DeleteCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("uuid", "", "lists only component with this uuid")
			cmd.Flags().String("type", "", "lists only component of this type")
			cmd.Flags().String("identifier", "", "lists only component with this identifier")
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *component) Get(id string) (*apiv2.Component, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.ComponentServiceGetRequest{Uuid: id}

	resp, err := c.c.Client.Adminv2().Component().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get component: %w", err)
	}

	return resp.Component, nil
}

func (c *component) Create(rq any) (*apiv2.Component, error) {
	panic("unimplemented")
}

func (c *component) Delete(id string) (*apiv2.Component, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.ComponentServiceDeleteRequest{Uuid: id}

	resp, err := c.c.Client.Adminv2().Component().Delete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete component: %w", err)
	}

	return resp.Component, nil
}
func (c *component) List() ([]*apiv2.Component, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	query := &apiv2.ComponentQuery{
		Uuid:       pointer.PointerOrNil(viper.GetString("uuid")),
		Identifier: pointer.PointerOrNil(viper.GetString("identifier")),
	}

	if viper.IsSet("type") {
		t, err := enum.GetEnum[apiv2.ComponentType](viper.GetString("type"))
		if err != nil {
			return nil, fmt.Errorf("unable to get component type of string %q %w", viper.GetString("type"), err)
		}
		query.Type = &t
	}

	req := &adminv2.ComponentServiceListRequest{Query: query}

	resp, err := c.c.Client.Adminv2().Component().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get components: %w", err)
	}

	return resp.Components, nil
}
func (c *component) Convert(r *apiv2.Component) (string, any, any, error) {
	panic("unimplemented")
}

func (c *component) Update(rq any) (*apiv2.Component, error) {
	panic("unimplemented")
}
