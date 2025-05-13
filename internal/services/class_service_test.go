package services

import (
	"glofox/internal/constants"
	"glofox/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClassRepo mocks the ClassRepo
type MockClassRepo struct {
	mock.Mock
}

func (m *MockClassRepo) Create(class models.Class) error {
	args := m.Called(class)
	return args.Error(0)
}

func (m *MockClassRepo) GetByName(name string) (models.Class, bool) {
	args := m.Called(name)
	class, _ := args.Get(0).(models.Class)
	exists, _ := args.Get(1).(bool)
	return class, exists
}

// MockBookingRepo is just a placeholder, not used in CreateClass
type MockBookingRepo struct {
	mock.Mock
}

func TestClassService_CreateClass(t *testing.T) {
	// Setup mocks
	mockClassRepo := new(MockClassRepo)
	mockBookingRepo := new(MockBookingRepo)
	service := NewClassService(mockClassRepo, mockBookingRepo)

	// Define test cases
	tests := []struct {
		name          string
		inputName     string
		startDateStr  string
		endDateStr    string
		capacity      int
		setupMock     func()
		expectedErr   error
		expectedClass *models.Class // Nil if no Create call expected
	}{
		{
			name:         "Happy Path",
			inputName:    "Yoga",
			startDateStr: "2025-06-01",
			endDateStr:   "2025-06-20",
			capacity:     10,
			setupMock: func() {
				startDate, _ := time.Parse(constants.DateFormat, "2025-06-01")
				endDate, _ := time.Parse(constants.DateFormat, "2025-06-20")
				mockClassRepo.On("Create", models.Class{
					Name:      "Yoga",
					StartDate: startDate,
					EndDate:   endDate,
					Capacity:  10,
				}).Return(nil)
			},
			expectedErr: nil,
			expectedClass: &models.Class{
				Name:      "Yoga",
				StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC),
				Capacity:  10,
			},
		},
		{
			name:          "Invalid Start Date Format",
			inputName:     "Yoga",
			startDateStr:  "2025-06-01T00:00:00Z",
			endDateStr:    "2025-06-20",
			capacity:      10,
			setupMock:     func() {},
			expectedErr:   constants.ErrInvalidStartDate,
			expectedClass: nil,
		},
		{
			name:          "Invalid End Date Format",
			inputName:     "Yoga",
			startDateStr:  "2025-06-01",
			endDateStr:    "2025/06/20",
			capacity:      10,
			setupMock:     func() {},
			expectedErr:   constants.ErrInvalidEndDate,
			expectedClass: nil,
		},
		{
			name:          "Start Date After End Date",
			inputName:     "Yoga",
			startDateStr:  "2025-06-21",
			endDateStr:    "2025-06-01",
			capacity:      10,
			setupMock:     func() {},
			expectedErr:   constants.ErrInvalidStartEndDate,
			expectedClass: nil,
		},
		{
			name:         "Same Start and End Date",
			inputName:    "Yoga",
			startDateStr: "2025-06-01",
			endDateStr:   "2025-06-01",
			capacity:     10,
			setupMock: func() {
				startDate, _ := time.Parse(constants.DateFormat, "2025-06-01")
				endDate, _ := time.Parse(constants.DateFormat, "2025-06-01")
				mockClassRepo.On("Create", models.Class{
					Name:      "Yoga",
					StartDate: startDate,
					EndDate:   endDate,
					Capacity:  10,
				}).Return(nil)
			},
			expectedErr: nil,
			expectedClass: &models.Class{
				Name:      "Yoga",
				StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				Capacity:  10,
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockClassRepo.Calls = nil
			mockClassRepo.ExpectedCalls = nil

			// Setup mock
			tt.setupMock()

			// Call CreateClass
			err := service.CreateClass(tt.inputName, tt.startDateStr, tt.endDateStr, tt.capacity)

			// Assert error
			assert.Equal(t, tt.expectedErr, err, "Expected error %v, got %v", tt.expectedErr, err)

			// Assert mock calls
			if tt.expectedClass != nil {
				mockClassRepo.AssertCalled(t, "Create", mock.Anything)
				// Verify the class passed to Create
				for _, call := range mockClassRepo.Calls {
					class := call.Arguments[0].(models.Class)
					assert.Equal(t, tt.expectedClass.Name, class.Name)
					assert.True(t, tt.expectedClass.StartDate.Equal(class.StartDate), "StartDate mismatch: expected %v, got %v", tt.expectedClass.StartDate, class.StartDate)
					assert.True(t, tt.expectedClass.EndDate.Equal(class.EndDate), "EndDate mismatch: expected %v, got %v", tt.expectedClass.EndDate, class.EndDate)
					assert.Equal(t, tt.expectedClass.Capacity, class.Capacity)
				}
			} else {
				mockClassRepo.AssertNotCalled(t, "Create")
			}
		})
	}
}
