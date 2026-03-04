package model

// CreateUserRequest is the request body for creating a user.
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// CreateTaskRequest is the request body for creating a task.
type CreateTaskRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
	UserID int    `json:"userId"`
}

// UpdateTaskRequest is the request body for updating a task (all fields optional).
type UpdateTaskRequest struct {
	Title  *string `json:"title,omitempty"`
	Status *string `json:"status,omitempty"`
	UserID *int    `json:"userId,omitempty"`
}
