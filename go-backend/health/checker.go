package health

import "go-backend/model"

// Checker interface for health check implementations
type Checker interface {
	Name() string
	Check() model.HealthCheckResult
}
