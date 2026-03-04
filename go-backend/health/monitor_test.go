package health

import (
	"testing"
	"time"

	"go-backend/model"
)

// stubChecker implements Checker for testing
type stubChecker struct {
	name   string
	result model.HealthCheckResult
}

func (s *stubChecker) Name() string { return s.name }
func (s *stubChecker) Check() model.HealthCheckResult { return s.result }

func TestNewMonitor(t *testing.T) {
	m := NewMonitor("1.0.0")
	if m == nil {
		t.Fatal("NewMonitor returned nil")
	}
	if m.version != "1.0.0" {
		t.Errorf("version = %q, want 1.0.0", m.version)
	}
	if len(m.checkers) != 0 {
		t.Errorf("expected 0 checkers, got %d", len(m.checkers))
	}
}

func TestMonitor_AddChecker(t *testing.T) {
	m := NewMonitor("test")
	c := &stubChecker{name: "stub", result: model.HealthCheckResult{Status: model.HealthStatusHealthy}}
	m.AddChecker(c)
	if len(m.checkers) != 1 {
		t.Errorf("expected 1 checker, got %d", len(m.checkers))
	}
	// chaining
	m.AddChecker(c).AddChecker(c)
	if len(m.checkers) != 3 {
		t.Errorf("expected 3 checkers after chaining, got %d", len(m.checkers))
	}
}

func TestMonitor_CheckHealth_Healthy(t *testing.T) {
	m := NewMonitor("v1")
	m.AddChecker(&stubChecker{name: "a", result: model.HealthCheckResult{Status: model.HealthStatusHealthy, Name: "a"}})
	m.AddChecker(&stubChecker{name: "b", result: model.HealthCheckResult{Status: model.HealthStatusHealthy, Name: "b"}})
	report := m.CheckHealth()
	if report.Status != "healthy" {
		t.Errorf("Status = %q, want healthy", report.Status)
	}
	if !report.Readiness.Ready {
		t.Error("Readiness.Ready should be true")
	}
	if !report.Liveness.Alive {
		t.Error("Liveness.Alive should be true")
	}
	if len(report.Checks) != 2 {
		t.Errorf("expected 2 checks, got %d", len(report.Checks))
	}
	if report.Version != "v1" {
		t.Errorf("Version = %q, want v1", report.Version)
	}
	if report.Uptime < 0 {
		t.Error("Uptime should be non-negative")
	}
}

func TestMonitor_CheckHealth_Degraded(t *testing.T) {
	m := NewMonitor("v1")
	m.AddChecker(&stubChecker{name: "ok", result: model.HealthCheckResult{Status: model.HealthStatusHealthy, Name: "ok"}})
	m.AddChecker(&stubChecker{name: "bad", result: model.HealthCheckResult{Status: model.HealthStatusDegraded, Name: "bad"}})
	report := m.CheckHealth()
	if report.Status != "degraded" {
		t.Errorf("Status = %q, want degraded", report.Status)
	}
	if report.Readiness.Ready {
		t.Error("Readiness.Ready should be false when degraded")
	}
	if report.Readiness.Reason != "System degraded" {
		t.Errorf("Readiness.Reason = %q, want System degraded", report.Readiness.Reason)
	}
	if !report.Liveness.Alive {
		t.Error("Liveness.Alive should be true when degraded (only unhealthy kills liveness)")
	}
}

func TestMonitor_CheckHealth_Unhealthy(t *testing.T) {
	m := NewMonitor("v1")
	m.AddChecker(&stubChecker{name: "critical", result: model.HealthCheckResult{Status: model.HealthStatusUnhealthy, Name: "critical"}})
	report := m.CheckHealth()
	if report.Status != "unhealthy" {
		t.Errorf("Status = %q, want unhealthy", report.Status)
	}
	if report.Readiness.Ready {
		t.Error("Readiness.Ready should be false when unhealthy")
	}
	if report.Readiness.Reason != "System unhealthy" {
		t.Errorf("Readiness.Reason = %q, want System unhealthy", report.Readiness.Reason)
	}
	if report.Liveness.Alive {
		t.Error("Liveness.Alive should be false when unhealthy")
	}
	if report.Liveness.Reason != "Critical component failure" {
		t.Errorf("Liveness.Reason = %q", report.Liveness.Reason)
	}
}

func TestMonitor_CheckHealth_UnhealthyOverridesDegraded(t *testing.T) {
	m := NewMonitor("v1")
	m.AddChecker(&stubChecker{name: "d", result: model.HealthCheckResult{Status: model.HealthStatusDegraded, Name: "d"}})
	m.AddChecker(&stubChecker{name: "u", result: model.HealthCheckResult{Status: model.HealthStatusUnhealthy, Name: "u"}})
	report := m.CheckHealth()
	if report.Status != "unhealthy" {
		t.Errorf("Status = %q, want unhealthy (unhealthy overrides degraded)", report.Status)
	}
}

func TestMonitor_CheckReadiness(t *testing.T) {
	m := NewMonitor("v1")
	m.AddChecker(&stubChecker{name: "x", result: model.HealthCheckResult{Status: model.HealthStatusHealthy, Name: "x"}})
	rs := m.CheckReadiness()
	if !rs.Ready {
		t.Error("CheckReadiness should be ready when healthy")
	}
}

func TestMonitor_CheckLiveness(t *testing.T) {
	m := NewMonitor("v1")
	m.AddChecker(&stubChecker{name: "x", result: model.HealthCheckResult{Status: model.HealthStatusHealthy, Name: "x"}})
	ls := m.CheckLiveness()
	if !ls.Alive {
		t.Error("CheckLiveness should be alive when healthy")
	}
}

func TestMonitor_NoCheckers(t *testing.T) {
	m := NewMonitor("v1")
	report := m.CheckHealth()
	if report.Status != "healthy" {
		t.Errorf("Status = %q, want healthy when no checkers", report.Status)
	}
	if len(report.Checks) != 0 {
		t.Errorf("expected 0 checks, got %d", len(report.Checks))
	}
}

func TestMonitor_ReportFields(t *testing.T) {
	m := NewMonitor("2.0.0")
	m.AddChecker(&stubChecker{name: "c", result: model.HealthCheckResult{
		Name: "c", Status: model.HealthStatusHealthy, Message: "ok", Latency: time.Second,
	}})
	report := m.CheckHealth()
	if report.Message != "System health check completed" {
		t.Errorf("Message = %q", report.Message)
	}
	if report.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
	if rc, ok := report.Checks["c"]; !ok {
		t.Error("check c not in report")
	} else if rc.Message != "ok" || rc.Latency != time.Second {
		t.Errorf("check c: Message=%q Latency=%v", rc.Message, rc.Latency)
	}
}
