package model

import "time"

// HealthStatus represents health check status.
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
)

// HealthCheckResult represents the result of a health check.
type HealthCheckResult struct {
	Name    string        `json:"name"`
	Status  HealthStatus  `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency"`
	Details interface{}   `json:"details,omitempty"`
}

// ReadinessStatus represents readiness status.
type ReadinessStatus struct {
	Ready  bool   `json:"ready"`
	Reason string `json:"reason,omitempty"`
}

// LivenessStatus represents liveness status.
type LivenessStatus struct {
	Alive  bool   `json:"alive"`
	Reason string `json:"reason,omitempty"`
}

// HealthReport represents a comprehensive health report.
type HealthReport struct {
	Status    string                       `json:"status"`
	Message   string                       `json:"message"`
	Version   string                       `json:"version"`
	Timestamp time.Time                    `json:"timestamp"`
	Uptime    time.Duration                `json:"uptime"`
	Checks    map[string]HealthCheckResult `json:"checks"`
	Readiness ReadinessStatus              `json:"readiness"`
	Liveness  LivenessStatus               `json:"liveness"`
}
