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

func main() {
	runMain(nil)
}

// fatalFunc is used for exit on error; can be overridden in tests.
var fatalFunc = log.Fatal

// runMain is the testable entry point. When shutdownCh is nil, uses signal-based shutdown.
func runMain(shutdownCh <-chan struct{}) {
	if err := run(shutdownCh); err != nil {
		fatalFunc(err)
	}
}

// run configures and runs the server. Extracted for testability.
// If shutdownCh is non-nil, it is used for shutdown; otherwise signal-based shutdown is used.
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

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// newHandler returns an http.Handler with all routes registered.
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
// shutdownCh: when closed/received, triggers shutdown. shutdownTimeout: max wait for shutdown; if <=0, uses 30s.
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
