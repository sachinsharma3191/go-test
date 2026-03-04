package service

import (
	"fmt"
	"strconv"

	"go-backend/errors"
	"go-backend/model"
	"go-backend/repository"
	"go-backend/validation"
)

// TaskService contains business logic for task-related operations.
type TaskService struct {
	taskRepository *repository.TaskRepository
	userRepository *repository.UserRepository
}

// NewTaskService creates a TaskService with internal repositories.
func NewTaskService(taskRepo *repository.TaskRepository, userRepo *repository.UserRepository) *TaskService {
	return &TaskService{
		taskRepository: taskRepo,
		userRepository: userRepo,
	}
}

// FindUserTasks returns all tasks for a user (no status filtering here).
func (s *TaskService) FindUserTasks(userID int) ([]model.Task, error) {
	allTasks, err := s.FindAllTasks()
	if err != nil {
		return nil, err
	}

	var tasks []model.Task
	for _, t := range allTasks {
		if t.UserID == userID {
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

func (s *TaskService) FindAllTasks() ([]model.Task, error) {
	return s.taskRepository.FindAll()
}

// FindUserTasksByStatus filters tasks by status and optionally by userID.
// When userID is empty, returns all tasks. When status is empty, returns all matching user's tasks.
func (s *TaskService) FindUserTasksByStatus(status, userID string) ([]model.Task, error) {
	var tasks []model.Task
	var err error
	if userID == "" {
		tasks, err = s.FindAllTasks()
	} else {
		id, parseErr := strconv.Atoi(userID)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid userID %q: %w", userID, parseErr)
		}
		tasks, err = s.FindUserTasks(id)
	}
	if err != nil {
		return nil, err
	}
	if status == "" {
		return tasks, nil
	}
	var filtered []model.Task
	for _, t := range tasks {
		if t.Status == status {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}

func (s *TaskService) CreateTask(title, status string, userID int) (*model.Task, error) {
	// Syntactic validation (required, status enum, userId present)
	if err := validation.ValidateTask(title, status, userID); err != nil {
		return nil, err
	}

	// Business rule: user must exist
	u, _ := s.userRepository.FindByID(userID)
	if u == nil {
		return nil, errors.NewNotFoundError("user", nil)
	}

	t := &model.Task{Title: title, Status: status, UserID: userID}
	return s.taskRepository.Create(t)
}

func (s *TaskService) FindTaskByID(id int) (*model.Task, error) {
	t, err := s.taskRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.NewNotFoundError("Task", nil)
	}
	return t, nil
}

func (s *TaskService) UpdateTask(id int, title, status *string, userID *int) (*model.Task, error) {
	// Validation for optional update fields
	if err := validation.ValidateTaskUpdate(title, status, userID); err != nil {
		return nil, err
	}

	// If userID is provided, ensure the user exists
	if userID != nil {
		u, _ := s.userRepository.FindByID(*userID)
		if u == nil {
			return nil, errors.NewNotFoundError("user", nil)
		}
	}

	patch := &model.Task{ID: id}
	if title != nil {
		patch.Title = *title
	}
	if status != nil {
		patch.Status = *status
	}
	if userID != nil {
		patch.UserID = *userID
	}
	return s.taskRepository.Update(patch)
}
