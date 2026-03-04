package store

import (
	"go-backend/model"
)

// AppData represents the application data model
// This can be extended with new fields without breaking store implementations
type AppData struct {
	Users []model.User `json:"users" db:"users"`
	Tasks []model.Task `json:"tasks" db:"tasks"`
}

// Ensure AppData implements the Data interface
var _ Data = (*AppData)(nil)

// NewAppData creates a new empty AppData instance
func NewAppData() AppData {
	return AppData{
		Users: make([]model.User, 0),
		Tasks: make([]model.Task, 0),
	}
}

// Clone creates a deep copy of AppData
func (d AppData) Clone() AppData {
	// Clone users
	users := make([]model.User, len(d.Users))
	copy(users, d.Users)

	// Clone tasks
	tasks := make([]model.Task, len(d.Tasks))
	copy(tasks, d.Tasks)

	return AppData{
		Users: users,
		Tasks: tasks,
	}
}

// IsEmpty returns true if the data contains no users or tasks
func (d AppData) IsEmpty() bool {
	return len(d.Users) == 0 && len(d.Tasks) == 0
}

// GetUserByID finds a user by ID
func (d AppData) GetUserByID(id int) *model.User {
	for _, user := range d.Users {
		if user.ID == id {
			return &user
		}
	}
	return nil
}

// GetTaskByID finds a task by ID
func (d AppData) GetTaskByID(id int) *model.Task {
	for _, task := range d.Tasks {
		if task.ID == id {
			return &task
		}
	}
	return nil
}

// AddUser adds a user to the data
func (d *AppData) AddUser(user model.User) {
	d.Users = append(d.Users, user)
}

// AddTask adds a task to the data
func (d *AppData) AddTask(task model.Task) {
	d.Tasks = append(d.Tasks, task)
}

// UpdateUser updates a user by ID
func (d *AppData) UpdateUser(id int, user model.User) bool {
	for i, existing := range d.Users {
		if existing.ID == id {
			d.Users[i] = user
			return true
		}
	}
	return false
}

// UpdateTask updates a task by ID
func (d *AppData) UpdateTask(id int, task model.Task) bool {
	for i, existing := range d.Tasks {
		if existing.ID == id {
			d.Tasks[i] = task
			return true
		}
	}
	return false
}

// DeleteUser removes a user by ID
func (d *AppData) DeleteUser(id int) bool {
	for i, user := range d.Users {
		if user.ID == id {
			d.Users = append(d.Users[:i], d.Users[i+1:]...)
			return true
		}
	}
	return false
}

// DeleteTask removes a task by ID
func (d *AppData) DeleteTask(id int) bool {
	for i, task := range d.Tasks {
		if task.ID == id {
			d.Tasks = append(d.Tasks[:i], d.Tasks[i+1:]...)
			return true
		}
	}
	return false
}
