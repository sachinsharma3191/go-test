package model

// Task represents a task entity (API and persistence).
type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	UserID int    `json:"userId"`
}
