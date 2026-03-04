package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-backend/errors"
	"go-backend/middleware"
	"go-backend/model"
	"go-backend/validation"
)

// TaskService is the interface for task operations (allows mocking in tests).
type TaskService interface {
	FindTaskByID(id int) (*model.Task, error)
	FindUserTasksByStatus(status, userID string) ([]model.Task, error)
	CreateTask(title, status string, userID int) (*model.Task, error)
	UpdateTask(id int, title, status *string, userID *int) (*model.Task, error)
}

// TaskHandler handles task-related HTTP endpoints.
type TaskHandler struct {
	service TaskService
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (c *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.HandleList(w, r)
	case http.MethodPost:
		c.HandleCreate(w, r)
	default:
		middleware.SendResponse(w, r, middleware.Error(errors.NewInvalidMethodError(r.Method), r))
	}
}

func (c *TaskHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	userID := r.URL.Query().Get("userId")
	tasks, err := c.service.FindUserTasksByStatus(status, userID)
	if err != nil {
		middleware.SendResponse(w, r, middleware.Error(err, r))
		return
	}
	middleware.SendResponse(w, r, middleware.Data(http.StatusOK, model.TasksResponse{Tasks: tasks, Count: len(tasks)}))
}

func (c *TaskHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendResponse(w, r, middleware.Error(errors.NewInvalidJSONError(err), r))
		return
	}
	task, err := c.service.CreateTask(req.Title, req.Status, req.UserID)
	if err != nil {
		if errors.IsErrorCode(err, errors.ErrCodeNotFound) {
			middleware.SendResponse(w, r, middleware.Error(errors.NewValidationError("User not found", err), r))
			return
		}
		if errors.IsErrorCode(err, errors.ErrCodeValidation) {
			middleware.SendResponse(w, r, middleware.Error(err, r))
			return
		}
		middleware.SendResponse(w, r, middleware.Error(errors.NewDataStoreError("Failed to create task", err), r))
		return
	}
	middleware.SendResponse(w, r, middleware.Data(http.StatusCreated, task))
}

func (c *TaskHandler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.HandleGetByID(w, r)
	case http.MethodPut, http.MethodPatch:
		c.HandleUpdate(w, r)
	default:
		middleware.SendResponse(w, r, middleware.Error(errors.NewInvalidMethodError(r.Method), r))
	}
}

func (c *TaskHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	id, err := validation.ValidateID(idStr)
	if err != nil {
		middleware.SendResponse(w, r, middleware.Error(err, r))
		return
	}
	task, err := c.service.FindTaskByID(id)
	if err != nil {
		middleware.SendResponse(w, r, middleware.Error(err, r))
		return
	}
	if task == nil {
		middleware.SendResponse(w, r, middleware.Error(errors.NewNotFoundError("Task", nil), r))
		return
	}
	middleware.SendResponse(w, r, middleware.Data(http.StatusOK, task))
}

func (c *TaskHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := validation.ValidateID(r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:])
	if err != nil {
		middleware.SendResponse(w, r, middleware.Error(err, r))
		return
	}
	var req model.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendResponse(w, r, middleware.Error(errors.NewInvalidJSONError(err), r))
		return
	}
	task, err := c.service.UpdateTask(id, req.Title, req.Status, req.UserID)
	if err != nil {
		if errors.IsErrorCode(err, errors.ErrCodeNotFound) {
			middleware.SendResponse(w, r, middleware.Error(err, r))
			return
		}
		if errors.IsErrorCode(err, errors.ErrCodeValidation) {
			middleware.SendResponse(w, r, middleware.Error(err, r))
			return
		}
		middleware.SendResponse(w, r, middleware.Error(errors.NewDataStoreError("Failed to update task", err), r))
		return
	}
	middleware.SendResponse(w, r, middleware.Data(http.StatusOK, task))
}
