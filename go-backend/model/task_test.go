package model

import (
	"encoding/json"
	"testing"
)

// TestTask_JSONSerialization tests JSON marshaling and unmarshaling
func TestTask_JSONSerialization(t *testing.T) {
	task := Task{
		ID:     1,
		Title:  "Implement authentication",
		Status: "pending",
		UserID: 1,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledTask Task
	err = json.Unmarshal(jsonData, &unmarshaledTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal task from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaledTask.ID != task.ID {
		t.Errorf("Expected ID %d, got %d", task.ID, unmarshaledTask.ID)
	}
	if unmarshaledTask.Title != task.Title {
		t.Errorf("Expected Title '%s', got '%s'", task.Title, unmarshaledTask.Title)
	}
	if unmarshaledTask.Status != task.Status {
		t.Errorf("Expected Status '%s', got '%s'", task.Status, unmarshaledTask.Status)
	}
	if unmarshaledTask.UserID != task.UserID {
		t.Errorf("Expected UserID %d, got %d", task.UserID, unmarshaledTask.UserID)
	}
}

// TestTask_JSONFields tests that JSON field names are correct
func TestTask_JSONFields(t *testing.T) {
	task := Task{
		ID:     456,
		Title:  "Test Task",
		Status: "in-progress",
		UserID: 789,
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task to JSON: %v", err)
	}

	// Parse JSON to verify field names
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to map: %v", err)
	}

	// Verify expected field names exist
	expectedFields := []string{"id", "title", "status", "userId"}
	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("Expected field '%s' not found in JSON", field)
		}
	}

	// Verify field values
	if result["id"] != float64(456) {
		t.Errorf("Expected id 456, got %v", result["id"])
	}
	if result["title"] != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %v", result["title"])
	}
	if result["status"] != "in-progress" {
		t.Errorf("Expected status 'in-progress', got %v", result["status"])
	}
	if result["userId"] != float64(789) {
		t.Errorf("Expected userId 789, got %v", result["userId"])
	}
}

// TestTask_ZeroValues tests task with zero values
func TestTask_ZeroValues(t *testing.T) {
	task := Task{}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal zero-value task: %v", err)
	}

	var unmarshaledTask Task
	err = json.Unmarshal(jsonData, &unmarshaledTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal zero-value task: %v", err)
	}

	// Verify zero values are preserved
	if unmarshaledTask.ID != 0 {
		t.Errorf("Expected ID 0, got %d", unmarshaledTask.ID)
	}
	if unmarshaledTask.Title != "" {
		t.Errorf("Expected empty Title, got '%s'", unmarshaledTask.Title)
	}
	if unmarshaledTask.Status != "" {
		t.Errorf("Expected empty Status, got '%s'", unmarshaledTask.Status)
	}
	if unmarshaledTask.UserID != 0 {
		t.Errorf("Expected UserID 0, got %d", unmarshaledTask.UserID)
	}
}

// TestTask_PartialJSON tests partial JSON data
func TestTask_PartialJSON(t *testing.T) {
	// Test JSON with missing fields
	jsonData := []byte(`{"id": 1, "title": "Partial Task"}`)

	var task Task
	err := json.Unmarshal(jsonData, &task)
	if err != nil {
		t.Fatalf("Failed to unmarshal partial JSON: %v", err)
	}

	// Verify populated fields
	if task.ID != 1 {
		t.Errorf("Expected ID 1, got %d", task.ID)
	}
	if task.Title != "Partial Task" {
		t.Errorf("Expected Title 'Partial Task', got '%s'", task.Title)
	}

	// Verify missing fields are zero values
	if task.Status != "" {
		t.Errorf("Expected empty Status, got '%s'", task.Status)
	}
	if task.UserID != 0 {
		t.Errorf("Expected UserID 0, got %d", task.UserID)
	}
}

// TestTask_InvalidJSON tests invalid JSON handling
func TestTask_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"id": "not-a-number", "title": 123}`)

	var task Task
	err := json.Unmarshal(invalidJSON, &task)
	if err == nil {
		t.Error("Expected error when unmarshaling invalid JSON")
	}
}

// TestTask_EmptyJSON tests empty JSON handling
func TestTask_EmptyJSON(t *testing.T) {
	emptyJSON := []byte(`{}`)

	var task Task
	err := json.Unmarshal(emptyJSON, &task)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty JSON: %v", err)
	}

	// All fields should be zero values
	if task.ID != 0 {
		t.Errorf("Expected ID 0, got %d", task.ID)
	}
	if task.Title != "" {
		t.Errorf("Expected empty Title, got '%s'", task.Title)
	}
	if task.Status != "" {
		t.Errorf("Expected empty Status, got '%s'", task.Status)
	}
	if task.UserID != 0 {
		t.Errorf("Expected UserID 0, got %d", task.UserID)
	}
}

// TestTask_StatusValues tests various status values
func TestTask_StatusValues(t *testing.T) {
	statuses := []string{"pending", "in-progress", "completed", "cancelled", "on-hold"}

	for i, status := range statuses {
		task := Task{
			ID:     i + 1,
			Title:  "Task " + string(rune('A'+i)),
			Status: status,
			UserID: i + 1,
		}

		jsonData, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("Failed to marshal task with status '%s': %v", status, err)
		}

		var unmarshaledTask Task
		err = json.Unmarshal(jsonData, &unmarshaledTask)
		if err != nil {
			t.Fatalf("Failed to unmarshal task with status '%s': %v", status, err)
		}

		if unmarshaledTask.Status != status {
			t.Errorf("Expected Status '%s', got '%s'", status, unmarshaledTask.Status)
		}
	}
}

// TestTask_LargeValues tests handling of large values
func TestTask_LargeValues(t *testing.T) {
	longString := "a"
	for i := 0; i < 1000; i++ {
		longString += "a"
	}

	task := Task{
		ID:     999999999,
		Title:  longString,
		Status: longString,
		UserID: 999999999,
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task with large values: %v", err)
	}

	var unmarshaledTask Task
	err = json.Unmarshal(jsonData, &unmarshaledTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal task with large values: %v", err)
	}

	// Verify large values are preserved
	if unmarshaledTask.ID != task.ID {
		t.Errorf("Expected ID %d, got %d", task.ID, unmarshaledTask.ID)
	}
	if unmarshaledTask.Title != task.Title {
		t.Errorf("Expected Title length %d, got %d", len(task.Title), len(unmarshaledTask.Title))
	}
	if unmarshaledTask.Status != task.Status {
		t.Errorf("Expected Status length %d, got %d", len(task.Status), len(unmarshaledTask.Status))
	}
	if unmarshaledTask.UserID != task.UserID {
		t.Errorf("Expected UserID %d, got %d", task.UserID, unmarshaledTask.UserID)
	}
}

// TestTask_SpecialCharacters tests handling of special characters
func TestTask_SpecialCharacters(t *testing.T) {
	task := Task{
		ID:     1,
		Title:  "Task \"Special\" Chars\nNew Line\tTab",
		Status: "in-progress/urgent",
		UserID: 1,
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task with special characters: %v", err)
	}

	var unmarshaledTask Task
	err = json.Unmarshal(jsonData, &unmarshaledTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal task with special characters: %v", err)
	}

	// Verify special characters are preserved
	if unmarshaledTask.Title != task.Title {
		t.Errorf("Expected Title with special characters, got '%s'", unmarshaledTask.Title)
	}
	if unmarshaledTask.Status != task.Status {
		t.Errorf("Expected Status with special characters, got '%s'", unmarshaledTask.Status)
	}
}

// TestTask_NegativeValues tests negative values handling
func TestTask_NegativeValues(t *testing.T) {
	task := Task{
		ID:     -1,
		Title:  "Negative ID Task",
		Status: "pending",
		UserID: -1,
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task with negative values: %v", err)
	}

	var unmarshaledTask Task
	err = json.Unmarshal(jsonData, &unmarshaledTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal task with negative values: %v", err)
	}

	if unmarshaledTask.ID != -1 {
		t.Errorf("Expected ID -1, got %d", unmarshaledTask.ID)
	}
	if unmarshaledTask.UserID != -1 {
		t.Errorf("Expected UserID -1, got %d", unmarshaledTask.UserID)
	}
}

// TestTask_ZeroUserID tests task with zero UserID
func TestTask_ZeroUserID(t *testing.T) {
	task := Task{
		ID:     1,
		Title:  "Unassigned Task",
		Status: "pending",
		UserID: 0, // Zero UserID indicates unassigned
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task with zero UserID: %v", err)
	}

	var unmarshaledTask Task
	err = json.Unmarshal(jsonData, &unmarshaledTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal task with zero UserID: %v", err)
	}

	if unmarshaledTask.UserID != 0 {
		t.Errorf("Expected UserID 0, got %d", unmarshaledTask.UserID)
	}
}

// TestTask_Unicode tests handling of Unicode characters
func TestTask_Unicode(t *testing.T) {
	task := Task{
		ID:     1,
		Title:  "Tâsk with ünïçødé 🚀",
		Status: "pêndïng",
		UserID: 1,
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task with Unicode: %v", err)
	}

	var unmarshaledTask Task
	err = json.Unmarshal(jsonData, &unmarshaledTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal task with Unicode: %v", err)
	}

	// Verify Unicode characters are preserved
	if unmarshaledTask.Title != task.Title {
		t.Errorf("Expected Unicode Title, got '%s'", unmarshaledTask.Title)
	}
	if unmarshaledTask.Status != task.Status {
		t.Errorf("Expected Unicode Status, got '%s'", unmarshaledTask.Status)
	}
}
