package main

import (
	"glofox/internal/constants"
	"glofox/internal/handlers"
	"glofox/internal/repository"
	"glofox/internal/services"
	"log"
)

func main() {

	// Initialize repositories
	classRepo := repository.NewClassRepo()
	bookingRepo := repository.NewBookingRepo()

	// Initialize service
	service := services.NewClassService(classRepo, bookingRepo)

	// Initialize handler
	handler := handlers.NewClassHandler(service)

	// Set up router with handler
	router := handlers.SetupRouter(handler)

	if err := router.Run(constants.APIServerPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
