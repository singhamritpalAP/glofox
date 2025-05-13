package constants

import "errors"

var (
	ErrInternalServer      = errors.New("internal server error")
	ErrInvalidDate         = errors.New("invalid date format, expected YYYY-MM-DD")
	ErrInvalidStartEndDate = errors.New("start date cannot be after end date")
	ErrInvalidStartDate    = errors.New("invalid start date format, expected YYYY-MM-DD")
	ErrInvalidEndDate      = errors.New("invalid end date format, expected YYYY-MM-DD")
	ErrClassNotFound       = errors.New("class not found")
	ErrClassAlreadyExists  = errors.New("class already exists")
)
