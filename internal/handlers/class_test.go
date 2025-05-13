package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"glofox/internal/constants"
	"glofox/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

// CreateClass mocks the CreateClass method
func (m *MockClassService) CreateClass(name, startDateStr, endDateStr string, capacity int) error {
	args := m.Called(name, startDateStr, endDateStr, capacity)
	return args.Error(0)
}

func TestClassHandler_CreateClass(t *testing.T) {
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
			jsonInput: `{"name":"Yoga","start_date":"2025-06-01","end_date":"2025-06-20","capacity":10}`,
			setupMock: func(m *MockClassService) {
				m.On("CreateClass", "Yoga", "2025-06-01", "2025-06-20", 10).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: models.Response{
				Status:  constants.SuccessMsg,
				Message: "Class Yoga created successfully",
			},
			expectService: true,
		},
		{
			name:      "Invalid Start Date Format",
			jsonInput: `{"name":"Yoga","start_date":"2025-06-01T00:00:00Z","end_date":"2025-06-20","capacity":10}`,
			setupMock: func(m *MockClassService) {
				m.On("CreateClass", "Yoga", "2025-06-01T00:00:00Z", "2025-06-20", 10).Return(constants.ErrInvalidStartDate)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: models.Response{
				Status:  "error",
				Message: constants.ErrInvalidStartDate.Error(),
			},
			expectService: true,
		},
		{
			name:      "Invalid End Date Format",
			jsonInput: `{"name":"Yoga","start_date":"2025-06-01","end_date":"2025/06/20","capacity":10}`,
			setupMock: func(m *MockClassService) {
				m.On("CreateClass", "Yoga", "2025-06-01", "2025/06/20", 10).Return(constants.ErrInvalidEndDate)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: models.Response{
				Status:  "error",
				Message: constants.ErrInvalidEndDate.Error(),
			},
			expectService: true,
		},
		{
			name:      "Start Date After End Date",
			jsonInput: `{"name":"Yoga","start_date":"2025-06-21","end_date":"2025-06-01","capacity":10}`,
			setupMock: func(m *MockClassService) {
				m.On("CreateClass", "Yoga", "2025-06-21", "2025-06-01", 10).Return(constants.ErrInvalidStartEndDate)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: models.Response{
				Status:  "error",
				Message: constants.ErrInvalidStartEndDate.Error(),
			},
			expectService: true,
		},
		{
			name:      "Class Already Exists",
			jsonInput: `{"name":"Yoga","start_date":"2025-06-01","end_date":"2025-06-20","capacity":10}`,
			setupMock: func(m *MockClassService) {
				m.On("CreateClass", "Yoga", "2025-06-01", "2025-06-20", 10).Return(constants.ErrClassAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedBody: models.Response{
				Status:  "error",
				Message: constants.ErrClassAlreadyExists.Error(),
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
			req, _ := http.NewRequest("POST", "/classes", bytes.NewBufferString(tt.jsonInput))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req

			// Call handler
			handler.CreateClass(ctx)

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
				mockService.AssertCalled(t, "CreateClass", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			} else {
				mockService.AssertNotCalled(t, "CreateClass")
			}
		})
	}
}
