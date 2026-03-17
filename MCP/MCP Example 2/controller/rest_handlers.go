package controller

import (
	"net/http"

	"MCPExample2/service"

	"github.com/gin-gonic/gin"
)

// HandleRoomPrice is the REST endpoint for ChatGPT GPT Actions.
// POST /api/v1/room-price
func HandleRoomPrice(c *gin.Context) {
	var input service.RoomPriceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek: " + err.Error()})
		return
	}

	resp, err := service.GetRoomPrice(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if len(resp.Rooms) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Seçtiğiniz tarihler için müsait oda bulunamadı.",
			"data":    resp,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}
