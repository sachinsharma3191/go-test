package handler

import (
	"context"
	"errors"
	"os"
	"testing"

	"go-backend/health"
	"go-backend/model"
	"go-backend/repository"
	"go-backend/service"
	"go-backend/store"

	_ "go-backend/store/json"
)

var errFailingStore = errors.New("failing store")

// Ensure unhealthyStubChecker implements health.Checker
var _ health.Checker = (*unhealthyStubChecker)(nil)

// setupTestHandlers creates test handlers with a temporary store.
func setupTestHandlers(t *testing.T) (*UserHandler, *TaskHandler, *HealthHandler) {
	t.Helper()
	f, err := os.CreateTemp("", "handler_*.json")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	tempFile := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(tempFile) })

	jsonStore, err := store.QuickCreate("json", tempFile)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	data := store.AppData{
		Users: []model.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", Role: "developer"},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Role: "designer"},
		},
		Tasks: []model.Task{
			{ID: 1, Title: "Task 1", Status: "pending", UserID: 1},
			{ID: 2, Title: "Task 2", Status: "completed", UserID: 2},
		},
	}

	if err := jsonStore.Save(context.Background(), data); err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	userRepo := repository.NewUserRepository(jsonStore)
	taskRepo := repository.NewTaskRepository(jsonStore)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)

	monitor := health.NewMonitor("test")
	monitor.AddChecker(health.NewMemoryChecker(100.0))
	healthService := service.NewHealthService(monitor)

	return NewUserHandler(userService, taskService),
		NewTaskHandler(taskService),
		NewHealthHandler(healthService)
}

// setupTestHandlersWithUnhealthyMonitor creates handlers with a monitor that reports unhealthy (for Ready 503 path).
func setupTestHandlersWithUnhealthyMonitor(t *testing.T) *HealthHandler {
	t.Helper()
	f, _ := os.CreateTemp("", "unhealthy_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })

	st, _ := store.QuickCreate("json", path)
	t.Cleanup(func() { _ = st.Close() })
	userRepo := repository.NewUserRepository(st)
	taskRepo := repository.NewTaskRepository(st)
	_ = repository.SeedIfEmpty(userRepo, taskRepo)

	monitor := health.NewMonitor("test")
	monitor.AddChecker(&unhealthyStubChecker{})
	healthService := service.NewHealthService(monitor)
	return NewHealthHandler(healthService)
}

type unhealthyStubChecker struct{}

func (c *unhealthyStubChecker) Name() string { return "unhealthy" }
func (c *unhealthyStubChecker) Check() model.HealthCheckResult {
	return model.HealthCheckResult{Name: "unhealthy", Status: model.HealthStatusUnhealthy}
}

// failingStore implements store.Store and returns errors to trigger handler error paths.
type failingStore struct{}

func (s *failingStore) Load(_ context.Context) (store.AppData, error) {
	return store.AppData{}, errFailingStore
}
func (s *failingStore) Save(_ context.Context, _ store.AppData) error {
	return errFailingStore
}
func (s *failingStore) Health(_ context.Context) error {
	return nil
}
func (s *failingStore) Close() error {
	return nil
}

// loadOkSaveFailingStore returns valid data on Load but fails on Save (to trigger DataStoreError on create/update).
type loadOkSaveFailingStore struct {
	data store.AppData
}

func (s *loadOkSaveFailingStore) Load(_ context.Context) (store.AppData, error) {
	return s.data, nil
}
func (s *loadOkSaveFailingStore) Save(_ context.Context, _ store.AppData) error {
	return errFailingStore
}
func (s *loadOkSaveFailingStore) Health(_ context.Context) error {
	return nil
}
func (s *loadOkSaveFailingStore) Close() error {
	return nil
}

func setupTestHandlersWithFailingStore(t *testing.T) (*UserHandler, *TaskHandler) {
	t.Helper()
	st := &failingStore{}
	userRepo := repository.NewUserRepository(st)
	taskRepo := repository.NewTaskRepository(st)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	return NewUserHandler(userService, taskService), NewTaskHandler(taskService)
}

// setupTestHandlersWithSaveFailingStore uses a store that Load succeeds but Save fails (for create/update DataStoreError).
func setupTestHandlersWithSaveFailingStore(t *testing.T) (*UserHandler, *TaskHandler) {
	t.Helper()
	st := &loadOkSaveFailingStore{
		data: store.AppData{
			Users: []model.User{{ID: 1, Name: "User1", Email: "u1@example.com", Role: "developer"}},
			Tasks: []model.Task{{ID: 1, Title: "T1", Status: "pending", UserID: 1}},
		},
	}
	userRepo := repository.NewUserRepository(st)
	taskRepo := repository.NewTaskRepository(st)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	return NewUserHandler(userService, taskService), NewTaskHandler(taskService)
}

// secondLoadFailingStore returns data on first Load, error on subsequent Loads (Stats: users ok, tasks fail).
type secondLoadFailingStore struct {
	data     store.AppData
	loadNum  int
}

func (s *secondLoadFailingStore) Load(_ context.Context) (store.AppData, error) {
	s.loadNum++
	if s.loadNum > 1 {
		return store.AppData{}, errFailingStore
	}
	return s.data, nil
}
func (s *secondLoadFailingStore) Save(_ context.Context, d store.AppData) error {
	s.data = d
	return nil
}
func (s *secondLoadFailingStore) Health(_ context.Context) error { return nil }
func (s *secondLoadFailingStore) Close() error                  { return nil }

func setupTestHandlersWithSecondLoadFailing(t *testing.T) (*UserHandler, *TaskHandler) {
	t.Helper()
	st := &secondLoadFailingStore{
		data: store.AppData{
			Users: []model.User{{ID: 1, Name: "U1", Email: "u1@x.com", Role: "developer"}},
			Tasks: []model.Task{},
		},
	}
	userRepo := repository.NewUserRepository(st)
	taskRepo := repository.NewTaskRepository(st)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	return NewUserHandler(userService, taskService), NewTaskHandler(taskService)
}
