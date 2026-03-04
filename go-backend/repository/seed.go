package repository

import (
	"fmt"
	"go-backend/model"
)

func defaultUsers() []model.User {
	return []model.User{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Role: "developer"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Role: "designer"},
		{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", Role: "manager"},
	}
}

func defaultTasks() []model.Task {
	return []model.Task{
		{ID: 1, Title: "Implement authentication", Status: "pending", UserID: 1},
		{ID: 2, Title: "Design user interface", Status: "in-progress", UserID: 2},
		{ID: 3, Title: "Review code changes", Status: "completed", UserID: 3},
	}
}

// SeedIfEmpty seeds default users and tasks only when the store is empty
func SeedIfEmpty(userRepo *UserRepository, taskRepo *TaskRepository) error {
	// Check for nil repositories
	if userRepo == nil || taskRepo == nil {
		return fmt.Errorf("repositories cannot be nil")
	}

	// Check if users exist
	users, err := userRepo.FindAll()
	if err != nil {
		return err
	}

	// Check if tasks exist
	tasks, err := taskRepo.FindAll()
	if err != nil {
		return err
	}

	if len(users) > 0 || len(tasks) > 0 {
		return nil // data already exists, skip seeding
	}

	// Add default users
	for _, user := range defaultUsers() {
		if _, err := userRepo.Create(&user); err != nil {
			return err
		}
	}

	// Add default tasks
	for _, task := range defaultTasks() {
		if _, err := taskRepo.Create(&task); err != nil {
			return err
		}
	}

	return nil
}
