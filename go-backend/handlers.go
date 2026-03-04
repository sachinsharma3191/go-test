package main

import (
	"go-backend/handler"
	"go-backend/health"
	"go-backend/repository"
	"go-backend/service"
)

// Handlers contains all HTTP handlers for efficient parameter passing
type Handlers struct {
	UserHandler   *handler.UserHandler
	TaskHandler   *handler.TaskHandler
	HealthHandler *handler.HealthHandler
}

// NewHandlers creates and initializes all handlers with their dependencies
func NewHandlers(userRepo *repository.UserRepository, taskRepo *repository.TaskRepository, monitor *health.Monitor) *Handlers {
	// Initialize services
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	healthService := service.NewHealthService(monitor)

	// Initialize controllers
	userHandler := handler.NewUserHandler(userService, taskService)
	taskHandler := handler.NewTaskHandler(taskService)
	healthHandler := handler.NewHealthHandler(healthService)

	return &Handlers{
		UserHandler:   userHandler,
		TaskHandler:   taskHandler,
		HealthHandler: healthHandler,
	}
}
