package v1

import (
	"fmt"

	"connectrpc.com/connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type user struct {
	c *config.Config
}

func newUserCmd(c *config.Config) *cobra.Command {
	w := &user{
		c: c,
	}

	gcli := genericcli.NewGenericCLI(w).WithFS(c.Fs)

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.User]{
		BinaryName:      config.BinaryName,
		GenericCLI:      gcli,
		Singular:        "user",
		Plural:          "users",
		Description:     "manage api users for accessing the metal-stack.io api",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd),
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.RunE = func(cmd *cobra.Command, args []string) error {
				return gcli.DescribeAndPrint("", w.c.DescribePrinter)
			}
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *user) Get(id string) (*apiv2.User, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.UserServiceGetRequest{}

	resp, err := c.c.Client.Apiv2().User().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp.Msg.GetUser(), nil
}

func (c *user) List() ([]*apiv2.User, error) {
	panic("unimplemented")
}

func (c *user) Create(rq any) (*apiv2.User, error) {
	panic("unimplemented")
}

func (c *user) Delete(id string) (*apiv2.User, error) {
	panic("unimplemented")
}

func (t *user) Convert(r *apiv2.User) (string, any, any, error) {
	panic("unimplemented")
}

func (t *user) Update(rq any) (*apiv2.User, error) {
	panic("unimplemented")
}
