package model

import (
	"encoding/json"
	"testing"
)

// TestCreateUserRequest_JSONSerialization tests JSON marshaling and unmarshaling
func TestCreateUserRequest_JSONSerialization(t *testing.T) {
	request := CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Role:  "developer",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal CreateUserRequest to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledRequest CreateUserRequest
	err = json.Unmarshal(jsonData, &unmarshaledRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal CreateUserRequest from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaledRequest.Name != request.Name {
		t.Errorf("Expected Name '%s', got '%s'", request.Name, unmarshaledRequest.Name)
	}
	if unmarshaledRequest.Email != request.Email {
		t.Errorf("Expected Email '%s', got '%s'", request.Email, unmarshaledRequest.Email)
	}
	if unmarshaledRequest.Role != request.Role {
		t.Errorf("Expected Role '%s', got '%s'", request.Role, unmarshaledRequest.Role)
	}
}

// TestCreateUserRequest_JSONFields tests that JSON field names are correct
func TestCreateUserRequest_JSONFields(t *testing.T) {
	request := CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "tester",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal CreateUserRequest to JSON: %v", err)
	}

	// Parse JSON to verify field names
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to map: %v", err)
	}

	// Verify expected field names exist
	expectedFields := []string{"name", "email", "role"}
	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("Expected field '%s' not found in JSON", field)
		}
	}

	// Verify field values
	if result["name"] != "Test User" {
		t.Errorf("Expected name 'Test User', got %v", result["name"])
	}
	if result["email"] != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %v", result["email"])
	}
	if result["role"] != "tester" {
		t.Errorf("Expected role 'tester', got %v", result["role"])
	}
}

// TestCreateUserRequest_EmptyFields tests empty field handling
func TestCreateUserRequest_EmptyFields(t *testing.T) {
	request := CreateUserRequest{
		Name:  "",
		Email: "",
		Role:  "",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal empty CreateUserRequest: %v", err)
	}

	var unmarshaledRequest CreateUserRequest
	err = json.Unmarshal(jsonData, &unmarshaledRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty CreateUserRequest: %v", err)
	}

	// Verify empty fields are preserved
	if unmarshaledRequest.Name != "" {
		t.Errorf("Expected empty Name, got '%s'", unmarshaledRequest.Name)
	}
	if unmarshaledRequest.Email != "" {
		t.Errorf("Expected empty Email, got '%s'", unmarshaledRequest.Email)
	}
	if unmarshaledRequest.Role != "" {
		t.Errorf("Expected empty Role, got '%s'", unmarshaledRequest.Role)
	}
}

// TestCreateTaskRequest_JSONSerialization tests JSON marshaling and unmarshaling
func TestCreateTaskRequest_JSONSerialization(t *testing.T) {
	request := CreateTaskRequest{
		Title:  "Implement authentication",
		Status: "pending",
		UserID: 1,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal CreateTaskRequest to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledRequest CreateTaskRequest
	err = json.Unmarshal(jsonData, &unmarshaledRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal CreateTaskRequest from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaledRequest.Title != request.Title {
		t.Errorf("Expected Title '%s', got '%s'", request.Title, unmarshaledRequest.Title)
	}
	if unmarshaledRequest.Status != request.Status {
		t.Errorf("Expected Status '%s', got '%s'", request.Status, unmarshaledRequest.Status)
	}
	if unmarshaledRequest.UserID != request.UserID {
		t.Errorf("Expected UserID %d, got %d", request.UserID, unmarshaledRequest.UserID)
	}
}

// TestCreateTaskRequest_JSONFields tests that JSON field names are correct
func TestCreateTaskRequest_JSONFields(t *testing.T) {
	request := CreateTaskRequest{
		Title:  "Test Task",
		Status: "in-progress",
		UserID: 123,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal CreateTaskRequest to JSON: %v", err)
	}

	// Parse JSON to verify field names
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to map: %v", err)
	}

	// Verify expected field names exist
	expectedFields := []string{"title", "status", "userId"}
	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("Expected field '%s' not found in JSON", field)
		}
	}

	// Verify field values
	if result["title"] != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %v", result["title"])
	}
	if result["status"] != "in-progress" {
		t.Errorf("Expected status 'in-progress', got %v", result["status"])
	}
	if result["userId"] != float64(123) {
		t.Errorf("Expected userId 123, got %v", result["userId"])
	}
}

// TestUpdateTaskRequest_JSONSerialization tests JSON marshaling and unmarshaling
func TestUpdateTaskRequest_JSONSerialization(t *testing.T) {
	title := "Updated Title"
	status := "completed"
	userID := 2

	request := UpdateTaskRequest{
		Title:  &title,
		Status: &status,
		UserID: &userID,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal UpdateTaskRequest to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledRequest UpdateTaskRequest
	err = json.Unmarshal(jsonData, &unmarshaledRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal UpdateTaskRequest from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaledRequest.Title == nil || *unmarshaledRequest.Title != title {
		t.Errorf("Expected Title '%s', got %v", title, unmarshaledRequest.Title)
	}
	if unmarshaledRequest.Status == nil || *unmarshaledRequest.Status != status {
		t.Errorf("Expected Status '%s', got %v", status, unmarshaledRequest.Status)
	}
	if unmarshaledRequest.UserID == nil || *unmarshaledRequest.UserID != userID {
		t.Errorf("Expected UserID %d, got %v", userID, unmarshaledRequest.UserID)
	}
}

// TestUpdateTaskRequest_PartialUpdates tests partial update scenarios
func TestUpdateTaskRequest_PartialUpdates(t *testing.T) {
	// Test with only title
	title := "Only Title"
	request := UpdateTaskRequest{
		Title: &title,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal partial UpdateTaskRequest: %v", err)
	}

	var unmarshaledRequest UpdateTaskRequest
	err = json.Unmarshal(jsonData, &unmarshaledRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal partial UpdateTaskRequest: %v", err)
	}

	if unmarshaledRequest.Title == nil || *unmarshaledRequest.Title != title {
		t.Errorf("Expected Title '%s', got %v", title, unmarshaledRequest.Title)
	}
	if unmarshaledRequest.Status != nil {
		t.Errorf("Expected nil Status, got %v", unmarshaledRequest.Status)
	}
	if unmarshaledRequest.UserID != nil {
		t.Errorf("Expected nil UserID, got %v", unmarshaledRequest.UserID)
	}

	// Test with only status
	status := "in-progress"
	request = UpdateTaskRequest{
		Status: &status,
	}

	jsonData, err = json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal status-only UpdateTaskRequest: %v", err)
	}

	unmarshaledRequest = UpdateTaskRequest{} // reset for clean unmarshal
	err = json.Unmarshal(jsonData, &unmarshaledRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal status-only UpdateTaskRequest: %v", err)
	}

	if unmarshaledRequest.Title != nil {
		t.Errorf("Expected nil Title, got %v", unmarshaledRequest.Title)
	}
	if unmarshaledRequest.Status == nil || *unmarshaledRequest.Status != status {
		t.Errorf("Expected Status '%s', got %v", status, unmarshaledRequest.Status)
	}
	if unmarshaledRequest.UserID != nil {
		t.Errorf("Expected nil UserID, got %v", unmarshaledRequest.UserID)
	}

	// Test with only UserID
	userID := 5
	request = UpdateTaskRequest{
		UserID: &userID,
	}

	jsonData, err = json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal userID-only UpdateTaskRequest: %v", err)
	}

	unmarshaledRequest = UpdateTaskRequest{} // reset for clean unmarshal
	err = json.Unmarshal(jsonData, &unmarshaledRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal userID-only UpdateTaskRequest: %v", err)
	}

	if unmarshaledRequest.Title != nil {
		t.Errorf("Expected nil Title, got %v", unmarshaledRequest.Title)
	}
	if unmarshaledRequest.Status != nil {
		t.Errorf("Expected nil Status, got %v", unmarshaledRequest.Status)
	}
	if unmarshaledRequest.UserID == nil || *unmarshaledRequest.UserID != userID {
		t.Errorf("Expected UserID %d, got %v", userID, unmarshaledRequest.UserID)
	}
}

// TestUpdateTaskRequest_EmptyRequest tests empty update request
func TestUpdateTaskRequest_EmptyRequest(t *testing.T) {
	request := UpdateTaskRequest{}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal empty UpdateTaskRequest: %v", err)
	}

	var unmarshaledRequest UpdateTaskRequest
	err = json.Unmarshal(jsonData, &unmarshaledRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty UpdateTaskRequest: %v", err)
	}

	// All fields should be nil
	if unmarshaledRequest.Title != nil {
		t.Errorf("Expected nil Title, got %v", unmarshaledRequest.Title)
	}
	if unmarshaledRequest.Status != nil {
		t.Errorf("Expected nil Status, got %v", unmarshaledRequest.Status)
	}
	if unmarshaledRequest.UserID != nil {
		t.Errorf("Expected nil UserID, got %v", unmarshaledRequest.UserID)
	}
}

// TestUpdateTaskRequest_ZeroValues tests zero values in pointers
func TestUpdateTaskRequest_ZeroValues(t *testing.T) {
	title := ""
	status := ""
	userID := 0

	request := UpdateTaskRequest{
		Title:  &title,
		Status: &status,
		UserID: &userID,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal zero-value UpdateTaskRequest: %v", err)
	}

	var unmarshaledRequest UpdateTaskRequest
	err = json.Unmarshal(jsonData, &unmarshaledRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal zero-value UpdateTaskRequest: %v", err)
	}

	// Verify zero values are preserved
	if unmarshaledRequest.Title == nil || *unmarshaledRequest.Title != "" {
		t.Errorf("Expected empty Title, got %v", unmarshaledRequest.Title)
	}
	if unmarshaledRequest.Status == nil || *unmarshaledRequest.Status != "" {
		t.Errorf("Expected empty Status, got %v", unmarshaledRequest.Status)
	}
	if unmarshaledRequest.UserID == nil || *unmarshaledRequest.UserID != 0 {
		t.Errorf("Expected UserID 0, got %v", unmarshaledRequest.UserID)
	}
}

// TestCreateUserRequest_ValidationScenarios tests various validation scenarios
func TestCreateUserRequest_ValidationScenarios(t *testing.T) {
	scenarios := []struct {
		name  string
		request CreateUserRequest
		valid bool
	}{
		{
			name: "Valid user",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  "developer",
			},
			valid: true,
		},
		{
			name: "Empty name",
			request: CreateUserRequest{
				Name:  "",
				Email: "john@example.com",
				Role:  "developer",
			},
			valid: false,
		},
		{
			name: "Empty email",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "",
				Role:  "developer",
			},
			valid: false,
		},
		{
			name: "Empty role",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  "",
			},
			valid: false,
		},
		{
			name: "Special characters",
			request: CreateUserRequest{
				Name:  "John \"The Rock\" Doe",
				Email: "john+special@example.com",
				Role:  "admin/super-user",
			},
			valid: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Test that the request can be marshaled and unmarshaled
			jsonData, err := json.Marshal(scenario.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			var unmarshaledRequest CreateUserRequest
			err = json.Unmarshal(jsonData, &unmarshaledRequest)
			if err != nil {
				t.Fatalf("Failed to unmarshal request: %v", err)
			}

			// Verify fields match
			if unmarshaledRequest.Name != scenario.request.Name {
				t.Errorf("Expected Name '%s', got '%s'", scenario.request.Name, unmarshaledRequest.Name)
			}
			if unmarshaledRequest.Email != scenario.request.Email {
				t.Errorf("Expected Email '%s', got '%s'", scenario.request.Email, unmarshaledRequest.Email)
			}
			if unmarshaledRequest.Role != scenario.request.Role {
				t.Errorf("Expected Role '%s', got '%s'", scenario.request.Role, unmarshaledRequest.Role)
			}
		})
	}
}

// TestCreateTaskRequest_ValidationScenarios tests various validation scenarios
func TestCreateTaskRequest_ValidationScenarios(t *testing.T) {
	scenarios := []struct {
		name  string
		request CreateTaskRequest
		valid bool
	}{
		{
			name: "Valid task",
			request: CreateTaskRequest{
				Title:  "Implement authentication",
				Status: "pending",
				UserID: 1,
			},
			valid: true,
		},
		{
			name: "Empty title",
			request: CreateTaskRequest{
				Title:  "",
				Status: "pending",
				UserID: 1,
			},
			valid: false,
		},
		{
			name: "Empty status",
			request: CreateTaskRequest{
				Title:  "Implement authentication",
				Status: "",
				UserID: 1,
			},
			valid: false,
		},
		{
			name: "Zero UserID",
			request: CreateTaskRequest{
				Title:  "Implement authentication",
				Status: "pending",
				UserID: 0,
			},
			valid: false,
		},
		{
			name: "Negative UserID",
			request: CreateTaskRequest{
				Title:  "Implement authentication",
				Status: "pending",
				UserID: -1,
			},
			valid: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Test that the request can be marshaled and unmarshaled
			jsonData, err := json.Marshal(scenario.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			var unmarshaledRequest CreateTaskRequest
			err = json.Unmarshal(jsonData, &unmarshaledRequest)
			if err != nil {
				t.Fatalf("Failed to unmarshal request: %v", err)
			}

			// Verify fields match
			if unmarshaledRequest.Title != scenario.request.Title {
				t.Errorf("Expected Title '%s', got '%s'", scenario.request.Title, unmarshaledRequest.Title)
			}
			if unmarshaledRequest.Status != scenario.request.Status {
				t.Errorf("Expected Status '%s', got '%s'", scenario.request.Status, unmarshaledRequest.Status)
			}
			if unmarshaledRequest.UserID != scenario.request.UserID {
				t.Errorf("Expected UserID %d, got %d", scenario.request.UserID, unmarshaledRequest.UserID)
			}
		})
	}
}

// TestRequestModels_SpecialCharacters tests special characters in all request models
func TestRequestModels_SpecialCharacters(t *testing.T) {
	specialString := "Special \"Characters\" \nNew Line\tTab"

	// Test CreateUserRequest
	userRequest := CreateUserRequest{
		Name:  specialString,
		Email: "special+chars@example.com",
		Role:  specialString,
	}

	jsonData, err := json.Marshal(userRequest)
	if err != nil {
		t.Fatalf("Failed to marshal CreateUserRequest with special chars: %v", err)
	}

	var unmarshaledUserRequest CreateUserRequest
	err = json.Unmarshal(jsonData, &unmarshaledUserRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal CreateUserRequest with special chars: %v", err)
	}

	if unmarshaledUserRequest.Name != specialString {
		t.Errorf("Expected Name with special chars, got '%s'", unmarshaledUserRequest.Name)
	}

	// Test CreateTaskRequest
	taskRequest := CreateTaskRequest{
		Title:  specialString,
		Status: specialString,
		UserID: 1,
	}

	jsonData, err = json.Marshal(taskRequest)
	if err != nil {
		t.Fatalf("Failed to marshal CreateTaskRequest with special chars: %v", err)
	}

	var unmarshaledTaskRequest CreateTaskRequest
	err = json.Unmarshal(jsonData, &unmarshaledTaskRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal CreateTaskRequest with special chars: %v", err)
	}

	if unmarshaledTaskRequest.Title != specialString {
		t.Errorf("Expected Title with special chars, got '%s'", unmarshaledTaskRequest.Title)
	}

	// Test UpdateTaskRequest
	title := specialString
	updateRequest := UpdateTaskRequest{
		Title: &title,
	}

	jsonData, err = json.Marshal(updateRequest)
	if err != nil {
		t.Fatalf("Failed to marshal UpdateTaskRequest with special chars: %v", err)
	}

	var unmarshaledUpdateRequest UpdateTaskRequest
	err = json.Unmarshal(jsonData, &unmarshaledUpdateRequest)
	if err != nil {
		t.Fatalf("Failed to unmarshal UpdateTaskRequest with special chars: %v", err)
	}

	if unmarshaledUpdateRequest.Title == nil || *unmarshaledUpdateRequest.Title != specialString {
		t.Errorf("Expected Title with special chars, got %v", unmarshaledUpdateRequest.Title)
	}
}
