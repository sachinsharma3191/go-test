package repository

import (
	"context"
	"os"
	"testing"

	"go-backend/model"
	"go-backend/store"
	_ "go-backend/store/json" // Import to register JSON factory
)

// setupTaskRepositoryTest creates a test store and task repository
func setupTaskRepositoryTest(t *testing.T) (*TaskRepository, store.Store) {
	t.Helper()
	tempFile := "test_task_repo.json"

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

	taskRepo := NewTaskRepository(jsonStore)
	return taskRepo, jsonStore
}

// TestTaskRepository_FindAll tests finding all tasks
func TestTaskRepository_FindAll(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	tasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].ID != 1 || tasks[1].ID != 2 {
		t.Error("Tasks not in expected order")
	}

	// Verify task data
	if tasks[0].Title != "Task 1" {
		t.Errorf("Expected first task title 'Task 1', got '%s'", tasks[0].Title)
	}
	if tasks[1].Status != "completed" {
		t.Errorf("Expected second task status 'completed', got '%s'", tasks[1].Status)
	}
}

// TestTaskRepository_FindAll_EmptyStore tests finding all tasks in empty store
func TestTaskRepository_FindAll_EmptyStore(t *testing.T) {
	tempFile := "test_empty_task.json"
	defer os.Remove(tempFile)

	jsonStore, err := store.QuickCreate("json", tempFile)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	taskRepo := NewTaskRepository(jsonStore)

	tasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks in empty store, got %d", len(tasks))
	}
}

// TestTaskRepository_FindByID tests finding a task by ID
func TestTaskRepository_FindByID(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	// Test existing task
	task, err := taskRepo.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if task == nil {
		t.Fatal("Expected to find task with ID 1")
	}

	if task.ID != 1 {
		t.Errorf("Expected task ID 1, got %d", task.ID)
	}
	if task.Title != "Task 1" {
		t.Errorf("Expected title 'Task 1', got '%s'", task.Title)
	}
	if task.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", task.Status)
	}
	if task.UserID != 1 {
		t.Errorf("Expected user ID 1, got %d", task.UserID)
	}

	// Test second task
	task, err = taskRepo.FindByID(2)
	if err != nil {
		t.Fatalf("FindByID for task 2 failed: %v", err)
	}

	if task == nil {
		t.Fatal("Expected to find task with ID 2")
	}
	if task.Title != "Task 2" {
		t.Errorf("Expected title 'Task 2', got '%s'", task.Title)
	}

	// Test non-existent task
	task, err = taskRepo.FindByID(999)
	if err != nil {
		t.Fatalf("FindByID for non-existent task failed: %v", err)
	}

	if task != nil {
		t.Error("Expected nil for non-existent task")
	}
}

// TestTaskRepository_Create tests creating a new task
func TestTaskRepository_Create(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	newTask := &model.Task{
		Title:  "New Task",
		Status: "in-progress",
		UserID: 1,
	}

	createdTask, err := taskRepo.Create(newTask)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if createdTask.ID == 0 {
		t.Error("Expected task ID to be set")
	}

	if createdTask.ID != 3 {
		t.Errorf("Expected task ID 3, got %d", createdTask.ID)
	}

	if createdTask.Title != "New Task" {
		t.Errorf("Expected title 'New Task', got '%s'", createdTask.Title)
	}
	if createdTask.Status != "in-progress" {
		t.Errorf("Expected status 'in-progress', got '%s'", createdTask.Status)
	}
	if createdTask.UserID != 1 {
		t.Errorf("Expected user ID 1, got %d", createdTask.UserID)
	}

	// Verify task was actually saved
	allTasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll after create failed: %v", err)
	}

	if len(allTasks) != 3 {
		t.Errorf("Expected 3 tasks after create, got %d", len(allTasks))
	}

	// Verify the new task is in the list
	found := false
	for _, task := range allTasks {
		if task.ID == 3 && task.Title == "New Task" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Created task not found in task list")
	}
}

// TestTaskRepository_Create_Multiple tests creating multiple tasks
func TestTaskRepository_Create_Multiple(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	// Create multiple tasks
	tasks := []*model.Task{
		{Title: "Task 3", Status: "pending", UserID: 1},
		{Title: "Task 4", Status: "in-progress", UserID: 2},
		{Title: "Task 5", Status: "completed", UserID: 1},
	}

	for i, task := range tasks {
		createdTask, err := taskRepo.Create(task)
		if err != nil {
			t.Fatalf("Create task %d failed: %v", i+1, err)
		}

		expectedID := i + 3
		if createdTask.ID != expectedID {
			t.Errorf("Expected task ID %d, got %d", expectedID, createdTask.ID)
		}
	}

	// Verify all tasks were created
	allTasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll after multiple creates failed: %v", err)
	}

	if len(allTasks) != 5 {
		t.Errorf("Expected 5 tasks after multiple creates, got %d", len(allTasks))
	}
}

// TestTaskRepository_Create_DifferentUsers tests creating tasks for different users
func TestTaskRepository_Create_DifferentUsers(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	// Create tasks for different users
	task1 := &model.Task{Title: "User 1 Task", Status: "pending", UserID: 1}
	task2 := &model.Task{Title: "User 2 Task", Status: "completed", UserID: 2}

	createdTask1, err := taskRepo.Create(task1)
	if err != nil {
		t.Fatalf("Create task for user 1 failed: %v", err)
	}

	createdTask2, err := taskRepo.Create(task2)
	if err != nil {
		t.Fatalf("Create task for user 2 failed: %v", err)
	}

	if createdTask1.UserID != 1 {
		t.Errorf("Expected user ID 1, got %d", createdTask1.UserID)
	}
	if createdTask2.UserID != 2 {
		t.Errorf("Expected user ID 2, got %d", createdTask2.UserID)
	}

	// Verify both tasks exist
	task1Found, err := taskRepo.FindByID(createdTask1.ID)
	if err != nil {
		t.Fatalf("Find task 1 failed: %v", err)
	}
	if task1Found.UserID != 1 {
		t.Errorf("Expected task 1 to belong to user 1, got %d", task1Found.UserID)
	}

	task2Found, err := taskRepo.FindByID(createdTask2.ID)
	if err != nil {
		t.Fatalf("Find task 2 failed: %v", err)
	}
	if task2Found.UserID != 2 {
		t.Errorf("Expected task 2 to belong to user 2, got %d", task2Found.UserID)
	}
}

// TestTaskRepository_Update tests updating an existing task
func TestTaskRepository_Update(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	// Test updating task 1
	updatedTask := &model.Task{
		ID:     1,
		Title:  "Updated Task 1",
		Status: "completed",
		UserID: 2,
	}

	result, err := taskRepo.Update(updatedTask)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if result.Title != "Updated Task 1" {
		t.Errorf("Expected title 'Updated Task 1', got '%s'", result.Title)
	}
	if result.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", result.Status)
	}
	if result.UserID != 2 {
		t.Errorf("Expected user ID 2, got %d", result.UserID)
	}

	// Verify update was saved
	task, err := taskRepo.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID after update failed: %v", err)
	}

	if task.Title != "Updated Task 1" {
		t.Errorf("Expected updated title 'Updated Task 1', got '%s'", task.Title)
	}

	// Test partial update (only title)
	partialUpdate := &model.Task{
		ID:    2,
		Title: "Updated Task 2",
	}

	result, err = taskRepo.Update(partialUpdate)
	if err != nil {
		t.Fatalf("Partial update failed: %v", err)
	}

	if result.Title != "Updated Task 2" {
		t.Errorf("Expected updated title 'Updated Task 2', got '%s'", result.Title)
	}

	// Original status and user ID should be preserved
	task, err = taskRepo.FindByID(2)
	if err != nil {
		t.Fatalf("FindByID after partial update failed: %v", err)
	}

	if task.Status != "completed" {
		t.Errorf("Expected original status 'completed', got '%s'", task.Status)
	}
	if task.UserID != 2 {
		t.Errorf("Expected original user ID 2, got %d", task.UserID)
	}
}

// TestTaskRepository_Update_NonExistent tests updating a non-existent task
func TestTaskRepository_Update_NonExistent(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	nonExistentTask := &model.Task{
		ID:     999,
		Title:  "Ghost Task",
		Status: "pending",
		UserID: 1,
	}

	_, err := taskRepo.Update(nonExistentTask)
	if err == nil {
		t.Error("Expected error when updating non-existent task")
	}

	// Verify error message contains expected content
	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

// TestTaskRepository_DeleteByID tests deleting a task
func TestTaskRepository_DeleteByID(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	// Test deleting existing task
	err := taskRepo.DeleteByID(1)
	if err != nil {
		t.Fatalf("DeleteByID failed: %v", err)
	}

	// Verify task was deleted
	task, err := taskRepo.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID after delete failed: %v", err)
	}

	if task != nil {
		t.Error("Expected nil for deleted task")
	}

	// Verify task count decreased
	allTasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll after delete failed: %v", err)
	}

	if len(allTasks) != 1 {
		t.Errorf("Expected 1 task after delete, got %d", len(allTasks))
	}

	// Verify the remaining task is the correct one
	if len(allTasks) > 0 && allTasks[0].ID != 2 {
		t.Errorf("Expected remaining task ID 2, got %d", allTasks[0].ID)
	}

	// Test deleting the remaining task
	err = taskRepo.DeleteByID(2)
	if err != nil {
		t.Fatalf("DeleteByID for second task failed: %v", err)
	}

	// Verify all tasks are deleted
	allTasks, err = taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll after deleting all tasks failed: %v", err)
	}

	if len(allTasks) != 0 {
		t.Errorf("Expected 0 tasks after deleting all, got %d", len(allTasks))
	}
}

// TestTaskRepository_DeleteByID_NonExistent tests deleting a non-existent task
func TestTaskRepository_DeleteByID_NonExistent(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	err := taskRepo.DeleteByID(999)
	if err == nil {
		t.Error("Expected error when deleting non-existent task")
	}

	// Verify error message contains expected content
	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

// TestTaskRepository_DeleteByID_Zero tests deleting with ID 0
func TestTaskRepository_DeleteByID_Zero(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	err := taskRepo.DeleteByID(0)
	if err == nil {
		t.Error("Expected error when deleting task with ID 0")
	}
}

// TestTaskRepository_ErrorCases tests error conditions
func TestTaskRepository_ErrorCases(t *testing.T) {
	// Test with nil store
	repo := NewTaskRepository(nil)

	_, err := repo.FindAll()
	if err == nil {
		t.Error("Expected error with nil store")
	}

	_, err = repo.FindByID(1)
	if err == nil {
		t.Error("Expected error with nil store")
	}

	_, err = repo.Create(&model.Task{Title: "Test", Status: "pending", UserID: 1})
	if err == nil {
		t.Error("Expected error with nil store")
	}

	_, err = repo.Update(&model.Task{ID: 1, Title: "Test", Status: "pending", UserID: 1})
	if err == nil {
		t.Error("Expected error with nil store")
	}

	err = repo.DeleteByID(1)
	if err == nil {
		t.Error("Expected error with nil store")
	}
}

// TestTaskRepository_Integration tests integration scenarios
func TestTaskRepository_Integration(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	// Create a new task
	newTask := &model.Task{
		Title:  "Integration Task",
		Status: "pending",
		UserID: 1,
	}

	createdTask, err := taskRepo.Create(newTask)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Find the task by ID
	foundTask, err := taskRepo.FindByID(createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to find task: %v", err)
	}

	if foundTask.Title != "Integration Task" {
		t.Errorf("Expected title 'Integration Task', got '%s'", foundTask.Title)
	}

	// Update the task
	updatedTask := &model.Task{
		ID:     createdTask.ID,
		Title:  "Updated Integration Task",
		Status: "completed",
		UserID: 2,
	}

	_, err = taskRepo.Update(updatedTask)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	// Verify the update
	updatedFoundTask, err := taskRepo.FindByID(createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to find updated task: %v", err)
	}

	if updatedFoundTask.Title != "Updated Integration Task" {
		t.Errorf("Expected updated title, got '%s'", updatedFoundTask.Title)
	}

	// Delete the task
	err = taskRepo.DeleteByID(createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	// Verify deletion
	_, err = taskRepo.FindByID(createdTask.ID)
	if err != nil {
		t.Fatalf("FindByID after delete failed: %v", err)
	}

	// Check final task count
	allTasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll after integration test failed: %v", err)
	}

	if len(allTasks) != 2 {
		t.Errorf("Expected 2 tasks after integration test, got %d", len(allTasks))
	}
}

// TestTaskRepository_UserTasks tests finding tasks by user
func TestTaskRepository_UserTasks(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	// Create additional tasks for different users
	task1 := &model.Task{Title: "User 1 Task 1", Status: "pending", UserID: 1}
	task2 := &model.Task{Title: "User 1 Task 2", Status: "in-progress", UserID: 1}
	task3 := &model.Task{Title: "User 2 Task 1", Status: "completed", UserID: 2}

	_, err := taskRepo.Create(task1)
	if err != nil {
		t.Fatalf("Failed to create task1: %v", err)
	}

	_, err = taskRepo.Create(task2)
	if err != nil {
		t.Fatalf("Failed to create task2: %v", err)
	}

	_, err = taskRepo.Create(task3)
	if err != nil {
		t.Fatalf("Failed to create task3: %v", err)
	}

	// Get all tasks
	allTasks, err := taskRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	// Count tasks by user
	user1Tasks := 0
	user2Tasks := 0
	for _, task := range allTasks {
		if task.UserID == 1 {
			user1Tasks++
		} else if task.UserID == 2 {
			user2Tasks++
		}
	}

	if user1Tasks != 3 {
		t.Errorf("Expected 3 tasks for user 1, got %d", user1Tasks)
	}
	if user2Tasks != 2 {
		t.Errorf("Expected 2 tasks for user 2, got %d", user2Tasks)
	}
}

// TestTaskRepository_StatusUpdates tests updating task status
func TestTaskRepository_StatusUpdates(t *testing.T) {
	taskRepo, _ := setupTaskRepositoryTest(t)

	// Create a new task
	newTask := &model.Task{
		Title:  "Status Test Task",
		Status: "pending",
		UserID: 1,
	}

	createdTask, err := taskRepo.Create(newTask)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Update status to in-progress
	statusUpdate := &model.Task{
		ID:     createdTask.ID,
		Status: "in-progress",
	}

	_, err = taskRepo.Update(statusUpdate)
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}

	// Verify status update
	task, err := taskRepo.FindByID(createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to find task: %v", err)
	}

	if task.Status != "in-progress" {
		t.Errorf("Expected status 'in-progress', got '%s'", task.Status)
	}

	// Original title should be preserved
	if task.Title != "Status Test Task" {
		t.Errorf("Expected original title 'Status Test Task', got '%s'", task.Title)
	}

	// Update status to completed
	completedUpdate := &model.Task{
		ID:     createdTask.ID,
		Status: "completed",
	}

	_, err = taskRepo.Update(completedUpdate)
	if err != nil {
		t.Fatalf("Failed to update status to completed: %v", err)
	}

	// Verify final status
	task, err = taskRepo.FindByID(createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to find task after final update: %v", err)
	}

	if task.Status != "completed" {
		t.Errorf("Expected final status 'completed', got '%s'", task.Status)
	}
}
