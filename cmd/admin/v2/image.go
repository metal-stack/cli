package v2

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
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("id", "", "image id")
			cmd.Flags().String("url", "", "image url")
			cmd.Flags().String("name", "", "image name")
			cmd.Flags().String("classification", "", "image classification")
			cmd.Flags().String("expires-in", "", "expires-in duration")
			cmd.Flags().String("description", "", "image description")
			cmd.Flags().StringSlice("features", nil, "image features can be machine and/or firewall")
		},
		CreateRequestFromCLI: w.createFromCLI,
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("id", "", "image id")
			cmd.Flags().String("url", "", "image url")
			cmd.Flags().String("name", "", "image name")
			cmd.Flags().String("classification", "", "image classification")
			cmd.Flags().String("expires-in", "", "expires-in duration")
			cmd.Flags().String("description", "", "image description")
			cmd.Flags().StringSlice("features", nil, "image features can be machine and/or firewall")
		},
		UpdateRequestFromCLI: w.updateFromCLI,
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

	resp, err := c.c.Client.Adminv2().Image().Create(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}

	return resp.Image, nil
}

func (c *image) createFromCLI() (*adminv2.ImageServiceCreateRequest, error) {
	var expiresAt *timestamppb.Timestamp
	if viper.IsSet("expires-in") {
		expiresAt = timestamppb.New(time.Now().Add(viper.GetDuration("expires-in")))
	}

	return &adminv2.ImageServiceCreateRequest{
		Image: &apiv2.Image{
			Id:             viper.GetString("id"),
			Url:            viper.GetString("url"),
			Classification: imageClassificationFromString(viper.GetString("classification")),
			Name:           pointer.PointerOrNil(viper.GetString("name")),
			Description:    pointer.PointerOrNil(viper.GetString("description")),
			ExpiresAt:      expiresAt,
			Features:       imageFeaturesFromString(viper.GetStringSlice("features")),
			Meta:           &apiv2.Meta{
				// TODO labels
			},
		},
	}, nil
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

	resp, err := c.c.Client.Adminv2().Image().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update image: %w", err)
	}

	return resp.Image, nil
}

func (c *image) updateFromCLI(args []string) (*adminv2.ImageServiceUpdateRequest, error) {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	req := &adminv2.ImageServiceUpdateRequest{
		Id:          id,
		Url:         pointer.PointerOrNil(viper.GetString("url")),
		Name:        pointer.PointerOrNil(viper.GetString("name")),
		Description: pointer.PointerOrNil(viper.GetString("description")),
		UpdateMeta: &apiv2.UpdateMeta{
			LockingStrategy: apiv2.OptimisticLockingStrategy_OPTIMISTIC_LOCKING_STRATEGY_SERVER,
		},
	}

	if viper.IsSet("expires-in") {
		req.ExpiresAt = timestamppb.New(time.Now().Add(viper.GetDuration("expires-in")))
	}
	if viper.IsSet("features") {
		req.Features = imageFeaturesFromString(viper.GetStringSlice("features"))
	}
	if viper.IsSet("classification") {
		req.Classification = imageClassificationFromString(viper.GetString("classification"))
	}

	return req, nil
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

func imageClassificationFromString(classification string) apiv2.ImageClassification {
	switch strings.ToLower(strings.TrimSpace(classification)) {
	case "preview":
		return apiv2.ImageClassification_IMAGE_CLASSIFICATION_PREVIEW
	case "supported":
		return apiv2.ImageClassification_IMAGE_CLASSIFICATION_SUPPORTED
	case "deprecated":
		return apiv2.ImageClassification_IMAGE_CLASSIFICATION_DEPRECATED
	}

	return apiv2.ImageClassification_IMAGE_CLASSIFICATION_UNSPECIFIED
}
