package completion

import (
	"strconv"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
)

func (c *Completion) AuditPhaseListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{apiv2.AuditPhase_AUDIT_PHASE_REQUEST.String(), apiv2.AuditPhase_AUDIT_PHASE_RESPONSE.String()}, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) AuditStatusCodesCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var result []string

	for i := range 16 {
		result = append(result, strconv.Itoa(i)+"\t"+codes.Code(uint32(i)).String())
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}
