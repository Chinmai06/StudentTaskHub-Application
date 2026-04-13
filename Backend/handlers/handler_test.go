package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"studenttaskhub/database"
	"studenttaskhub/models"
	"testing"

	"github.com/gorilla/mux"
)

// setupTestDB initializes an in-memory SQLite database for testing
func setupTestDB(t *testing.T) {
	database.InitDB(":memory:")
}

// createTestUser registers a user directly for test setup
func createTestUser(t *testing.T, username, email, password string) {
	body := map[string]string{
		"username": username,
		"email":    email,
		"ufid":     "12345678",
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	Register(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create test user %s: %d", username, rr.Code)
	}
}

// createTestTask creates a task and returns its ID
func createTestTask(t *testing.T, title, description, deadline, priority, createdBy string) int {
	body := models.CreateTaskRequest{
		Title:       title,
		Category:    "Study",
		Description: description,
		Location:    "Library",
		Deadline:    deadline,
		Priority:    priority,
		CreatedBy:   createdBy,
	}
	if priority == "" {
		body.Priority = "Medium"
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateTask(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create test task: %d - %s", rr.Code, rr.Body.String())
	}
	var task models.Task
	json.NewDecoder(rr.Body).Decode(&task)
	return task.ID
}

// helper to complete a task for feedback tests
func completeTask(t *testing.T, taskID int, creator, claimer string) {
	claimBody := `{"claimed_by":"` + claimer + `"}`
	claimReq := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(claimBody))
	claimReq.Header.Set("Content-Type", "application/json")
	claimReq = mux.SetURLVars(claimReq, map[string]string{"id": strconv.Itoa(taskID)})
	claimRR := httptest.NewRecorder()
	ClaimTask(claimRR, claimReq)

	doneBody := `{"status":"done"}`
	doneReq := httptest.NewRequest("PATCH", "/api/tasks/1/status?username="+claimer, bytes.NewBufferString(doneBody))
	doneReq.Header.Set("Content-Type", "application/json")
	doneReq = mux.SetURLVars(doneReq, map[string]string{"id": strconv.Itoa(taskID)})
	doneRR := httptest.NewRecorder()
	UpdateTaskStatus(doneRR, doneReq)
}

// ============================================================
// Register Tests
// ============================================================

func TestRegister_Success(t *testing.T) {
	setupTestDB(t)

	body := `{"username":"testuser","email":"test@ufl.edu","ufid":"12345678","password":"pass123"}`
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Register(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rr.Code)
	}

	var resp map[string]string
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp["username"] != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", resp["username"])
	}
}

func TestRegister_MissingFields(t *testing.T) {
	setupTestDB(t)

	body := `{"username":"testuser"}`
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Register(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	setupTestDB(t)

	body := `{"username":"testuser","email":"test@ufl.edu","password":"pass123"}`
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	Register(rr, req)

	body2 := `{"username":"testuser","email":"test2@ufl.edu","password":"pass123"}`
	req2 := httptest.NewRequest("POST", "/api/register", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	Register(rr2, req2)

	if rr2.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", rr2.Code)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	setupTestDB(t)

	body := `{"username":"testuser","email":"test@ufl.edu","password":"abc"}`
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Register(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestRegister_WithUFID(t *testing.T) {
	setupTestDB(t)

	body := `{"username":"testuser","email":"test@ufl.edu","ufid":"87654321","password":"pass123"}`
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Register(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rr.Code)
	}
}

// ============================================================
// Login Tests
// ============================================================

func TestLogin_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")

	body := `{"username":"chinmai","password":"pass123"}`
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var resp models.LoginResponse
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.Username != "chinmai" {
		t.Errorf("Expected username 'chinmai', got '%s'", resp.Username)
	}
	if resp.Email != "chinmai@ufl.edu" {
		t.Errorf("Expected email 'chinmai@ufl.edu', got '%s'", resp.Email)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")

	body := `{"username":"chinmai","password":"wrongpassword"}`
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestLogin_NonExistentUser(t *testing.T) {
	setupTestDB(t)

	body := `{"username":"nobody","password":"pass123"}`
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	setupTestDB(t)

	body := `{"username":"chinmai"}`
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Login(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

// ============================================================
// CreateTask Tests
// ============================================================

func TestCreateTask_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")

	body := `{"title":"ML Assignment","category":"Project","description":"Neural nets","location":"Marston Library","deadline":"2026-03-01","priority":"High","created_by":"chinmai"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateTask(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rr.Code)
	}

	var task models.Task
	json.NewDecoder(rr.Body).Decode(&task)
	if task.Title != "ML Assignment" {
		t.Errorf("Expected title 'ML Assignment', got '%s'", task.Title)
	}
	if task.Category != "Project" {
		t.Errorf("Expected category 'Project', got '%s'", task.Category)
	}
	if task.Location != "Marston Library" {
		t.Errorf("Expected location 'Marston Library', got '%s'", task.Location)
	}
	if task.Priority != "High" {
		t.Errorf("Expected priority 'High', got '%s'", task.Priority)
	}
	if task.Status != "open" {
		t.Errorf("Expected status 'open', got '%s'", task.Status)
	}
}

func TestCreateTask_MissingFields(t *testing.T) {
	setupTestDB(t)

	body := `{"title":"ML Assignment"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestCreateTask_InvalidDeadline(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")

	body := `{"title":"Test","deadline":"not-a-date","created_by":"chinmai"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestCreateTask_DefaultPriority(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")

	body := `{"title":"Read Chapter","deadline":"2026-04-01","created_by":"chinmai"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateTask(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rr.Code)
	}

	var task models.Task
	json.NewDecoder(rr.Body).Decode(&task)
	if task.Priority != "Medium" {
		t.Errorf("Expected default priority 'Medium', got '%s'", task.Priority)
	}
}

func TestCreateTask_UnregisteredUser(t *testing.T) {
	setupTestDB(t)

	body := `{"title":"Test","deadline":"2026-03-01","created_by":"nobody"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestCreateTask_WithCategoryAndLocation(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")

	body := `{"title":"Study Session","category":"Study","description":"Review for exam","location":"Reitz Union","deadline":"2026-04-01","priority":"Low","created_by":"chinmai"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateTask(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rr.Code)
	}

	var task models.Task
	json.NewDecoder(rr.Body).Decode(&task)
	if task.Category != "Study" {
		t.Errorf("Expected category 'Study', got '%s'", task.Category)
	}
	if task.Location != "Reitz Union" {
		t.Errorf("Expected location 'Reitz Union', got '%s'", task.Location)
	}
	if task.Priority != "Low" {
		t.Errorf("Expected priority 'Low', got '%s'", task.Priority)
	}
}

// ============================================================
// GetTasks Tests
// ============================================================

func TestGetTasks_Empty(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest("GET", "/api/tasks", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}

func TestGetTasks_FilterByStatus(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "Task 1", "Desc", "2026-03-01", "High", "chinmai")
	createTestTask(t, "Task 2", "Desc", "2026-04-01", "Medium", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks?status=open", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 2 {
		t.Errorf("Expected 2 open tasks, got %d", len(tasks))
	}
}

func TestGetTasks_FilterByCreatedUser(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	createTestTask(t, "Task 1", "Desc", "2026-03-01", "High", "chinmai")
	createTestTask(t, "Task 2", "Desc", "2026-04-01", "Medium", "chinmai")
	createTestTask(t, "Task 3", "Desc", "2026-05-01", "Low", "alice")

	req := httptest.NewRequest("GET", "/api/tasks?created_by=chinmai", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks by chinmai, got %d", len(tasks))
	}
}

func TestGetTasks_Search(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "ML Assignment", "Neural networks", "2026-03-01", "High", "chinmai")
	createTestTask(t, "Read Chapter 5", "Database textbook", "2026-04-01", "Medium", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks?search=ML", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task matching 'ML', got %d", len(tasks))
	}
}

func TestGetTasks_SearchCaseInsensitive(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "ML Assignment", "Neural networks", "2026-03-01", "High", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks?search=ml", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task matching 'ml' (case-insensitive), got %d", len(tasks))
	}
}

func TestGetTasks_FilterByPriority(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "High Task", "Desc", "2026-03-01", "High", "chinmai")
	createTestTask(t, "Low Task", "Desc", "2026-04-01", "Low", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks?priority=High", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 High priority task, got %d", len(tasks))
	}
}

func TestGetTasks_FilterByDeadlineRange(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "Task 1", "Desc", "2026-03-01", "High", "chinmai")
	createTestTask(t, "Task 2", "Desc", "2026-05-01", "Medium", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks?deadline_before=2026-04-01", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task before 2026-04-01, got %d", len(tasks))
	}
}

func TestGetTasks_SortByDeadline(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "Later Task", "Desc", "2026-05-01", "Medium", "chinmai")
	createTestTask(t, "Earlier Task", "Desc", "2026-03-01", "Medium", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks?sort=deadline", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) >= 2 && tasks[0].Title != "Earlier Task" {
		t.Errorf("Expected 'Earlier Task' first, got '%s'", tasks[0].Title)
	}
}

func TestGetTasks_SortByPriority(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "Low Task", "Desc", "2026-03-01", "Low", "chinmai")
	createTestTask(t, "High Task", "Desc", "2026-04-01", "High", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks?sort=priority", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) >= 2 && tasks[0].Title != "High Task" {
		t.Errorf("Expected 'High Task' first when sorted by priority, got '%s'", tasks[0].Title)
	}
}

// ============================================================
// GetTask Tests
// ============================================================

func TestGetTask_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	taskID := createTestTask(t, "ML Assignment", "Neural nets", "2026-03-01", "High", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(taskID)})
	rr := httptest.NewRecorder()

	GetTask(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest("GET", "/api/tasks/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()

	GetTask(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}

// ============================================================
// UpdateTask Tests
// ============================================================

func TestUpdateTask_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "Old Title", "Desc", "2026-03-01", "High", "chinmai")

	body := `{"title":"New Title","category":"Project","location":"New Location"}`
	req := httptest.NewRequest("PUT", "/api/tasks/1?username=chinmai", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	UpdateTask(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var task models.Task
	json.NewDecoder(rr.Body).Decode(&task)
	if task.Title != "New Title" {
		t.Errorf("Expected title 'New Title', got '%s'", task.Title)
	}
	if task.Category != "Project" {
		t.Errorf("Expected category 'Project', got '%s'", task.Category)
	}
}

func TestUpdateTask_Forbidden(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	createTestTask(t, "Chinmai Task", "Desc", "2026-03-01", "High", "chinmai")

	body := `{"title":"Hacked"}`
	req := httptest.NewRequest("PUT", "/api/tasks/1?username=alice", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	UpdateTask(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rr.Code)
	}
}

// ============================================================
// DeleteTask Tests
// ============================================================

func TestDeleteTask_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "Delete Me", "Desc", "2026-03-01", "Medium", "chinmai")

	req := httptest.NewRequest("DELETE", "/api/tasks/1?username=chinmai", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	DeleteTask(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestDeleteTask_Forbidden(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	createTestTask(t, "Chinmai Task", "Desc", "2026-03-01", "Medium", "chinmai")

	req := httptest.NewRequest("DELETE", "/api/tasks/1?username=alice", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	DeleteTask(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rr.Code)
	}
}

// ============================================================
// ClaimTask Tests
// ============================================================

func TestClaimTask_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	createTestTask(t, "Open Task", "Desc", "2026-03-01", "High", "chinmai")

	body := `{"claimed_by":"alice"}`
	req := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	ClaimTask(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var task models.Task
	json.NewDecoder(rr.Body).Decode(&task)
	if task.Status != "claimed" {
		t.Errorf("Expected status 'claimed', got '%s'", task.Status)
	}
	if task.ClaimedBy != "alice" {
		t.Errorf("Expected claimed_by 'alice', got '%s'", task.ClaimedBy)
	}
}

func TestClaimTask_AlreadyClaimed(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	createTestUser(t, "bob", "bob@ufl.edu", "pass123")
	createTestTask(t, "Open Task", "Desc", "2026-03-01", "High", "chinmai")

	body := `{"claimed_by":"alice"}`
	req := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	ClaimTask(rr, req)

	body2 := `{"claimed_by":"bob"}`
	req2 := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2 = mux.SetURLVars(req2, map[string]string{"id": "1"})
	rr2 := httptest.NewRecorder()
	ClaimTask(rr2, req2)

	if rr2.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", rr2.Code)
	}
}

// ============================================================
// UpdateTaskStatus Tests
// ============================================================

func TestUpdateTaskStatus_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	createTestTask(t, "Task", "Desc", "2026-03-01", "High", "chinmai")

	claimBody := `{"claimed_by":"alice"}`
	claimReq := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(claimBody))
	claimReq.Header.Set("Content-Type", "application/json")
	claimReq = mux.SetURLVars(claimReq, map[string]string{"id": "1"})
	claimRR := httptest.NewRecorder()
	ClaimTask(claimRR, claimReq)

	body := `{"status":"in_progress"}`
	req := httptest.NewRequest("PATCH", "/api/tasks/1/status?username=alice", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	UpdateTaskStatus(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var task models.Task
	json.NewDecoder(rr.Body).Decode(&task)
	if task.Status != "in_progress" {
		t.Errorf("Expected status 'in_progress', got '%s'", task.Status)
	}
}

func TestUpdateTaskStatus_InvalidStatus(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "Task", "Desc", "2026-03-01", "High", "chinmai")

	body := `{"status":"invalid_status"}`
	req := httptest.NewRequest("PATCH", "/api/tasks/1/status?username=chinmai", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	UpdateTaskStatus(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestUpdateTaskStatus_Forbidden(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	createTestUser(t, "bob", "bob@ufl.edu", "pass123")
	createTestTask(t, "Task", "Desc", "2026-03-01", "High", "chinmai")

	claimBody := `{"claimed_by":"alice"}`
	claimReq := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(claimBody))
	claimReq.Header.Set("Content-Type", "application/json")
	claimReq = mux.SetURLVars(claimReq, map[string]string{"id": "1"})
	claimRR := httptest.NewRecorder()
	ClaimTask(claimRR, claimReq)

	body := `{"status":"done"}`
	req := httptest.NewRequest("PATCH", "/api/tasks/1/status?username=bob", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	UpdateTaskStatus(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rr.Code)
	}
}

func TestMarkTaskAsCompleted(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	taskID := createTestTask(t, "Complete Me", "Desc", "2026-03-01", "High", "chinmai")
	completeTask(t, taskID, "chinmai", "alice")

	// Verify it's done
	req := httptest.NewRequest("GET", "/api/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(taskID)})
	rr := httptest.NewRecorder()
	GetTask(rr, req)

	var task models.Task
	json.NewDecoder(rr.Body).Decode(&task)
	if task.Status != "done" {
		t.Errorf("Expected status 'done', got '%s'", task.Status)
	}
}

// ============================================================
// Profile Tests (Sprint 3)
// ============================================================

func TestGetProfile_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")

	req := httptest.NewRequest("GET", "/api/profile/chinmai", nil)
	req = mux.SetURLVars(req, map[string]string{"username": "chinmai"})
	rr := httptest.NewRecorder()

	GetProfile(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest("GET", "/api/profile/nobody", nil)
	req = mux.SetURLVars(req, map[string]string{"username": "nobody"})
	rr := httptest.NewRecorder()

	GetProfile(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")

	body := `{"full_name":"Chinmai Reddy","bio":"CS student","major":"Computer Science","year":"Senior","skills":"Go, Python, React"}`
	req := httptest.NewRequest("PUT", "/api/profile/chinmai?username=chinmai", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"username": "chinmai"})
	rr := httptest.NewRecorder()

	UpdateProfile(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var profile models.Profile
	json.NewDecoder(rr.Body).Decode(&profile)
	if profile.FullName != "Chinmai Reddy" {
		t.Errorf("Expected full_name 'Chinmai Reddy', got '%s'", profile.FullName)
	}
	if profile.Major != "Computer Science" {
		t.Errorf("Expected major 'Computer Science', got '%s'", profile.Major)
	}
}

func TestUpdateProfile_Forbidden(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")

	body := `{"full_name":"Hacked"}`
	req := httptest.NewRequest("PUT", "/api/profile/chinmai?username=alice", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"username": "chinmai"})
	rr := httptest.NewRecorder()

	UpdateProfile(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rr.Code)
	}
}

// ============================================================
// Feedback Tests (Sprint 3)
// ============================================================

func TestAddFeedback_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	taskID := createTestTask(t, "Test Task", "Desc", "2026-03-01", "High", "chinmai")
	completeTask(t, taskID, "chinmai", "alice")

	body := `{"rating":5,"comment":"Great work!"}`
	req := httptest.NewRequest("POST", "/api/tasks/1/feedback?username=alice", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(taskID)})
	rr := httptest.NewRecorder()

	AddFeedback(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rr.Code)
	}

	var fb models.Feedback
	json.NewDecoder(rr.Body).Decode(&fb)
	if fb.Rating != 5 {
		t.Errorf("Expected rating 5, got %d", fb.Rating)
	}
}

func TestAddFeedback_InvalidRating(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	taskID := createTestTask(t, "Test Task", "Desc", "2026-03-01", "High", "chinmai")
	completeTask(t, taskID, "chinmai", "alice")

	body := `{"rating":6,"comment":"Too high"}`
	req := httptest.NewRequest("POST", "/api/tasks/1/feedback?username=alice", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(taskID)})
	rr := httptest.NewRecorder()

	AddFeedback(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestAddFeedback_TaskNotDone(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	taskID := createTestTask(t, "Open Task", "Desc", "2026-03-01", "High", "chinmai")

	body := `{"rating":4,"comment":"Not done yet"}`
	req := httptest.NewRequest("POST", "/api/tasks/1/feedback?username=chinmai", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(taskID)})
	rr := httptest.NewRecorder()

	AddFeedback(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestAddFeedback_DuplicateFeedback(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	taskID := createTestTask(t, "Test Task", "Desc", "2026-03-01", "High", "chinmai")
	completeTask(t, taskID, "chinmai", "alice")

	body := `{"rating":5,"comment":"Great!"}`
	req := httptest.NewRequest("POST", "/api/tasks/1/feedback?username=alice", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(taskID)})
	rr := httptest.NewRecorder()
	AddFeedback(rr, req)

	body2 := `{"rating":3,"comment":"Changed mind"}`
	req2 := httptest.NewRequest("POST", "/api/tasks/1/feedback?username=alice", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2 = mux.SetURLVars(req2, map[string]string{"id": strconv.Itoa(taskID)})
	rr2 := httptest.NewRecorder()
	AddFeedback(rr2, req2)

	if rr2.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", rr2.Code)
	}
}

func TestGetFeedback_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	taskID := createTestTask(t, "Test Task", "Desc", "2026-03-01", "High", "chinmai")
	completeTask(t, taskID, "chinmai", "alice")

	addBody := `{"rating":5,"comment":"Excellent!"}`
	addReq := httptest.NewRequest("POST", "/api/tasks/1/feedback?username=alice", bytes.NewBufferString(addBody))
	addReq.Header.Set("Content-Type", "application/json")
	addReq = mux.SetURLVars(addReq, map[string]string{"id": strconv.Itoa(taskID)})
	addRR := httptest.NewRecorder()
	AddFeedback(addRR, addReq)

	req := httptest.NewRequest("GET", "/api/tasks/1/feedback", nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(taskID)})
	rr := httptest.NewRecorder()

	GetFeedback(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var feedbacks []models.Feedback
	json.NewDecoder(rr.Body).Decode(&feedbacks)
	if len(feedbacks) != 1 {
		t.Errorf("Expected 1 feedback, got %d", len(feedbacks))
	}
}

func TestGetFeedback_Empty(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	taskID := createTestTask(t, "No Feedback", "Desc", "2026-03-01", "Medium", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks/1/feedback", nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(taskID)})
	rr := httptest.NewRecorder()

	GetFeedback(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var feedbacks []models.Feedback
	json.NewDecoder(rr.Body).Decode(&feedbacks)
	if len(feedbacks) != 0 {
		t.Errorf("Expected 0 feedbacks, got %d", len(feedbacks))
	}
}
