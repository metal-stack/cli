package v1

import (
	"fmt"
	"strings"
	"time"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type image struct {
	c *config.Config
}

func newImageCmd(c *config.Config) *cobra.Command {
	w := &image{
		c: c,
	}
	gcli := genericcli.NewGenericCLI(w).WithFS(c.Fs)

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.ImageServiceCreateRequest, *adminv2.ImageServiceUpdateRequest, *apiv2.Image]{
		BinaryName:      config.BinaryName,
		GenericCLI:      gcli,
		Singular:        "image",
		Plural:          "images",
		Description:     "manage images which are used to be installed on machines and firewalls",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.CreateCmd, genericcli.UpdateCmd, genericcli.DeleteCmd, genericcli.EditCmd),
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.RunE = func(cmd *cobra.Command, args []string) error {
				return gcli.DescribeAndPrint("", w.c.DescribePrinter)
			}
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *image) Get(id string) (*apiv2.Image, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.ImageServiceGetRequest{Id: id}

	resp, err := c.c.Client.Apiv2().Image().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return resp.Image, nil
}

func (c *image) Create(rq *adminv2.ImageServiceCreateRequest) (*apiv2.Image, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	var expiresAt *timestamppb.Timestamp
	if viper.IsSet("expires-in") {
		expiresAt = timestamppb.New(time.Now().Add(viper.GetDuration("expires-in")))
	}

	req := &adminv2.ImageServiceCreateRequest{
		Image: &apiv2.Image{
			Id:          viper.GetString("id"),
			Url:         viper.GetString("url"),
			Description: pointer.PointerOrNil(viper.GetString("description")),
			ExpiresAt:   expiresAt,
			Features:    imageFeaturesFromString(viper.GetStringSlice("features")),
			Meta:        &apiv2.Meta{
				// TODO labels
			},
		},
	}

	resp, err := c.c.Client.Adminv2().Image().Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return resp.Image, nil
}

func (c *image) Delete(id string) (*apiv2.Image, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.ImageServiceDeleteRequest{Id: id}

	resp, err := c.c.Client.Adminv2().Image().Delete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete image: %w", err)
	}

	return resp.Image, nil
}
func (c *image) List() ([]*apiv2.Image, error) {
	panic("unimplemented")

}
func (c *image) Convert(r *apiv2.Image) (string, *adminv2.ImageServiceCreateRequest, *adminv2.ImageServiceUpdateRequest, error) {

	return r.Id, &adminv2.ImageServiceCreateRequest{
			Image: &apiv2.Image{
				Id:             r.Id,
				Url:            r.Url,
				Name:           r.Name,
				Description:    r.Description,
				Features:       r.Features,
				Meta:           r.Meta,
				Classification: r.Classification,
				ExpiresAt:      r.ExpiresAt,
			},
		}, &adminv2.ImageServiceUpdateRequest{
			Id:          r.Id,
			Url:         &r.Url,
			Name:        r.Name,
			Description: r.Description,
			Features:    r.Features,
			UpdateMeta: &apiv2.UpdateMeta{
				LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_CLIENT,
				UpdatedAt:       r.Meta.UpdatedAt,
			},
			Classification: r.Classification,
			ExpiresAt:      r.ExpiresAt,
		}, nil

}

func (c *image) Update(rq *adminv2.ImageServiceUpdateRequest) (*apiv2.Image, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	var expiresAt *timestamppb.Timestamp
	if viper.IsSet("expires-in") {
		expiresAt = timestamppb.New(time.Now().Add(viper.GetDuration("expires-in")))
	}

	req := &adminv2.ImageServiceUpdateRequest{
		Id:          viper.GetString("id"),
		Url:         pointer.Pointer(viper.GetString("url")),
		Description: pointer.PointerOrNil(viper.GetString("description")),
		ExpiresAt:   expiresAt,
		Features:    imageFeaturesFromString(viper.GetStringSlice("features")),
		UpdateMeta: &apiv2.UpdateMeta{
			LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_CLIENT,
			UpdatedAt:       rq.UpdateMeta.GetUpdatedAt(),
		},
	}

	resp, err := c.c.Client.Adminv2().Image().Update(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return resp.Image, nil
}

func imageFeaturesFromString(features []string) []apiv2.ImageFeature {
	if len(features) == 0 {
		return nil
	}

	var result []apiv2.ImageFeature
	for _, f := range features {
		switch strings.ToLower(f) {
		case "machine":
			result = append(result, apiv2.ImageFeature_IMAGE_FEATURE_MACHINE)
		case "firewall":
			result = append(result, apiv2.ImageFeature_IMAGE_FEATURE_FIREWALL)
		}
	}
	return result
}
