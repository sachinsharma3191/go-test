package repository

import (
	"context"
	"fmt"
	"go-backend/errors"
	"go-backend/model"
	"go-backend/store"
)

// TaskRepository implements persistence for tasks using the generic store
type TaskRepository struct {
	*BaseRepository
}

// NewTaskRepository returns a new TaskRepository backed by the given store
func NewTaskRepository(store store.Store) *TaskRepository {
	return &TaskRepository{
		BaseRepository: NewBaseRepository(store),
	}
}

// FindAll returns all tasks
func (r *TaskRepository) FindAll() ([]model.Task, error) {
	if r.store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	data, err := r.store.Load(context.Background())
	if err != nil {
		return nil, err
	}

	tasks := make([]model.Task, len(data.Tasks))
	copy(tasks, data.Tasks)
	return tasks, nil
}

// FindByID returns the task with the given ID, or nil if not found
func (r *TaskRepository) FindByID(id int) (*model.Task, error) {
	if r.store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	data, err := r.store.Load(context.Background())
	if err != nil {
		return nil, err
	}

	return data.GetTaskByID(id), nil
}

// Create adds a new task, persists, and returns it
func (r *TaskRepository) Create(t *model.Task) (*model.Task, error) {
	if r.store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	data, err := r.store.Load(context.Background())
	if err != nil {
		return nil, err
	}

	newID := 1
	if len(data.Tasks) > 0 {
		newID = data.Tasks[len(data.Tasks)-1].ID + 1
	}
	newTask := model.Task{ID: newID, Title: t.Title, Status: t.Status, UserID: t.UserID}

	data.AddTask(newTask)
	if err := r.store.Save(context.Background(), data); err != nil {
		return nil, errors.NewDataStoreError("Failed to save task data", err)
	}

	return &newTask, nil
}

// Update updates an existing task by ID
func (r *TaskRepository) Update(t *model.Task) (*model.Task, error) {
	if r.store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	data, err := r.store.Load(context.Background())
	if err != nil {
		return nil, err
	}

	// Get existing task
	existing := data.GetTaskByID(t.ID)
	if existing == nil {
		return nil, errors.NewNotFoundError("task", nil)
	}

	// Apply partial updates
	updated := *existing
	if t.Title != "" {
		updated.Title = t.Title
	}
	if t.Status != "" {
		updated.Status = t.Status
	}
	if t.UserID != 0 {
		updated.UserID = t.UserID
	}

	if !data.UpdateTask(t.ID, updated) {
		return nil, errors.NewNotFoundError("task", nil)
	}

	if err := r.store.Save(context.Background(), data); err != nil {
		return nil, errors.NewDataStoreError("Failed to update task data", err)
	}

	return &updated, nil
}

// DeleteByID removes a task by ID
func (r *TaskRepository) DeleteByID(id int) error {
	if r.store == nil {
		return fmt.Errorf("store is nil")
	}
	data, err := r.store.Load(context.Background())
	if err != nil {
		return err
	}

	if !data.DeleteTask(id) {
		return errors.NewNotFoundError("task", nil)
	}

	return r.store.Save(context.Background(), data)
}
