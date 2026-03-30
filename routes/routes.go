package routes

import (
	"inventory-service/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes defines all API endpoints and maps them to their respective handlers
func SetupRoutes(r *gin.Engine) {
	r.GET("/inventory", handlers.GetInventory)
	r.POST("/inventory/update-stock", handlers.UpdateStock)
	r.GET("/inventory/low-stock", handlers.GetLowStock)
}
