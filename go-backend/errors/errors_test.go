package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	// Test error without internal error
	err := NewValidationError("test message", nil)
	expected := "VALIDATION_ERROR: test message"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}

	// Test error with internal error
	internalErr := errors.New("internal error")
	err = NewDataStoreError("data store failed", internalErr)
	expected = "DATA_STORE_ERROR: data store failed (internal: internal error)"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

func TestAppError_Unwrap(t *testing.T) {
	internalErr := errors.New("internal error")
	err := NewDataStoreError("data store failed", internalErr)

	unwrapped := err.Unwrap()
	if !errors.Is(internalErr, unwrapped) {
		t.Errorf("Expected unwrapped error %v, got %v", internalErr, unwrapped)
	}
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name           string
		constructor    func(...interface{}) *AppError
		expectedCode   ErrorCode
		expectedStatus int
	}{
		{
			name:           "ValidationError",
			constructor:    func(args ...interface{}) *AppError { return NewValidationError("test", nil) },
			expectedCode:   ErrCodeValidation,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "NotFoundError",
			constructor:    func(args ...interface{}) *AppError { return NewNotFoundError("user", nil) },
			expectedCode:   ErrCodeNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "DuplicateError",
			constructor:    func(args ...interface{}) *AppError { return NewDuplicateError("user", "email", nil) },
			expectedCode:   ErrCodeDuplicate,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InternalError",
			constructor:    func(args ...interface{}) *AppError { return NewInternalError("test", errors.New("internal")) },
			expectedCode:   ErrCodeInternal,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor()
			if err.Code != tt.expectedCode {
				t.Errorf("Expected error code %s, got %s", tt.expectedCode, err.Code)
			}
			if err.HTTPStatus != tt.expectedStatus {
				t.Errorf("Expected HTTP status %d, got %d", tt.expectedStatus, err.HTTPStatus)
			}
		})
	}
}

func TestIsErrorCode(t *testing.T) {
	err := NewValidationError("test", nil)

	if !IsErrorCode(err, ErrCodeValidation) {
		t.Error("Expected error to be validation error")
	}

	if IsErrorCode(err, ErrCodeNotFound) {
		t.Error("Expected error to not be not found error")
	}
}

func TestGetHTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "AppError",
			err:      NewValidationError("test", nil),
			expected: http.StatusBadRequest,
		},
		{
			name:     "GenericError",
			err:      errors.New("generic error"),
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := GetHTTPStatus(tt.err)
			if status != tt.expected {
				t.Errorf("Expected HTTP status %d, got %d", tt.expected, status)
			}
		})
	}
}

func TestRequestTooLargeError(t *testing.T) {
	err := NewRequestTooLargeError(1024, 2048)

	if err.Code != ErrCodeRequestTooLarge {
		t.Errorf("Expected error code %s, got %s", ErrCodeRequestTooLarge, err.Code)
	}

	if err.HTTPStatus != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected HTTP status %d, got %d", http.StatusRequestEntityTooLarge, err.HTTPStatus)
	}

	expectedMessage := "Request size 2048 exceeds maximum allowed size 1024"
	if err.Message != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, err.Message)
	}
}

func TestInvalidMethodError(t *testing.T) {
	err := NewInvalidMethodError("POST")

	if err.Code != ErrCodeInvalidMethod {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidMethod, err.Code)
	}

	if err.HTTPStatus != http.StatusMethodNotAllowed {
		t.Errorf("Expected HTTP status %d, got %d", http.StatusMethodNotAllowed, err.HTTPStatus)
	}

	expectedMessage := "Method POST not allowed"
	if err.Message != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, err.Message)
	}
}

func TestNewInvalidJSONError(t *testing.T) {
	err := NewInvalidJSONError(nil)
	if err.Code != ErrCodeInvalidJSON {
		t.Errorf("Expected ErrCodeInvalidJSON, got %s", err.Code)
	}
	if err.Internal != nil {
		t.Error("Expected nil internal for nil arg")
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("task", nil)
	if err.Code != ErrCodeNotFound {
		t.Errorf("Expected ErrCodeNotFound, got %s", err.Code)
	}
	if err.Message != "task not found" {
		t.Errorf("Expected message, got %q", err.Message)
	}
}

func TestIsErrorCode_NonAppError(t *testing.T) {
	if IsErrorCode(errors.New("other"), ErrCodeValidation) {
		t.Error("IsErrorCode should return false for non-AppError")
	}
}
