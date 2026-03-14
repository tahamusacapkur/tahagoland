package controller

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewTestServer() *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "test", Version: "v1.0.0"},
		nil,
	)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "greet",
		Description: "Test: Bir kişiye merhaba der",
	}, testGreet)
	return server
}

func testGreet(ctx context.Context, req *mcp.CallToolRequest, input GreetInput) (*mcp.CallToolResult, GreetOutput, error) {
	return nil, GreetOutput{Greeting: "🧪 [Test] Merhaba " + input.Name + "!"}, nil
}
