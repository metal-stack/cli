package v2

import (
	"fmt"
	"time"

	"github.com/metal-stack/api/go/errorutil"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/pkg/helpers"
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
	var (
		w = &image{
			c: c,
		}

		queryFlags = func(cmd *cobra.Command) {
			cmd.Flags().StringP("id", "", "", "image id to filter for")
			cmd.Flags().StringP("os", "", "", "image os to filter for")
			cmd.Flags().StringP("version", "", "", "image version to filter for")
			cmd.Flags().StringP("name", "", "", "image name to filter for")
			cmd.Flags().StringP("description", "", "", "image description to filter for")
			cmd.Flags().StringP("feature", "", "", "image feature to filter for, can be either machine|firewall")
		}
	)

	cmdsConfig := &genericcli.CmdsConfig[*adminv2.ImageServiceCreateRequest, *adminv2.ImageServiceUpdateRequest, *apiv2.Image]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "image",
		Plural:          "images",
		Description:     "manage images which are used to be installed on machines and firewalls",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: queryFlags,
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("id", "", "image id")
			cmd.Flags().String("url", "", "image url")
			cmd.Flags().String("name", "", "image name")
			cmd.Flags().String("classification", "", "image classification")
			cmd.Flags().String("expires-in", "", "expires-in duration")
			cmd.Flags().String("description", "", "image description")
			cmd.Flags().StringSlice("features", nil, "image features can be machine and/or firewall")
			cmd.Flags().StringSlice("labels", nil, "labels to add to the image")
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
			cmd.Flags().StringSlice("labels", nil, "labels to replace for the image")
			cmd.Flags().StringSlice("add-labels", nil, "labels to add to the image")
			cmd.Flags().StringSlice("remove-labels", nil, "labels to remove to the image")
		},
		UpdateRequestFromCLI: w.updateFromCLI,
	}

	usageCmd := &cobra.Command{
		Use:   "usage",
		Short: "show image usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.usage()
		},
	}

	queryFlags(usageCmd)

	return genericcli.NewCmds(cmdsConfig, usageCmd)
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
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, fmt.Errorf("failed to create image: %w", err)
	}

	return resp.Image, nil
}

func (c *image) createFromCLI() (*adminv2.ImageServiceCreateRequest, error) {
	var expiresAt *timestamppb.Timestamp
	if viper.IsSet("expires-in") {
		expiresAt = timestamppb.New(time.Now().Add(viper.GetDuration("expires-in")))
	}

	labels, err := helpers.LabelsFromSlice(viper.GetStringSlice("labels"))
	if err != nil {
		return nil, err
	}

	return &adminv2.ImageServiceCreateRequest{
		Image: &apiv2.Image{
			Id:             viper.GetString("id"),
			Url:            viper.GetString("url"),
			Classification: helpers.ImageClassificationFromString(viper.GetString("classification")),
			Name:           pointer.PointerOrNil(viper.GetString("name")),
			Description:    pointer.PointerOrNil(viper.GetString("description")),
			ExpiresAt:      expiresAt,
			Features:       helpers.ImageFeaturesFromString(viper.GetStringSlice("features")),
			Meta: &apiv2.Meta{
				Labels: labels,
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
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Image().List(ctx, &apiv2.ImageServiceListRequest{Query: &apiv2.ImageQuery{
		Id:          pointer.PointerOrNil(viper.GetString("id")),
		Os:          pointer.PointerOrNil(viper.GetString("os")),
		Version:     pointer.PointerOrNil(viper.GetString("version")),
		Name:        pointer.PointerOrNil(viper.GetString("name")),
		Description: pointer.PointerOrNil(viper.GetString("description")),
		Feature:     helpers.ImageFeatureFromString(viper.GetString("feature")),
	}})
	if err != nil {
		return nil, fmt.Errorf("failed to get images: %w", err)
	}

	return resp.Images, nil
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
			UpdateMeta:     helpers.UpdateMetaFromMeta(r.Meta),
			Labels:         helpers.UpdateLabelsFromMeta(r.Meta),
			Id:             r.Id,
			Url:            &r.Url,
			Name:           r.Name,
			Description:    r.Description,
			Features:       r.Features,
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

	updateLabels, err := helpers.UpdateLabelsFromCLI()
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
		Features: helpers.ImageFeaturesFromString(viper.GetStringSlice("features")),
		Labels:   updateLabels,
	}

	if viper.IsSet("expires-in") {
		req.ExpiresAt = timestamppb.New(time.Now().Add(viper.GetDuration("expires-in")))
	}
	if viper.IsSet("classification") {
		req.Classification = helpers.ImageClassificationFromString(viper.GetString("classification"))
	}

	return req, nil
}

func (c *image) usage() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Adminv2().Image().Usage(ctx, &adminv2.ImageServiceUsageRequest{
		Query: &apiv2.ImageQuery{
			Id:          pointer.PointerOrNil(viper.GetString("id")),
			Os:          pointer.PointerOrNil(viper.GetString("os")),
			Version:     pointer.PointerOrNil(viper.GetString("version")),
			Name:        pointer.PointerOrNil(viper.GetString("name")),
			Description: pointer.PointerOrNil(viper.GetString("description")),
			Feature:     helpers.ImageFeatureFromString(viper.GetString("feature")),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get images: %w", err)
	}

	return c.c.ListPrinter.Print(resp.ImageUsage)
}
