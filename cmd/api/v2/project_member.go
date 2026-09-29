package v2

import (
	"fmt"

	"github.com/metal-stack/api/go/errorutil"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	cliv2 "github.com/metal-stack/api/go/metalstack/cli/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
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

	cmdsConfig := &genericcli.CmdsConfig[*apiv2.ProjectServiceAddMemberRequest, *apiv2.ProjectServiceUpdateMemberRequest, *cliv2.ProjectMember]{
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

	return genericcli.NewCmds(cmdsConfig)
}

func (c *projectMember) Get(id string) (*cliv2.ProjectMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	project := c.c.GetProject()

	resp, err := c.c.Client.Apiv2().Project().Get(ctx, &apiv2.ProjectServiceGetRequest{
		Project: project,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	for _, member := range resp.GetProjectMembers() {
		if member.Id == id {
			return &cliv2.ProjectMember{
				Project:       project,
				ProjectMember: member,
			}, nil
		}
	}

	return nil, fmt.Errorf("unable to find project members %q in project %q", id, project)
}

func (c *projectMember) List() ([]*cliv2.ProjectMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	project := c.c.GetProject()

	resp, err := c.c.Client.Apiv2().Project().Get(ctx, &apiv2.ProjectServiceGetRequest{
		Project: project,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	var res []*cliv2.ProjectMember

	for _, member := range resp.GetProjectMembers() {
		res = append(res, &cliv2.ProjectMember{
			Project:       project,
			ProjectMember: member,
		})
	}

	return res, nil
}

func (c *projectMember) Create(rq *apiv2.ProjectServiceAddMemberRequest) (*cliv2.ProjectMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Project().AddMember(ctx, rq)
	if err != nil {
		if errorutil.IsConflict(err) {
			return nil, genericcli.AlreadyExistsError()
		}

		return nil, err
	}

	return &cliv2.ProjectMember{
		Project:       rq.Project,
		ProjectMember: resp.ProjectMember,
	}, nil
}

func (c *projectMember) Delete(id string) (*cliv2.ProjectMember, error) {
	project := c.c.GetProject()

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Project().RemoveMember(ctx, &apiv2.ProjectServiceRemoveMemberRequest{
		Project: project,
		Member:  id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to remove member from project: %w", err)
	}

	return &cliv2.ProjectMember{
		Project:       project,
		ProjectMember: resp.ProjectMember,
	}, nil
}

func (c *projectMember) Convert(r *cliv2.ProjectMember) (string, *apiv2.ProjectServiceAddMemberRequest, *apiv2.ProjectServiceUpdateMemberRequest, error) {
	return r.ProjectMember.Id, &apiv2.ProjectServiceAddMemberRequest{
			Project: r.Project,
			Member:  r.ProjectMember.Id,
			Role:    r.ProjectMember.Role,
		},
		&apiv2.ProjectServiceUpdateMemberRequest{
			Project: r.Project,
			Member:  r.ProjectMember.Id,
			Role:    r.ProjectMember.Role,
		},
		nil
}

func (c *projectMember) Update(rq *apiv2.ProjectServiceUpdateMemberRequest) (*cliv2.ProjectMember, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv2().Project().UpdateMember(ctx, rq)
	if err != nil {
		return nil, fmt.Errorf("failed to update member: %w", err)
	}

	return &cliv2.ProjectMember{
		Project:       rq.Project,
		ProjectMember: resp.ProjectMember,
	}, nil
}

func (c *projectMember) updateRequestFromCLI(args []string) (*apiv2.ProjectServiceUpdateMemberRequest, error) {
	member, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	project := c.c.GetProject()

	return &apiv2.ProjectServiceUpdateMemberRequest{
		Project: project,
		Member:  member,
		Role:    apiv2.ProjectRole(apiv2.ProjectRole_value[viper.GetString("role")]),
	}, nil
}

func (c *projectMember) createRequestFromCLI() (*apiv2.ProjectServiceAddMemberRequest, error) {
	project := c.c.GetProject()

	return &apiv2.ProjectServiceAddMemberRequest{
		Project: project,
		Member:  viper.GetString("id"),
		Role:    apiv2.ProjectRole(apiv2.ProjectRole_value[viper.GetString("role")]),
	}, nil
}
