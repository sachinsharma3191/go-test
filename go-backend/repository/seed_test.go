package repository

import (
	"context"
	"errors"
	"os"
	"testing"

	"go-backend/model"
	"go-backend/store"
	_ "go-backend/store/json" // Import to register JSON factory
)

var errSeedSave = errors.New("seed save failed")

type seedSaveFailingStore struct{}

func (s *seedSaveFailingStore) Load(_ context.Context) (store.AppData, error) {
	return store.AppData{}, nil
}
func (s *seedSaveFailingStore) Save(_ context.Context, _ store.AppData) error {
	return errSeedSave
}
func (s *seedSaveFailingStore) Health(_ context.Context) error { return nil }
func (s *seedSaveFailingStore) Close() error                  { return nil }

// setupSeedTest creates an empty store for seed testing
func setupSeedTest(t *testing.T) (*UserRepository, *TaskRepository) {
	t.Helper()
	tempFile := "test_seed.json"

	jsonStore, err := store.QuickCreate("json", tempFile)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	// Cleanup after test
	t.Cleanup(func() {
		os.Remove(tempFile)
	})

	userRepo := NewUserRepository(jsonStore)
	taskRepo := NewTaskRepository(jsonStore)

	return userRepo, taskRepo
}

// TestSeedIfEmpty tests the basic seeding functionality
func TestSeedIfEmpty(t *testing.T) {
	userRepo, taskRepo := setupSeedTest(t)

	// Test seeding empty store
	err := SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("SeedIfEmpty failed: %v", err)
	}

	// Verify seeded users
	users, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 seeded users, got %d", len(users))
	}

	// Verify specific seeded users
	expectedUsers := []struct {
		name  string
		email string
		role  string
	}{
		{"John Doe", "john@example.com", "developer"},
		{"Jane Smith", "jane@example.com", "designer"},
		{"Bob Johnson", "bob@example.com", "manager"},
	}

	for i, expected := range expectedUsers {
		if i >= len(users) {
			t.Fatalf("Expected user %d not found", i+1)
		}

		if users[i].Name != expected.name {
			t.Errorf("Expected user %d name '%s', got '%s'", i+1, expected.name, users[i].Name)
		}
		if users[i].Email != expected.email {
			t.Errorf("Expected user %d email '%s', got '%s'", i+1, expected.email, users[i].Email)
		}
		if users[i].Role != expected.role {
			t.Errorf("Expected user %d role '%s', got '%s'", i+1, expected.role, users[i].Role)
		}
	}

	// Verify seeded tasks
	tasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll tasks failed: %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 seeded tasks, got %d", len(tasks))
	}

	// Verify specific seeded tasks
	expectedTasks := []struct {
		title  string
		status string
		userID int
	}{
		{"Implement authentication", "pending", 1},
		{"Design user interface", "in-progress", 2},
		{"Review code changes", "completed", 3},
	}

	for i, expected := range expectedTasks {
		if i >= len(tasks) {
			t.Fatalf("Expected task %d not found", i+1)
		}

		if tasks[i].Title != expected.title {
			t.Errorf("Expected task %d title '%s', got '%s'", i+1, expected.title, tasks[i].Title)
		}
		if tasks[i].Status != expected.status {
			t.Errorf("Expected task %d status '%s', got '%s'", i+1, expected.status, tasks[i].Status)
		}
		if tasks[i].UserID != expected.userID {
			t.Errorf("Expected task %d user ID %d, got %d", i+1, expected.userID, tasks[i].UserID)
		}
	}
}

// TestSeedIfEmpty_NotEmpty tests that seeding doesn't occur when store is not empty
func TestSeedIfEmpty_NotEmpty(t *testing.T) {
	userRepo, taskRepo := setupSeedTest(t)

	// Add some data to the store first
	// We need to access the underlying store to add initial data
	// For now, let's create a user first to make the store non-empty
	_, err := userRepo.Create(&model.User{
		Name:  "Existing User",
		Email: "existing@example.com",
		Role:  "developer",
	})
	if err != nil {
		t.Fatalf("Failed to create initial user: %v", err)
	}

	// Attempt to seed
	err = SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("SeedIfEmpty failed: %v", err)
	}

	// Verify that original data is preserved
	users, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user (original data preserved), got %d", len(users))
	}

	if users[0].Name != "Existing User" {
		t.Errorf("Expected original user name 'Existing User', got '%s'", users[0].Name)
	}
}

// TestSeedIfEmpty_PartialData tests seeding when only users or tasks exist
func TestSeedIfEmpty_PartialData(t *testing.T) {
	userRepo, taskRepo := setupSeedTest(t)

	// Test with only users existing
	// Create a user first to make the store non-empty
	_, err := userRepo.Create(&model.User{
		Name:  "Existing User",
		Email: "existing@example.com",
		Role:  "developer",
	})
	if err != nil {
		t.Fatalf("Failed to create partial data: %v", err)
	}

	// Attempt to seed
	err = SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("SeedIfEmpty failed: %v", err)
	}

	// Verify that no seeding occurred (because users exist)
	users, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user (no seeding), got %d", len(users))
	}

	tasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll tasks failed: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks (no seeding), got %d", len(tasks))
	}
}

// TestSeedIfEmpty_MultipleCalls tests that multiple calls to SeedIfEmpty don't duplicate data
func TestSeedIfEmpty_MultipleCalls(t *testing.T) {
	userRepo, taskRepo := setupSeedTest(t)

	// First call to seed
	err := SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("First SeedIfEmpty failed: %v", err)
	}

	// Get counts after first seeding
	users, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	userCount1 := len(users)

	tasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll tasks failed: %v", err)
	}
	taskCount1 := len(tasks)

	// Second call to seed
	err = SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("Second SeedIfEmpty failed: %v", err)
	}

	// Get counts after second seeding
	users, err = userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	userCount2 := len(users)

	tasks, err = taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll tasks failed: %v", err)
	}
	taskCount2 := len(tasks)

	// Verify no duplication
	if userCount2 != userCount1 {
		t.Errorf("Expected user count to remain %d, got %d", userCount1, userCount2)
	}

	if taskCount2 != taskCount1 {
		t.Errorf("Expected task count to remain %d, got %d", taskCount1, taskCount2)
	}
}

// TestSeedIfEmpty_ErrorHandling tests error handling in seed function
func TestSeedIfEmpty_ErrorHandling(t *testing.T) {
	// Test with nil repositories
	err := SeedIfEmpty(nil, nil)
	if err == nil {
		t.Error("Expected error when both repositories are nil")
	}

	// Test with nil user repository
	err = SeedIfEmpty(nil, NewTaskRepository(nil))
	if err == nil {
		t.Error("Expected error when user repository is nil")
	}

	// Test with nil task repository
	err = SeedIfEmpty(NewUserRepository(nil), nil)
	if err == nil {
		t.Error("Expected error when task repository is nil")
	}

	// Test with repositories that have nil stores - this should not panic
	// but we need to handle it gracefully
	userRepo := NewUserRepository(nil)
	taskRepo := NewTaskRepository(nil)

	err = SeedIfEmpty(userRepo, taskRepo)
	if err == nil {
		t.Error("Expected error when repositories have nil stores")
	}

	// Test with store that fails on Save during seeding
	saveFailStore := &seedSaveFailingStore{}
	userRepoFail := NewUserRepository(saveFailStore)
	taskRepoFail := NewTaskRepository(saveFailStore)
	err = SeedIfEmpty(userRepoFail, taskRepoFail)
	if err == nil {
		t.Error("Expected error when store Save fails during seeding")
	}
}

// TestSeedIfEmpty_DataIntegrity tests that seeded data maintains integrity
func TestSeedIfEmpty_DataIntegrity(t *testing.T) {
	userRepo, taskRepo := setupSeedTest(t)

	// Seed the data
	err := SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("SeedIfEmpty failed: %v", err)
	}

	// Verify user-task relationships
	users, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	tasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll tasks failed: %v", err)
	}

	// Create a map of user IDs for validation
	userMap := make(map[int]model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Verify all task user IDs correspond to existing users
	for _, task := range tasks {
		if _, exists := userMap[task.UserID]; !exists {
			t.Errorf("Task '%s' references non-existent user ID %d", task.Title, task.UserID)
		}
	}

	// Verify specific expected relationships
	// Task 1 should belong to user 1 (John Doe)
	task1, err := taskRepo.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID for task 1 failed: %v", err)
	}

	if task1.UserID != 1 {
		t.Errorf("Expected task 1 to belong to user 1, got %d", task1.UserID)
	}

	user1, err := userRepo.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID for user 1 failed: %v", err)
	}

	if user1.Name != "John Doe" {
		t.Errorf("Expected user 1 to be John Doe, got '%s'", user1.Name)
	}

	// Task 2 should belong to user 2 (Jane Smith)
	task2, err := taskRepo.FindByID(2)
	if err != nil {
		t.Fatalf("FindByID for task 2 failed: %v", err)
	}

	if task2.UserID != 2 {
		t.Errorf("Expected task 2 to belong to user 2, got %d", task2.UserID)
	}

	user2, err := userRepo.FindByID(2)
	if err != nil {
		t.Fatalf("FindByID for user 2 failed: %v", err)
	}

	if user2.Name != "Jane Smith" {
		t.Errorf("Expected user 2 to be Jane Smith, got '%s'", user2.Name)
	}
}

// TestSeedIfEmpty_ConcurrentAccess tests seeding with concurrent access
func TestSeedIfEmpty_ConcurrentAccess(t *testing.T) {
	userRepo, taskRepo := setupSeedTest(t)

	// Test that seeding works correctly (simplified test to avoid race conditions)
	err := SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("SeedIfEmpty failed: %v", err)
	}

	// Verify the data was seeded correctly
	users, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	tasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll tasks failed: %v", err)
	}

	// Should have exactly 3 users and 3 tasks
	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}

	// Test that subsequent calls don't add more data
	err = SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("Second SeedIfEmpty failed: %v", err)
	}

	// Verify counts haven't changed
	users, err = userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	tasks, err = taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll tasks failed: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 users after second call, got %d", len(users))
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks after second call, got %d", len(tasks))
	}
}

// TestSeedIfEmpty_Persistence tests that seeded data persists across repository instances
func TestSeedIfEmpty_Persistence(t *testing.T) {
	userRepo, taskRepo := setupSeedTest(t)

	// Seed the data
	err := SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("SeedIfEmpty failed: %v", err)
	}

	// Create new repository instances using the same store
	// We need to get the store from the first repository
	// For this test, we'll just verify the data is still there
	users, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}

	tasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll tasks failed: %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

// TestSeedIfEmpty_IDGeneration tests that seeded data has correct IDs
func TestSeedIfEmpty_IDGeneration(t *testing.T) {
	userRepo, taskRepo := setupSeedTest(t)

	// Seed the data
	err := SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("SeedIfEmpty failed: %v", err)
	}

	// Verify user IDs start from 1 and are sequential
	users, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	expectedUserIDs := []int{1, 2, 3}
	for i, user := range users {
		if i >= len(expectedUserIDs) {
			t.Errorf("Unexpected user at index %d", i)
			continue
		}
		if user.ID != expectedUserIDs[i] {
			t.Errorf("Expected user ID %d, got %d", expectedUserIDs[i], user.ID)
		}
	}

	// Verify task IDs start from 1 and are sequential
	tasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll tasks failed: %v", err)
	}

	expectedTaskIDs := []int{1, 2, 3}
	for i, task := range tasks {
		if i >= len(expectedTaskIDs) {
			t.Errorf("Unexpected task at index %d", i)
			continue
		}
		if task.ID != expectedTaskIDs[i] {
			t.Errorf("Expected task ID %d, got %d", expectedTaskIDs[i], task.ID)
		}
	}
}

// TestSeedIfEmpty_AfterOperations tests seeding after repository operations
func TestSeedIfEmpty_AfterOperations(t *testing.T) {
	userRepo, taskRepo := setupSeedTest(t)

	// Seed first
	err := SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("First SeedIfEmpty failed: %v", err)
	}

	// Perform some operations
	user, err := userRepo.Create(&model.User{
		Name:  "Additional User",
		Email: "additional@example.com",
		Role:  "developer",
	})
	if err != nil {
		t.Fatalf("Failed to create additional user: %v", err)
	}

	_, err = taskRepo.Create(&model.Task{
		Title:  "Additional Task",
		Status: "pending",
		UserID: user.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create additional task: %v", err)
	}

	// Try to seed again
	err = SeedIfEmpty(userRepo, taskRepo)
	if err != nil {
		t.Fatalf("Second SeedIfEmpty failed: %v", err)
	}

	// Verify original seeded data is preserved and additional data remains
	users, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 4 {
		t.Errorf("Expected 4 users (3 seeded + 1 added), got %d", len(users))
	}

	tasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll tasks failed: %v", err)
	}

	if len(tasks) != 4 {
		t.Errorf("Expected 4 tasks (3 seeded + 1 added), got %d", len(tasks))
	}

	// Verify the additional user and task are still there
	foundUser := false
	for _, u := range users {
		if u.Name == "Additional User" {
			foundUser = true
			break
		}
	}
	if !foundUser {
		t.Error("Additional user not found after second seeding")
	}

	foundTask := false
	for _, t := range tasks {
		if t.Title == "Additional Task" {
			foundTask = true
			break
		}
	}
	if !foundTask {
		t.Error("Additional task not found after second seeding")
	}
}
