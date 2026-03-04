package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "go-backend/errors"
)

func TestSend_ErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	err := apperrors.NewValidationError("bad input", nil)
	SendResponse(rr, req, Error(err, req))

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["code"] != string(apperrors.ErrCodeValidation) {
		t.Errorf("expected code %s, got %v", apperrors.ErrCodeValidation, body["code"])
	}
}

func TestSend_ValidationError_Writes400(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", nil)

	SendResponse(rr, req, Validation(map[string]string{"field": "error"}, req))

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestLoggingMiddleware_CallsNextAndSetsStatus(t *testing.T) {
	called := false
	handler := LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/logging", nil)
	handler.ServeHTTP(rr, req)

	if !called {
		t.Fatalf("expected wrapped handler to be called")
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestErrorMiddleware_RecoversFromPanic(t *testing.T) {
	handler := ErrorMiddleware(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestErrorMiddleware_NoPanic(t *testing.T) {
	handler := ErrorMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestCORSMiddleware_OPTIONS(t *testing.T) {
	called := false
	handler := CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/cors", nil)
	handler.ServeHTTP(rr, req)
	if called {
		t.Error("OPTIONS should not call next handler")
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestCORSMiddleware_GET(t *testing.T) {
	called := false
	handler := CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/cors", nil)
	handler.ServeHTTP(rr, req)
	if !called {
		t.Error("GET should call next handler")
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS headers")
	}
}

func TestSendResponse_Data(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	SendResponse(rr, req, Data(http.StatusCreated, map[string]string{"id": "1"}))
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}

func TestError_NonAppError(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	SendResponse(rr, req, Error(errors.New("generic"), req))
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for non-AppError, got %d", rr.Code)
	}
}
