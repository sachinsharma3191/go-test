package service

import (
	"go-backend/health"
	"go-backend/model"
)

// HealthService exposes health, readiness, and liveness via the injected monitor.
type HealthService struct {
	monitor *health.Monitor
}

// NewHealthService creates a HealthService with the given monitor.
func NewHealthService(monitor *health.Monitor) *HealthService {
	return &HealthService{monitor: monitor}
}

// GetHealthReport returns the full health report (all checkers, readiness, liveness).
func (s *HealthService) GetHealthReport() model.HealthReport {
	return s.monitor.CheckHealth()
}

// GetReadiness returns the readiness status (ready for traffic or not).
func (s *HealthService) GetReadiness() model.ReadinessStatus {
	return s.monitor.CheckReadiness()
}

// GetLiveness returns the liveness status (service is alive).
func (s *HealthService) GetLiveness() model.LivenessStatus {
	return s.monitor.CheckLiveness()
}
