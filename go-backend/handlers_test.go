// handlers_test.go
package main

import (
	"testing"

	"go-backend/health"
	"go-backend/repository"
)

// TestNewHandlers_InitializesAllHandlers verifies handler initialization
func TestNewHandlers_InitializesAllHandlers(t *testing.T) {
	// Arrange
	userRepo := &repository.UserRepository{}
	taskRepo := &repository.TaskRepository{}
	monitor := health.NewMonitor("test-version")

	// Act
	handlers := NewHandlers(userRepo, taskRepo, monitor)

	// Assert
	if handlers == nil {
		t.Fatal("expected handlers to be initialized, got nil")
	}

	if handlers.UserHandler == nil {
		t.Error("expected UserHandler to be initialized")
	}

	if handlers.TaskHandler == nil {
		t.Error("expected TaskHandler to be initialized")
	}

	if handlers.HealthHandler == nil {
		t.Error("expected HealthHandler to be initialized")
	}
}

// TestNewHandlers_DoesNotShareInstances verifies each call creates new instances
func TestNewHandlers_DoesNotShareInstances(t *testing.T) {
	userRepo := &repository.UserRepository{}
	taskRepo := &repository.TaskRepository{}
	monitor := health.NewMonitor("test-version")

	h1 := NewHandlers(userRepo, taskRepo, monitor)
	h2 := NewHandlers(userRepo, taskRepo, monitor)

	if h1 == h2 {
		t.Fatal("expected different handler instances")
	}
}
