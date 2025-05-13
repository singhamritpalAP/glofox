package services

import (
	"fmt"
	"glofox/internal/constants"
	"glofox/internal/utils"
	"log"
	"runtime/debug"
	"time"
)

// BookClass creates a booking
func (service *ClassService) BookClass(className, memberName, dateStr string) (err error) {
	// Recover from panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v\nStack trace:\n%s", r, debug.Stack())
			err = constants.ErrInternalServer
		}
	}()

	// specific date format validation
	date, err := time.Parse(constants.DateFormat, dateStr)
	if err != nil {
		return constants.ErrInvalidDate
	}

	// Check if class exists
	class, exists := service.classRepo.GetByName(className)
	if !exists {
		return constants.ErrClassNotFound
	}

	// Check if date is valid for the class
	if !utils.IsDateInRange(date, class.StartDate, class.EndDate) {
		return fmt.Errorf("date %s is not valid for class %s", dateStr, className)
	}

	// Create booking
	return service.bookingRepo.Create(className, memberName, date)
}
