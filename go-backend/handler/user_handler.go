package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-backend/errors"
	"go-backend/middleware"
	"go-backend/model"
	"go-backend/service"
	"go-backend/validation"
)

// UserHandler handles user-related HTTP endpoints.
type UserHandler struct {
	userService *service.UserService
	taskService *service.TaskService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *service.UserService, taskService *service.TaskService) *UserHandler {
	return &UserHandler{
		userService: userService,
		taskService: taskService,
	}
}

// Users handles /users for listing and creating users.
func (h *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		middleware.SendResponse(w, r, middleware.Error(errors.NewInvalidMethodError(r.Method), r))
	}
}

func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListUsers()
	if err != nil {
		middleware.SendResponse(w, r, middleware.Error(errors.NewDataStoreError("failed to list users", err), r))
		return
	}
	middleware.SendResponse(w, r, middleware.Data(http.StatusOK, model.UsersResponse{
		Users: users,
		Count: len(users),
	}))
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendResponse(w, r, middleware.Error(errors.NewInvalidJSONError(err), r))
		return
	}

	user, err := h.userService.CreateUser(req.Name, req.Email, req.Role)
	if err != nil {
		if errors.IsErrorCode(err, errors.ErrCodeDuplicate) || errors.IsErrorCode(err, errors.ErrCodeValidation) {
			middleware.SendResponse(w, r, middleware.Error(err, r))
			return
		}
		middleware.SendResponse(w, r, middleware.Error(errors.NewDataStoreError("failed to create user", err), r))
		return
	}

	middleware.SendResponse(w, r, middleware.Data(http.StatusCreated, user))
}

// Stats handles /stats (or /users/stats depending on routing) and returns platform stats.
func (h *UserHandler) Stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendResponse(w, r, middleware.Error(errors.NewInvalidMethodError(r.Method), r))
		return
	}

	users, err := h.userService.ListUsers()
	if err != nil {
		middleware.SendResponse(w, r, middleware.Error(errors.NewDataStoreError("failed to list users", err), r))
		return
	}

	tasks, err := h.taskService.FindAllTasks()
	if err != nil {
		middleware.SendResponse(w, r, middleware.Error(errors.NewDataStoreError("failed to list tasks", err), r))
		return
	}

	var stats model.StatsResponse
	stats.Users.Total = len(users)
	stats.Tasks.Total = len(tasks)

	for _, task := range tasks {
		switch task.Status {
		case "pending":
			stats.Tasks.Pending++
		case "in-progress":
			stats.Tasks.InProgress++
		case "completed":
			stats.Tasks.Completed++
		}
	}

	middleware.SendResponse(w, r, middleware.Data(http.StatusOK, stats))
}

// UserByID handles /users/{id}
func (h *UserHandler) UserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendResponse(w, r, middleware.Error(errors.NewInvalidMethodError(r.Method), r))
		return
	}

	id, err := validation.ValidateID(r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:])
	if err != nil {
		middleware.SendResponse(w, r, middleware.Error(err, r))
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if errors.IsErrorCode(err, errors.ErrCodeNotFound) {
			middleware.SendResponse(w, r, middleware.Error(err, r))
			return
		}
		middleware.SendResponse(w, r, middleware.Error(errors.NewDataStoreError("failed to get user", err), r))
		return
	}

	middleware.SendResponse(w, r, middleware.Data(http.StatusOK, user))
}
