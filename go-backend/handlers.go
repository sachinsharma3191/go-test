// Package main provides HTTP handler initialization and dependency injection.
//
// This package contains the main handler container and factory functions for creating
// all HTTP handlers with their proper dependencies. It follows the dependency injection
// pattern to ensure testability and separation of concerns.
//
// Handler Architecture:
//   - UserHandler: Manages user-related HTTP endpoints
//   - TaskHandler: Manages task-related HTTP endpoints  
//   - HealthHandler: Manages health check endpoints
//
// Dependency Flow:
//   Store → Repository → Service → Handler
//
// Each handler receives its required services through constructor injection,
// making them easy to test and maintain.
package main

import (
	"go-backend/handler"
	"go-backend/health"
	"go-backend/repository"
	"go-backend/service"
)

// Handlers contains all HTTP handlers for efficient parameter passing.
// This struct acts as a container for all route handlers, allowing them to be
// created once and passed together to the router setup.
//
// Fields:
//   - UserHandler: Handles user CRUD operations and statistics
//   - TaskHandler: Handles task CRUD operations and management
//   - HealthHandler: Handles health check and monitoring endpoints
type Handlers struct {
	UserHandler   *handler.UserHandler
	TaskHandler   *handler.TaskHandler
	HealthHandler *handler.HealthHandler
}

// NewHandlers creates and initializes all handlers with their dependencies.
// This function implements the dependency injection pattern by creating the
// service layer first, then injecting those services into the appropriate handlers.
//
// Creation Order:
//   1. Services are created with repository dependencies
//   2. Handlers are created with service dependencies
//   3. All handlers are packaged into a Handlers struct
//
// Parameters:
//   - userRepo: Repository for user data access
//   - taskRepo: Repository for task data access
//   - monitor: Health monitoring system
//
// Returns:
//   - *Handlers: Container with all initialized HTTP handlers
//
// Example:
//   handlers := NewHandlers(userRepo, taskRepo, healthMonitor)
//   // Use handlers to set up HTTP routes
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
