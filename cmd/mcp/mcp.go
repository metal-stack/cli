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

var long = fmt.Sprintf(`
Use %q to serve as mcp server. You must configure your coding agent to make use of the mcp server.

Example opencode.json:

"mcp": {
  "metal": {
    "type": "local",
    "command": [
      %q,
      "mcp"
    ],
    "enabled": true,
  }
}

Then login with %q and start you coding agent and ask questions like:

- list all metal partitions
- give me all available metal sizes and images
- create a ip address, ask me questions
`, config.BinaryName, config.BinaryName, config.BinaryName)

func NewMCPCmd(c *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "start mcp server",
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, mcps := gosdk.NewServer("metal-stack.io mcp server", v.V.String())

			forwardToClients(c, mcps)

			if err := raw.Run(cmd.Context(), &mcp.StdioTransport{}); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func forwardToClients(c *config.Config, mcps runtime.MCPServer) {
	// Be prepared for later opts.
	var opts []runtime.Option
	// API Clients
	apiv2mcp.ForwardToAuditServiceClient(mcps, c.Client.Apiv2().Audit(), opts...)
	apiv2mcp.ForwardToFilesystemServiceClient(mcps, c.Client.Apiv2().Filesystem(), opts...)
	apiv2mcp.ForwardToHealthServiceClient(mcps, c.Client.Apiv2().Health(), opts...)
	apiv2mcp.ForwardToImageServiceClient(mcps, c.Client.Apiv2().Image(), opts...)
	apiv2mcp.ForwardToIPServiceClient(mcps, c.Client.Apiv2().IP(), opts...)
	apiv2mcp.ForwardToMachineServiceClient(mcps, c.Client.Apiv2().Machine(), opts...)
	apiv2mcp.ForwardToMethodServiceClient(mcps, c.Client.Apiv2().Method(), opts...)
	apiv2mcp.ForwardToNetworkServiceClient(mcps, c.Client.Apiv2().Network(), opts...)
	apiv2mcp.ForwardToPartitionServiceClient(mcps, c.Client.Apiv2().Partition(), opts...)
	apiv2mcp.ForwardToProjectServiceClient(mcps, c.Client.Apiv2().Project(), opts...)
	apiv2mcp.ForwardToSizeImageConstraintServiceClient(mcps, c.Client.Apiv2().SizeImageConstraint(), opts...)
	apiv2mcp.ForwardToSizeReservationServiceClient(mcps, c.Client.Apiv2().SizeReservation(), opts...)
	apiv2mcp.ForwardToSizeServiceClient(mcps, c.Client.Apiv2().Size(), opts...)
	apiv2mcp.ForwardToTenantServiceClient(mcps, c.Client.Apiv2().Tenant(), opts...)
	apiv2mcp.ForwardToTokenServiceClient(mcps, c.Client.Apiv2().Token(), opts...)
	apiv2mcp.ForwardToUserServiceClient(mcps, c.Client.Apiv2().User(), opts...)
	apiv2mcp.ForwardToVersionServiceClient(mcps, c.Client.Apiv2().Version(), opts...)

	// Admin Clients
	adminv2mcp.ForwardToAuditServiceClient(mcps, c.Client.Adminv2().Audit(), opts...)
	adminv2mcp.ForwardToComponentServiceClient(mcps, c.Client.Adminv2().Component(), opts...)
	adminv2mcp.ForwardToFilesystemServiceClient(mcps, c.Client.Adminv2().Filesystem(), opts...)
	adminv2mcp.ForwardToImageServiceClient(mcps, c.Client.Adminv2().Image(), opts...)
	adminv2mcp.ForwardToIPServiceClient(mcps, c.Client.Adminv2().IP(), opts...)
	adminv2mcp.ForwardToMachineServiceClient(mcps, c.Client.Adminv2().Machine(), opts...)
	adminv2mcp.ForwardToNetworkServiceClient(mcps, c.Client.Adminv2().Network(), opts...)
	adminv2mcp.ForwardToPartitionServiceClient(mcps, c.Client.Adminv2().Partition(), opts...)
	adminv2mcp.ForwardToProjectServiceClient(mcps, c.Client.Adminv2().Project(), opts...)
	adminv2mcp.ForwardToSizeImageConstraintServiceClient(mcps, c.Client.Adminv2().SizeImageConstraint(), opts...)
	adminv2mcp.ForwardToSizeReservationServiceClient(mcps, c.Client.Adminv2().SizeReservation(), opts...)
	adminv2mcp.ForwardToSizeServiceClient(mcps, c.Client.Adminv2().Size(), opts...)
	adminv2mcp.ForwardToSwitchServiceClient(mcps, c.Client.Adminv2().Switch(), opts...)
	adminv2mcp.ForwardToTaskServiceClient(mcps, c.Client.Adminv2().Task(), opts...)
	adminv2mcp.ForwardToTenantServiceClient(mcps, c.Client.Adminv2().Tenant(), opts...)
	adminv2mcp.ForwardToTokenServiceClient(mcps, c.Client.Adminv2().Token(), opts...)
	adminv2mcp.ForwardToVPNServiceClient(mcps, c.Client.Adminv2().VPN(), opts...)
}
