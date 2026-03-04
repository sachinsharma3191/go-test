package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-backend/errors"
	"go-backend/model"
)

// mockTaskServiceNilTask returns (nil, nil) for FindTaskByID to cover task==nil branch.
type mockTaskServiceNilTask struct{}

func (m *mockTaskServiceNilTask) FindTaskByID(id int) (*model.Task, error) {
	return nil, nil
}

func (m *mockTaskServiceNilTask) FindUserTasksByStatus(_, _ string) ([]model.Task, error) {
	return nil, nil
}

func (m *mockTaskServiceNilTask) CreateTask(_, _ string, _ int) (*model.Task, error) {
	return nil, errors.NewNotFoundError("user", nil)
}

func (m *mockTaskServiceNilTask) UpdateTask(_ int, _, _ *string, _ *int) (*model.Task, error) {
	return nil, errors.NewNotFoundError("task", nil)
}

func TestTaskHandler_HandleTasks_List(t *testing.T) {
	_, taskHandler, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	taskHandler.HandleTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.TasksResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if len(response.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(response.Tasks))
	}
	if response.Count != 2 {
		t.Errorf("Expected count 2, got %d", response.Count)
	}
}

func TestTaskHandler_HandleTasks_Create(t *testing.T) {
	_, taskHandler, _ := setupTestHandlers(t)

	taskReq := model.CreateTaskRequest{
		Title:  "New Task",
		Status: "pending",
		UserID: 1,
	}

	body, _ := json.Marshal(taskReq)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	taskHandler.HandleTasks(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response model.Task
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response.Title != "New Task" {
		t.Errorf("Expected title 'New Task', got '%s'", response.Title)
	}
}

func TestTaskHandler_HandleTasks_CreateValidation(t *testing.T) {
	_, taskHandler, _ := setupTestHandlers(t)

	testCases := []struct {
		name    string
		request model.CreateTaskRequest
		status  int
	}{
		{"empty title", model.CreateTaskRequest{Title: "", Status: "pending", UserID: 1}, http.StatusBadRequest},
		{"invalid status", model.CreateTaskRequest{Title: "Test", Status: "invalid-status", UserID: 1}, http.StatusBadRequest},
		{"non-existent user", model.CreateTaskRequest{Title: "Test", Status: "pending", UserID: 999}, http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.request)
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			taskHandler.HandleTasks(w, req)

			if w.Code != tc.status {
				t.Errorf("Expected status %d, got %d", tc.status, w.Code)
			}
		})
	}
}

func TestTaskHandler_HandleTasks_InvalidMethod(t *testing.T) {
	_, taskHandler, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodPut, "/tasks", nil)
	w := httptest.NewRecorder()

	taskHandler.HandleTasks(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestTaskHandler_HandleTaskByID_Get(t *testing.T) {
	_, taskHandler, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	w := httptest.NewRecorder()

	taskHandler.HandleTaskByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.Task
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response.ID != 1 {
		t.Errorf("Expected task ID 1, got %d", response.ID)
	}
}

func TestTaskHandler_HandleTaskByID_GetNotFound(t *testing.T) {
	_, taskHandler, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	w := httptest.NewRecorder()

	taskHandler.HandleTaskByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestTaskHandler_HandleGetByID_TaskNil(t *testing.T) {
	// Mock that returns (nil, nil) to cover defensive task==nil branch
	h := NewTaskHandler(&mockTaskServiceNilTask{})
	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	w := httptest.NewRecorder()
	h.HandleTaskByID(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404 when task is nil, got %d", w.Code)
	}
}

func TestTaskHandler_HandleTaskByID_Update(t *testing.T) {
	_, taskHandler, _ := setupTestHandlers(t)

	title, status := "Updated Task", "completed"
	updateReq := model.UpdateTaskRequest{Title: &title, Status: &status}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	taskHandler.HandleTaskByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.Task
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response.Title != "Updated Task" {
		t.Errorf("Expected title 'Updated Task', got '%s'", response.Title)
	}
}

func TestTaskHandler_HandleTaskByID_UpdateNotFound(t *testing.T) {
	_, taskHandler, _ := setupTestHandlers(t)

	title, status := "Updated", "completed"
	updateReq := model.UpdateTaskRequest{Title: &title, Status: &status}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/tasks/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	taskHandler.HandleTaskByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestTaskHandler_HandleUpdate_ErrCodeValidation(t *testing.T) {
	_, taskHandler, _ := setupTestHandlers(t)

	status := "invalid-status"
	updateReq := model.UpdateTaskRequest{Status: &status}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	taskHandler.HandleTaskByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid status, got %d", w.Code)
	}
}

func TestTaskHandler_HandleTaskByID_DeleteReturns405(t *testing.T) {
	_, taskHandler, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	w := httptest.NewRecorder()

	taskHandler.HandleTaskByID(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d (DELETE not implemented), got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestTaskHandler_HandleTaskByID_InvalidMethod(t *testing.T) {
	_, taskHandler, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodPost, "/tasks/1", nil)
	w := httptest.NewRecorder()

	taskHandler.HandleTaskByID(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestTaskHandler_Integration(t *testing.T) {
	userHandler, taskHandler, _ := setupTestHandlers(t)

	// Create user
	userReq := model.CreateUserRequest{
		Name:  "Integration User",
		Email: "integration@example.com",
		Role:  "developer",
	}
	body, _ := json.Marshal(userReq)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	userHandler.Users(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Create user: expected 201, got %d", w.Code)
	}

	var createdUser model.User
	if err := json.Unmarshal(w.Body.Bytes(), &createdUser); err != nil {
		t.Fatalf("Parse user: %v", err)
	}

	// Create task
	taskReq := model.CreateTaskRequest{
		Title:  "Integration Task",
		Status: "pending",
		UserID: createdUser.ID,
	}
	body, _ = json.Marshal(taskReq)
	req = httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	taskHandler.HandleTasks(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Create task: expected 201, got %d", w.Code)
	}

	var createdTask model.Task
	if err := json.Unmarshal(w.Body.Bytes(), &createdTask); err != nil {
		t.Fatalf("Parse task: %v", err)
	}

	// Update task
	title, status := "Updated Integration Task", "completed"
	updateReq := model.UpdateTaskRequest{Title: &title, Status: &status}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", createdTask.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	taskHandler.HandleTaskByID(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Update task: expected 200, got %d", w.Code)
	}

	var updatedTask model.Task
	if err := json.Unmarshal(w.Body.Bytes(), &updatedTask); err != nil {
		t.Fatalf("Parse updated task: %v", err)
	}
	if updatedTask.Title != "Updated Integration Task" {
		t.Errorf("Expected updated title, got %q", updatedTask.Title)
	}

	// DELETE returns 405
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%d", createdTask.ID), nil)
	w = httptest.NewRecorder()
	taskHandler.HandleTaskByID(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("DELETE: expected 405, got %d", w.Code)
	}
}

func TestTaskHandler_HandleList_StoreError(t *testing.T) {
	_, taskHandler := setupTestHandlersWithFailingStore(t)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	taskHandler.HandleTasks(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d on store error, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTaskHandler_HandleCreate_DataStoreError(t *testing.T) {
	_, taskHandler := setupTestHandlersWithSaveFailingStore(t)

	taskReq := model.CreateTaskRequest{Title: "Test", Status: "pending", UserID: 1}
	body, _ := json.Marshal(taskReq)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	taskHandler.HandleTasks(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d on DataStoreError, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTaskHandler_HandleGetByID_StoreError(t *testing.T) {
	_, taskHandler := setupTestHandlersWithFailingStore(t)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	w := httptest.NewRecorder()

	taskHandler.HandleTaskByID(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d on store error, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTaskHandler_HandleUpdate_DataStoreError(t *testing.T) {
	_, taskHandler := setupTestHandlersWithSaveFailingStore(t)

	title, status := "Updated", "completed"
	updateReq := model.UpdateTaskRequest{Title: &title, Status: &status}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	taskHandler.HandleTaskByID(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d on DataStoreError, got %d", http.StatusInternalServerError, w.Code)
	}
}
