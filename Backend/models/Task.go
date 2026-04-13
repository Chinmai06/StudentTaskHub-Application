package models

import "time"

// Task priority levels (matching frontend)
const (
	PriorityHigh   = "High"
	PriorityMedium = "Medium"
	PriorityLow    = "Low"
)

// Task status values
const (
	StatusOpen       = "open"
	StatusClaimed    = "claimed"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

// Task represents a task in the system
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Priority    string    `json:"priority"`
	Deadline    string    `json:"deadline"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	ClaimedBy   string    `json:"claimed_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTaskRequest is the JSON body for creating a task
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Priority    string `json:"priority"`
	Deadline    string `json:"deadline"`
	CreatedBy   string `json:"created_by"`
}

// UpdateTaskRequest is the JSON body for editing a task
type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Priority    string `json:"priority"`
	Deadline    string `json:"deadline"`
}

// ClaimTaskRequest is the JSON body for claiming a task
type ClaimTaskRequest struct {
	ClaimedBy string `json:"claimed_by"`
}

// UpdateStatusRequest is the JSON body for updating task status
type UpdateStatusRequest struct {
	Status string `json:"status"`
}
