package service

import (
	"testing"

	"go-backend/health"
	"go-backend/model"
)

func TestNewHealthService(t *testing.T) {
	mon := health.NewMonitor("1.0.0")
	svc := NewHealthService(mon)
	if svc == nil {
		t.Fatal("NewHealthService returned nil")
	}
	if svc.monitor != mon {
		t.Error("monitor not set correctly")
	}
}

func TestHealthService_GetHealthReport(t *testing.T) {
	mon := health.NewMonitor("2.0.0")
	mon.AddChecker(health.NewMemoryChecker(1e12)) // stay healthy
	svc := NewHealthService(mon)
	report := svc.GetHealthReport()
	if report.Status != "healthy" {
		t.Errorf("GetHealthReport Status = %q, want healthy", report.Status)
	}
	if report.Version != "2.0.0" {
		t.Errorf("Version = %q", report.Version)
	}
}

func TestHealthService_GetReadiness(t *testing.T) {
	mon := health.NewMonitor("test")
	mon.AddChecker(health.NewMemoryChecker(1e12))
	svc := NewHealthService(mon)
	rs := svc.GetReadiness()
	if !rs.Ready {
		t.Error("GetReadiness should be ready when healthy")
	}
}

func TestHealthService_GetLiveness(t *testing.T) {
	mon := health.NewMonitor("test")
	mon.AddChecker(health.NewMemoryChecker(1e12))
	svc := NewHealthService(mon)
	ls := svc.GetLiveness()
	if !ls.Alive {
		t.Error("GetLiveness should be alive when healthy")
	}
}

// stubChecker for testing
type healthStubChecker struct {
	result model.HealthCheckResult
}

func (c *healthStubChecker) Name() string { return "stub" }
func (c *healthStubChecker) Check() model.HealthCheckResult { return c.result }

func TestHealthService_GetLiveness_Unhealthy(t *testing.T) {
	mon := health.NewMonitor("test")
	mon.AddChecker(&healthStubChecker{result: model.HealthCheckResult{Status: model.HealthStatusUnhealthy}})
	svc := NewHealthService(mon)
	ls := svc.GetLiveness()
	if ls.Alive {
		t.Error("GetLiveness should be not alive when unhealthy")
	}
}
