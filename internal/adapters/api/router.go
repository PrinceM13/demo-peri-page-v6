package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures the Gin router with all routes and middleware.
func SetupRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", handler.HealthCheck)

	// Print endpoint
	router.POST("/print", handler.Print)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
