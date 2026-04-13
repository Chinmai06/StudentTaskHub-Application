package models

import "time"

// Feedback represents a review/rating on a task
type Feedback struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	Username  string    `json:"username"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateFeedbackRequest is the JSON body for submitting feedback
type CreateFeedbackRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}
