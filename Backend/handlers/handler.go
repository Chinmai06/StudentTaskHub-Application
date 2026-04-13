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

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, map[string]string{"error": message})
}

func getTaskByID(id int) (*models.Task, error) {
	var t models.Task
	err := database.DB.QueryRow(
		`SELECT id, title, category, description, location, priority, deadline, status, created_by, claimed_by, created_at, updated_at
		 FROM tasks WHERE id = ?`, id,
	).Scan(&t.ID, &t.Title, &t.Category, &t.Description, &t.Location, &t.Priority,
		&t.Deadline, &t.Status, &t.CreatedBy, &t.ClaimedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ============================================================
// User handlers
// ============================================================

// POST /api/register
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		sendError(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	if len(req.Password) < 6 {
		sendError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	_, err = database.DB.Exec(
		"INSERT INTO users (username, email, ufid, password) VALUES (?, ?, ?, ?)",
		req.Username, req.Email, req.UFID, string(hashedPassword),
	)
	if err != nil {
		sendError(w, http.StatusConflict, "Username or email already exists")
		return
	}

	// Auto-create profile
	database.DB.Exec("INSERT INTO profiles (username) VALUES (?)", req.Username)

	sendJSON(w, http.StatusCreated, map[string]string{
		"message":  "User registered successfully",
		"username": req.Username,
	})
}

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

	var storedPassword, email, ufid string
	err := database.DB.QueryRow(
		"SELECT password, email, ufid FROM users WHERE username = ?", req.Username,
	).Scan(&storedPassword, &email, &ufid)

	if err != nil {
		sendError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	sendJSON(w, http.StatusOK, models.LoginResponse{
		Message:  "Login successful",
		Username: req.Username,
		Email:    email,
		UFID:     ufid,
	})
}

// ============================================================
// Profile handlers
// ============================================================

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

// PUT /api/profile/{username}?username=xxx
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

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
		_, err = database.DB.Exec(
			`INSERT INTO profiles (username, full_name, bio, major, year, skills, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			username, req.FullName, req.Bio, req.Major, req.Year, req.Skills, now,
		)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Failed to create profile")
			return
		}
	}

	var p models.Profile
	database.DB.QueryRow(
		`SELECT id, username, full_name, bio, major, year, skills, created_at, updated_at FROM profiles WHERE username = ?`, username,
	).Scan(&p.ID, &p.Username, &p.FullName, &p.Bio, &p.Major, &p.Year, &p.Skills, &p.CreatedAt, &p.UpdatedAt)

	sendJSON(w, http.StatusOK, p)
}

// ============================================================
// Task handlers
// ============================================================

// POST /api/tasks
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Deadline == "" || req.CreatedBy == "" {
		sendError(w, http.StatusBadRequest, "Title, deadline, and created_by are required")
		return
	}

	// Default priority to Medium if not provided or invalid
	validPriorities := map[string]bool{"High": true, "Medium": true, "Low": true}
	if !validPriorities[req.Priority] {
		req.Priority = models.PriorityMedium
	}

	// Default category
	if req.Category == "" {
		req.Category = "Study"
	}

	if _, err := time.Parse("2006-01-02", req.Deadline); err != nil {
		sendError(w, http.StatusBadRequest, "Deadline must be in YYYY-MM-DD format")
		return
	}

	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.CreatedBy).Scan(&exists)
	if err != nil || exists == 0 {
		sendError(w, http.StatusBadRequest, "User does not exist. Please register first.")
		return
	}

	now := time.Now()
	result, err := database.DB.Exec(
		`INSERT INTO tasks (title, category, description, location, priority, deadline, status, created_by, claimed_by, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, 'open', ?, '', ?, ?)`,
		req.Title, req.Category, req.Description, req.Location, req.Priority, req.Deadline, req.CreatedBy, now, now,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	id, _ := result.LastInsertId()

	task := models.Task{
		ID:          int(id),
		Title:       req.Title,
		Category:    req.Category,
		Description: req.Description,
		Location:    req.Location,
		Priority:    req.Priority,
		Deadline:    req.Deadline,
		Status:      models.StatusOpen,
		CreatedBy:   req.CreatedBy,
		ClaimedBy:   "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	sendJSON(w, http.StatusCreated, task)
}

// GET /api/tasks?status=open&sort=deadline&search=ML&priority=High&category=Study
func GetTasks(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, title, category, description, location, priority, deadline, status, created_by, claimed_by, created_at, updated_at FROM tasks WHERE 1=1"
	var args []interface{}

	status := r.URL.Query().Get("status")
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	createdBy := r.URL.Query().Get("created_by")
	if createdBy != "" {
		query += " AND created_by = ?"
		args = append(args, createdBy)
	}

	claimedBy := r.URL.Query().Get("claimed_by")
	if claimedBy != "" {
		query += " AND claimed_by = ?"
		args = append(args, claimedBy)
	}

	priority := r.URL.Query().Get("priority")
	if priority != "" {
		query += " AND priority = ?"
		args = append(args, priority)
	}

	category := r.URL.Query().Get("category")
	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	search := r.URL.Query().Get("search")
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query += " AND (LOWER(title) LIKE ? OR LOWER(description) LIKE ?)"
		args = append(args, searchTerm, searchTerm)
	}

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

	sort := r.URL.Query().Get("sort")
	switch sort {
	case "deadline":
		query += " ORDER BY deadline ASC"
	case "priority":
		query += " ORDER BY CASE priority WHEN 'High' THEN 1 WHEN 'Medium' THEN 2 WHEN 'Low' THEN 3 END ASC"
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
		err := rows.Scan(&t.ID, &t.Title, &t.Category, &t.Description, &t.Location, &t.Priority,
			&t.Deadline, &t.Status, &t.CreatedBy, &t.ClaimedBy, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			continue
		}
		tasks = append(tasks, t)
	}

	sendJSON(w, http.StatusOK, tasks)
}

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

// PUT /api/tasks/{id}?username=xxx
func UpdateTask(w http.ResponseWriter, r *http.Request) {
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

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

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
	if req.Priority != "" {
		validPriorities := map[string]bool{"High": true, "Medium": true, "Low": true}
		if !validPriorities[req.Priority] {
			sendError(w, http.StatusBadRequest, "Priority must be 'High', 'Medium', or 'Low'")
			return
		}
	}

	if req.Deadline != "" {
		if _, err := time.Parse("2006-01-02", req.Deadline); err != nil {
			sendError(w, http.StatusBadRequest, "Deadline must be in YYYY-MM-DD format")
			return
		}
	}

	// Keep existing values for empty fields
	if req.Title == "" {
		req.Title = task.Title
	}
	if req.Category == "" {
		req.Category = task.Category
	}
	if req.Description == "" {
		req.Description = task.Description
	}
	if req.Location == "" {
		req.Location = task.Location
	}
	if req.Deadline == "" {
		req.Deadline = task.Deadline
	}
	if req.Priority == "" {
		req.Priority = task.Priority
	}

	now := time.Now()
	_, err = database.DB.Exec(
		`UPDATE tasks SET title = ?, category = ?, description = ?, location = ?, priority = ?, deadline = ?, updated_at = ? WHERE id = ?`,
		req.Title, req.Category, req.Description, req.Location, req.Priority, req.Deadline, now, id,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	updated, _ := getTaskByID(id)
	sendJSON(w, http.StatusOK, updated)
}

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

	username := r.URL.Query().Get("username")
	if username == "" {
		sendError(w, http.StatusBadRequest, "Username query parameter is required")
		return
	}
	if task.CreatedBy != username {
		sendError(w, http.StatusForbidden, "Only the task creator can delete this task")
		return
	}

	database.DB.Exec("DELETE FROM feedback WHERE task_id = ?", id)
	_, err = database.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

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

	validStatuses := map[string]bool{
		models.StatusOpen: true, models.StatusClaimed: true,
		models.StatusInProgress: true, models.StatusDone: true,
	}
	if !validStatuses[req.Status] {
		sendError(w, http.StatusBadRequest, "Status must be 'open', 'claimed', 'in_progress', or 'done'")
		return
	}

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
	_, err = database.DB.Exec(`UPDATE tasks SET status = ?, updated_at = ? WHERE id = ?`, req.Status, now, id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to update status")
		return
	}

	updated, _ := getTaskByID(id)
	sendJSON(w, http.StatusOK, updated)
}

// ============================================================
// Feedback handlers
// ============================================================

// POST /api/tasks/{id}/feedback?username=xxx
func AddFeedback(w http.ResponseWriter, r *http.Request) {
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

	if task.Status != models.StatusDone {
		sendError(w, http.StatusBadRequest, "Feedback can only be added to completed tasks")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		sendError(w, http.StatusBadRequest, "Username query parameter is required")
		return
	}

	if task.CreatedBy != username && task.ClaimedBy != username {
		sendError(w, http.StatusForbidden, "Only the task creator or claimer can leave feedback")
		return
	}

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
	sendJSON(w, http.StatusCreated, models.Feedback{
		ID: int(feedbackID), TaskID: id, Username: username,
		Rating: req.Rating, Comment: req.Comment, CreatedAt: now,
	})
}

// GET /api/tasks/{id}/feedback
func GetFeedback(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

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
		rows.Scan(&f.ID, &f.TaskID, &f.Username, &f.Rating, &f.Comment, &f.CreatedAt)
		feedbacks = append(feedbacks, f)
	}

	sendJSON(w, http.StatusOK, feedbacks)
}
