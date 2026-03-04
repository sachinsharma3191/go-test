package handler

import (
	"net/http"

	"go-backend/errors"
	"go-backend/middleware"
	"go-backend/service"
)

// HealthHandler handles health-related HTTP endpoints.
type HealthHandler struct {
	healthService *service.HealthService
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(healthService *service.HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

// Health returns the full health report.
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendResponse(w, r, middleware.Error(errors.NewInvalidMethodError(r.Method), r))
		return
	}

	report := h.healthService.GetHealthReport()
	middleware.SendResponse(w, r, middleware.Data(http.StatusOK, report))
}

// Ready returns readiness status (503 if not ready).
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendResponse(w, r, middleware.Error(errors.NewInvalidMethodError(r.Method), r))
		return
	}

	status := h.healthService.GetReadiness()
	code := http.StatusOK
	if !status.Ready {
		code = http.StatusServiceUnavailable
	}

	middleware.SendResponse(w, r, middleware.Data(code, status))
}

// Live returns liveness status.
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendResponse(w, r, middleware.Error(errors.NewInvalidMethodError(r.Method), r))
		return
	}

	status := h.healthService.GetLiveness()
	middleware.SendResponse(w, r, middleware.Data(http.StatusOK, status))
}
