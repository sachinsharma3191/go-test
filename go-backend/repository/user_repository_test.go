package repository

import (
	"context"
	"os"
	"testing"

	"go-backend/model"
	"go-backend/store"
	_ "go-backend/store/json" // Import to register JSON factory
)

// setupUserRepositoryTest creates a test store and user repository
func setupUserRepositoryTest(t *testing.T) (*UserRepository, store.Store) {
	t.Helper()
	tempFile := "test_user_repo.json"

	jsonStore, err := store.QuickCreate("json", tempFile)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	// Add test data
	data := store.AppData{
		Users: []model.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", Role: "developer"},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Role: "designer"},
		},
		Tasks: []model.Task{
			{ID: 1, Title: "Task 1", Status: "pending", UserID: 1},
			{ID: 2, Title: "Task 2", Status: "completed", UserID: 2},
		},
	}

	err = jsonStore.Save(context.Background(), data)
	if err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// Cleanup after test
	t.Cleanup(func() {
		os.Remove(tempFile)
	})

	userRepo := NewUserRepository(jsonStore)
	return userRepo, jsonStore
}

// TestUserRepository_FindAll tests finding all users
func TestUserRepository_FindAll(t *testing.T) {
	userRepo, _ := setupUserRepositoryTest(t)

	users, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	if users[0].ID != 1 || users[1].ID != 2 {
		t.Error("Users not in expected order")
	}

	// Verify user data
	if users[0].Name != "John Doe" {
		t.Errorf("Expected first user name 'John Doe', got '%s'", users[0].Name)
	}
	if users[1].Email != "jane@example.com" {
		t.Errorf("Expected second user email 'jane@example.com', got '%s'", users[1].Email)
	}
}

// TestUserRepository_FindAll_EmptyStore tests finding all users in empty store
func TestUserRepository_FindAll_EmptyStore(t *testing.T) {
	tempFile := "test_empty_user.json"
	defer os.Remove(tempFile)

	jsonStore, err := store.QuickCreate("json", tempFile)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	userRepo := NewUserRepository(jsonStore)

	users, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("Expected 0 users in empty store, got %d", len(users))
	}
}

// TestUserRepository_FindByID tests finding a user by ID
func TestUserRepository_FindByID(t *testing.T) {
	userRepo, _ := setupUserRepositoryTest(t)

	// Test existing user
	user, err := userRepo.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if user == nil {
		t.Fatal("Expected to find user with ID 1")
	}

	if user.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", user.ID)
	}
	if user.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", user.Email)
	}
	if user.Role != "developer" {
		t.Errorf("Expected role 'developer', got '%s'", user.Role)
	}

	// Test second user
	user, err = userRepo.FindByID(2)
	if err != nil {
		t.Fatalf("FindByID for user 2 failed: %v", err)
	}

	if user == nil {
		t.Fatal("Expected to find user with ID 2")
	}
	if user.Name != "Jane Smith" {
		t.Errorf("Expected name 'Jane Smith', got '%s'", user.Name)
	}

	// Test non-existent user
	user, err = userRepo.FindByID(999)
	if err != nil {
		t.Fatalf("FindByID for non-existent user failed: %v", err)
	}

	if user != nil {
		t.Error("Expected nil for non-existent user")
	}
}

// TestUserRepository_Create tests creating a new user
func TestUserRepository_Create(t *testing.T) {
	userRepo, _ := setupUserRepositoryTest(t)

	newUser := &model.User{
		Name:  "New User",
		Email: "newuser@example.com",
		Role:  "manager",
	}

	createdUser, err := userRepo.Create(newUser)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if createdUser.ID == 0 {
		t.Error("Expected user ID to be set")
	}

	if createdUser.ID != 3 {
		t.Errorf("Expected user ID 3, got %d", createdUser.ID)
	}

	if createdUser.Name != "New User" {
		t.Errorf("Expected name 'New User', got '%s'", createdUser.Name)
	}
	if createdUser.Email != "newuser@example.com" {
		t.Errorf("Expected email 'newuser@example.com', got '%s'", createdUser.Email)
	}
	if createdUser.Role != "manager" {
		t.Errorf("Expected role 'manager', got '%s'", createdUser.Role)
	}

	// Verify user was actually saved
	allUsers, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll after create failed: %v", err)
	}

	if len(allUsers) != 3 {
		t.Errorf("Expected 3 users after create, got %d", len(allUsers))
	}

	// Verify the new user is in the list
	found := false
	for _, user := range allUsers {
		if user.ID == 3 && user.Name == "New User" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Created user not found in user list")
	}
}

// TestUserRepository_Create_Multiple tests creating multiple users
func TestUserRepository_Create_Multiple(t *testing.T) {
	userRepo, _ := setupUserRepositoryTest(t)

	// Create multiple users
	users := []*model.User{
		{Name: "User 3", Email: "user3@example.com", Role: "developer"},
		{Name: "User 4", Email: "user4@example.com", Role: "designer"},
		{Name: "User 5", Email: "user5@example.com", Role: "manager"},
	}

	for i, user := range users {
		createdUser, err := userRepo.Create(user)
		if err != nil {
			t.Fatalf("Create user %d failed: %v", i+1, err)
		}

		expectedID := i + 3
		if createdUser.ID != expectedID {
			t.Errorf("Expected user ID %d, got %d", expectedID, createdUser.ID)
		}
	}

	// Verify all users were created
	allUsers, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll after multiple creates failed: %v", err)
	}

	if len(allUsers) != 5 {
		t.Errorf("Expected 5 users after multiple creates, got %d", len(allUsers))
	}
}

// TestUserRepository_Update tests updating an existing user
func TestUserRepository_Update(t *testing.T) {
	userRepo, _ := setupUserRepositoryTest(t)

	// Test updating user 1
	updatedUser := &model.User{
		ID:    1,
		Name:  "Updated John",
		Email: "updated@example.com",
		Role:  "senior-developer",
	}

	result, err := userRepo.Update(updatedUser)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if result.Name != "Updated John" {
		t.Errorf("Expected name 'Updated John', got '%s'", result.Name)
	}
	if result.Email != "updated@example.com" {
		t.Errorf("Expected email 'updated@example.com', got '%s'", result.Email)
	}
	if result.Role != "senior-developer" {
		t.Errorf("Expected role 'senior-developer', got '%s'", result.Role)
	}

	// Verify update was saved
	user, err := userRepo.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID after update failed: %v", err)
	}

	if user.Name != "Updated John" {
		t.Errorf("Expected updated name 'Updated John', got '%s'", user.Name)
	}

	// Test partial update (only name)
	partialUpdate := &model.User{
		ID:   2,
		Name: "Updated Jane",
	}

	result, err = userRepo.Update(partialUpdate)
	if err != nil {
		t.Fatalf("Partial update failed: %v", err)
	}

	if result.Name != "Updated Jane" {
		t.Errorf("Expected updated name 'Updated Jane', got '%s'", result.Name)
	}

	// Note: The repository Update method completely replaces the user
	// So email and role will be empty if not provided
	user, err = userRepo.FindByID(2)
	if err != nil {
		t.Fatalf("FindByID after partial update failed: %v", err)
	}

	if user.Name != "Updated Jane" {
		t.Errorf("Expected updated name 'Updated Jane', got '%s'", user.Name)
	}
	// Email and role will be empty since they weren't provided in the update
	if user.Email != "" {
		t.Errorf("Expected empty email after partial update, got '%s'", user.Email)
	}
	if user.Role != "" {
		t.Errorf("Expected empty role after partial update, got '%s'", user.Role)
	}
}

// TestUserRepository_Update_NonExistent tests updating a non-existent user
func TestUserRepository_Update_NonExistent(t *testing.T) {
	userRepo, _ := setupUserRepositoryTest(t)

	nonExistentUser := &model.User{
		ID:    999,
		Name:  "Ghost",
		Email: "ghost@example.com",
		Role:  "developer",
	}

	_, err := userRepo.Update(nonExistentUser)
	if err == nil {
		t.Error("Expected error when updating non-existent user")
	}

	// Verify error message contains expected content
	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

// TestUserRepository_DeleteByID tests deleting a user
func TestUserRepository_DeleteByID(t *testing.T) {
	userRepo, _ := setupUserRepositoryTest(t)

	// Test deleting existing user
	err := userRepo.DeleteByID(1)
	if err != nil {
		t.Fatalf("DeleteByID failed: %v", err)
	}

	// Verify user was deleted
	user, err := userRepo.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID after delete failed: %v", err)
	}

	if user != nil {
		t.Error("Expected nil for deleted user")
	}

	// Verify user count decreased
	allUsers, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll after delete failed: %v", err)
	}

	if len(allUsers) != 1 {
		t.Errorf("Expected 1 user after delete, got %d", len(allUsers))
	}

	// Verify the remaining user is the correct one
	if len(allUsers) > 0 && allUsers[0].ID != 2 {
		t.Errorf("Expected remaining user ID 2, got %d", allUsers[0].ID)
	}

	// Test deleting the remaining user
	err = userRepo.DeleteByID(2)
	if err != nil {
		t.Fatalf("DeleteByID for second user failed: %v", err)
	}

	// Verify all users are deleted
	allUsers, err = userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll after deleting all users failed: %v", err)
	}

	if len(allUsers) != 0 {
		t.Errorf("Expected 0 users after deleting all, got %d", len(allUsers))
	}
}

// TestUserRepository_DeleteByID_NonExistent tests deleting a non-existent user
func TestUserRepository_DeleteByID_NonExistent(t *testing.T) {
	userRepo, _ := setupUserRepositoryTest(t)

	err := userRepo.DeleteByID(999)
	if err == nil {
		t.Error("Expected error when deleting non-existent user")
	}

	// Verify error message contains expected content
	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

// TestUserRepository_DeleteByID_Zero tests deleting with ID 0
func TestUserRepository_DeleteByID_Zero(t *testing.T) {
	userRepo, _ := setupUserRepositoryTest(t)

	err := userRepo.DeleteByID(0)
	if err == nil {
		t.Error("Expected error when deleting user with ID 0")
	}
}

// TestUserRepository_ErrorCases tests error conditions
func TestUserRepository_ErrorCases(t *testing.T) {
	// Test with nil store
	repo := NewUserRepository(nil)

	_, err := repo.FindAll()
	if err == nil {
		t.Error("Expected error with nil store")
	}

	_, err = repo.FindByID(1)
	if err == nil {
		t.Error("Expected error with nil store")
	}

	_, err = repo.Create(&model.User{Name: "Test", Email: "test@example.com", Role: "developer"})
	if err == nil {
		t.Error("Expected error with nil store")
	}

	_, err = repo.Update(&model.User{ID: 1, Name: "Test", Email: "test@example.com", Role: "developer"})
	if err == nil {
		t.Error("Expected error with nil store")
	}

	err = repo.DeleteByID(1)
	if err == nil {
		t.Error("Expected error with nil store")
	}
}

// TestUserRepository_Integration tests integration scenarios
func TestUserRepository_Integration(t *testing.T) {
	userRepo, _ := setupUserRepositoryTest(t)

	// Create a new user
	newUser := &model.User{
		Name:  "Integration User",
		Email: "integration@example.com",
		Role:  "developer",
	}

	createdUser, err := userRepo.Create(newUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Find the user by ID
	foundUser, err := userRepo.FindByID(createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if foundUser.Name != "Integration User" {
		t.Errorf("Expected name 'Integration User', got '%s'", foundUser.Name)
	}

	// Update the user
	updatedUser := &model.User{
		ID:    createdUser.ID,
		Name:  "Updated Integration User",
		Email: "updated@example.com",
		Role:  "senior-developer",
	}

	_, err = userRepo.Update(updatedUser)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify the update
	updatedFoundUser, err := userRepo.FindByID(createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to find updated user: %v", err)
	}

	if updatedFoundUser.Name != "Updated Integration User" {
		t.Errorf("Expected updated name, got '%s'", updatedFoundUser.Name)
	}

	// Delete the user
	err = userRepo.DeleteByID(createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify deletion
	_, err = userRepo.FindByID(createdUser.ID)
	if err != nil {
		t.Fatalf("FindByID after delete failed: %v", err)
	}

	// Check final user count
	allUsers, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll after integration test failed: %v", err)
	}

	if len(allUsers) != 2 {
		t.Errorf("Expected 2 users after integration test, got %d", len(allUsers))
	}
}
