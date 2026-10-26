package v2

import (
	"fmt"

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

type projectMember struct {
	c *config.Config
}

func newProjectMemberCmd(c *config.Config) *cobra.Command {
	w := &projectMember{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv2.ProjectMemberServiceCreateRequest, *apiv2.ProjectMemberServiceUpdateRequest, *apiv2.ProjectMember]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "member",
		Plural:          "members",
		Description:     "manage api project members",
		Sorter:          sorters.ProjectMemberSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("project", "", "the project to describe the project members of, defaults to project of the default project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
		},
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("project", "", "the project to list the project members of, defaults to project of the default project")
			cmd.Flags().String("role", "", "the role of the member")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("role", c.Completion.ProjectRole))
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("project", "", "the project for which to create a project member, defaults to project of the default project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("project", "", "the project for which to update a project member, defaults to project of the default project")
			cmd.Flags().String("role", "", "the role of the member")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("role", c.Completion.ProjectRole))
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("project", "", "the project from which to delete a project member, defaults to project of the default project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
		},
		EditCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("project", "", "the project from which to edit a project member, defaults to project of the default project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.Project))
		},
		CreateRequestFromCLI: w.createRequestFromCLI,
		UpdateRequestFromCLI: w.updateRequestFromCLI,
		ValidArgsFn:          w.c.Completion.ProjectMember,
	}

	leaveCmd := &cobra.Command{
		Use:   "leave",
		Short: "leave the given project",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return w.leave()
		},
	}
	leaveCmd.Flags().String("project", "", "the project to leave, defaults to project of the default project")
	genericcli.Must(leaveCmd.RegisterFlagCompletionFunc("project", c.Completion.Project))

	return genericcli.NewCmds(cmdsConfig, leaveCmd)
}

func (c *projectMember) Get(id string) (*apiv2.ProjectMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	project := c.c.GetProject()

	resp, err := c.c.Client.Apiv2().ProjectMember().Get(ctx, &apiv2.ProjectMemberServiceGetRequest{
		Project: project,
		Member:  id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Member, nil
}

func (c *projectMember) List() ([]*apiv2.ProjectMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	project := c.c.GetProject()

	resp, err := c.c.Client.Apiv2().ProjectMember().List(ctx, &apiv2.ProjectMemberServiceListRequest{
		Project: project,
	})
	if err != nil {
		return nil, err
	}

	return resp.Members, nil
}

func (c *projectMember) Create(rq *apiv2.ProjectMemberServiceCreateRequest) (*apiv2.ProjectMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().ProjectMember().Create(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, err
	}

	return resp.Member, nil
}

func (c *projectMember) Delete(id string) (*apiv2.ProjectMember, error) {
	project := c.c.GetProject()

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().ProjectMember().Delete(ctx, &apiv2.ProjectMemberServiceDeleteRequest{
		Project: project,
		Member:  id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to remove member from project: %w", err)
	}

	return resp.Member, nil
}

func (c *projectMember) Convert(r *apiv2.ProjectMember) (string, *apiv2.ProjectMemberServiceCreateRequest, *apiv2.ProjectMemberServiceUpdateRequest, error) {
	return r.Member, &apiv2.ProjectMemberServiceCreateRequest{
			Project: r.Project,
			Member:  r.Member,
			Role:    r.Role,
			Labels:  pointer.SafeDeref(r.Meta).Labels,
		},
		&apiv2.ProjectMemberServiceUpdateRequest{
			Project:    r.Project,
			Member:     r.Member,
			Role:       pointer.PointerOrNil(r.Role),
			UpdateMeta: helpers.UpdateMetaFromMeta(r.Meta),
			Labels:     helpers.UpdateLabelsFromMeta(r.Meta),
		},
		nil
}

func (c *projectMember) Update(rq *apiv2.ProjectMemberServiceUpdateRequest) (*apiv2.ProjectMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().ProjectMember().Update(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update member: %w", err)
	}

	return resp.Member, nil
}

func (c *projectMember) updateRequestFromCLI(args []string) (*apiv2.ProjectMemberServiceUpdateRequest, error) {
	member, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	project := c.c.GetProject()

	return &apiv2.ProjectMemberServiceUpdateRequest{
		Project: project,
		Member:  member,
		Role:    pointer.PointerOrNil(apiv2.ProjectRole(apiv2.ProjectRole_value[viper.GetString("role")])),
	}, nil
}

func (c *projectMember) createRequestFromCLI() (*apiv2.ProjectMemberServiceCreateRequest, error) {
	project := c.c.GetProject()

	return &apiv2.ProjectMemberServiceCreateRequest{
		Project: project,
		Member:  viper.GetString("id"),
		Role:    apiv2.ProjectRole(apiv2.ProjectRole_value[viper.GetString("role")]),
	}, nil
}

func (c *projectMember) leave() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().ProjectMember().Leave(ctx, &apiv2.ProjectMemberServiceLeaveRequest{
		Project: c.c.GetProject(),
	})
	if err != nil {
		return fmt.Errorf("failed to leave project: %w", err)
	}

	return c.c.DescribePrinter.Print(resp.Member)
}
