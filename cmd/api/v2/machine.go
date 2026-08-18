package v2

import (
	"github.com/metal-stack/api/go/errorutil"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/cli/pkg/helpers"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type machine struct {
	c *config.Config
}

func newMachineCmd(c *config.Config) *cobra.Command {
	w := &machine{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv2.MachineServiceCreateRequest, *apiv2.MachineServiceUpdateRequest, *apiv2.Machine]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "machine",
		Aliases:         []string{"ms"},
		Plural:          "machines",
		Description:     "allocate a machine",
		Sorter:          sorters.MachineSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		CreateRequestFromCLI: func() (*apiv2.MachineServiceCreateRequest, error) {
			return helpers.MachineCreateRequestFromCLI(c)
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			helpers.AddMachineCreateFlags(cmd, "machine", c.Completion)
			cmd.Aliases = []string{"allocate"}
		},
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("uuid", "", "allocation uuid of machine which should be listed")
			cmd.Flags().String("name", "", "name from machines which should be listed")
			cmd.Flags().String("hostname", "", "hostname from machines which should be listed")
			cmd.Flags().String("size", "", "size from machines which should be listed")
			cmd.Flags().String("image", "", "image")
			cmd.Flags().StringP("project", "p", "", "project from where machines should be listed")
			cmd.Flags().StringP("partition", "", "", "partition from where machines should be listed")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("size", c.Completion.Size))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("image", c.Completion.Image))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.Partition))

			cmd.Long = cmd.Short + "\n" + helpers.MachineListEmojiHelpText()
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project from where machines should be listed")
			cmd.Flags().String("description", "", "description of the machine")
			cmd.Flags().StringSlice("labels", nil, "labels to replace for the machine")
			cmd.Flags().StringSlice("add-labels", nil, "labels to add to the machine")
			cmd.Flags().StringSlice("remove-labels", nil, "labels to remove to the machine")
			cmd.Flags().StringP("ssh-public-key", "i", "",
				`SSH public key for access via ssh and console. [optional]
Can be either the public key as string, or pointing to the public key file to use e.g.: "@~/.ssh/id_rsa.pub".
If ~/.ssh/[id_ed25519.pub | id_rsa.pub | id_dsa.pub] is present it will be picked as default, matching the first one in this order.`)

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
		},
		UpdateRequestFromCLI: func(args []string) (*apiv2.MachineServiceUpdateRequest, error) {
			return helpers.MachineUpdateRequestFromCLI(c, args)
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the machine")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the machine")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
		},
		ValidArgsFn: c.Completion.Machine,
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *machine) Create(rq *apiv2.MachineServiceCreateRequest) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Machine().Create(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, err
	}

	return resp.Machine, nil
}

func (c *machine) Delete(id string) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv2.MachineServiceDeleteRequest{
		Uuid:    id,
		Project: c.c.GetProject(),
	}

	if viper.IsSet("file") {
		var err error
		req.Uuid, req.Project, err = helpers.DecodeProject(id)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.c.Client.Apiv2().Machine().Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Machine, nil
}

func (c *machine) Get(id string) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Machine().Get(ctx, &apiv2.MachineServiceGetRequest{
		Project: c.c.GetProject(),
		Uuid:    id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Machine, nil
}

func (c *machine) List() ([]*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	var allocation *apiv2.MachineAllocationQuery

	if viper.IsSet("hostname") || viper.IsSet("image") {
		allocation = &apiv2.MachineAllocationQuery{
			Hostname: pointer.PointerOrNil(viper.GetString("hostname")),
			Image:    pointer.PointerOrNil(viper.GetString("image")),
		}
	}

	resp, err := c.c.Client.Apiv2().Machine().List(ctx, &apiv2.MachineServiceListRequest{
		Project: c.c.GetProject(),
		Query: &apiv2.MachineQuery{
			Uuid:       pointer.PointerOrNil(viper.GetString("id")),
			Name:       pointer.PointerOrNil(viper.GetString("name")),
			Partition:  pointer.PointerOrNil(viper.GetString("partition")),
			Size:       pointer.PointerOrNil(viper.GetString("size")),
			Allocation: allocation,
		},
	})
	if err != nil {
		return nil, err
	}

	return resp.Machines, nil
}

func (c *machine) Update(rq *apiv2.MachineServiceUpdateRequest) (*apiv2.Machine, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Machine().Update(ctx, rq)
	if err != nil {
		return nil, err
	}

	return resp.Machine, nil
}

func (c *machine) Convert(r *apiv2.Machine) (string, *apiv2.MachineServiceCreateRequest, *apiv2.MachineServiceUpdateRequest, error) {
	update, err := helpers.MachineResponseToUpdate(r)
	if err != nil {
		return "", nil, nil, err
	}

	create, err := helpers.MachineResponseToCreate(r)
	if err != nil {
		return "", nil, nil, err
	}

	return helpers.EncodeProject(r.Uuid, r.Allocation.Project), create, update, err
}
