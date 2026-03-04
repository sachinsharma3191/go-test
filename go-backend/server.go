// Package main provides the entry point and server configuration for the Go backend API.
//
// This server implements a RESTful API for user and task management with the following features:
//   - JSON file-based storage with atomic operations
//   - In-memory caching with TTL support
//   - Health monitoring and readiness probes
//   - Structured logging and error handling
//   - Graceful shutdown with signal handling
//
// Architecture:
//   - Handler Layer: HTTP request/response processing
//   - Service Layer: Business logic and validation
//   - Repository Layer: Data access and persistence
//   - Store Layer: File-based storage abstraction
//
// Environment Variables:
//   - STORE_BACKEND: Storage backend type (default: "json")
//   - DATA_FILE: Path to JSON data file (default: "data/data.json")
//
// Default Configuration:
//   - Port: 8080
//   - Read/Write Timeout: 15 seconds
//   - Idle Timeout: 60 seconds
//   - Shutdown Timeout: 30 seconds
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-backend/health"
	"go-backend/repository"
	"go-backend/store"
	_ "go-backend/store/json" // Import to register JSON factory
)

const (
	version = "1.0.0"
	port    = "8080"
)

// main is the application entry point.
// It initializes and starts the HTTP server with default signal-based shutdown handling.
// This function calls runMain with a nil shutdown channel to enable signal handling.
func main() {
	runMain(nil)
}

// fatalFunc is used for exit on error; can be overridden in tests.
var fatalFunc = log.Fatal

// runMain is the testable entry point. When shutdownCh is nil, uses signal-based shutdown.
// This function allows for dependency injection of the shutdown channel in tests.
//
// Parameters:
//   - shutdownCh: Channel for triggering shutdown. If nil, signal-based shutdown is used.
func runMain(shutdownCh <-chan struct{}) {
	if err := run(shutdownCh); err != nil {
		fatalFunc(err)
	}
}

// run configures and runs the server. Extracted for testability.
// If shutdownCh is non-nil, it is used for shutdown; otherwise signal-based shutdown is used.
//
// This function performs the following steps:
//   1. Loads environment configuration
//   2. Creates and initializes the data store
//   3. Sets up repositories and services
//   4. Seeds initial data if the store is empty
//   5. Configures health monitoring
//   6. Creates and registers HTTP handlers
//   7. Starts the HTTP server
//
// Parameters:
//   - shutdownCh: Channel for triggering shutdown. If nil, signal-based shutdown is used.
//
// Returns:
//   - error: Any error encountered during server setup
func run(shutdownCh <-chan struct{}) error {
	backend := getEnvOrDefault("STORE_BACKEND", "json")
	dataFile := getEnvOrDefault("DATA_FILE", "data/data.json")

	st, err := store.QuickCreate(backend, dataFile)
	if err != nil {
		return err
	}
	defer st.Close()

	userRepo := repository.NewUserRepository(st)
	taskRepo := repository.NewTaskRepository(st)
	if err := repository.SeedIfEmpty(userRepo, taskRepo); err != nil {
		log.Printf("Warning: seed data: %v", err)
	}

	monitor := health.NewMonitor(version)
	monitor.AddChecker(health.NewMemoryChecker(100.0))
	handlers := NewHandlers(userRepo, taskRepo, monitor)

	ch := shutdownCh
	if ch == nil {
		ch = makeSignalShutdownChan()
	}
	RunServer(newHandler(handlers), ":"+port, ch, 0)
	return nil
}

// makeSignalShutdownChan returns a channel that closes when SIGINT or SIGTERM is received.
// This enables graceful shutdown when the process receives termination signals.
//
// Returns:
//   - <-chan struct{}: Channel that closes when shutdown signal is received
func makeSignalShutdownChan() <-chan struct{} {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})
	go func() {
		<-quit
		close(done)
	}()
	return done
}

// getEnvOrDefault returns the value of an environment variable or a default value.
// This function provides a convenient way to handle optional environment configuration.
//
// Parameters:
//   - key: The environment variable name to look up
//   - def: The default value to return if the environment variable is not set
//
// Returns:
//   - string: The environment variable value or the default
func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// newHandler returns an http.Handler with all routes registered.
// This function sets up the HTTP routing for all API endpoints.
//
// Registered Routes:
//   - GET /health: Comprehensive health report
//   - GET /health/ready: Readiness probe (Kubernetes)
//   - GET /health/live: Liveness probe (Kubernetes)
//   - GET|POST /api/users: User listing and creation
//   - GET|PUT|DELETE /api/users/{id}: User operations by ID
//   - GET|POST /api/tasks: Task listing and creation
//   - GET|PUT|DELETE /api/tasks/{id}: Task operations by ID
//   - GET /api/stats: System statistics
//
// Parameters:
//   - handlers: Container for all HTTP handlers
//
// Returns:
//   - http.Handler: Configured HTTP handler with all routes registered
func newHandler(handlers *Handlers) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler.Health)
	mux.HandleFunc("/health/ready", handlers.HealthHandler.Ready)
	mux.HandleFunc("/health/live", handlers.HealthHandler.Live)
	mux.HandleFunc("/api/users", handlers.UserHandler.Users)
	mux.HandleFunc("/api/users/", handlers.UserHandler.UserByID)
	mux.HandleFunc("/api/tasks", handlers.TaskHandler.HandleTasks)
	mux.HandleFunc("/api/tasks/", handlers.TaskHandler.HandleTaskByID)
	mux.HandleFunc("/api/stats", handlers.UserHandler.Stats)
	return mux
}

// RunServer runs the HTTP server until the shutdown channel receives a value.
// This function provides a configurable HTTP server with graceful shutdown capabilities.
//
// Server Configuration:
//   - Read Timeout: 15 seconds (time to read request headers and body)
//   - Write Timeout: 15 seconds (time to write response)
//   - Idle Timeout: 60 seconds (time to keep idle connections open)
//
// Shutdown Process:
//   1. Wait for signal on shutdownCh
//   2. Log shutdown initiation
//   3. Create context with timeout (default 30 seconds)
//   4. Call server.Shutdown() to gracefully stop accepting new connections
//   5. Wait for existing connections to finish or timeout
//   6. Log completion
//
// Parameters:
//   - handler: HTTP handler containing all route handlers
//   - addr: Server address in format ":port" (e.g., ":8080")
//   - shutdownCh: Channel that triggers shutdown when closed/received
//   - shutdownTimeout: Maximum time to wait for graceful shutdown.
//     If <= 0, defaults to 30 seconds.
//
// Example Usage:
//   RunServer(handler, ":8080", makeShutdownChan(), 30*time.Second)
func RunServer(handler http.Handler, addr string, shutdownCh <-chan struct{}, shutdownTimeout time.Duration) {
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on http://localhost%s", addr)
		log.Printf("Version: %s", version)
		log.Printf("Health check: http://localhost%s/health", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server failed to start: %v", err)
		}
	}()

	<-shutdownCh
	log.Println("Shutting down server...")

	timeout := shutdownTimeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("Server shutdown complete")
}
