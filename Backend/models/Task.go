package models

import "time"

// Task priority levels
const (
	PriorityHigh   = "high"
	PriorityNormal = "normal"
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
	Description string    `json:"description"`
	Deadline    string    `json:"deadline"`              // format: "2025-12-31"
	Priority    string    `json:"priority"`              // high or normal
	Status      string    `json:"status"`                // open, claimed, in_progress, done
	CreatedBy   string    `json:"created_by"`            // username of the creator
	ClaimedBy   string    `json:"claimed_by,omitempty"`  // username of who claimed it
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTaskRequest is the JSON body for creating a task
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Priority    string `json:"priority"`    // "high" if selected, otherwise defaults to "normal"
	CreatedBy   string `json:"created_by"`
}

// UpdateTaskRequest is the JSON body for editing a task
type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Priority    string `json:"priority"`
}

// ClaimTaskRequest is the JSON body for claiming a task
type ClaimTaskRequest struct {
	ClaimedBy string `json:"claimed_by"`
}

// UpdateStatusRequest is the JSON body for updating task status
type UpdateStatusRequest struct {
	Status string `json:"status"`
}