package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"glofox/internal/constants"
	"glofox/internal/models"
	"glofox/internal/utils"
	"net/http"
)

// CreateClass handles POST /classes
func (h *ClassHandler) CreateClass(ctx *gin.Context) {
	var req models.ClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.HandleErrorResp(ctx, http.StatusBadRequest, err, constants.ErrInvalidReq)
		return
	}

	err := h.service.CreateClass(req.Name, req.StartDate, req.EndDate, req.Capacity)
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, constants.ErrClassAlreadyExists) {
			statusCode = http.StatusConflict
		}
		utils.HandleErrorResp(ctx, statusCode, err, "")
		return
	}

	ctx.JSON(http.StatusCreated, models.Response{
		Status:  constants.SuccessMsg,
		Message: fmt.Sprintf("Class %s created successfully", req.Name),
	})
}
