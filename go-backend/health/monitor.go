package health

import (
	"time"

	"go-backend/model"
)

// Monitor runs health checkers and aggregates results into reports.
type Monitor struct {
	checkers  []Checker
	startTime time.Time
	version   string
}

// NewMonitor creates a new Monitor with the given version string.
func NewMonitor(version string) *Monitor {
	return &Monitor{
		checkers:  make([]Checker, 0),
		startTime: time.Now(),
		version:   version,
	}
}

// AddChecker registers a checker. Returns the monitor for chaining.
func (m *Monitor) AddChecker(c Checker) *Monitor {
	m.checkers = append(m.checkers, c)
	return m
}

// CheckHealth runs all checkers and returns an aggregated health report.
func (m *Monitor) CheckHealth() model.HealthReport {
	checks := make(map[string]model.HealthCheckResult)
	status := model.HealthStatusHealthy
	
	for _, c := range m.checkers {
		result := c.Check()
		checks[c.Name()] = result
		switch result.Status {
		case model.HealthStatusUnhealthy:
			status = model.HealthStatusUnhealthy
		case model.HealthStatusDegraded:
			if status == model.HealthStatusHealthy {
				status = model.HealthStatusDegraded
			}
		}
	}
	
	readiness, liveness := m.deriveReadinessLiveness(status)
	return model.HealthReport{
		Status:    string(status),
		Message:   "System health check completed",
		Version:   m.version,
		Timestamp: time.Now().UTC(),
		Uptime:    time.Since(m.startTime),
		Checks:    checks,
		Readiness: readiness,
		Liveness:  liveness,
	}
}

func (m *Monitor) deriveReadinessLiveness(status model.HealthStatus) (model.ReadinessStatus, model.LivenessStatus) {
	readiness := model.ReadinessStatus{Ready: true}
	liveness := model.LivenessStatus{Alive: true, Reason: "Service is responding to requests"}
	
	switch status {
	case model.HealthStatusUnhealthy:
		readiness.Ready = false
		readiness.Reason = "System unhealthy"
		liveness.Alive = false
		liveness.Reason = "Critical component failure"
	case model.HealthStatusDegraded:
		readiness.Ready = false
		readiness.Reason = "System degraded"
	}
	
	return readiness, liveness
}

// CheckReadiness returns readiness status (whether the service can accept traffic).
func (m *Monitor) CheckReadiness() model.ReadinessStatus {
	return m.CheckHealth().Readiness
}

// CheckLiveness returns liveness status (whether the process is alive and healthy).
func (m *Monitor) CheckLiveness() model.LivenessStatus {
	return m.CheckHealth().Liveness
}
