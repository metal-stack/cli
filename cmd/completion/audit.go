package completion

import (
	"strconv"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
)

func (c *Completion) AuditPhase(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return c.genericEnums(apiv2.AuditPhase_name)
}

func (c *Completion) AuditStatusCodes(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var result []string

	for i := range 16 {
		result = append(result, strconv.Itoa(i)+"\t"+codes.Code(uint32(i)).String())
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}
