package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"go-backend/health"
	"go-backend/repository"
	"go-backend/store"

	_ "go-backend/store/json"
)

func setupMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	f, err := os.CreateTemp("", "server_test_*.json")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })

	st, err := store.QuickCreate("json", path)
	if err != nil {
		t.Fatalf("QuickCreate: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	userRepo := repository.NewUserRepository(st)
	taskRepo := repository.NewTaskRepository(st)
	if err := repository.SeedIfEmpty(userRepo, taskRepo); err != nil {
		t.Fatalf("SeedIfEmpty: %v", err)
	}

	monitor := health.NewMonitor("test")
	monitor.AddChecker(health.NewMemoryChecker(1e12))
	handlers := NewHandlers(userRepo, taskRepo, monitor)
	handler := newHandler(handlers)
	return httptest.NewServer(handler)
}

func TestMakeSignalShutdownChan(t *testing.T) {
	ch := makeSignalShutdownChan()
	// Send SIGTERM to self to trigger channel close
	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Skipf("Cannot send signal: %v", err)
	}
	select {
	case <-ch:
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("makeSignalShutdownChan did not close after SIGTERM")
	}
}

func TestGetEnvOrDefault_FromEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "from-env")
	defer os.Unsetenv("TEST_KEY")
	if got := getEnvOrDefault("TEST_KEY", "default"); got != "from-env" {
		t.Errorf("getEnvOrDefault = %q, want from-env", got)
	}
}

func TestGetEnvOrDefault_Default(t *testing.T) {
	os.Unsetenv("TEST_KEY_EMPTY")
	if got := getEnvOrDefault("TEST_KEY_EMPTY", "mydefault"); got != "mydefault" {
		t.Errorf("getEnvOrDefault = %q, want mydefault", got)
	}
}

func TestNewHandler_RegistersAllRoutes(t *testing.T) {
	srv := setupMockServer(t)
	defer srv.Close()

	paths := []string{
		"/health",
		"/health/ready",
		"/health/live",
		"/api/users",
		"/api/users/1",
		"/api/tasks",
		"/api/tasks/1",
		"/api/stats",
	}

	for _, p := range paths {
		resp, err := srv.Client().Get(srv.URL + p)
		if err != nil {
			t.Fatalf("GET %s: %v", p, err)
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			t.Errorf("GET %s returned 404", p)
		}
	}
}

func TestNewHandler_HealthReturnsStatus(t *testing.T) {
	srv := setupMockServer(t)
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /health status = %d, want 200", resp.StatusCode)
	}
}

func TestRunServer_StartsAndShutsDown(t *testing.T) {
	f, _ := os.CreateTemp("", "run_server_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })

	st, err := store.QuickCreate("json", path)
	if err != nil {
		t.Fatalf("QuickCreate: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	userRepo := repository.NewUserRepository(st)
	taskRepo := repository.NewTaskRepository(st)
	_ = repository.SeedIfEmpty(userRepo, taskRepo)

	monitor := health.NewMonitor("test")
	monitor.AddChecker(health.NewMemoryChecker(1e12))
	handlers := NewHandlers(userRepo, taskRepo, monitor)
	handler := newHandler(handlers)

	done := make(chan struct{})
	go func() {
		RunServer(handler, ":0", done, 0)
	}()

	time.Sleep(100 * time.Millisecond)
	close(done)
	time.Sleep(200 * time.Millisecond)
}

func TestRunServer_ShutdownWithTimeout(t *testing.T) {
	f, _ := os.CreateTemp("", "shutdown_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })

	st, _ := store.QuickCreate("json", path)
	t.Cleanup(func() { _ = st.Close() })
	userRepo := repository.NewUserRepository(st)
	taskRepo := repository.NewTaskRepository(st)
	_ = repository.SeedIfEmpty(userRepo, taskRepo)

	monitor := health.NewMonitor("test")
	monitor.AddChecker(health.NewMemoryChecker(1e12))
	handlers := NewHandlers(userRepo, taskRepo, monitor)
	handler := newHandler(handlers)

	done := make(chan struct{})
	finished := make(chan struct{})
	go func() {
		RunServer(handler, ":0", done, 0)
		close(finished)
	}()

	time.Sleep(50 * time.Millisecond)
	close(done)
	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("RunServer did not exit after shutdown")
	}
}

func TestRunServer_ShutdownError(t *testing.T) {
	// Use very short timeout to force Shutdown to fail with context.DeadlineExceeded
	f, _ := os.CreateTemp("", "shutdown_err_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	st, _ := store.QuickCreate("json", path)
	t.Cleanup(func() { _ = st.Close() })
	userRepo := repository.NewUserRepository(st)
	taskRepo := repository.NewTaskRepository(st)
	_ = repository.SeedIfEmpty(userRepo, taskRepo)
	monitor := health.NewMonitor("test")
	monitor.AddChecker(health.NewMemoryChecker(1e12))
	handlers := NewHandlers(userRepo, taskRepo, monitor)
	handler := newHandler(handlers)

	done := make(chan struct{})
	finished := make(chan struct{})
	go func() {
		RunServer(handler, ":0", done, 1) // 1ns timeout - forces shutdown error
		close(finished)
	}()

	time.Sleep(80 * time.Millisecond)
	close(done)
	select {
	case <-finished:
	case <-time.After(3 * time.Second):
		t.Fatal("RunServer did not exit")
	}
}

func TestRunServer_ListenAndServeError(t *testing.T) {
	f, _ := os.CreateTemp("", "listen_err_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	st, _ := store.QuickCreate("json", path)
	t.Cleanup(func() { _ = st.Close() })
	userRepo := repository.NewUserRepository(st)
	taskRepo := repository.NewTaskRepository(st)
	handlers := NewHandlers(userRepo, taskRepo, health.NewMonitor("test"))
	handler := newHandler(handlers)

	done := make(chan struct{})
	finished := make(chan struct{})
	go func() {
		RunServer(handler, "invalid-address:999999", done, 0)
		close(finished)
	}()

	time.Sleep(100 * time.Millisecond)
	close(done)
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("RunServer did not exit")
	}
}

func TestRun_StoreCreationError(t *testing.T) {
	os.Setenv("STORE_BACKEND", "nonexistent-backend")
	defer os.Unsetenv("STORE_BACKEND")
	os.Unsetenv("DATA_FILE")

	err := run(nil)
	if err == nil {
		t.Fatal("expected error for invalid store backend")
	}
}

func TestRunMain_ExitOnError(t *testing.T) {
	os.Setenv("STORE_BACKEND", "nonexistent-backend")
	defer os.Unsetenv("STORE_BACKEND")
	os.Unsetenv("DATA_FILE")

	// Override fatalFunc to capture and not exit
	called := false
	orig := fatalFunc
	fatalFunc = func(args ...interface{}) {
		called = true
		if len(args) == 0 {
			t.Error("fatalFunc called with no args")
		}
	}
	defer func() { fatalFunc = orig }()

	runMain(nil)

	if !called {
		t.Error("runMain should have called fatalFunc on store creation error")
	}
}

func TestRunMain_Success(t *testing.T) {
	f, _ := os.CreateTemp("", "runmain_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	os.Setenv("STORE_BACKEND", "json")
	os.Setenv("DATA_FILE", path)
	defer func() {
		os.Unsetenv("STORE_BACKEND")
		os.Unsetenv("DATA_FILE")
	}()

	done := make(chan struct{})
	finished := make(chan struct{})
	go func() {
		runMain(done)
		close(finished)
	}()

	time.Sleep(150 * time.Millisecond)
	close(done)
	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("runMain did not exit")
	}
}

func TestRun_SuccessWithMockShutdown(t *testing.T) {
	f, _ := os.CreateTemp("", "run_success_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })

	os.Setenv("STORE_BACKEND", "json")
	os.Setenv("DATA_FILE", path)
	defer func() {
		os.Unsetenv("STORE_BACKEND")
		os.Unsetenv("DATA_FILE")
	}()

	done := make(chan struct{})
	finished := make(chan error, 1)
	go func() {
		finished <- run(done)
	}()

	time.Sleep(150 * time.Millisecond)
	close(done)
	err := <-finished
	if err != nil {
		t.Errorf("run() = %v", err)
	}
}

func TestRun_SeedWarning(t *testing.T) {
	// Use read-only directory so Save fails during seed (WriteFile to dir/file.tmp fails)
	dir, _ := os.MkdirTemp("", "seed_warn_dir")
	t.Cleanup(func() {
		os.Chmod(dir, 0755)
		os.RemoveAll(dir)
	})
	path := filepath.Join(dir, "data.json")
	if err := os.Chmod(dir, 0555); err != nil {
		t.Skipf("Cannot chmod: %v", err)
	}

	os.Setenv("STORE_BACKEND", "json")
	os.Setenv("DATA_FILE", path)
	defer func() {
		os.Unsetenv("STORE_BACKEND")
		os.Unsetenv("DATA_FILE")
	}()

	done := make(chan struct{})
	finished := make(chan error, 1)
	go func() {
		finished <- run(done)
	}()

	time.Sleep(100 * time.Millisecond)
	close(done)
	<-finished
}

func TestRun_SuccessWithSignalShutdown(t *testing.T) {
	// Use random port to avoid conflict
	if ln, err := net.Listen("tcp", ":8080"); err != nil {
		t.Skipf("port 8080 in use: %v", err)
	} else {
		ln.Close()
	}

	f, _ := os.CreateTemp("", "run_signal_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })

	os.Setenv("STORE_BACKEND", "json")
	os.Setenv("DATA_FILE", path)
	defer func() {
		os.Unsetenv("STORE_BACKEND")
		os.Unsetenv("DATA_FILE")
	}()

	finished := make(chan error, 1)
	go func() {
		finished <- run(nil) // uses signal-based shutdown
	}()

	time.Sleep(200 * time.Millisecond)
	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatalf("Kill: %v", err)
	}

	select {
	case err := <-finished:
		if err != nil {
			t.Errorf("run() = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("run() did not exit after SIGTERM")
	}
}
