package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"glofox/internal/constants"
	"glofox/internal/models"
	"glofox/internal/services"
	"glofox/internal/utils"
	"net/http"
)

// ClassHandler handles HTTP requests
type ClassHandler struct {
	service services.IService
}

func NewClassHandler(service services.IService) IHandler {
	return &ClassHandler{service: service}
}

// CreateBooking handles POST /bookings
func (h *ClassHandler) CreateBooking(ctx *gin.Context) {
	var req models.BookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.HandleErrorResp(ctx, http.StatusBadRequest, err, constants.ErrInvalidReq)
		return
	}

	err := h.service.BookClass(req.ClassName, req.MemberName, req.Date)
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, constants.ErrClassNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.HandleErrorResp(ctx, statusCode, err, "")
		return
	}

	ctx.JSON(http.StatusCreated, models.Response{
		Status:  constants.SuccessMsg,
		Message: fmt.Sprintf("Booking created for %s on %s for class %s", req.MemberName, req.Date, req.ClassName),
	})
}
