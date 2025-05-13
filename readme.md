# Glofox Class Booking API

This is a RESTful API for managing boutiques, studios, and gyms classes and bookings, built with Go and the Gin framework. It supports creating classes and booking members into classes.

## Prerequisites

- Go 1.22
- Git
- Curl (for testing)

## Setup and Running

Follow these steps to clone and run the application locally.
- **Clone the Repository**:
   ```bash
   git clone https://github.com/singhamritpalAP/glofox.git
   ```
- **Navigate to the cmd Directory:**
   ```bash
   cd glofox/cmd/
   ```
- **Run the Application:**
    ```bash
    go run .
   ```
## Run Happy Flow Tests
- Open a new terminal and execute the `curl` commands from the `README.md`:
   - Create a Class
     ```bash
     curl -X POST http://localhost:8080/classes -H "Content-Type: application/json" -d '{"name":"Yoga","start_date":"2025-06-01","end_date":"2025-06-20","capacity":10}'
     ```
     Expected Response (HTTP 201):
     {
     "status": "success",
     "message": "Class Pilates created successfully"
     }
  
   - Create a Booking
   - ```bash   
     curl -X POST http://localhost:8080/bookings -H "Content-Type: application/json" -d '{"class_name":"Yoga","name":"Amrit","date":"2025-06-10"}'
     ```
     Expected Response (HTTP 201):
     {
     "status": "success",
     "message": "Booking created for Amrit on 2025-06-10 for class Yoga"
     }