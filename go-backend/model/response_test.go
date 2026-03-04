package model

import (
	"encoding/json"
	"testing"
)

// TestSuccessResponse_JSONSerialization tests JSON marshaling and unmarshaling
func TestSuccessResponse_JSONSerialization(t *testing.T) {
	response := SuccessResponse{
		Success: true,
		Data: map[string]interface{}{
			"id":   1,
			"name": "Test User",
		},
		Message: "Operation successful",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal SuccessResponse to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledResponse SuccessResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal SuccessResponse from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaledResponse.Success != response.Success {
		t.Errorf("Expected Success %v, got %v", response.Success, unmarshaledResponse.Success)
	}
	if unmarshaledResponse.Message != response.Message {
		t.Errorf("Expected Message '%s', got '%s'", response.Message, unmarshaledResponse.Message)
	}

	// For Data field, we need to handle the interface{} carefully
	dataMap, ok := unmarshaledResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected Data to be a map, got %T", unmarshaledResponse.Data)
	}

	if dataMap["id"] != float64(1) {
		t.Errorf("Expected Data.id 1, got %v", dataMap["id"])
	}
	if dataMap["name"] != "Test User" {
		t.Errorf("Expected Data.name 'Test User', got %v", dataMap["name"])
	}
}

// TestSuccessResponse_OptionalFields tests optional fields
func TestSuccessResponse_OptionalFields(t *testing.T) {
	// Test with minimal fields
	response := SuccessResponse{
		Success: true,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal minimal SuccessResponse: %v", err)
	}

	var unmarshaledResponse SuccessResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal minimal SuccessResponse: %v", err)
	}

	if unmarshaledResponse.Success != true {
		t.Errorf("Expected Success true, got %v", unmarshaledResponse.Success)
	}
	if unmarshaledResponse.Data != nil {
		t.Errorf("Expected nil Data, got %v", unmarshaledResponse.Data)
	}
	if unmarshaledResponse.Message != "" {
		t.Errorf("Expected empty Message, got '%s'", unmarshaledResponse.Message)
	}
}

// TestErrorResponse_JSONSerialization tests JSON marshaling and unmarshaling
func TestErrorResponse_JSONSerialization(t *testing.T) {
	response := ErrorResponse{
		Error:   "Something went wrong",
		Code:    "INTERNAL_ERROR",
		Details: "Database connection failed",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorResponse to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledResponse ErrorResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal ErrorResponse from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaledResponse.Error != response.Error {
		t.Errorf("Expected Error '%s', got '%s'", response.Error, unmarshaledResponse.Error)
	}
	if unmarshaledResponse.Code != response.Code {
		t.Errorf("Expected Code '%s', got '%s'", response.Code, unmarshaledResponse.Code)
	}
	if unmarshaledResponse.Details != response.Details {
		t.Errorf("Expected Details '%s', got '%s'", response.Details, unmarshaledResponse.Details)
	}
}

// TestErrorResponse_OptionalFields tests optional fields
func TestErrorResponse_OptionalFields(t *testing.T) {
	// Test with minimal fields
	response := ErrorResponse{
		Error: "Basic error",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal minimal ErrorResponse: %v", err)
	}

	var unmarshaledResponse ErrorResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal minimal ErrorResponse: %v", err)
	}

	if unmarshaledResponse.Error != "Basic error" {
		t.Errorf("Expected Error 'Basic error', got '%s'", unmarshaledResponse.Error)
	}
	if unmarshaledResponse.Code != "" {
		t.Errorf("Expected empty Code, got '%s'", unmarshaledResponse.Code)
	}
	if unmarshaledResponse.Details != "" {
		t.Errorf("Expected empty Details, got '%s'", unmarshaledResponse.Details)
	}
}

// TestValidationError_JSONSerialization tests JSON marshaling and unmarshaling
func TestValidationError_JSONSerialization(t *testing.T) {
	response := ValidationError{
		Error: "Validation failed",
		Fields: map[string]string{
			"name":  "Name is required",
			"email": "Invalid email format",
		},
		Code: "VALIDATION_ERROR",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal ValidationError to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledResponse ValidationError
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal ValidationError from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaledResponse.Error != response.Error {
		t.Errorf("Expected Error '%s', got '%s'", response.Error, unmarshaledResponse.Error)
	}
	if unmarshaledResponse.Code != response.Code {
		t.Errorf("Expected Code '%s', got '%s'", response.Code, unmarshaledResponse.Code)
	}

	// Verify Fields
	if len(unmarshaledResponse.Fields) != len(response.Fields) {
		t.Errorf("Expected %d fields, got %d", len(response.Fields), len(unmarshaledResponse.Fields))
	}

	for field, message := range response.Fields {
		unmarshaledMessage, exists := unmarshaledResponse.Fields[field]
		if !exists {
			t.Errorf("Expected field '%s' not found", field)
			continue
		}
		if unmarshaledMessage != message {
			t.Errorf("Expected field '%s' message '%s', got '%s'", field, message, unmarshaledMessage)
		}
	}
}

// TestValidationError_EmptyFields tests empty fields
func TestValidationError_EmptyFields(t *testing.T) {
	response := ValidationError{
		Error: "Validation failed",
		Code:  "VALIDATION_ERROR",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal ValidationError with empty fields: %v", err)
	}

	var unmarshaledResponse ValidationError
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal ValidationError with empty fields: %v", err)
	}

	if unmarshaledResponse.Error != "Validation failed" {
		t.Errorf("Expected Error 'Validation failed', got '%s'", unmarshaledResponse.Error)
	}
	if unmarshaledResponse.Code != "VALIDATION_ERROR" {
		t.Errorf("Expected Code 'VALIDATION_ERROR', got '%s'", unmarshaledResponse.Code)
	}
	if unmarshaledResponse.Fields != nil {
		t.Errorf("Expected nil Fields, got %v", unmarshaledResponse.Fields)
	}
}

// TestUsersResponse_JSONSerialization tests JSON marshaling and unmarshaling
func TestUsersResponse_JSONSerialization(t *testing.T) {
	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Role: "developer"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Role: "designer"},
	}

	response := UsersResponse{
		Users: users,
		Count: len(users),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal UsersResponse to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledResponse UsersResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal UsersResponse from JSON: %v", err)
	}

	// Verify fields match
	if unmarshaledResponse.Count != response.Count {
		t.Errorf("Expected Count %d, got %d", response.Count, unmarshaledResponse.Count)
	}

	if len(unmarshaledResponse.Users) != len(response.Users) {
		t.Errorf("Expected %d users, got %d", len(response.Users), len(unmarshaledResponse.Users))
	}

	for i, user := range response.Users {
		if i >= len(unmarshaledResponse.Users) {
			t.Errorf("User %d not found in unmarshaled response", i)
			continue
		}
		unmarshaledUser := unmarshaledResponse.Users[i]
		if unmarshaledUser.ID != user.ID {
			t.Errorf("Expected user %d ID %d, got %d", i, user.ID, unmarshaledUser.ID)
		}
		if unmarshaledUser.Name != user.Name {
			t.Errorf("Expected user %d Name '%s', got '%s'", i, user.Name, unmarshaledUser.Name)
		}
		if unmarshaledUser.Email != user.Email {
			t.Errorf("Expected user %d Email '%s', got '%s'", i, user.Email, unmarshaledUser.Email)
		}
		if unmarshaledUser.Role != user.Role {
			t.Errorf("Expected user %d Role '%s', got '%s'", i, user.Role, unmarshaledUser.Role)
		}
	}
}

// TestUsersResponse_EmptyUsers tests empty users list
func TestUsersResponse_EmptyUsers(t *testing.T) {
	response := UsersResponse{
		Users: []User{},
		Count: 0,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal empty UsersResponse: %v", err)
	}

	var unmarshaledResponse UsersResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty UsersResponse: %v", err)
	}

	if unmarshaledResponse.Count != 0 {
		t.Errorf("Expected Count 0, got %d", unmarshaledResponse.Count)
	}
	if len(unmarshaledResponse.Users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(unmarshaledResponse.Users))
	}
}

// TestTasksResponse_JSONSerialization tests JSON marshaling and unmarshaling
func TestTasksResponse_JSONSerialization(t *testing.T) {
	tasks := []Task{
		{ID: 1, Title: "Task 1", Status: "pending", UserID: 1},
		{ID: 2, Title: "Task 2", Status: "completed", UserID: 2},
	}

	response := TasksResponse{
		Tasks: tasks,
		Count: len(tasks),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal TasksResponse to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledResponse TasksResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal TasksResponse from JSON: %v", err)
	}

	// Verify fields match
	if unmarshaledResponse.Count != response.Count {
		t.Errorf("Expected Count %d, got %d", response.Count, unmarshaledResponse.Count)
	}

	if len(unmarshaledResponse.Tasks) != len(response.Tasks) {
		t.Errorf("Expected %d tasks, got %d", len(response.Tasks), len(unmarshaledResponse.Tasks))
	}

	for i, task := range response.Tasks {
		if i >= len(unmarshaledResponse.Tasks) {
			t.Errorf("Task %d not found in unmarshaled response", i)
			continue
		}
		unmarshaledTask := unmarshaledResponse.Tasks[i]
		if unmarshaledTask.ID != task.ID {
			t.Errorf("Expected task %d ID %d, got %d", i, task.ID, unmarshaledTask.ID)
		}
		if unmarshaledTask.Title != task.Title {
			t.Errorf("Expected task %d Title '%s', got '%s'", i, task.Title, unmarshaledTask.Title)
		}
		if unmarshaledTask.Status != task.Status {
			t.Errorf("Expected task %d Status '%s', got '%s'", i, task.Status, unmarshaledTask.Status)
		}
		if unmarshaledTask.UserID != task.UserID {
			t.Errorf("Expected task %d UserID %d, got %d", i, task.UserID, unmarshaledTask.UserID)
		}
	}
}

// TestStatsResponse_JSONSerialization tests JSON marshaling and unmarshaling
func TestStatsResponse_JSONSerialization(t *testing.T) {
	response := StatsResponse{
		Users: struct {
			Total int `json:"total"`
		}{
			Total: 10,
		},
		Tasks: struct {
			Total      int `json:"total"`
			Pending    int `json:"pending"`
			InProgress int `json:"inProgress"`
			Completed  int `json:"completed"`
		}{
			Total:      25,
			Pending:    5,
			InProgress: 10,
			Completed:  10,
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal StatsResponse to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledResponse StatsResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal StatsResponse from JSON: %v", err)
	}

	// Verify fields match
	if unmarshaledResponse.Users.Total != response.Users.Total {
		t.Errorf("Expected Users.Total %d, got %d", response.Users.Total, unmarshaledResponse.Users.Total)
	}
	if unmarshaledResponse.Tasks.Total != response.Tasks.Total {
		t.Errorf("Expected Tasks.Total %d, got %d", response.Tasks.Total, unmarshaledResponse.Tasks.Total)
	}
	if unmarshaledResponse.Tasks.Pending != response.Tasks.Pending {
		t.Errorf("Expected Tasks.Pending %d, got %d", response.Tasks.Pending, unmarshaledResponse.Tasks.Pending)
	}
	if unmarshaledResponse.Tasks.InProgress != response.Tasks.InProgress {
		t.Errorf("Expected Tasks.InProgress %d, got %d", response.Tasks.InProgress, unmarshaledResponse.Tasks.InProgress)
	}
	if unmarshaledResponse.Tasks.Completed != response.Tasks.Completed {
		t.Errorf("Expected Tasks.Completed %d, got %d", response.Tasks.Completed, unmarshaledResponse.Tasks.Completed)
	}
}

// TestCacheStats_JSONSerialization tests JSON marshaling and unmarshaling
func TestCacheStats_JSONSerialization(t *testing.T) {
	response := CacheStats{
		Hits:         1000,
		Misses:       200,
		Evictions:    50,
		TotalEntries: 150,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal CacheStats to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledResponse CacheStats
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal CacheStats from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaledResponse.Hits != response.Hits {
		t.Errorf("Expected Hits %d, got %d", response.Hits, unmarshaledResponse.Hits)
	}
	if unmarshaledResponse.Misses != response.Misses {
		t.Errorf("Expected Misses %d, got %d", response.Misses, unmarshaledResponse.Misses)
	}
	if unmarshaledResponse.Evictions != response.Evictions {
		t.Errorf("Expected Evictions %d, got %d", response.Evictions, unmarshaledResponse.Evictions)
	}
	if unmarshaledResponse.TotalEntries != response.TotalEntries {
		t.Errorf("Expected TotalEntries %d, got %d", response.TotalEntries, unmarshaledResponse.TotalEntries)
	}
}

// TestCacheStats_ZeroValues tests zero values
func TestCacheStats_ZeroValues(t *testing.T) {
	response := CacheStats{}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal zero-value CacheStats: %v", err)
	}

	var unmarshaledResponse CacheStats
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal zero-value CacheStats: %v", err)
	}

	if unmarshaledResponse.Hits != 0 {
		t.Errorf("Expected Hits 0, got %d", unmarshaledResponse.Hits)
	}
	if unmarshaledResponse.Misses != 0 {
		t.Errorf("Expected Misses 0, got %d", unmarshaledResponse.Misses)
	}
	if unmarshaledResponse.Evictions != 0 {
		t.Errorf("Expected Evictions 0, got %d", unmarshaledResponse.Evictions)
	}
	if unmarshaledResponse.TotalEntries != 0 {
		t.Errorf("Expected TotalEntries 0, got %d", unmarshaledResponse.TotalEntries)
	}
}

// TestResponseModels_JSONFields tests that JSON field names are correct
func TestResponseModels_JSONFields(t *testing.T) {
	// Test SuccessResponse
	successResp := SuccessResponse{
		Success: true,
		Data:    "test data",
		Message: "success",
	}

	jsonData, err := json.Marshal(successResp)
	if err != nil {
		t.Fatalf("Failed to marshal SuccessResponse: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal SuccessResponse to map: %v", err)
	}

	expectedFields := []string{"success", "data", "message"}
	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("SuccessResponse: Expected field '%s' not found", field)
		}
	}

	// Test ErrorResponse
	errorResp := ErrorResponse{
		Error:   "test error",
		Code:    "TEST_ERROR",
		Details: "test details",
	}

	jsonData, err = json.Marshal(errorResp)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorResponse: %v", err)
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal ErrorResponse to map: %v", err)
	}

	expectedFields = []string{"error", "code", "details"}
	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("ErrorResponse: Expected field '%s' not found", field)
		}
	}

	// Test UsersResponse
	usersResp := UsersResponse{
		Users: []User{{ID: 1, Name: "Test", Email: "test@example.com", Role: "tester"}},
		Count: 1,
	}

	jsonData, err = json.Marshal(usersResp)
	if err != nil {
		t.Fatalf("Failed to marshal UsersResponse: %v", err)
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal UsersResponse to map: %v", err)
	}

	expectedFields = []string{"users", "count"}
	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("UsersResponse: Expected field '%s' not found", field)
		}
	}

	// Test TasksResponse
	tasksResp := TasksResponse{
		Tasks: []Task{{ID: 1, Title: "Test", Status: "pending", UserID: 1}},
		Count: 1,
	}

	jsonData, err = json.Marshal(tasksResp)
	if err != nil {
		t.Fatalf("Failed to marshal TasksResponse: %v", err)
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal TasksResponse to map: %v", err)
	}

	expectedFields = []string{"tasks", "count"}
	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("TasksResponse: Expected field '%s' not found", field)
		}
	}
}

// TestResponseModels_LargeData tests handling of large data sets
func TestResponseModels_LargeData(t *testing.T) {
	// Create large user list
	users := make([]User, 100)
	for i := 0; i < 100; i++ {
		users[i] = User{
			ID:    i + 1,
			Name:  "User " + string(rune('A'+i%26)),
			Email: "user@example.com",
			Role:  "tester",
		}
	}

	usersResp := UsersResponse{
		Users: users,
		Count: len(users),
	}

	jsonData, err := json.Marshal(usersResp)
	if err != nil {
		t.Fatalf("Failed to marshal large UsersResponse: %v", err)
	}

	var unmarshaledResponse UsersResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal large UsersResponse: %v", err)
	}

	if unmarshaledResponse.Count != 100 {
		t.Errorf("Expected Count 100, got %d", unmarshaledResponse.Count)
	}
	if len(unmarshaledResponse.Users) != 100 {
		t.Errorf("Expected 100 users, got %d", len(unmarshaledResponse.Users))
	}

	// Create large task list
	tasks := make([]Task, 50)
	for i := 0; i < 50; i++ {
		tasks[i] = Task{
			ID:     i + 1,
			Title:  "Task " + string(rune('A'+i%26)),
			Status: "pending",
			UserID: (i % 5) + 1,
		}
	}

	tasksResp := TasksResponse{
		Tasks: tasks,
		Count: len(tasks),
	}

	jsonData, err = json.Marshal(tasksResp)
	if err != nil {
		t.Fatalf("Failed to marshal large TasksResponse: %v", err)
	}

	var unmarshaledTasksResponse TasksResponse
	err = json.Unmarshal(jsonData, &unmarshaledTasksResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal large TasksResponse: %v", err)
	}

	if unmarshaledTasksResponse.Count != 50 {
		t.Errorf("Expected Count 50, got %d", unmarshaledTasksResponse.Count)
	}
	if len(unmarshaledTasksResponse.Tasks) != 50 {
		t.Errorf("Expected 50 tasks, got %d", len(unmarshaledTasksResponse.Tasks))
	}
}
