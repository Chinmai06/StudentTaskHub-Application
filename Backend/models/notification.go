package models

import "time"

// Notification represents an in-app notification for a user
type Notification struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	TaskID    int       `json:"task_id"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
