package controller

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewHotelContentServer() *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "hotel-content", Version: "v1.0.0"},
		nil,
	)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "greet",
		Description: "Hotel Content: Bir kişiye merhaba der",
	}, hotelContentGreet)
	return server
}

func hotelContentGreet(ctx context.Context, req *mcp.CallToolRequest, input GreetInput) (*mcp.CallToolResult, GreetOutput, error) {
	return nil, GreetOutput{Greeting: "🏨 [Hotel Content] Merhaba " + input.Name + "!"}, nil
}
