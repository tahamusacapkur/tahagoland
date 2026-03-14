package controller

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewUpsellServer() *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "upsell", Version: "v1.0.0"},
		nil,
	)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "greet",
		Description: "Upsell: Bir kişiye merhaba der",
	}, upsellGreet)
	return server
}

func upsellGreet(ctx context.Context, req *mcp.CallToolRequest, input GreetInput) (*mcp.CallToolResult, GreetOutput, error) {
	return nil, GreetOutput{Greeting: "💎 [Upsell] Merhaba " + input.Name + "!"}, nil
}
