package handlers

import "github.com/gin-gonic/gin"

// IHandler defines functions in handlers
type IHandler interface {
	CreateClass(ctx *gin.Context)
	CreateBooking(ctx *gin.Context)
}
