package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-backend/errors"
	"go-backend/model"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	startTime  time.Time
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		startTime:      time.Now(),
	}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := NewResponseWriter(w)
		next(rw, r)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	}
}

func ErrorMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC %s %s - %v", r.Method, r.URL.Path, err)
				SendResponse(w, r, Error(errors.NewInternalError("Internal server error", fmt.Errorf("panic: %v", err)), r))
			}
		}()
		next(w, r)
	}
}

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// Response describes an HTTP response (success or failure).
type Response struct {
	StatusCode int
	Body       interface{}
}

// SendResponse writes a single response (success or failure) to w. Use Data, Error, or Validation to build responses.
func SendResponse(w http.ResponseWriter, r *http.Request, res Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	json.NewEncoder(w).Encode(res.Body)
}

// Data returns a success response with the given status code and payload.
func Data(statusCode int, body interface{}) Response {
	return Response{StatusCode: statusCode, Body: body}
}

// Error returns an error response from err, logs it, and includes request context in details.
func Error(err error, r *http.Request) Response {
	statusCode := http.StatusInternalServerError
	errorCode := "INTERNAL_ERROR"
	message := "Internal server error"
	if appErr, ok := err.(*errors.AppError); ok {
		statusCode = appErr.HTTPStatus
		errorCode = string(appErr.Code)
		message = appErr.Message
	}
	log.Printf("ERROR %s %s %d - code=%s message='%s'", r.Method, r.URL.Path, statusCode, errorCode, message)
	return Response{
		StatusCode: statusCode,
		Body: model.ErrorResponse{
			Error:   message,
			Code:    errorCode,
			Details: fmt.Sprintf("Method: %s, Path: %s", r.Method, r.URL.Path),
		},
	}
}

// Validation returns a 400 validation error response with field-level details.
func Validation(fields map[string]string, r *http.Request) Response {
	log.Printf("VALIDATION_ERROR %s %s 400 - fields=%v", r.Method, r.URL.Path, fields)
	return Response{
		StatusCode: http.StatusBadRequest,
		Body: model.ValidationError{
			Error:  "Validation failed",
			Fields: fields,
			Code:   "VALIDATION_ERROR",
		},
	}
}
