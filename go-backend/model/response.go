package model

// SuccessResponse is the standard API success envelope (data + message).
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse is the standard API error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// ValidationError is the API response for validation errors with field details.
type ValidationError struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
	Code   string            `json:"code"`
}

// UsersResponse is the response for listing users.
type UsersResponse struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

// TasksResponse is the response for listing tasks.
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
	Count int    `json:"count"`
}

// StatsResponse is the response for the stats endpoint.
type StatsResponse struct {
	Users struct {
		Total int `json:"total"`
	} `json:"users"`
	Tasks struct {
		Total      int `json:"total"`
		Pending    int `json:"pending"`
		InProgress int `json:"inProgress"`
		Completed  int `json:"completed"`
	} `json:"tasks"`
}

// CacheStats is the response for the cache stats endpoint.
type CacheStats struct {
	Hits         int `json:"hits"`
	Misses       int `json:"misses"`
	Evictions    int `json:"evictions"`
	TotalEntries int `json:"totalEntries"`
}
