package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"glofox/internal/constants"
	"glofox/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClassService is a mock for ClassService
type MockClassService struct {
	mock.Mock
}

// BookClass mocks the BookClass method
func (m *MockClassService) BookClass(className, memberName, dateStr string) error {
	args := m.Called(className, memberName, dateStr)
	return args.Error(0)
}

func TestClassHandler_CreateBooking(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Define test cases
	tests := []struct {
		name           string
		jsonInput      string
		setupMock      func(*MockClassService)
		expectedStatus int
		expectedBody   models.Response
		expectService  bool
	}{
		{
			name:      "Happy Path",
			jsonInput: `{"class_name":"Yoga","name":"Alice","date":"2025-06-10"}`,
			setupMock: func(m *MockClassService) {
				m.On("BookClass", "Yoga", "Alice", "2025-06-10").Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: models.Response{
				Status:  constants.SuccessMsg,
				Message: "Booking created for Alice on 2025-06-10 for class Yoga",
			},
			expectService: true,
		},
		{
			name:      "Invalid Date Format",
			jsonInput: `{"class_name":"Yoga","name":"Alice","date":"2025-06-10"}`,
			setupMock: func(m *MockClassService) {
				m.On("BookClass", "Yoga", "Alice", "2025-06-10").Return(constants.ErrInvalidDate)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: models.Response{
				Status:  "error",
				Message: constants.ErrInvalidDate.Error(),
			},
			expectService: true,
		},
		{
			name:      "Class Not Found",
			jsonInput: `{"class_name":"Yoga","name":"Alice","date":"2025-06-10"}`,
			setupMock: func(m *MockClassService) {
				m.On("BookClass", "Yoga", "Alice", "2025-06-10").Return(constants.ErrClassNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: models.Response{
				Status:  "error",
				Message: constants.ErrClassNotFound.Error(),
			},
			expectService: true,
		},
		{
			name:      "Invalid Date Range",
			jsonInput: `{"class_name":"Yoga","name":"Alice","date":"2025-06-21"}`,
			setupMock: func(m *MockClassService) {
				m.On("BookClass", "Yoga", "Alice", "2025-06-21").Return(fmt.Errorf("date 2025-06-21 is not valid for class Yoga"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: models.Response{
				Status:  "error",
				Message: "date 2025-06-21 is not valid for class Yoga",
			},
			expectService: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := new(MockClassService)
			tt.setupMock(mockService)
			handler := NewClassHandler(mockService)

			// Create test context and recorder
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			// Create HTTP request
			req, _ := http.NewRequest("POST", "/bookings", bytes.NewBufferString(tt.jsonInput))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req

			// Call handler
			handler.CreateBooking(ctx)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code, "Expected status %d, got %d", tt.expectedStatus, w.Code)

			// Assert response body
			var resp models.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err, "Failed to unmarshal response")
			assert.Equal(t, tt.expectedBody.Status, resp.Status, "Expected status %s, got %s", tt.expectedBody.Status, resp.Status)
			assert.Equal(t, tt.expectedBody.Message, resp.Message, "Expected message %s, got %s", tt.expectedBody.Message, resp.Message)

			// Assert service calls
			if tt.expectService {
				mockService.AssertCalled(t, "BookClass", mock.Anything, mock.Anything, mock.Anything)
			} else {
				mockService.AssertNotCalled(t, "BookClass")
			}
		})
	}
}
