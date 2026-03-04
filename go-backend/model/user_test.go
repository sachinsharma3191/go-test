package model

import (
	"encoding/json"
	"testing"
)

// TestUser_JSONSerialization tests JSON marshaling and unmarshaling
func TestUser_JSONSerialization(t *testing.T) {
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
		Role:  "developer",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledUser User
	err = json.Unmarshal(jsonData, &unmarshaledUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal user from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaledUser.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, unmarshaledUser.ID)
	}
	if unmarshaledUser.Name != user.Name {
		t.Errorf("Expected Name '%s', got '%s'", user.Name, unmarshaledUser.Name)
	}
	if unmarshaledUser.Email != user.Email {
		t.Errorf("Expected Email '%s', got '%s'", user.Email, unmarshaledUser.Email)
	}
	if unmarshaledUser.Role != user.Role {
		t.Errorf("Expected Role '%s', got '%s'", user.Role, unmarshaledUser.Role)
	}
}

// TestUser_JSONFields tests that JSON field names are correct
func TestUser_JSONFields(t *testing.T) {
	user := User{
		ID:    123,
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "tester",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	// Parse JSON to verify field names
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to map: %v", err)
	}

	// Verify expected field names exist
	expectedFields := []string{"id", "name", "email", "role"}
	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("Expected field '%s' not found in JSON", field)
		}
	}

	// Verify field values
	if result["id"] != float64(123) {
		t.Errorf("Expected id 123, got %v", result["id"])
	}
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

// TestUser_ZeroValues tests user with zero values
func TestUser_ZeroValues(t *testing.T) {
	user := User{}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal zero-value user: %v", err)
	}

	var unmarshaledUser User
	err = json.Unmarshal(jsonData, &unmarshaledUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal zero-value user: %v", err)
	}

	// Verify zero values are preserved
	if unmarshaledUser.ID != 0 {
		t.Errorf("Expected ID 0, got %d", unmarshaledUser.ID)
	}
	if unmarshaledUser.Name != "" {
		t.Errorf("Expected empty Name, got '%s'", unmarshaledUser.Name)
	}
	if unmarshaledUser.Email != "" {
		t.Errorf("Expected empty Email, got '%s'", unmarshaledUser.Email)
	}
	if unmarshaledUser.Role != "" {
		t.Errorf("Expected empty Role, got '%s'", unmarshaledUser.Role)
	}
}

// TestUser_PartialJSON tests partial JSON data
func TestUser_PartialJSON(t *testing.T) {
	// Test JSON with missing fields
	jsonData := []byte(`{"id": 1, "name": "Partial User"}`)

	var user User
	err := json.Unmarshal(jsonData, &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal partial JSON: %v", err)
	}

	// Verify populated fields
	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}
	if user.Name != "Partial User" {
		t.Errorf("Expected Name 'Partial User', got '%s'", user.Name)
	}

	// Verify missing fields are zero values
	if user.Email != "" {
		t.Errorf("Expected empty Email, got '%s'", user.Email)
	}
	if user.Role != "" {
		t.Errorf("Expected empty Role, got '%s'", user.Role)
	}
}

// TestUser_InvalidJSON tests invalid JSON handling
func TestUser_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"id": "not-a-number", "name": 123}`)

	var user User
	err := json.Unmarshal(invalidJSON, &user)
	if err == nil {
		t.Error("Expected error when unmarshaling invalid JSON")
	}
}

// TestUser_EmptyJSON tests empty JSON handling
func TestUser_EmptyJSON(t *testing.T) {
	emptyJSON := []byte(`{}`)

	var user User
	err := json.Unmarshal(emptyJSON, &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty JSON: %v", err)
	}

	// All fields should be zero values
	if user.ID != 0 {
		t.Errorf("Expected ID 0, got %d", user.ID)
	}
	if user.Name != "" {
		t.Errorf("Expected empty Name, got '%s'", user.Name)
	}
	if user.Email != "" {
		t.Errorf("Expected empty Email, got '%s'", user.Email)
	}
	if user.Role != "" {
		t.Errorf("Expected empty Role, got '%s'", user.Role)
	}
}

// TestUser_LargeValues tests handling of large values
func TestUser_LargeValues(t *testing.T) {
	longString := "a"
	for i := 0; i < 1000; i++ {
		longString += "a"
	}

	user := User{
		ID:    999999999,
		Name:  longString,
		Email: longString + "@example.com",
		Role:  longString,
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user with large values: %v", err)
	}

	var unmarshaledUser User
	err = json.Unmarshal(jsonData, &unmarshaledUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal user with large values: %v", err)
	}

	// Verify large values are preserved
	if unmarshaledUser.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, unmarshaledUser.ID)
	}
	if unmarshaledUser.Name != user.Name {
		t.Errorf("Expected Name length %d, got %d", len(user.Name), len(unmarshaledUser.Name))
	}
	if unmarshaledUser.Email != user.Email {
		t.Errorf("Expected Email length %d, got %d", len(user.Email), len(unmarshaledUser.Email))
	}
	if unmarshaledUser.Role != user.Role {
		t.Errorf("Expected Role length %d, got %d", len(user.Role), len(unmarshaledUser.Role))
	}
}

// TestUser_SpecialCharacters tests handling of special characters
func TestUser_SpecialCharacters(t *testing.T) {
	user := User{
		ID:    1,
		Name:  "John \"The Rock\" Doe\nNew Line\tTab",
		Email: "test+special@example.com",
		Role:  "admin/super-user",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user with special characters: %v", err)
	}

	var unmarshaledUser User
	err = json.Unmarshal(jsonData, &unmarshaledUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal user with special characters: %v", err)
	}

	// Verify special characters are preserved
	if unmarshaledUser.Name != user.Name {
		t.Errorf("Expected Name with special characters, got '%s'", unmarshaledUser.Name)
	}
	if unmarshaledUser.Email != user.Email {
		t.Errorf("Expected Email with special characters, got '%s'", unmarshaledUser.Email)
	}
	if unmarshaledUser.Role != user.Role {
		t.Errorf("Expected Role with special characters, got '%s'", unmarshaledUser.Role)
	}
}

// TestUser_NegativeID tests negative ID handling
func TestUser_NegativeID(t *testing.T) {
	user := User{
		ID:    -1,
		Name:  "Negative ID User",
		Email: "negative@example.com",
		Role:  "tester",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user with negative ID: %v", err)
	}

	var unmarshaledUser User
	err = json.Unmarshal(jsonData, &unmarshaledUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal user with negative ID: %v", err)
	}

	if unmarshaledUser.ID != -1 {
		t.Errorf("Expected ID -1, got %d", unmarshaledUser.ID)
	}
}
