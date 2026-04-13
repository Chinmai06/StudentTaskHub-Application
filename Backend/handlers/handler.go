package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"studenttaskhub/database"
	"studenttaskhub/models"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================
// Helper functions
// ============================================================

// sendJSON writes a JSON response
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// sendError writes an error JSON response
func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, map[string]string{"error": message})
}

// getTaskByID is a helper to fetch a task by its ID
func getTaskByID(id int) (*models.Task, error) {
	var t models.Task
	err := database.DB.QueryRow(
		`SELECT id, title, description, deadline, priority, status, created_by, claimed_by, created_at, updated_at
		 FROM tasks WHERE id = ?`, id,
	).Scan(&t.ID, &t.Title, &t.Description, &t.Deadline, &t.Priority,
		&t.Status, &t.CreatedBy, &t.ClaimedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

<<<<<<< HEAD
// ============================================================
// User handlers (Sprint 2)
// ============================================================

=======
>>>>>>> b7448005 (WIP: local changes before pull)
// Register creates a new user account
// POST /api/register
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" {
		sendError(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	// Validate password length
	if len(req.Password) < 6 {
		sendError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Insert user into database
	_, err = database.DB.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		req.Username, req.Email, string(hashedPassword),
	)
	if err != nil {
		sendError(w, http.StatusConflict, "Username or email already exists")
		return
	}

	// Sprint 3: Auto-create an empty profile for the new user
	_, err = database.DB.Exec(
		"INSERT INTO profiles (username) VALUES (?)",
		req.Username,
	)
	if err != nil {
		// Profile creation failed but user was created - log but don't fail
		// Profile can be created later via PUT /api/profile
	}

	sendJSON(w, http.StatusCreated, map[string]string{
		"message":  "User registered successfully",
		"username": req.Username,
	})
}

// Login authenticates a user
// POST /api/login
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		sendError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	// Get stored password hash from database
	var storedPassword string
	err := database.DB.QueryRow(
		"SELECT password FROM users WHERE username = ?", req.Username,
	).Scan(&storedPassword)

	if err != nil {
		sendError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	sendJSON(w, http.StatusOK, models.LoginResponse{
		Message:  "Login successful",
		Username: req.Username,
	})
}

// ============================================================
// Profile handlers (Sprint 3 - NEW)
// ============================================================

// GetProfile returns a user's profile
// GET /api/profile/{username}
func GetProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	var p models.Profile
	err := database.DB.QueryRow(
		`SELECT id, username, full_name, bio, major, year, skills, created_at, updated_at
		 FROM profiles WHERE username = ?`, username,
	).Scan(&p.ID, &p.Username, &p.FullName, &p.Bio, &p.Major, &p.Year, &p.Skills, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		sendError(w, http.StatusNotFound, "Profile not found")
		return
	}

	sendJSON(w, http.StatusOK, p)
}

// UpdateProfile updates a user's profile
// PUT /api/profile/{username}
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	// Check who is making the request
	reqUsername := r.URL.Query().Get("username")
	if reqUsername == "" {
		sendError(w, http.StatusBadRequest, "Username query parameter is required")
		return
	}
	if reqUsername != username {
		sendError(w, http.StatusForbidden, "You can only edit your own profile")
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	now := time.Now()

	// Try to update existing profile, or insert if it doesn't exist
	result, err := database.DB.Exec(
		`UPDATE profiles SET full_name = ?, bio = ?, major = ?, year = ?, skills = ?, updated_at = ? WHERE username = ?`,
		req.FullName, req.Bio, req.Major, req.Year, req.Skills, now, username,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// Profile doesn't exist, create it
		_, err = database.DB.Exec(
			`INSERT INTO profiles (username, full_name, bio, major, year, skills, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			username, req.FullName, req.Bio, req.Major, req.Year, req.Skills, now,
		)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Failed to create profile")
			return
		}
	}

	// Return updated profile
	var p models.Profile
	database.DB.QueryRow(
		`SELECT id, username, full_name, bio, major, year, skills, created_at, updated_at FROM profiles WHERE username = ?`, username,
	).Scan(&p.ID, &p.Username, &p.FullName, &p.Bio, &p.Major, &p.Year, &p.Skills, &p.CreatedAt, &p.UpdatedAt)

	sendJSON(w, http.StatusOK, p)
}

// ============================================================
// Task handlers (Sprint 1 + Sprint 2 search/filter)
// ============================================================

// CreateTask creates a new task
// POST /api/tasks
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Title == "" || req.Deadline == "" || req.CreatedBy == "" {
		sendError(w, http.StatusBadRequest, "Title, deadline, and created_by are required")
		return
	}

	// Default priority to "normal" if not provided or not "high"
	if req.Priority != models.PriorityHigh {
		req.Priority = models.PriorityNormal
	}

	// Validate deadline format
	if _, err := time.Parse("2006-01-02", req.Deadline); err != nil {
		sendError(w, http.StatusBadRequest, "Deadline must be in YYYY-MM-DD format")
		return
	}

	// Sprint 2: Verify the creator is a registered user
	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.CreatedBy).Scan(&exists)
	if err != nil || exists == 0 {
		sendError(w, http.StatusBadRequest, "User does not exist. Please register first.")
		return
	}

	now := time.Now()
	result, err := database.DB.Exec(
		`INSERT INTO tasks (title, description, deadline, priority, status, created_by, claimed_by, created_at, updated_at)
		 VALUES (?, ?, ?, ?, 'open', ?, '', ?, ?)`,
		req.Title, req.Description, req.Deadline, req.Priority, req.CreatedBy, now, now,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	id, _ := result.LastInsertId()

	task := models.Task{
		ID:          int(id),
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Priority:    req.Priority,
		Status:      models.StatusOpen,
		CreatedBy:   req.CreatedBy,
		ClaimedBy:   "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	sendJSON(w, http.StatusCreated, task)
}

// GetTasks returns all tasks, with optional filters and search
// GET /api/tasks?status=open&sort=deadline&search=ML&priority=high
func GetTasks(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, title, description, deadline, priority, status, created_by, claimed_by, created_at, updated_at FROM tasks WHERE 1=1"
	var args []interface{}

	// Filter by status if provided
	status := r.URL.Query().Get("status")
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	// Filter by created_by if provided
	createdBy := r.URL.Query().Get("created_by")
	if createdBy != "" {
		query += " AND created_by = ?"
		args = append(args, createdBy)
	}

	// Filter by claimed_by if provided
	claimedBy := r.URL.Query().Get("claimed_by")
	if claimedBy != "" {
		query += " AND claimed_by = ?"
		args = append(args, claimedBy)
	}

	// Sprint 2: Filter by priority if provided
	priority := r.URL.Query().Get("priority")
	if priority != "" {
		query += " AND priority = ?"
		args = append(args, priority)
	}

	// Sprint 2: Search by title or description (case-insensitive)
	search := r.URL.Query().Get("search")
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query += " AND (LOWER(title) LIKE ? OR LOWER(description) LIKE ?)"
		args = append(args, searchTerm, searchTerm)
	}

	// Sprint 2: Filter by deadline range
	deadlineBefore := r.URL.Query().Get("deadline_before")
	if deadlineBefore != "" {
		query += " AND deadline <= ?"
		args = append(args, deadlineBefore)
	}
	deadlineAfter := r.URL.Query().Get("deadline_after")
	if deadlineAfter != "" {
		query += " AND deadline >= ?"
		args = append(args, deadlineAfter)
	}

	// Sort by deadline or priority
	sort := r.URL.Query().Get("sort")
	switch sort {
	case "deadline":
		query += " ORDER BY deadline ASC"
	case "priority":
		query += " ORDER BY CASE priority WHEN 'high' THEN 1 WHEN 'normal' THEN 2 END ASC"
	case "newest":
		query += " ORDER BY created_at DESC"
	case "oldest":
		query += " ORDER BY created_at ASC"
	default:
		query += " ORDER BY created_at DESC"
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Deadline, &t.Priority,
			&t.Status, &t.CreatedBy, &t.ClaimedBy, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			continue
		}
		tasks = append(tasks, t)
	}

	sendJSON(w, http.StatusOK, tasks)
}

// GetTask returns a single task by ID
// GET /api/tasks/{id}
func GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := getTaskByID(id)
	if err != nil {
		sendError(w, http.StatusNotFound, "Task not found")
		return
	}

	sendJSON(w, http.StatusOK, task)
}

// UpdateTask edits a task (only the creator can edit)
// PUT /api/tasks/{id}?username=xxx
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Get the existing task
	task, err := getTaskByID(id)
	if err != nil {
		sendError(w, http.StatusNotFound, "Task not found")
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check who is making the request
	username := r.URL.Query().Get("username")
	if username == "" {
		sendError(w, http.StatusBadRequest, "Username query parameter is required")
		return
	}
	if task.CreatedBy != username {
		sendError(w, http.StatusForbidden, "Only the task creator can edit this task")
		return
	}

	// Validate priority if provided
	if req.Priority != "" && req.Priority != models.PriorityHigh && req.Priority != models.PriorityNormal {
		sendError(w, http.StatusBadRequest, "Priority must be 'high' or 'normal'")
		return
	}

	// Validate deadline if provided
	if req.Deadline != "" {
		if _, err := time.Parse("2006-01-02", req.Deadline); err != nil {
			sendError(w, http.StatusBadRequest, "Deadline must be in YYYY-MM-DD format")
			return
		}
	}

	// Use existing values if fields are empty
	if req.Title == "" {
		req.Title = task.Title
	}
	if req.Description == "" {
		req.Description = task.Description
	}
	if req.Deadline == "" {
		req.Deadline = task.Deadline
	}
	if req.Priority == "" {
		req.Priority = task.Priority
	}

	now := time.Now()
	_, err = database.DB.Exec(
		`UPDATE tasks SET title = ?, description = ?, deadline = ?, priority = ?, updated_at = ? WHERE id = ?`,
		req.Title, req.Description, req.Deadline, req.Priority, now, id,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	// Return updated task
	updated, _ := getTaskByID(id)
	sendJSON(w, http.StatusOK, updated)
}

// DeleteTask deletes a task (only the creator can delete)
// DELETE /api/tasks/{id}?username=xxx
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := getTaskByID(id)
	if err != nil {
		sendError(w, http.StatusNotFound, "Task not found")
		return
	}

	// Check who is making the request
	username := r.URL.Query().Get("username")
	if username == "" {
		sendError(w, http.StatusBadRequest, "Username query parameter is required")
		return
	}
	if task.CreatedBy != username {
		sendError(w, http.StatusForbidden, "Only the task creator can delete this task")
		return
	}

	// Sprint 3: Also delete related feedback when task is deleted
	database.DB.Exec("DELETE FROM feedback WHERE task_id = ?", id)

	_, err = database.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

// ClaimTask allows a user to claim an open task
// POST /api/tasks/{id}/claim
func ClaimTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := getTaskByID(id)
	if err != nil {
		sendError(w, http.StatusNotFound, "Task not found")
		return
	}

	// Only open tasks can be claimed
	if task.Status != models.StatusOpen {
		sendError(w, http.StatusConflict, "Task is not open for claiming")
		return
	}

	var req models.ClaimTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ClaimedBy == "" {
		sendError(w, http.StatusBadRequest, "claimed_by is required")
		return
	}

	// Sprint 2: Verify the claimer is a registered user
	var exists int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.ClaimedBy).Scan(&exists)
	if err != nil || exists == 0 {
		sendError(w, http.StatusBadRequest, "User does not exist. Please register first.")
		return
	}

	now := time.Now()
	_, err = database.DB.Exec(
		`UPDATE tasks SET status = 'claimed', claimed_by = ?, updated_at = ? WHERE id = ?`,
		req.ClaimedBy, now, id,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to claim task")
		return
	}

	updated, _ := getTaskByID(id)
	sendJSON(w, http.StatusOK, updated)
}

// UpdateTaskStatus updates the status of a task
// PATCH /api/tasks/{id}/status?username=xxx
func UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := getTaskByID(id)
	if err != nil {
		sendError(w, http.StatusNotFound, "Task not found")
		return
	}

	var req models.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		models.StatusOpen:       true,
		models.StatusClaimed:    true,
		models.StatusInProgress: true,
		models.StatusDone:       true,
	}
	if !validStatuses[req.Status] {
		sendError(w, http.StatusBadRequest, "Status must be 'open', 'claimed', 'in_progress', or 'done'")
		return
	}

	// Check authorization: only creator or claimer can update status
	username := r.URL.Query().Get("username")
	if username == "" {
		sendError(w, http.StatusBadRequest, "Username query parameter is required")
		return
	}
	if task.CreatedBy != username && task.ClaimedBy != username {
		sendError(w, http.StatusForbidden, "Only the task creator or claimer can update status")
		return
	}

	now := time.Now()
	_, err = database.DB.Exec(
		`UPDATE tasks SET status = ?, updated_at = ? WHERE id = ?`,
		req.Status, now, id,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to update status")
		return
	}

	updated, _ := getTaskByID(id)
	sendJSON(w, http.StatusOK, updated)
}

// ============================================================
// Feedback handlers (Sprint 3 - NEW)
// ============================================================

// AddFeedback adds a review/rating to a completed task
// POST /api/tasks/{id}/feedback?username=xxx
func AddFeedback(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Verify task exists
	task, err := getTaskByID(id)
	if err != nil {
		sendError(w, http.StatusNotFound, "Task not found")
		return
	}

	// Only completed tasks can receive feedback
	if task.Status != models.StatusDone {
		sendError(w, http.StatusBadRequest, "Feedback can only be added to completed tasks")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		sendError(w, http.StatusBadRequest, "Username query parameter is required")
		return
	}

	// Only creator or claimer can leave feedback
	if task.CreatedBy != username && task.ClaimedBy != username {
		sendError(w, http.StatusForbidden, "Only the task creator or claimer can leave feedback")
		return
	}

	// Check if user already left feedback on this task
	var feedbackExists int
	database.DB.QueryRow("SELECT COUNT(*) FROM feedback WHERE task_id = ? AND username = ?", id, username).Scan(&feedbackExists)
	if feedbackExists > 0 {
		sendError(w, http.StatusConflict, "You have already submitted feedback for this task")
		return
	}

	var req models.CreateFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		sendError(w, http.StatusBadRequest, "Rating must be between 1 and 5")
		return
	}

	now := time.Now()
	result, err := database.DB.Exec(
		`INSERT INTO feedback (task_id, username, rating, comment, created_at) VALUES (?, ?, ?, ?, ?)`,
		id, username, req.Rating, req.Comment, now,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to add feedback")
		return
	}

	feedbackID, _ := result.LastInsertId()

	feedback := models.Feedback{
		ID:        int(feedbackID),
		TaskID:    id,
		Username:  username,
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: now,
	}

	sendJSON(w, http.StatusCreated, feedback)
}

// GetFeedback returns all feedback for a task
// GET /api/tasks/{id}/feedback
func GetFeedback(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Verify task exists
	_, err = getTaskByID(id)
	if err != nil {
		sendError(w, http.StatusNotFound, "Task not found")
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, task_id, username, rating, comment, created_at FROM feedback WHERE task_id = ? ORDER BY created_at DESC`, id,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch feedback")
		return
	}
	defer rows.Close()

	feedbacks := []models.Feedback{}
	for rows.Next() {
		var f models.Feedback
		err := rows.Scan(&f.ID, &f.TaskID, &f.Username, &f.Rating, &f.Comment, &f.CreatedAt)
		if err != nil {
			continue
		}
		feedbacks = append(feedbacks, f)
	}

	sendJSON(w, http.StatusOK, feedbacks)
}
