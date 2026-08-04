package mcp

import (
	"fmt"

	"github.com/metal-stack/api/go/metalstack/admin/v2/adminv2mcp"
	"github.com/metal-stack/api/go/metalstack/api/v2/apiv2mcp"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/v"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/redpanda-data/protoc-gen-go-mcp/pkg/runtime"
	"github.com/redpanda-data/protoc-gen-go-mcp/pkg/runtime/gosdk"

	"github.com/spf13/cobra"
)

func NewMCPCmd(c *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "start mcp server",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, mcps := gosdk.NewServer("metal-stack.io mcp server", v.V.String())

			registerApiv2Clients(c, mcps)
			registerAdminv2Clients(c, mcps)

			if err := raw.Run(cmd.Context(), &mcp.StdioTransport{}); err != nil {
				fmt.Printf("Server error: %v\n", err)
			}

			return nil
		},
	}
	return cmd
}

func registerApiv2Clients(c *config.Config, mcps runtime.MCPServer) {

	// Forward to MCP Clients
	apiv2mcp.ForwardToAuditServiceClient(mcps, c.Client.Apiv2().Audit())
	apiv2mcp.ForwardToFilesystemServiceClient(mcps, c.Client.Apiv2().Filesystem())
	apiv2mcp.ForwardToHealthServiceClient(mcps, c.Client.Apiv2().Health())
	apiv2mcp.ForwardToImageServiceClient(mcps, c.Client.Apiv2().Image())
	apiv2mcp.ForwardToIPServiceClient(mcps, c.Client.Apiv2().IP())
	apiv2mcp.ForwardToMachineServiceClient(mcps, c.Client.Apiv2().Machine())
	apiv2mcp.ForwardToMethodServiceClient(mcps, c.Client.Apiv2().Method())
	apiv2mcp.ForwardToNetworkServiceClient(mcps, c.Client.Apiv2().Network())
	apiv2mcp.ForwardToPartitionServiceClient(mcps, c.Client.Apiv2().Partition())
	apiv2mcp.ForwardToProjectServiceClient(mcps, c.Client.Apiv2().Project())
	apiv2mcp.ForwardToSizeImageConstraintServiceClient(mcps, c.Client.Apiv2().SizeImageConstraint())
	apiv2mcp.ForwardToSizeReservationServiceClient(mcps, c.Client.Apiv2().SizeReservation())
	apiv2mcp.ForwardToSizeServiceClient(mcps, c.Client.Apiv2().Size())
	apiv2mcp.ForwardToTenantServiceClient(mcps, c.Client.Apiv2().Tenant())
	apiv2mcp.ForwardToTokenServiceClient(mcps, c.Client.Apiv2().Token())
	apiv2mcp.ForwardToUserServiceClient(mcps, c.Client.Apiv2().User())
	apiv2mcp.ForwardToVersionServiceClient(mcps, c.Client.Apiv2().Version())
}

func registerAdminv2Clients(c *config.Config, mcps runtime.MCPServer) {

	// Forwart to Admin MCP Clients
	adminv2mcp.ForwardToAuditServiceClient(mcps, c.Client.Adminv2().Audit())
	adminv2mcp.ForwardToComponentServiceClient(mcps, c.Client.Adminv2().Component())
	adminv2mcp.ForwardToFilesystemServiceClient(mcps, c.Client.Adminv2().Filesystem())
	adminv2mcp.ForwardToImageServiceClient(mcps, c.Client.Adminv2().Image())
	adminv2mcp.ForwardToIPServiceClient(mcps, c.Client.Adminv2().IP())
	adminv2mcp.ForwardToMachineServiceClient(mcps, c.Client.Adminv2().Machine())
	adminv2mcp.ForwardToNetworkServiceClient(mcps, c.Client.Adminv2().Network())
	adminv2mcp.ForwardToPartitionServiceClient(mcps, c.Client.Adminv2().Partition())
	adminv2mcp.ForwardToProjectServiceClient(mcps, c.Client.Adminv2().Project())
	adminv2mcp.ForwardToSizeImageConstraintServiceClient(mcps, c.Client.Adminv2().SizeImageConstraint())
	adminv2mcp.ForwardToSizeReservationServiceClient(mcps, c.Client.Adminv2().SizeReservation())
	adminv2mcp.ForwardToSizeServiceClient(mcps, c.Client.Adminv2().Size())
	adminv2mcp.ForwardToSwitchServiceClient(mcps, c.Client.Adminv2().Switch())
	adminv2mcp.ForwardToTaskServiceClient(mcps, c.Client.Adminv2().Task())
	adminv2mcp.ForwardToTenantServiceClient(mcps, c.Client.Adminv2().Tenant())
	adminv2mcp.ForwardToTokenServiceClient(mcps, c.Client.Adminv2().Token())
	adminv2mcp.ForwardToVPNServiceClient(mcps, c.Client.Adminv2().VPN())
}
