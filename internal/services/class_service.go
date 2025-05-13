package services

import (
	"glofox/internal/constants"
	"glofox/internal/models"
	"glofox/internal/repository"
	"glofox/internal/utils"
	"log"
	"runtime/debug"
	"time"
)

type ClassService struct {
	classRepo   repository.ClassRepository
	bookingRepo repository.BookingRepository
}

func NewClassService(classRepo repository.ClassRepository, bookingRepo repository.BookingRepository) *ClassService {
	return &ClassService{
		classRepo:   classRepo,
		bookingRepo: bookingRepo,
	}
}

// CreateClass adds a new class
func (service *ClassService) CreateClass(name, startDateStr, endDateStr string, capacity int) (err error) {
	// Recover from panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v\nStack trace:\n%s", r, debug.Stack())
			err = constants.ErrInternalServer
		}
	}()

	// specific date format validation
	startDate, err := time.Parse(constants.DateFormat, startDateStr)
	if err != nil {
		return constants.ErrInvalidStartDate
	}

	endDate, err := time.Parse(constants.DateFormat, endDateStr)
	if err != nil {
		return constants.ErrInvalidEndDate
	}

	err = utils.IsValidDate(startDate, endDate)
	if err != nil {
		return err
	}

	class := models.Class{
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		Capacity:  capacity,
	}
	return service.classRepo.Create(class)
}
