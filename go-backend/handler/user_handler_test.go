package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-backend/errors"
	"go-backend/model"
)

func TestUserHandler_Users_List(t *testing.T) {
	userHandler, _, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	userHandler.Users(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.UsersResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if len(response.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(response.Users))
	}
	if response.Count != 2 {
		t.Errorf("Expected count 2, got %d", response.Count)
	}
}

func TestUserHandler_Users_Create(t *testing.T) {
	userHandler, _, _ := setupTestHandlers(t)

	userReq := model.CreateUserRequest{
		Name:  "New User",
		Email: "newuser@example.com",
		Role:  "developer",
	}

	body, _ := json.Marshal(userReq)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	userHandler.Users(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response model.User
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response.Name != "New User" {
		t.Errorf("Expected name 'New User', got '%s'", response.Name)
	}
}

func TestUserHandler_Users_CreateValidation(t *testing.T) {
	userHandler, _, _ := setupTestHandlers(t)

	testCases := []struct {
		name    string
		request model.CreateUserRequest
		status  int
	}{
		{"empty name", model.CreateUserRequest{Name: "", Email: "test@example.com", Role: "developer"}, http.StatusBadRequest},
		{"invalid email", model.CreateUserRequest{Name: "Test", Email: "invalid-email", Role: "developer"}, http.StatusBadRequest},
		{"invalid role", model.CreateUserRequest{Name: "Test", Email: "test@example.com", Role: "invalid-role"}, http.StatusBadRequest},
		{"duplicate email", model.CreateUserRequest{Name: "Other", Email: "john@example.com", Role: "developer"}, http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.request)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			userHandler.Users(w, req)

			if w.Code != tc.status {
				t.Errorf("Expected status %d, got %d", tc.status, w.Code)
			}
		})
	}
}

func TestUserHandler_Users_CreateInvalidJSON(t *testing.T) {
	userHandler, _, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	userHandler.Users(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUserHandler_Users_InvalidMethod(t *testing.T) {
	userHandler, _, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodPut, "/users", nil)
	w := httptest.NewRecorder()

	userHandler.Users(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	var errorResponse map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Fatalf("Failed to parse error: %v", err)
	}
	if errorResponse["code"] != string(errors.ErrCodeInvalidMethod) {
		t.Errorf("Expected error code %s, got %v", errors.ErrCodeInvalidMethod, errorResponse["code"])
	}
}

func TestUserHandler_UserByID_Found(t *testing.T) {
	userHandler, _, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()

	userHandler.UserByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.User
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", response.ID)
	}
}

func TestUserHandler_UserByID_NotFound(t *testing.T) {
	userHandler, _, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	w := httptest.NewRecorder()

	userHandler.UserByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUserHandler_UserByID_InvalidID(t *testing.T) {
	userHandler, _, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
	w := httptest.NewRecorder()

	userHandler.UserByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUserHandler_UserByID_InvalidMethod(t *testing.T) {
	userHandler, _, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodPost, "/users/1", nil)
	w := httptest.NewRecorder()

	userHandler.UserByID(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestUserHandler_Stats(t *testing.T) {
	userHandler, _, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()

	userHandler.Stats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.StatsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response.Users.Total != 2 || response.Tasks.Total != 2 {
		t.Errorf("Expected 2 users and 2 tasks, got %d users, %d tasks", response.Users.Total, response.Tasks.Total)
	}
}

func TestUserHandler_ListUsers_StoreError(t *testing.T) {
	userHandler, _ := setupTestHandlersWithFailingStore(t)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	userHandler.Users(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d on store error, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUserHandler_CreateUser_DataStoreError(t *testing.T) {
	userHandler, _ := setupTestHandlersWithSaveFailingStore(t)

	userReq := model.CreateUserRequest{Name: "Test", Email: "test@example.com", Role: "developer"}
	body, _ := json.Marshal(userReq)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	userHandler.Users(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d on DataStoreError, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUserHandler_Stats_UsersStoreError(t *testing.T) {
	userHandler, _ := setupTestHandlersWithFailingStore(t)

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()

	userHandler.Stats(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d on store error, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUserHandler_Stats_TasksStoreError(t *testing.T) {
	userHandler, _ := setupTestHandlersWithSecondLoadFailing(t)

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()

	userHandler.Stats(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d when tasks load fails, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUserHandler_UserByID_DataStoreError(t *testing.T) {
	userHandler, _ := setupTestHandlersWithFailingStore(t)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()

	userHandler.UserByID(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d on DataStoreError, got %d", http.StatusInternalServerError, w.Code)
	}
}
