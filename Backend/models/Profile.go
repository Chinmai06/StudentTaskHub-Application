package models

import "time"

// Profile represents a user's profile
type Profile struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Bio       string    `json:"bio"`
	Major     string    `json:"major"`
	Year      string    `json:"year"`
	Skills    string    `json:"skills"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateProfileRequest is the JSON body for updating a profile
type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	Major    string `json:"major"`
	Year     string `json:"year"`
	Skills   string `json:"skills"`
}
