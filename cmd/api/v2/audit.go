package v2

import (
	"fmt"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/cli/cmd/sorters"
	helpersaudit "github.com/metal-stack/cli/pkg/helpers/audit"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type audit struct {
	c *config.Config
}

func newAuditCmd(c *config.Config) *cobra.Command {
	w := &audit{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv2.AuditTrace]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:        "audit",
		Plural:          "audits",
		Description:     "read api audit traces of a tenant",
		Sorter:          sorters.AuditSorter(),
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DescribeCmd),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("tenant", "", "tenant of the audit trace.")

			cmd.Flags().String("phase", "", "the audit trace phase.")

			cmd.Flags().Bool("prettify-body", true, "attempts to interpret the body as json and prettifies it.")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("phase", c.Completion.AuditPhaseListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))
		},
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("request-id", "", "request id of the audit trace.")

			cmd.Flags().String("from", "", "start of range of the audit traces. e.g. 1h, 10m, 2006-01-02 15:04:05")
			cmd.Flags().String("to", "", "end of range of the audit traces. e.g. 1h, 10m, 2006-01-02 15:04:05")

			cmd.Flags().String("user", "", "user of the audit trace.")
			cmd.Flags().String("tenant", "", "tenant of the audit trace.")

			cmd.Flags().String("project", "", "project id of the audit trace")

			cmd.Flags().String("phase", "", "the audit trace phase.")
			cmd.Flags().String("method", "", "api method of the audit trace.")
			cmd.Flags().Int32("result-code", 0, "gRPC result status code of the audit trace.")
			cmd.Flags().String("source-ip", "", "source-ip of the audit trace.")

			cmd.Flags().String("body", "", "filters audit trace body payloads for the given text (full-text search).")

			cmd.Flags().Int64("limit", 0, "limit the number of audit traces.")

			cmd.Flags().Bool("prettify-body", true, "attempts to interpret the body as json and prettifies it.")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("phase", c.Completion.AuditPhaseListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("result-code", c.Completion.AuditStatusCodesCompletion))
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *audit) Get(id string) (*apiv2.AuditTrace, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Client.Apiv2().Audit().Get(ctx, &apiv2.AuditServiceGetRequest{
		Login: tenant,
		Uuid:  id,
		Phase: helpersaudit.ToPhase(viper.GetString("phase")),
	},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit trace: %w", err)
	}

	if viper.GetBool("prettify-body") {
		helpersaudit.TryPrettifyBody(resp.Trace)
	}

	return resp.Trace, nil
}

func (c *audit) List() ([]*apiv2.AuditTrace, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	fromDateTime, err := helpersaudit.EventuallyRelativeDateTime(viper.GetString("from"))
	if err != nil {
		return nil, err
	}
	toDateTime, err := helpersaudit.EventuallyRelativeDateTime(viper.GetString("to"))
	if err != nil {
		return nil, err
	}

	tenant, err := c.c.GetTenant()
	if err != nil {
		return nil, fmt.Errorf("tenant is required: %w", err)
	}

	var code *int32
	if viper.IsSet("result-code") {
		code = new(viper.GetInt32("result-code"))
	}

	resp, err := c.c.Client.Apiv2().Audit().List(ctx, &apiv2.AuditServiceListRequest{
		Login: tenant,
		Query: &apiv2.AuditQuery{
			Uuid:       pointer.PointerOrNil(viper.GetString("request-id")),
			From:       fromDateTime,
			To:         toDateTime,
			User:       pointer.PointerOrNil(viper.GetString("user")),
			Project:    pointer.PointerOrNil(viper.GetString("project")),
			Method:     pointer.PointerOrNil(viper.GetString("method")),
			ResultCode: code,
			Body:       pointer.PointerOrNil(viper.GetString("body")),
			SourceIp:   pointer.PointerOrNil(viper.GetString("source-ip")),
			Limit:      pointer.PointerOrNil(viper.GetInt32("limit")),
			Phase:      helpersaudit.ToPhase(viper.GetString("phase")),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list audit traces: %w", err)
	}

	if viper.GetBool("prettify-body") {
		for _, trace := range resp.Traces {
			helpersaudit.TryPrettifyBody(trace)
		}
	}

	return resp.Traces, nil
}

func (c *audit) Create(rq any) (*apiv2.AuditTrace, error) {
	panic("unimplemented")
}

func (c *audit) Delete(id string) (*apiv2.AuditTrace, error) {
	panic("unimplemented")
}

func (t *audit) Convert(r *apiv2.AuditTrace) (string, any, any, error) {
	panic("unimplemented")
}

func (t *audit) Update(rq any) (*apiv2.AuditTrace, error) {
	panic("unimplemented")
}
