package v2

import (
	"fmt"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type size struct {
	c *config.Config
}

func newSizeCmd(c *config.Config) *cobra.Command {
	w := &size{
		c: c,
	}

	gcli := genericcli.NewGenericCLI(w).WithFS(c.Fs)

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.Size]{
		BinaryName:      config.BinaryName,
		GenericCLI:      gcli,
		Singular:        "size",
		Plural:          "sizes",
		Description:     "manage sizes which defines the cpu, gpu, memory and storage properties of machines",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.SizeListCompletion,
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("id", "", "", "size id to filter for")
			cmd.Flags().StringP("name", "", "", "size name to filter for")
			cmd.Flags().StringP("description", "", "", "size description to filter for")
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *size) Get(id string) (*apiv2.Size, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.SizeServiceGetRequest{Id: id}

	resp, err := c.c.Client.Apiv2().Size().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}

	return resp.Size, nil
}

func (c *size) List() ([]*apiv2.Size, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.SizeServiceListRequest{Query: &apiv2.SizeQuery{
		Id:          pointer.PointerOrNil(viper.GetString("id")),
		Name:        pointer.PointerOrNil(viper.GetString("name")),
		Description: pointer.PointerOrNil(viper.GetString("description")),
	}}

	resp, err := c.c.Client.Apiv2().Size().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get sizes: %w", err)
	}

	return resp.Sizes, nil
}

func (c *size) Create(rq any) (*apiv2.Size, error) {
	panic("unimplemented")
}

func (c *size) Delete(id string) (*apiv2.Size, error) {
	panic("unimplemented")
}

func (t *size) Convert(r *apiv2.Size) (string, any, any, error) {
	panic("unimplemented")
}

func (t *size) Update(rq any) (*apiv2.Size, error) {
	panic("unimplemented")
}
