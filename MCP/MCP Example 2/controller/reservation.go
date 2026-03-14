package controller

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewReservationServer() *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "reservation", Version: "v1.0.0"},
		nil,
	)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "greet",
		Description: "Reservation: Bir kişiye merhaba der",
	}, reservationGreet)
	return server
}

func reservationGreet(ctx context.Context, req *mcp.CallToolRequest, input GreetInput) (*mcp.CallToolResult, GreetOutput, error) {
	return nil, GreetOutput{Greeting: "📅 [Reservation] Merhaba " + input.Name + "!"}, nil
}
