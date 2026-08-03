package mcp

import (
	"context"
	"fmt"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/api/go/metalstack/api/v2/apiv2connect"
	"github.com/metal-stack/api/go/metalstack/api/v2/apiv2mcp"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/v"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/redpanda-data/protoc-gen-go-mcp/pkg/runtime"
	"github.com/redpanda-data/protoc-gen-go-mcp/pkg/runtime/gosdk"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NewMCPCmd(c *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "start mcp server",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, mcps := gosdk.NewServer("metal-stack.io mcp server", v.V.String())

			registerApiv2Clients(c, mcps)

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
	apiv2mcp.ForwardToAuditServiceClient(mcps, &auditServiceGRPCAdapter{client: c.Client.Apiv2().Audit()})
	// apiv2mcp.ForwardToFilesystemServiceClient(mcps, filesystemService)
	// apiv2mcp.ForwardToHealthServiceClient(mcps, healthService)
	// apiv2mcp.ForwardToImageServiceClient(mcps, imageService)
	// apiv2mcp.ForwardToIPServiceClient(mcps, ipService)
	// apiv2mcp.ForwardToMachineServiceClient(mcps, machineService)
	// apiv2mcp.ForwardToMethodServiceClient(mcps, methodService)
	// apiv2mcp.ForwardToNetworkServiceClient(mcps, networkService)
	// apiv2mcp.ForwardToPartitionServiceClient(mcps, partitionService)
	// apiv2mcp.ForwardToProjectServiceClient(mcps, projectService)
	// apiv2mcp.ForwardToSizeImageConstraintServiceClient(mcps, sizeImageConstraintService)
	// apiv2mcp.ForwardToSizeReservationServiceClient(mcps, sizeReservationService)
	// apiv2mcp.ForwardToSizeServiceClient(mcps, sizeService)
	// apiv2mcp.ForwardToTenantServiceClient(mcps, tenantService)
	// apiv2mcp.ForwardToTokenServiceClient(mcps, tokenService)
	// apiv2mcp.ForwardToUserServiceClient(mcps, userService)
	// apiv2mcp.ForwardToVersionServiceClient(mcps, versionService)
}

// Created a PR to fix this in the upstream grpc-mcp compiler: https://github.com/redpanda-data/protoc-gen-go-mcp/pull/51
type auditServiceGRPCAdapter struct {
	client apiv2connect.AuditServiceClient
}

func (a *auditServiceGRPCAdapter) Get(ctx context.Context, req *apiv2.AuditServiceGetRequest, _ ...grpc.CallOption) (*apiv2.AuditServiceGetResponse, error) {
	return a.client.Get(ctx, req)
}

func (a *auditServiceGRPCAdapter) List(ctx context.Context, req *apiv2.AuditServiceListRequest, _ ...grpc.CallOption) (*apiv2.AuditServiceListResponse, error) {
	return a.client.List(ctx, req)
}
