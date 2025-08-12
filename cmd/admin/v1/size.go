package v1

import (
	"fmt"

	"connectrpc.com/connect"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
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

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.SizeServiceCreateRequest, *adminv2.SizeServiceUpdateRequest, *apiv2.Size]{
		BinaryName:      config.BinaryName,
		GenericCLI:      gcli,
		Singular:        "size",
		Plural:          "sizes",
		Description:     "manage sizes which defines the cpu, gpu, memory and storage properties of machines",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ValidArgsFn:     c.Completion.SizeListCompletion,
		OnlyCmds:        genericcli.OnlyCmds(genericcli.CreateCmd, genericcli.UpdateCmd, genericcli.DeleteCmd, genericcli.EditCmd),
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *size) Get(id string) (*apiv2.Size, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.SizeServiceGetRequest{Id: id}

	resp, err := c.c.Client.Apiv2().Size().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return resp.Msg.Size, nil
}

func (c *size) Create(rq *adminv2.SizeServiceCreateRequest) (*apiv2.Size, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Size().Create(ctx, connect.NewRequest(rq))
	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}

	return resp.Msg.Size, nil
}

func (c *size) Delete(id string) (*apiv2.Size, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.SizeServiceDeleteRequest{Id: id}

	resp, err := c.c.Client.Adminv2().Size().Delete(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to delete size: %w", err)
	}

	return resp.Msg.Size, nil
}
func (c *size) List() ([]*apiv2.Size, error) {
	panic("unimplemented")

}
func (c *size) Convert(r *apiv2.Size) (string, *adminv2.SizeServiceCreateRequest, *adminv2.SizeServiceUpdateRequest, error) {

	return r.Id, &adminv2.SizeServiceCreateRequest{
			Size: &apiv2.Size{
				Id:          r.Id,
				Name:        r.Name,
				Description: r.Description,
				Meta:        r.Meta,
				Constraints: r.Constraints,
			},
		}, &adminv2.SizeServiceUpdateRequest{
			Id:          r.Id,
			Name:        r.Name,
			Description: r.Description,
			Constraints: r.Constraints,
			// FIXME
			Labels: &apiv2.UpdateLabels{},
		}, nil

}

func (c *size) Update(rq *adminv2.SizeServiceUpdateRequest) (*apiv2.Size, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.SizeServiceUpdateRequest{
		Id:          viper.GetString("id"),
		Name:        pointer.PointerOrNil(viper.GetString("name")),
		Description: pointer.PointerOrNil(viper.GetString("description")),
		Constraints: rq.Constraints,
		Labels:      &apiv2.UpdateLabels{}, // FIXME
	}

	resp, err := c.c.Client.Adminv2().Size().Update(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return resp.Msg.Size, nil
}
