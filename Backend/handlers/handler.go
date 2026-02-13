package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"studenttaskhub/database"
	"studenttaskhub/models"
	"time"

	"github.com/gorilla/mux"
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

// ============================================================
// Task handlers
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

// GetTasks returns all tasks, with optional filters
// GET /api/tasks?status=open&sort=deadline
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

	// Sort by deadline or priority
	sort := r.URL.Query().Get("sort")
	switch sort {
	case "deadline":
		query += " ORDER BY deadline ASC"
	case "priority":
		query += " ORDER BY CASE priority WHEN 'high' THEN 1 WHEN 'normal' THEN 2 END ASC"
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

	// Check who is making the request (passed as query param for Sprint 1)
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
