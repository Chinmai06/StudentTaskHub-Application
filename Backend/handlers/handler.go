package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"studenttaskhub/database"
	"studenttaskhub/models"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

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

// Sprint 4: Helper to create a notification
func createNotification(username, message string, taskID int) {
	database.DB.Exec(
		`INSERT INTO notifications (username, message, task_id, is_read, created_at) VALUES (?, ?, ?, 0, ?)`,
		username, message, taskID, time.Now(),
	)
}

// Sprint 4: Helper to sanitize input - trims whitespace and limits length
func sanitize(input string, maxLen int) string {
	s := strings.TrimSpace(input)
	if len(s) > maxLen {
		return s[:maxLen]
	}
	return s
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

	req.Username = sanitize(req.Username, 50)
	req.Email = sanitize(req.Email, 100)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		sendError(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	// Sprint 4: Better validation
	if len(req.Username) < 3 {
		sendError(w, http.StatusBadRequest, "Username must be at least 3 characters")
		return
	}

	if !strings.Contains(req.Email, "@") {
		sendError(w, http.StatusBadRequest, "Please enter a valid email address")
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

	req.Username = strings.TrimSpace(req.Username)

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

	// Sprint 4: Sanitize profile inputs
	req.FullName = sanitize(req.FullName, 100)
	req.Bio = sanitize(req.Bio, 500)
	req.Major = sanitize(req.Major, 100)
	req.Year = sanitize(req.Year, 20)
	req.Skills = sanitize(req.Skills, 500)

	now := time.Now()
	result, err := database.DB.Exec(
		`UPDATE profiles SET full_name=?, bio=?, major=?, year=?, skills=?, updated_at=? WHERE username=?`,
		req.FullName, req.Bio, req.Major, req.Year, req.Skills, now, username,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		database.DB.Exec(
			`INSERT INTO profiles (username, full_name, bio, major, year, skills, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			username, req.FullName, req.Bio, req.Major, req.Year, req.Skills, now,
		)
	}

	var p models.Profile
	database.DB.QueryRow(
		`SELECT id, username, full_name, bio, major, year, skills, created_at, updated_at FROM profiles WHERE username=?`, username,
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

	// Sprint 4: Sanitize inputs
	req.Title = sanitize(req.Title, 200)
	req.Description = sanitize(req.Description, 1000)
	req.Location = sanitize(req.Location, 200)
	req.Category = sanitize(req.Category, 50)

	if req.Title == "" || req.Deadline == "" || req.CreatedBy == "" {
		sendError(w, http.StatusBadRequest, "Title, deadline, and created_by are required")
		return
	}

	validPriorities := map[string]bool{"High": true, "Medium": true, "Low": true}
	if !validPriorities[req.Priority] {
		req.Priority = models.PriorityMedium
	}

	validCategories := map[string]bool{"Study": true, "Project": true, "Errand": true, "Event": true}
	if !validCategories[req.Category] {
		req.Category = "Study"
	}

	if _, err := time.Parse("2006-01-02", req.Deadline); err != nil {
		sendError(w, http.StatusBadRequest, "Deadline must be in YYYY-MM-DD format")
		return
	}

	var exists int
	database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", req.CreatedBy).Scan(&exists)
	if exists == 0 {
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
		ID: int(id), Title: req.Title, Category: req.Category, Description: req.Description,
		Location: req.Location, Priority: req.Priority, Deadline: req.Deadline,
		Status: models.StatusOpen, CreatedBy: req.CreatedBy, ClaimedBy: "", CreatedAt: now, UpdatedAt: now,
	}
	sendJSON(w, http.StatusCreated, task)
}

// GET /api/tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, title, category, description, location, priority, deadline, status, created_by, claimed_by, created_at, updated_at FROM tasks WHERE 1=1"
	var args []interface{}

	if v := r.URL.Query().Get("status"); v != "" {
		query += " AND status = ?"
		args = append(args, v)
	}
	if v := r.URL.Query().Get("created_by"); v != "" {
		query += " AND created_by = ?"
		args = append(args, v)
	}
	if v := r.URL.Query().Get("claimed_by"); v != "" {
		query += " AND claimed_by = ?"
		args = append(args, v)
	}
	if v := r.URL.Query().Get("priority"); v != "" {
		query += " AND priority = ?"
		args = append(args, v)
	}
	if v := r.URL.Query().Get("category"); v != "" {
		query += " AND category = ?"
		args = append(args, v)
	}
	if v := r.URL.Query().Get("search"); v != "" {
		searchTerm := "%" + strings.ToLower(v) + "%"
		query += " AND (LOWER(title) LIKE ? OR LOWER(description) LIKE ?)"
		args = append(args, searchTerm, searchTerm)
	}
	if v := r.URL.Query().Get("deadline_before"); v != "" {
		query += " AND deadline <= ?"
		args = append(args, v)
	}
	if v := r.URL.Query().Get("deadline_after"); v != "" {
		query += " AND deadline >= ?"
		args = append(args, v)
	}

	switch r.URL.Query().Get("sort") {
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
		rows.Scan(&t.ID, &t.Title, &t.Category, &t.Description, &t.Location, &t.Priority,
			&t.Deadline, &t.Status, &t.CreatedBy, &t.ClaimedBy, &t.CreatedAt, &t.UpdatedAt)
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

	if req.Priority != "" {
		valid := map[string]bool{"High": true, "Medium": true, "Low": true}
		if !valid[req.Priority] {
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
	database.DB.Exec(
		`UPDATE tasks SET title=?, category=?, description=?, location=?, priority=?, deadline=?, updated_at=? WHERE id=?`,
		req.Title, req.Category, req.Description, req.Location, req.Priority, req.Deadline, now, id,
	)
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
	database.DB.Exec("DELETE FROM notifications WHERE task_id = ?", id)
	database.DB.Exec("DELETE FROM feedback WHERE task_id = ?", id)
	database.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
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
	database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", req.ClaimedBy).Scan(&exists)
	if exists == 0 {
		sendError(w, http.StatusBadRequest, "User does not exist. Please register first.")
		return
	}

	now := time.Now()
	database.DB.Exec(`UPDATE tasks SET status='claimed', claimed_by=?, updated_at=? WHERE id=?`, req.ClaimedBy, now, id)

	// Sprint 4: Notify task creator that someone claimed their task
	createNotification(task.CreatedBy,
		fmt.Sprintf("%s claimed your task: %s", req.ClaimedBy, task.Title), id)

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

	valid := map[string]bool{models.StatusOpen: true, models.StatusClaimed: true, models.StatusInProgress: true, models.StatusDone: true}
	if !valid[req.Status] {
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
	database.DB.Exec(`UPDATE tasks SET status=?, updated_at=? WHERE id=?`, req.Status, now, id)

	// Sprint 4: Notify the other party about status change
	if username == task.ClaimedBy && task.CreatedBy != "" {
		createNotification(task.CreatedBy,
			fmt.Sprintf("Task '%s' status changed to %s by %s", task.Title, req.Status, username), id)
	} else if username == task.CreatedBy && task.ClaimedBy != "" {
		createNotification(task.ClaimedBy,
			fmt.Sprintf("Task '%s' status changed to %s by %s", task.Title, req.Status, username), id)
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

	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM feedback WHERE task_id=? AND username=?", id, username).Scan(&count)
	if count > 0 {
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

	// Sprint 4: Sanitize comment
	req.Comment = sanitize(req.Comment, 500)

	now := time.Now()
	result, _ := database.DB.Exec(
		`INSERT INTO feedback (task_id, username, rating, comment, created_at) VALUES (?, ?, ?, ?, ?)`,
		id, username, req.Rating, req.Comment, now,
	)
	fbID, _ := result.LastInsertId()

	// Sprint 4: Notify the other party about feedback
	if username == task.CreatedBy && task.ClaimedBy != "" {
		createNotification(task.ClaimedBy,
			fmt.Sprintf("%s left feedback on task '%s'", username, task.Title), id)
	} else if username == task.ClaimedBy && task.CreatedBy != "" {
		createNotification(task.CreatedBy,
			fmt.Sprintf("%s left feedback on task '%s'", username, task.Title), id)
	}

	sendJSON(w, http.StatusCreated, models.Feedback{
		ID: int(fbID), TaskID: id, Username: username, Rating: req.Rating, Comment: req.Comment, CreatedAt: now,
	})
}

// GET /api/tasks/{id}/feedback
func GetFeedback(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	if _, err := getTaskByID(id); err != nil {
		sendError(w, http.StatusNotFound, "Task not found")
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, task_id, username, rating, comment, created_at FROM feedback WHERE task_id=? ORDER BY created_at DESC`, id,
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

// ============================================================
// Notification handlers (Sprint 4 - NEW)
// ============================================================

// GET /api/notifications?username=xxx
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		sendError(w, http.StatusBadRequest, "Username query parameter is required")
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, username, message, task_id, is_read, created_at FROM notifications WHERE username=? ORDER BY created_at DESC LIMIT 50`,
		username,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}
	defer rows.Close()

	notifications := []models.Notification{}
	for rows.Next() {
		var n models.Notification
		var isRead int
		rows.Scan(&n.ID, &n.Username, &n.Message, &n.TaskID, &isRead, &n.CreatedAt)
		n.IsRead = isRead == 1
		notifications = append(notifications, n)
	}
	sendJSON(w, http.StatusOK, notifications)
}

// PATCH /api/notifications/{id}/read?username=xxx
func MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		sendError(w, http.StatusBadRequest, "Username query parameter is required")
		return
	}

	// Sprint 4: Privacy - only the notification owner can mark it as read
	var notifUsername string
	err = database.DB.QueryRow("SELECT username FROM notifications WHERE id=?", id).Scan(&notifUsername)
	if err != nil {
		sendError(w, http.StatusNotFound, "Notification not found")
		return
	}
	if notifUsername != username {
		sendError(w, http.StatusForbidden, "You can only read your own notifications")
		return
	}

	database.DB.Exec("UPDATE notifications SET is_read=1 WHERE id=?", id)
	sendJSON(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

// PATCH /api/notifications/read-all?username=xxx
func MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		sendError(w, http.StatusBadRequest, "Username query parameter is required")
		return
	}

	database.DB.Exec("UPDATE notifications SET is_read=1 WHERE username=?", username)
	sendJSON(w, http.StatusOK, map[string]string{"message": "All notifications marked as read"})
}

// GET /api/notifications/unread-count?username=xxx
func GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		sendError(w, http.StatusBadRequest, "Username query parameter is required")
		return
	}

	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM notifications WHERE username=? AND is_read=0", username).Scan(&count)
	sendJSON(w, http.StatusOK, map[string]int{"unread_count": count})
}
