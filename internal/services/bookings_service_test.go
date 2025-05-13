package services

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"glofox/internal/constants"
	"glofox/internal/models"
	"glofox/internal/utils"
	"testing"
	"time"
)

func (m *MockBookingRepo) Create(className, memberName string, date time.Time) error {
	args := m.Called(className, memberName, date)
	return args.Error(0)
}

func TestClassService_BookClass(t *testing.T) {
	// Setup mocks
	mockClassRepo := new(MockClassRepo)
	mockBookingRepo := new(MockBookingRepo)
	service := NewClassService(mockClassRepo, mockBookingRepo)

	// Define test cases
	tests := []struct {
		name            string
		className       string
		memberName      string
		dateStr         string
		setupMock       func()
		expectedErr     error
		expectedBooking *struct {
			className  string
			memberName string
			date       time.Time
		}
	}{
		{
			name:       "Happy Path",
			className:  "Yoga",
			memberName: "Alice",
			dateStr:    "2025-06-10",
			setupMock: func() {
				startDate, _ := time.Parse(constants.DateFormat, "2025-06-01")
				endDate, _ := time.Parse(constants.DateFormat, "2025-06-20")
				date, _ := time.Parse(constants.DateFormat, "2025-06-10")
				mockClassRepo.On("GetByName", "Yoga").Return(models.Class{
					Name:      "Yoga",
					StartDate: startDate,
					EndDate:   endDate,
					Capacity:  10,
				}, true)
				mockBookingRepo.On("Create", "Yoga", "Alice", utils.ToMidnightUTC(date)).Return(nil)
			},
			expectedErr: nil,
			expectedBooking: &struct {
				className  string
				memberName string
				date       time.Time
			}{
				className:  "Yoga",
				memberName: "Alice",
				date:       time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:            "Invalid Date Format",
			className:       "Yoga",
			memberName:      "Alice",
			dateStr:         "2025/06/10",
			setupMock:       func() {},
			expectedErr:     constants.ErrInvalidDate,
			expectedBooking: nil,
		},
		{
			name:       "Class Not Found",
			className:  "Yoga",
			memberName: "Alice",
			dateStr:    "2025-06-10",
			setupMock: func() {
				mockClassRepo.On("GetByName", "Yoga").Return(models.Class{}, false)
			},
			expectedErr:     constants.ErrClassNotFound,
			expectedBooking: nil,
		},
		{
			name:       "Date Before Start",
			className:  "Yoga",
			memberName: "Alice",
			dateStr:    "2025-05-31",
			setupMock: func() {
				startDate, _ := time.Parse(constants.DateFormat, "2025-06-01")
				endDate, _ := time.Parse(constants.DateFormat, "2025-06-20")
				mockClassRepo.On("GetByName", "Yoga").Return(models.Class{
					Name:      "Yoga",
					StartDate: startDate,
					EndDate:   endDate,
					Capacity:  10,
				}, true)
			},
			expectedErr:     fmt.Errorf("date 2025-05-31 is not valid for class Yoga"),
			expectedBooking: nil,
		},
		{
			name:       "Date After End",
			className:  "Yoga",
			memberName: "Alice",
			dateStr:    "2025-06-21",
			setupMock: func() {
				startDate, _ := time.Parse(constants.DateFormat, "2025-06-01")
				endDate, _ := time.Parse(constants.DateFormat, "2025-06-20")
				mockClassRepo.On("GetByName", "Yoga").Return(models.Class{
					Name:      "Yoga",
					StartDate: startDate,
					EndDate:   endDate,
					Capacity:  10,
				}, true)
			},
			expectedErr:     fmt.Errorf("date 2025-06-21 is not valid for class Yoga"),
			expectedBooking: nil,
		},
		{
			name:       "Date Equals Start",
			className:  "Yoga",
			memberName: "Alice",
			dateStr:    "2025-06-01",
			setupMock: func() {
				startDate, _ := time.Parse(constants.DateFormat, "2025-06-01")
				endDate, _ := time.Parse(constants.DateFormat, "2025-06-20")
				date, _ := time.Parse(constants.DateFormat, "2025-06-01")
				mockClassRepo.On("GetByName", "Yoga").Return(models.Class{
					Name:      "Yoga",
					StartDate: startDate,
					EndDate:   endDate,
					Capacity:  10,
				}, true)
				mockBookingRepo.On("Create", "Yoga", "Alice", utils.ToMidnightUTC(date)).Return(nil)
			},
			expectedErr: nil,
			expectedBooking: &struct {
				className  string
				memberName string
				date       time.Time
			}{
				className:  "Yoga",
				memberName: "Alice",
				date:       time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:       "Date Equals End",
			className:  "Yoga",
			memberName: "Alice",
			dateStr:    "2025-06-20",
			setupMock: func() {
				startDate, _ := time.Parse(constants.DateFormat, "2025-06-01")
				endDate, _ := time.Parse(constants.DateFormat, "2025-06-20")
				date, _ := time.Parse(constants.DateFormat, "2025-06-20")
				mockClassRepo.On("GetByName", "Yoga").Return(models.Class{
					Name:      "Yoga",
					StartDate: startDate,
					EndDate:   endDate,
					Capacity:  10,
				}, true)
				mockBookingRepo.On("Create", "Yoga", "Alice", utils.ToMidnightUTC(date)).Return(nil)
			},
			expectedErr: nil,
			expectedBooking: &struct {
				className  string
				memberName string
				date       time.Time
			}{
				className:  "Yoga",
				memberName: "Alice",
				date:       time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockClassRepo.Calls = nil
			mockClassRepo.ExpectedCalls = nil
			mockBookingRepo.Calls = nil
			mockBookingRepo.ExpectedCalls = nil

			// Setup mock
			tt.setupMock()

			// Call BookClass
			err := service.BookClass(tt.className, tt.memberName, tt.dateStr)

			// Assert error
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "Expected error %v, got %v", tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			// Assert mock calls
			if tt.expectedBooking != nil {
				mockClassRepo.AssertCalled(t, "GetByName", tt.className)
				mockBookingRepo.AssertCalled(t, "Create", tt.expectedBooking.className, tt.expectedBooking.memberName, tt.expectedBooking.date)
			} else {
				if errors.Is(err, constants.ErrInvalidDate) || errors.Is(err, constants.ErrClassNotFound) {
					mockBookingRepo.AssertNotCalled(t, "Create")
				}
			}
		})
	}
}
