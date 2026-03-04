package health

import (
	"testing"

	"go-backend/model"
)

func TestNewMemoryChecker(t *testing.T) {
	c := NewMemoryChecker(100.0)
	if c == nil {
		t.Fatal("NewMemoryChecker returned nil")
	}
	if c.MaxMemoryMB != 100.0 {
		t.Errorf("MaxMemoryMB = %v, want 100", c.MaxMemoryMB)
	}
}

func TestMemoryChecker_Name(t *testing.T) {
	c := NewMemoryChecker(50.0)
	if c.Name() != "memory" {
		t.Errorf("Name() = %q, want memory", c.Name())
	}
}

func TestMemoryChecker_Check_Healthy(t *testing.T) {
	// Use a very high threshold so we stay healthy (unless machine has huge alloc)
	c := NewMemoryChecker(1e12) // 1TB
	result := c.Check()
	if result.Status != model.HealthStatusHealthy {
		t.Errorf("Status = %v, want healthy (alloc should be under 1TB)", result.Status)
	}
	if result.Name != "memory" {
		t.Errorf("Name = %q", result.Name)
	}
	if result.Latency <= 0 {
		t.Error("Latency should be positive")
	}
	details, ok := result.Details.(map[string]interface{})
	if !ok {
		t.Fatalf("Details should be map, got %T", result.Details)
	}
	if _, ok := details["alloc_mb"]; !ok {
		t.Error("Details should have alloc_mb")
	}
	if _, ok := details["sys_mb"]; !ok {
		t.Error("Details should have sys_mb")
	}
	if _, ok := details["num_gc"]; !ok {
		t.Error("Details should have num_gc")
	}
	if _, ok := details["goroutines"]; !ok {
		t.Error("Details should have goroutines")
	}
}

func TestMemoryChecker_Check_Degraded(t *testing.T) {
	// Use 0 so any alloc is over threshold
	c := NewMemoryChecker(0)
	result := c.Check()
	if result.Status != model.HealthStatusDegraded {
		t.Errorf("Status = %v, want degraded (alloc > 0)", result.Status)
	}
}
