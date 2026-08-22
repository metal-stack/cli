package mcp

import (
	"fmt"
	"log/slog"

	"github.com/metal-stack/api/go/client"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/v"
	"github.com/modelcontextprotocol/go-sdk/mcp"
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
			// Disable logging because it would interfere with the stdout mcp protocol
			slog.SetDefault(slog.New(slog.DiscardHandler))

			raw, mcps := gosdk.NewServer("metal-stack.io mcp server", v.V.String())

			client.ForwardToAdminv2(c.Client.Adminv2(), mcps)
			client.ForwardToApiv2(c.Client.Apiv2(), mcps)

			if err := raw.Run(cmd.Context(), &mcp.StdioTransport{}); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
