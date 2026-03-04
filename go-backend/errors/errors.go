package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents different types of application errors
type ErrorCode string

const (
	// Validation errors (400)
	ErrCodeValidation  ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidJSON ErrorCode = "INVALID_JSON"

	// Data store errors (404, 400, 500)
	ErrCodeNotFound  ErrorCode = "NOT_FOUND"
	ErrCodeDuplicate ErrorCode = "DUPLICATE"
	ErrCodeDataStore ErrorCode = "DATA_STORE_ERROR"

	// System errors (500, 413, 405)
	ErrCodeInternal        ErrorCode = "INTERNAL_ERROR"
	ErrCodeRequestTooLarge ErrorCode = "REQUEST_TOO_LARGE"
	ErrCodeInvalidMethod   ErrorCode = "INVALID_METHOD"
)

// HTTPStatusMapping maps error codes to HTTP status codes
var HTTPStatusMapping = map[ErrorCode]int{
	ErrCodeValidation:      http.StatusBadRequest,
	ErrCodeInvalidJSON:     http.StatusBadRequest,
	ErrCodeNotFound:        http.StatusNotFound,
	ErrCodeDuplicate:       http.StatusBadRequest,
	ErrCodeDataStore:       http.StatusInternalServerError,
	ErrCodeInternal:        http.StatusInternalServerError,
	ErrCodeRequestTooLarge: http.StatusRequestEntityTooLarge,
	ErrCodeInvalidMethod:   http.StatusMethodNotAllowed,
}

// AppError represents an application error with context
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	HTTPStatus int       `json:"-"`
	Internal   error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %s (internal: %v)", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the internal error for error unwrapping
func (e *AppError) Unwrap() error {
	return e.Internal
}

// Core Error Constructors (these are the ones we actually use)
func NewValidationError(message string, internal error) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Internal:   internal,
	}
}

func NewInvalidJSONError(internal error) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidJSON,
		Message:    "Invalid JSON format",
		HTTPStatus: http.StatusBadRequest,
		Internal:   internal,
	}
}

func NewNotFoundError(resource string, internal error) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		HTTPStatus: http.StatusNotFound,
		Internal:   internal,
	}
}

func NewDuplicateError(resource string, field string, internal error) *AppError {
	return &AppError{
		Code:       ErrCodeDuplicate,
		Message:    fmt.Sprintf("%s with this %s already exists", resource, field),
		HTTPStatus: http.StatusBadRequest,
		Internal:   internal,
	}
}

func NewDataStoreError(message string, internal error) *AppError {
	return &AppError{
		Code:       ErrCodeDataStore,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Internal:   internal,
	}
}

func NewInternalError(message string, internal error) *AppError {
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Internal:   internal,
	}
}

func NewRequestTooLargeError(maxSize int64, actualSize int64) *AppError {
	return &AppError{
		Code:       ErrCodeRequestTooLarge,
		Message:    fmt.Sprintf("Request size %d exceeds maximum allowed size %d", actualSize, maxSize),
		HTTPStatus: http.StatusRequestEntityTooLarge,
	}
}

func NewInvalidMethodError(method string) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidMethod,
		Message:    fmt.Sprintf("Method %s not allowed", method),
		HTTPStatus: http.StatusMethodNotAllowed,
	}
}

// Utility Functions
func IsErrorCode(err error, code ErrorCode) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}

func GetHTTPStatus(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}
