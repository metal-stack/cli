package v2

import (
	"fmt"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type adminSizeImageConstraint struct {
	c *config.Config
}

func newAdminSizeImageConstraintCmd(c *config.Config) *cobra.Command {
	w := &adminSizeImageConstraint{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, *adminv2.SizeImageConstraintServiceUpdateRequest, *apiv2.SizeImageConstraint]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "size-image-constraint",
		Plural:          "size-image-constraints",
		Description:     "manage size image constraints",
		Sorter:          sorters.SizeImageConstraintSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd, genericcli.DeleteCmd),
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *adminSizeImageConstraint) Get(id string) (*apiv2.SizeImageConstraint, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.SizeImageConstraintServiceGetRequest{Size: id}

	resp, err := c.c.Client.Adminv2().SizeImageConstraint().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get size image constraint: %w", err)
	}

	return resp.SizeImageConstraint, nil
}

func (c *adminSizeImageConstraint) List() ([]*apiv2.SizeImageConstraint, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.SizeImageConstraintServiceListRequest{}

	resp, err := c.c.Client.Adminv2().SizeImageConstraint().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list size image constraints: %w", err)
	}

	return resp.SizeImageConstraints, nil
}

func (c *adminSizeImageConstraint) Create(rq any) (*apiv2.SizeImageConstraint, error) {
	panic("unimplemented")
}

func (c *adminSizeImageConstraint) Delete(id string) (*apiv2.SizeImageConstraint, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().SizeImageConstraint().Delete(ctx, &adminv2.SizeImageConstraintServiceDeleteRequest{Size: id})
	if err != nil {
		return nil, fmt.Errorf("failed to delete size image constraint: %w", err)
	}

	return resp.SizeImageConstraint, nil
}

func (c *adminSizeImageConstraint) Convert(r *apiv2.SizeImageConstraint) (string, any, *adminv2.SizeImageConstraintServiceUpdateRequest, error) {
	return r.Size, nil, &adminv2.SizeImageConstraintServiceUpdateRequest{
		Size:             r.Size,
		ImageConstraints: r.ImageConstraints,
		Name:             r.Name,
		Description:      r.Description,
	}, nil
}

func (c *adminSizeImageConstraint) Update(rq *adminv2.SizeImageConstraintServiceUpdateRequest) (*apiv2.SizeImageConstraint, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().SizeImageConstraint().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update size image constraint: %w", err)
	}

	return resp.SizeImageConstraint, nil
}
