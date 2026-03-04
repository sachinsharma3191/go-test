package health

import (
	"fmt"
	"runtime"
	"time"

	"go-backend/model"
)

// MemoryChecker checks memory usage
type MemoryChecker struct {
	MaxMemoryMB float64
}

// NewMemoryChecker creates a MemoryChecker.
func NewMemoryChecker(maxMemoryMB float64) *MemoryChecker {
	return &MemoryChecker{MaxMemoryMB: maxMemoryMB}
}

func (c *MemoryChecker) Name() string {
	return "memory"
}

func (c *MemoryChecker) Check() model.HealthCheckResult {
	start := time.Now()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	latency := time.Since(start)
	memoryMB := float64(memStats.Alloc) / 1024 / 1024
	status := model.HealthStatusHealthy
	message := fmt.Sprintf("Memory usage: %.2f MB", memoryMB)
	if memoryMB > c.MaxMemoryMB {
		status = model.HealthStatusDegraded
		message = fmt.Sprintf("High memory usage: %.2f MB (threshold: %.2f MB)", memoryMB, c.MaxMemoryMB)
	}
	return model.HealthCheckResult{
		Name:    c.Name(),
		Status:  status,
		Message: message,
		Latency: latency,
		Details: map[string]interface{}{
			"alloc_mb":   memoryMB,
			"sys_mb":     float64(memStats.Sys) / 1024 / 1024,
			"num_gc":     memStats.NumGC,
			"goroutines": runtime.NumGoroutine(),
		},
	}
}
