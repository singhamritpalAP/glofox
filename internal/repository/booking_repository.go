package repository

import (
	"glofox/internal/utils"
	"sync"
	"time"
)

type BookingRepository interface {
	Create(className, memberName string, date time.Time) error
}

// BookingRepo manages the in-memory booking data
type BookingRepo struct {
	// Key: class name, Sub-key: date, Value: list of member names
	bookings map[string]map[time.Time][]string
	mu       sync.RWMutex
}

// NewBookingRepo creates a new BookingRepo
func NewBookingRepo() *BookingRepo {
	return &BookingRepo{
		bookings: make(map[string]map[time.Time][]string),
	}
}

// Create for creating a new booking
func (bookingRepo *BookingRepo) Create(className, memberName string, date time.Time) error {
	bookingRepo.mu.Lock()
	defer bookingRepo.mu.Unlock()

	// Normalize date to midnight
	date = utils.ToMidnightUTC(date)

	if _, exists := bookingRepo.bookings[className]; !exists {
		bookingRepo.bookings[className] = make(map[time.Time][]string)
	}
	bookingRepo.bookings[className][date] = append(bookingRepo.bookings[className][date], memberName)
	return nil
}
