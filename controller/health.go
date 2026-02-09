package controller

import (
	"net/http"

	"projectwebcurhat/contract"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	service *contract.Service
}

func (h *HealthController) GetPrefix() string {
	return ""
}

func (h *HealthController) InitService(service *contract.Service) {
	h.service = service
}

func (h *HealthController) InitRoute(app *gin.RouterGroup) {
	app.GET("/health", h.HandleHealth)
}

func (h *HealthController) HandleHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Server is running",
		"data": gin.H{
			"status":     "healthy",
			"room_count": h.service.Room.GetRoomCount(),
		},
	})
}
