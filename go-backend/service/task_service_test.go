package service

import (
	"context"
	"os"
	"testing"

	"go-backend/model"
	"go-backend/repository"
	"go-backend/store"

	_ "go-backend/store/json"
)

func setupTaskService(t *testing.T) *TaskService {
	t.Helper()
	f, err := os.CreateTemp("", "task_svc_*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })

	s, err := store.QuickCreate("json", path)
	if err != nil {
		t.Fatalf("QuickCreate: %v", err)
	}
	data := store.AppData{
		Users: []model.User{
			{ID: 1, Name: "Alice", Email: "alice@example.com", Role: "developer"},
			{ID: 2, Name: "Bob", Email: "bob@example.com", Role: "designer"},
		},
		Tasks: []model.Task{
			{ID: 1, Title: "Task 1", Status: "pending", UserID: 1},
			{ID: 2, Title: "Task 2", Status: "completed", UserID: 2},
		},
	}
	if err := s.Save(context.Background(), data); err != nil {
		t.Fatalf("Save: %v", err)
	}
	userRepo := repository.NewUserRepository(s)
	taskRepo := repository.NewTaskRepository(s)
	return NewTaskService(taskRepo, userRepo)
}

func TestTaskService_FindAllTasks(t *testing.T) {
	svc := setupTaskService(t)
	tasks, err := svc.FindAllTasks()
	if err != nil {
		t.Fatalf("FindAllTasks: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("len(tasks) = %d, want 2", len(tasks))
	}
}

func TestTaskService_FindUserTasks(t *testing.T) {
	svc := setupTaskService(t)
	tasks, err := svc.FindUserTasks(1)
	if err != nil {
		t.Fatalf("FindUserTasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("len(tasks) = %d, want 1", len(tasks))
	}
	if tasks[0].UserID != 1 {
		t.Errorf("task UserID = %d", tasks[0].UserID)
	}
}

func TestTaskService_FindUserTasksByStatus_All(t *testing.T) {
	svc := setupTaskService(t)
	tasks, err := svc.FindUserTasksByStatus("", "")
	if err != nil {
		t.Fatalf("FindUserTasksByStatus: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("len(tasks) = %d, want 2", len(tasks))
	}
}

func TestTaskService_FindUserTasksByStatus_ByUser(t *testing.T) {
	svc := setupTaskService(t)
	tasks, err := svc.FindUserTasksByStatus("", "1")
	if err != nil {
		t.Fatalf("FindUserTasksByStatus: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("len(tasks) = %d, want 1", len(tasks))
	}
}

func TestTaskService_FindUserTasksByStatus_ByStatus(t *testing.T) {
	svc := setupTaskService(t)
	tasks, err := svc.FindUserTasksByStatus("completed", "")
	if err != nil {
		t.Fatalf("FindUserTasksByStatus: %v", err)
	}
	if len(tasks) != 1 || tasks[0].Status != "completed" {
		t.Errorf("tasks = %+v", tasks)
	}
}

func TestTaskService_FindUserTasksByStatus_InvalidUserID(t *testing.T) {
	svc := setupTaskService(t)
	_, err := svc.FindUserTasksByStatus("", "abc")
	if err == nil {
		t.Fatal("expected error for invalid userId")
	}
}

func TestTaskService_CreateTask(t *testing.T) {
	svc := setupTaskService(t)
	task, err := svc.CreateTask("New Task", "pending", 1)
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if task.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if task.Title != "New Task" || task.UserID != 1 {
		t.Errorf("task = %+v", task)
	}
}

func TestTaskService_CreateTask_UserNotFound(t *testing.T) {
	svc := setupTaskService(t)
	_, err := svc.CreateTask("Task", "pending", 999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestTaskService_CreateTask_ValidationError(t *testing.T) {
	svc := setupTaskService(t)
	_, err := svc.CreateTask("", "pending", 1)
	if err == nil {
		t.Fatal("expected error for empty title")
	}
	_, err = svc.CreateTask("x", "invalid-status", 1)
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
}

func TestTaskService_FindTaskByID(t *testing.T) {
	svc := setupTaskService(t)
	task, err := svc.FindTaskByID(1)
	if err != nil {
		t.Fatalf("FindTaskByID: %v", err)
	}
	if task == nil || task.Title != "Task 1" {
		t.Errorf("task = %+v", task)
	}
}

func TestTaskService_FindTaskByID_NotFound(t *testing.T) {
	svc := setupTaskService(t)
	task, err := svc.FindTaskByID(999)
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}
	if task != nil {
		t.Error("expected nil task")
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	svc := setupTaskService(t)
	title, status := "Updated Title", "in-progress"
	task, err := svc.UpdateTask(1, &title, &status, nil)
	if err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}
	if task.Title != "Updated Title" || task.Status != "in-progress" {
		t.Errorf("task = %+v", task)
	}
}

func TestTaskService_UpdateTask_Partial(t *testing.T) {
	svc := setupTaskService(t)
	status := "completed"
	task, err := svc.UpdateTask(1, nil, &status, nil)
	if err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}
	if task.Status != "completed" {
		t.Errorf("task.Status = %q", task.Status)
	}
}

func TestTaskService_UpdateTask_UserNotFound(t *testing.T) {
	svc := setupTaskService(t)
	uid := 999
	_, err := svc.UpdateTask(1, nil, nil, &uid)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}
