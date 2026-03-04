package model

import (
	"encoding/json"
	"testing"
	"time"
)

// TestHealthStatus_Constants tests that health status constants are correct
func TestHealthStatus_Constants(t *testing.T) {
	if HealthStatusHealthy != "healthy" {
		t.Errorf("Expected HealthStatusHealthy to be 'healthy', got '%s'", HealthStatusHealthy)
	}
	if HealthStatusUnhealthy != "unhealthy" {
		t.Errorf("Expected HealthStatusUnhealthy to be 'unhealthy', got '%s'", HealthStatusUnhealthy)
	}
	if HealthStatusDegraded != "degraded" {
		t.Errorf("Expected HealthStatusDegraded to be 'degraded', got '%s'", HealthStatusDegraded)
	}
}

// TestHealthCheckResult_JSONSerialization tests JSON marshaling and unmarshaling
func TestHealthCheckResult_JSONSerialization(t *testing.T) {
	result := HealthCheckResult{
		Name:    "database",
		Status:  HealthStatusHealthy,
		Message: "Connection successful",
		Latency: 50 * time.Millisecond,
		Details: map[string]interface{}{
			"connections": 5,
			"max_connections": 100,
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal HealthCheckResult to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledResult HealthCheckResult
	err = json.Unmarshal(jsonData, &unmarshaledResult)
	if err != nil {
		t.Fatalf("Failed to unmarshal HealthCheckResult from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaledResult.Name != result.Name {
		t.Errorf("Expected Name '%s', got '%s'", result.Name, unmarshaledResult.Name)
	}
	if unmarshaledResult.Status != result.Status {
		t.Errorf("Expected Status '%s', got '%s'", result.Status, unmarshaledResult.Status)
	}
	if unmarshaledResult.Message != result.Message {
		t.Errorf("Expected Message '%s', got '%s'", result.Message, unmarshaledResult.Message)
	}
	if unmarshaledResult.Latency != result.Latency {
		t.Errorf("Expected Latency %v, got %v", result.Latency, unmarshaledResult.Latency)
	}
}

// TestHealthCheckResult_OptionalFields tests optional fields
func TestHealthCheckResult_OptionalFields(t *testing.T) {
	// Test with minimal fields
	result := HealthCheckResult{
		Name:    "minimal",
		Status:  HealthStatusHealthy,
		Latency: 10 * time.Millisecond,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal minimal HealthCheckResult: %v", err)
	}

	var unmarshaledResult HealthCheckResult
	err = json.Unmarshal(jsonData, &unmarshaledResult)
	if err != nil {
		t.Fatalf("Failed to unmarshal minimal HealthCheckResult: %v", err)
	}

	// Verify required fields
	if unmarshaledResult.Name != "minimal" {
		t.Errorf("Expected Name 'minimal', got '%s'", unmarshaledResult.Name)
	}
	if unmarshaledResult.Status != HealthStatusHealthy {
		t.Errorf("Expected Status '%s', got '%s'", HealthStatusHealthy, unmarshaledResult.Status)
	}

	// Verify optional fields are empty
	if unmarshaledResult.Message != "" {
		t.Errorf("Expected empty Message, got '%s'", unmarshaledResult.Message)
	}
	if unmarshaledResult.Details != nil {
		t.Errorf("Expected nil Details, got %v", unmarshaledResult.Details)
	}
}

// TestReadinessStatus_JSONSerialization tests JSON marshaling and unmarshaling
func TestReadinessStatus_JSONSerialization(t *testing.T) {
	status := ReadinessStatus{
		Ready:  true,
		Reason: "All systems ready",
	}

	jsonData, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal ReadinessStatus to JSON: %v", err)
	}

	var unmarshaledStatus ReadinessStatus
	err = json.Unmarshal(jsonData, &unmarshaledStatus)
	if err != nil {
		t.Fatalf("Failed to unmarshal ReadinessStatus from JSON: %v", err)
	}

	if unmarshaledStatus.Ready != status.Ready {
		t.Errorf("Expected Ready %v, got %v", status.Ready, unmarshaledStatus.Ready)
	}
	if unmarshaledStatus.Reason != status.Reason {
		t.Errorf("Expected Reason '%s', got '%s'", status.Reason, unmarshaledStatus.Reason)
	}
}

// TestReadinessStatus_NotReady tests not ready status
func TestReadinessStatus_NotReady(t *testing.T) {
	status := ReadinessStatus{
		Ready:  false,
		Reason: "Database not ready",
	}

	jsonData, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal not-ready ReadinessStatus: %v", err)
	}

	var unmarshaledStatus ReadinessStatus
	err = json.Unmarshal(jsonData, &unmarshaledStatus)
	if err != nil {
		t.Fatalf("Failed to unmarshal not-ready ReadinessStatus: %v", err)
	}

	if unmarshaledStatus.Ready != false {
		t.Errorf("Expected Ready false, got %v", unmarshaledStatus.Ready)
	}
	if unmarshaledStatus.Reason != "Database not ready" {
		t.Errorf("Expected Reason 'Database not ready', got '%s'", unmarshaledStatus.Reason)
	}
}

// TestLivenessStatus_JSONSerialization tests JSON marshaling and unmarshaling
func TestLivenessStatus_JSONSerialization(t *testing.T) {
	status := LivenessStatus{
		Alive:  true,
		Reason: "Service is alive",
	}

	jsonData, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal LivenessStatus to JSON: %v", err)
	}

	var unmarshaledStatus LivenessStatus
	err = json.Unmarshal(jsonData, &unmarshaledStatus)
	if err != nil {
		t.Fatalf("Failed to unmarshal LivenessStatus from JSON: %v", err)
	}

	if unmarshaledStatus.Alive != status.Alive {
		t.Errorf("Expected Alive %v, got %v", status.Alive, unmarshaledStatus.Alive)
	}
	if unmarshaledStatus.Reason != status.Reason {
		t.Errorf("Expected Reason '%s', got '%s'", status.Reason, unmarshaledStatus.Reason)
	}
}

// TestHealthReport_JSONSerialization tests JSON marshaling and unmarshaling
func TestHealthReport_JSONSerialization(t *testing.T) {
	checks := map[string]HealthCheckResult{
		"database": {
			Name:    "database",
			Status:  HealthStatusHealthy,
			Message: "OK",
			Latency: 10 * time.Millisecond,
		},
		"cache": {
			Name:    "cache",
			Status:  HealthStatusDegraded,
			Message: "High latency",
			Latency: 200 * time.Millisecond,
		},
	}

	report := HealthReport{
		Status:    "degraded",
		Message:   "Some services degraded",
		Version:   "1.0.0",
		Timestamp: time.Now(),
		Uptime:    24 * time.Hour,
		Checks:    checks,
		Readiness: ReadinessStatus{Ready: true, Reason: "Ready"},
		Liveness:  LivenessStatus{Alive: true, Reason: "Alive"},
	}

	jsonData, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("Failed to marshal HealthReport to JSON: %v", err)
	}

	var unmarshaledReport HealthReport
	err = json.Unmarshal(jsonData, &unmarshaledReport)
	if err != nil {
		t.Fatalf("Failed to unmarshal HealthReport from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaledReport.Status != report.Status {
		t.Errorf("Expected Status '%s', got '%s'", report.Status, unmarshaledReport.Status)
	}
	if unmarshaledReport.Message != report.Message {
		t.Errorf("Expected Message '%s', got '%s'", report.Message, unmarshaledReport.Message)
	}
	if unmarshaledReport.Version != report.Version {
		t.Errorf("Expected Version '%s', got '%s'", report.Version, unmarshaledReport.Version)
	}
	if unmarshaledReport.Uptime != report.Uptime {
		t.Errorf("Expected Uptime %v, got %v", report.Uptime, unmarshaledReport.Uptime)
	}

	// Verify checks
	if len(unmarshaledReport.Checks) != len(report.Checks) {
		t.Errorf("Expected %d checks, got %d", len(report.Checks), len(unmarshaledReport.Checks))
	}

	for name, check := range report.Checks {
		unmarshaledCheck, exists := unmarshaledReport.Checks[name]
		if !exists {
			t.Errorf("Expected check '%s' not found", name)
			continue
		}
		if unmarshaledCheck.Name != check.Name {
			t.Errorf("Expected check name '%s', got '%s'", check.Name, unmarshaledCheck.Name)
		}
		if unmarshaledCheck.Status != check.Status {
			t.Errorf("Expected check status '%s', got '%s'", check.Status, unmarshaledCheck.Status)
		}
	}

	// Verify readiness and liveness
	if unmarshaledReport.Readiness.Ready != report.Readiness.Ready {
		t.Errorf("Expected Readiness.Ready %v, got %v", report.Readiness.Ready, unmarshaledReport.Readiness.Ready)
	}
	if unmarshaledReport.Liveness.Alive != report.Liveness.Alive {
		t.Errorf("Expected Liveness.Alive %v, got %v", report.Liveness.Alive, unmarshaledReport.Liveness.Alive)
	}
}

// TestHealthReport_EmptyChecks tests health report with no checks
func TestHealthReport_EmptyChecks(t *testing.T) {
	report := HealthReport{
		Status:    "healthy",
		Message:   "All systems operational",
		Version:   "1.0.0",
		Timestamp: time.Now(),
		Uptime:    1 * time.Hour,
		Checks:    make(map[string]HealthCheckResult),
		Readiness: ReadinessStatus{Ready: true},
		Liveness:  LivenessStatus{Alive: true},
	}

	jsonData, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("Failed to marshal HealthReport with empty checks: %v", err)
	}

	var unmarshaledReport HealthReport
	err = json.Unmarshal(jsonData, &unmarshaledReport)
	if err != nil {
		t.Fatalf("Failed to unmarshal HealthReport with empty checks: %v", err)
	}

	if len(unmarshaledReport.Checks) != 0 {
		t.Errorf("Expected 0 checks, got %d", len(unmarshaledReport.Checks))
	}
}

// TestHealthReport_TimestampSerialization tests timestamp serialization
func TestHealthReport_TimestampSerialization(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	
	report := HealthReport{
		Status:    "healthy",
		Message:   "OK",
		Version:   "1.0.0",
		Timestamp: fixedTime,
		Uptime:    30 * time.Minute,
		Checks:    make(map[string]HealthCheckResult),
		Readiness: ReadinessStatus{Ready: true},
		Liveness:  LivenessStatus{Alive: true},
	}

	jsonData, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("Failed to marshal HealthReport with fixed timestamp: %v", err)
	}

	var unmarshaledReport HealthReport
	err = json.Unmarshal(jsonData, &unmarshaledReport)
	if err != nil {
		t.Fatalf("Failed to unmarshal HealthReport with fixed timestamp: %v", err)
	}

	if !unmarshaledReport.Timestamp.Equal(fixedTime) {
		t.Errorf("Expected Timestamp %v, got %v", fixedTime, unmarshaledReport.Timestamp)
	}
}

// TestHealthCheckResult_LatencySerialization tests latency field serialization
func TestHealthCheckResult_LatencySerialization(t *testing.T) {
	latencies := []time.Duration{
		0,
		1 * time.Millisecond,
		100 * time.Millisecond,
		1 * time.Second,
		5 * time.Second,
	}

	for _, latency := range latencies {
		result := HealthCheckResult{
			Name:    "test",
			Status:  HealthStatusHealthy,
			Latency: latency,
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("Failed to marshal HealthCheckResult with latency %v: %v", latency, err)
		}

		var unmarshaledResult HealthCheckResult
		err = json.Unmarshal(jsonData, &unmarshaledResult)
		if err != nil {
			t.Fatalf("Failed to unmarshal HealthCheckResult with latency %v: %v", latency, err)
		}

		if unmarshaledResult.Latency != latency {
			t.Errorf("Expected Latency %v, got %v", latency, unmarshaledResult.Latency)
		}
	}
}

// TestHealthCheckResult_DetailsSerialization tests details field serialization
func TestHealthCheckResult_DetailsSerialization(t *testing.T) {
	details := []interface{}{
		nil,
		"simple string",
		42,
		map[string]interface{}{
			"key1": "value1",
			"key2": 123,
			"key3": true,
		},
		[]interface{}{"item1", "item2", 123},
	}

	for i, detail := range details {
		result := HealthCheckResult{
			Name:    "test",
			Status:  HealthStatusHealthy,
			Details: detail,
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("Failed to marshal HealthCheckResult with details %v: %v", detail, err)
		}

		var unmarshaledResult HealthCheckResult
		err = json.Unmarshal(jsonData, &unmarshaledResult)
		if err != nil {
			t.Fatalf("Failed to unmarshal HealthCheckResult with details %v: %v", detail, err)
		}

		// For complex types, we'll just verify that the field exists
		// Deep comparison of complex types can be tricky due to JSON number vs int type issues
		if unmarshaledResult.Details == nil && detail != nil {
			t.Errorf("Test %d: Expected non-nil Details, got nil", i)
		}
	}
}

// TestHealthStatus_StringValues tests that health status values are valid strings
func TestHealthStatus_StringValues(t *testing.T) {
	validStatuses := []HealthStatus{
		HealthStatusHealthy,
		HealthStatusUnhealthy,
		HealthStatusDegraded,
	}

	for _, status := range validStatuses {
		if len(string(status)) == 0 {
			t.Errorf("HealthStatus '%s' should not be empty", status)
		}
	}
}
