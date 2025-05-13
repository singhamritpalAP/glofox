package models

import (
	"time"
)

// Class represents a class with its details
type Class struct {
	Name      string
	StartDate time.Time
	EndDate   time.Time
	Capacity  int
}

// Booking represents a booking for a class on a specific date
type Booking struct {
	MemberName string    `json:"name"`
	Date       time.Time `json:"date"`
}

// ClassRequest represents the JSON request for /classes
type ClassRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Capacity  int    `json:"capacity" binding:"required,gt=0"`
}

// BookingRequest represents the JSON request for /bookings
type BookingRequest struct {
	ClassName  string `json:"class_name" binding:"required"`
	MemberName string `json:"name" binding:"required"`
	Date       string `json:"date" binding:"required"`
}

// Response represents the JSON response
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
