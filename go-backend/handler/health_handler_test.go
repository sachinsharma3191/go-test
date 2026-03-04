package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_Health(t *testing.T) {
	_, _, healthHandler := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%v'", response["status"])
	}
	if response["version"] != "test" {
		t.Errorf("Expected version 'test', got '%v'", response["version"])
	}
}

func TestHealthHandler_InvalidMethod(t *testing.T) {
	_, _, healthHandler := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler.Health(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestHealthHandler_Ready(t *testing.T) {
	_, _, healthHandler := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()

	healthHandler.Ready(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHealthHandler_Ready_NotReady_503(t *testing.T) {
	healthHandler := setupTestHandlersWithUnhealthyMonitor(t)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()

	healthHandler.Ready(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d when not ready, got %d", http.StatusServiceUnavailable, w.Code)
	}
}

func TestHealthHandler_Ready_InvalidMethod(t *testing.T) {
	_, _, healthHandler := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodPost, "/health/ready", nil)
	w := httptest.NewRecorder()

	healthHandler.Ready(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestHealthHandler_Live(t *testing.T) {
	_, _, healthHandler := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()

	healthHandler.Live(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHealthHandler_Live_InvalidMethod(t *testing.T) {
	_, _, healthHandler := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodPost, "/health/live", nil)
	w := httptest.NewRecorder()

	healthHandler.Live(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
