package handlers

import (
	"github.com/gin-gonic/gin"
	"glofox/internal/constants"
)

// SetupRouter configures the Gin router with handlers
func SetupRouter(handler IHandler) *gin.Engine {
	router := gin.Default()
	// Middleware to handle panics and recover
	router.Use(gin.Recovery())
	// Middleware for logging requests
	router.Use(gin.Logger())

	// Define API endpoints
	router.POST(constants.ClassEndpoint, handler.CreateClass)
	router.POST(constants.BookingEndpoint, handler.CreateBooking)

	return router
}
