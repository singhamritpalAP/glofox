package utils

import (
	"github.com/gin-gonic/gin"
	"glofox/internal/constants"
	"glofox/internal/models"
	"log"
	"time"
)

// HandleErrorResp logs and writes error on the context
func HandleErrorResp(ctx *gin.Context, statusCode int, err error, errMsg string) {
	log.Println(errMsg, err)
	ctx.JSON(statusCode, models.Response{
		Status:  "error",
		Message: errMsg + err.Error(),
	})
}

// ToMidnightUTC normalizes a time.Time to midnight UTC (00:00:00.000).
func ToMidnightUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// IsDateInRange validates if the date is in specified range
func IsDateInRange(date, start, end time.Time) bool {
	date = ToMidnightUTC(date)
	start = ToMidnightUTC(start)
	end = ToMidnightUTC(end)
	return !date.Before(start) && !date.After(end)
}

// IsValidDate validates if the date is valid or not
func IsValidDate(startDate, endDate time.Time) error {
	if startDate.After(endDate) {
		return constants.ErrInvalidStartEndDate
	}
	return nil
}
