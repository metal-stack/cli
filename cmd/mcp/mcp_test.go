package mcp

import (
	"log"
	"os/exec"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

func TestMCPServer(t *testing.T) {
	ctx := t.Context()
	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// Connect to a server over stdin/stdout.
	transport := &mcp.CommandTransport{Command: exec.Command("/home/stefan/bin/metalctlv2", "mcp")}
	session, err := client.Connect(ctx, transport, nil)
	require.NoError(t, err)
	defer session.Close()

	tools, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	require.NoError(t, err)
	for _, tool := range tools.Tools {
		t.Logf("%s", tool.Name)
		t.Logf("%#v", tool)
	}

	// Call a tool on the server.
	params := &mcp.CallToolParams{
		Name:      "metalstack_api_v2_AuditService_Get",
		Arguments: map[string]any{"name": "you"},
	}
	res, err := session.CallTool(ctx, params)
	require.NoError(t, err)
	if res.IsError {
		t.Logf("tool failed %v", res.Content)
		t.Fail()
	}
	for _, c := range res.Content {
		log.Print(c.(*mcp.TextContent).Text)
	}
}
