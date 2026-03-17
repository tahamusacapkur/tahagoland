package router

import (
	"net/http"

	"MCPExample2/controller"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Setup() *gin.Engine {
	r := gin.Default()

	opts := &mcp.StreamableHTTPOptions{
		Stateless:    true,
		JSONResponse: false,
	}

	// MCP Server endpoint'leri
	register(r, "/test", controller.NewTestServer(), opts)
	register(r, "/reservation", controller.NewReservationServer(), opts)
	register(r, "/hotel-content", controller.NewHotelContentServer(), opts)
	register(r, "/upsell", controller.NewUpsellServer(), opts)

	// REST API endpoint'leri (ChatGPT GPT Actions için)
	api := r.Group("/api/v1")
	{
		api.POST("/room-price", controller.HandleRoomPrice)
	}

	// OpenAPI spec (ChatGPT Actions import için)
	r.StaticFile("/openapi.json", "./openapi.json")

	// Privacy policy (ChatGPT Actions gereksinimi)
	r.GET("/privacy", func(c *gin.Context) {
		c.String(200, "Privacy Policy for ICI Hotel API. Bu API otel oda fiyat sorgulama hizmeti sunar. Kişisel veri saklanmaz.")
	})

	// Sağlık kontrolü
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Ana sayfa
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "MCP + REST API Server çalışıyor",
			"mcp_endpoints": []string{
				"/test",
				"/reservation",
				"/hotel-content",
				"/upsell",
			},
			"rest_endpoints": []string{
				"POST /api/v1/room-price",
			},
			"openapi": "/openapi.json",
		})
	})

	return r
}

func register(r *gin.Engine, path string, server *mcp.Server, opts *mcp.StreamableHTTPOptions) {
	handler := mcp.NewStreamableHTTPHandler(
		func(req *http.Request) *mcp.Server { return server },
		opts,
	)
	r.Any(path, gin.WrapH(handler))
}
