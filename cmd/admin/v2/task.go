package v2

import (
	"fmt"
	"strings"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type task struct {
	c *config.Config
}

func newTaskCmd(c *config.Config) *cobra.Command {
	w := &task{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *adminv2.TaskInfo]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "task",
		Plural:          "tasks",
		Description:     "get task insights",
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("queue", "default", "the queue for which tasks should be described")
		},
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("queue", "", "the queue for which tasks should be listed")
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("queue", "default", "the queue of the task which should be delete")
		},
		OnlyCmds: genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DescribeCmd, genericcli.DeleteCmd),
	}

	queueCmd := &cobra.Command{
		Use:   "queues",
		Short: "list all queues",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.queues()
		},
	}

	return genericcli.NewCmds(cmdsConfig, queueCmd)
}

func (t *task) queues() error {
	ctx, cancel := t.c.NewRequestContext()
	defer cancel()

	req := &adminv2.TaskServiceQueuesRequest{}

	resp, err := t.c.Client.Adminv2().Task().Queues(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get task queues: %w", err)
	}

	_, _ = fmt.Fprint(t.c.Out, strings.Join(resp.Queues, "\n"))
	_, _ = fmt.Fprint(t.c.Out, "\n")
	return err
}

func (t *task) Get(id string) (*adminv2.TaskInfo, error) {
	ctx, cancel := t.c.NewRequestContext()
	defer cancel()

	req := &adminv2.TaskServiceGetRequest{TaskId: id, Queue: viper.GetString("queue")}

	resp, err := t.c.Client.Adminv2().Task().Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return resp.Task, nil
}
func (t *task) List() ([]*adminv2.TaskInfo, error) {
	ctx, cancel := t.c.NewRequestContext()
	defer cancel()

	req := &adminv2.TaskServiceListRequest{}
	if viper.IsSet("queue") {
		req.Queue = new(viper.GetString("queue"))
	}

	resp, err := t.c.Client.Adminv2().Task().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return resp.Tasks, nil
}

func (t *task) Create(rq any) (*adminv2.TaskInfo, error) {
	panic("unimplemented")
}

func (t *task) Delete(id string) (*adminv2.TaskInfo, error) {
	ctx, cancel := t.c.NewRequestContext()
	defer cancel()

	req := &adminv2.TaskServiceDeleteRequest{TaskId: id, Queue: viper.GetString("queue")}

	_, err := t.c.Client.Adminv2().Task().Delete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return nil, nil
}

func (t *task) Convert(r *adminv2.TaskInfo) (string, any, any, error) {
	panic("unimplemented")
}

func (t *task) Update(rq any) (*adminv2.TaskInfo, error) {
	panic("unimplemented")
}
