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
		Description: description,
		Deadline:    deadline,
		Priority:    priority,
		CreatedBy:   createdBy,
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateTask(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create test task: %d", rr.Code)
	}
	var task models.Task
	json.NewDecoder(rr.Body).Decode(&task)
	return task.ID
}

// ============================================================
// Register Tests
// ============================================================

func TestRegister_Success(t *testing.T) {
	setupTestDB(t)

	body := `{"username":"testuser","email":"test@ufl.edu","password":"pass123"}`
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

	// Try to register again with same username
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

	body := `{"title":"ML Assignment","description":"Neural nets","deadline":"2026-03-01","priority":"high","created_by":"chinmai"}`
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
	if task.Priority != "high" {
		t.Errorf("Expected priority 'high', got '%s'", task.Priority)
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
	if task.Priority != "normal" {
		t.Errorf("Expected default priority 'normal', got '%s'", task.Priority)
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

// ============================================================
// GetTasks Tests
// ============================================================
func TestGetTasks_FilterByCreatedUser(t *testing.T) {
	setupTestDB(t)

	createTestUser(t, "nikita", "aliminati.reddy@ufl.edu", "123456")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")

	createTestTask(t, "Task 1", "Desc", "2026-03-01", "high", "nikita")
	createTestTask(t, "Task 2", "Desc", "2026-04-01", "", "nikita")
	createTestTask(t, "Task 3", "Desc", "2026-05-01", "", "alice")

	req := httptest.NewRequest("GET", "/api/tasks?created_by=nikita", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks created by nikita, got %d", len(tasks))
	}

	for _, task := range tasks {
		if task.CreatedBy != "nikita" {
			t.Errorf("Expected created_by 'nikita', got '%s'", task.CreatedBy)
		}
	}
}

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
	createTestTask(t, "Task 1", "Desc", "2026-03-01", "high", "chinmai")
	createTestTask(t, "Task 2", "Desc", "2026-04-01", "", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks?status=open", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 2 {
		t.Errorf("Expected 2 open tasks, got %d", len(tasks))
	}
}

func TestGetTasks_Search(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "ML Assignment", "Neural networks", "2026-03-01", "high", "chinmai")
	createTestTask(t, "Read Chapter 5", "Database textbook", "2026-04-01", "", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks?search=ML", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task matching 'ML', got %d", len(tasks))
	}
	if len(tasks) > 0 && tasks[0].Title != "ML Assignment" {
		t.Errorf("Expected 'ML Assignment', got '%s'", tasks[0].Title)
	}
}
func TestGetTasks_FilterByClaimedUser(t *testing.T) {
	setupTestDB(t)

	createTestUser(t, "nikita", "aliminati.reddy@ufl.edu", "123456")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")

	createTestTask(t, "Task 1", "Desc", "2026-03-01", "high", "nikita")
	createTestTask(t, "Task 2", "Desc", "2026-04-01", "", "nikita")

	claimBody := `{"claimed_by":"alice"}`
	reqClaim := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(claimBody))
	reqClaim.Header.Set("Content-Type", "application/json")
	reqClaim = mux.SetURLVars(reqClaim, map[string]string{"id": "1"})
	rrClaim := httptest.NewRecorder()
	ClaimTask(rrClaim, reqClaim)

	req := httptest.NewRequest("GET", "/api/tasks?claimed_by=alice", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)

	if len(tasks) != 1 {
		t.Errorf("Expected 1 claimed task, got %d", len(tasks))
	}

	if len(tasks) > 0 {
		if tasks[0].ClaimedBy != "alice" {
			t.Errorf("Expected claimed_by 'alice', got '%s'", tasks[0].ClaimedBy)
		}
		if tasks[0].Status != "claimed" {
			t.Errorf("Expected status 'claimed', got '%s'", tasks[0].Status)
		}
	}
}

func TestSearchTasks_ByKeyword(t *testing.T) {
	setupTestDB(t)

	createTestUser(t, "nikita", "aliminati.reddy@ufl.edu", "123456")

	createTestTask(t, "Machine Learning Project", "Neural networks", "2026-03-01", "high", "nikita")
	createTestTask(t, "Database Assignment", "SQL queries", "2026-04-01", "", "nikita")
	createTestTask(t, "AI Homework", "ML basics", "2026-05-01", "", "nikita")

	req := httptest.NewRequest("GET", "/api/tasks?search=ML", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task matching 'ML', got %d", len(tasks))
	}

	found := false
	for _, task := range tasks {
		if task.Title == "AI Homework" {
			found = true
		}
	}

	if !found {
		t.Errorf("Expected matching task not found in results")
	}
}
func TestGetTasks_SearchCaseInsensitive(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "ML Assignment", "Neural networks", "2026-03-01", "high", "chinmai")

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
	createTestTask(t, "High Task", "Desc", "2026-03-01", "high", "chinmai")
	createTestTask(t, "Normal Task", "Desc", "2026-04-01", "", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks?priority=high", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 high priority task, got %d", len(tasks))
	}
}

func TestGetTasks_FilterByDeadlineRange(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestTask(t, "Task 1", "Desc", "2026-03-01", "high", "chinmai")
	createTestTask(t, "Task 2", "Desc", "2026-05-01", "", "chinmai")

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
	createTestTask(t, "Later Task", "Desc", "2026-05-01", "", "chinmai")
	createTestTask(t, "Earlier Task", "Desc", "2026-03-01", "", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks?sort=deadline", nil)
	rr := httptest.NewRecorder()

	GetTasks(rr, req)

	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) >= 2 && tasks[0].Title != "Earlier Task" {
		t.Errorf("Expected 'Earlier Task' first when sorted by deadline, got '%s'", tasks[0].Title)
	}
}

// ============================================================
// GetTask Tests
// ============================================================

func TestGetTask_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	taskID := createTestTask(t, "ML Assignment", "Neural nets", "2026-03-01", "high", "chinmai")

	req := httptest.NewRequest("GET", "/api/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	GetTask(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var task models.Task
	json.NewDecoder(rr.Body).Decode(&task)
	if task.ID != taskID {
		t.Errorf("Expected task ID %d, got %d", taskID, task.ID)
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
	createTestTask(t, "Old Title", "Desc", "2026-03-01", "high", "chinmai")

	body := `{"title":"New Title"}`
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
}

func TestUpdateTask_Forbidden(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	createTestTask(t, "Chinmai Task", "Desc", "2026-03-01", "high", "chinmai")

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
	createTestTask(t, "Delete Me", "Desc", "2026-03-01", "", "chinmai")

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
	createTestTask(t, "Chinmai Task", "Desc", "2026-03-01", "", "chinmai")

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
	createTestTask(t, "Open Task", "Desc", "2026-03-01", "high", "chinmai")

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
	createTestTask(t, "Open Task", "Desc", "2026-03-01", "high", "chinmai")

	// First claim
	body := `{"claimed_by":"alice"}`
	req := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	ClaimTask(rr, req)

	// Second claim should fail
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
	createTestTask(t, "Task", "Desc", "2026-03-01", "high", "chinmai")

	// Claim first
	claimBody := `{"claimed_by":"alice"}`
	claimReq := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(claimBody))
	claimReq.Header.Set("Content-Type", "application/json")
	claimReq = mux.SetURLVars(claimReq, map[string]string{"id": "1"})
	claimRR := httptest.NewRecorder()
	ClaimTask(claimRR, claimReq)

	// Update status
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
	createTestTask(t, "Task", "Desc", "2026-03-01", "high", "chinmai")

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
	createTestTask(t, "Task", "Desc", "2026-03-01", "high", "chinmai")

	// Claim as alice
	claimBody := `{"claimed_by":"alice"}`
	claimReq := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(claimBody))
	claimReq.Header.Set("Content-Type", "application/json")
	claimReq = mux.SetURLVars(claimReq, map[string]string{"id": "1"})
	claimRR := httptest.NewRecorder()
	ClaimTask(claimRR, claimReq)

	// Bob tries to update status - should fail
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

	createTestUser(t, "nikita", "aliminati.reddy@ufl.edu", "123456")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")

	createTestTask(t, "Complete Me", "Desc", "2026-03-01", "high", "nikita")

	claimBody := `{"claimed_by":"alice"}`
	reqClaim := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(claimBody))
	reqClaim.Header.Set("Content-Type", "application/json")
	reqClaim = mux.SetURLVars(reqClaim, map[string]string{"id": "1"})
	rrClaim := httptest.NewRecorder()
	ClaimTask(rrClaim, reqClaim)

	updateBody := `{"status":"done"}`
	req := httptest.NewRequest("PATCH", "/api/tasks/1/status?username=alice", bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	UpdateTaskStatus(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var task models.Task
	json.NewDecoder(rr.Body).Decode(&task)

	if task.Status != "done" {
		t.Errorf("Expected status 'done', got '%s'", task.Status)
	}
}

// ============================================================
// Sprint 3: Profile Tests
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

	var profile models.Profile
	json.NewDecoder(rr.Body).Decode(&profile)
	if profile.Username != "chinmai" {
		t.Errorf("Expected username 'chinmai', got '%s'", profile.Username)
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

	body := `{"full_name":"Chinmai Reddy","bio":"CS student at UF","major":"Computer Science","year":"Senior","skills":"Go, Python, React"}`
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
	if profile.Skills != "Go, Python, React" {
		t.Errorf("Expected skills 'Go, Python, React', got '%s'", profile.Skills)
	}
}

func TestUpdateProfile_Forbidden(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")

	body := `{"full_name":"Hacked Name"}`
	req := httptest.NewRequest("PUT", "/api/profile/chinmai?username=alice", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"username": "chinmai"})
	rr := httptest.NewRecorder()

	UpdateProfile(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rr.Code)
	}
}

func TestUpdateProfile_MissingUsername(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")

	body := `{"full_name":"Test"}`
	req := httptest.NewRequest("PUT", "/api/profile/chinmai", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"username": "chinmai"})
	rr := httptest.NewRecorder()

	UpdateProfile(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

// ============================================================
// Sprint 3: Feedback Tests
// ============================================================

// helper to complete a task for feedback tests
func completeTask(t *testing.T, taskID int, creator, claimer string) {
	// Claim
	claimBody := `{"claimed_by":"` + claimer + `"}`
	claimReq := httptest.NewRequest("POST", "/api/tasks/1/claim", bytes.NewBufferString(claimBody))
	claimReq.Header.Set("Content-Type", "application/json")
	claimReq = mux.SetURLVars(claimReq, map[string]string{"id": strconv.Itoa(taskID)})
	claimRR := httptest.NewRecorder()
	ClaimTask(claimRR, claimReq)

	// Mark as done
	doneBody := `{"status":"done"}`
	doneReq := httptest.NewRequest("PATCH", "/api/tasks/1/status?username="+claimer, bytes.NewBufferString(doneBody))
	doneReq.Header.Set("Content-Type", "application/json")
	doneReq = mux.SetURLVars(doneReq, map[string]string{"id": strconv.Itoa(taskID)})
	doneRR := httptest.NewRecorder()
	UpdateTaskStatus(doneRR, doneReq)
}

func TestAddFeedback_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	taskID := createTestTask(t, "Test Task", "Desc", "2026-03-01", "high", "chinmai")
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
	if fb.Comment != "Great work!" {
		t.Errorf("Expected comment 'Great work!', got '%s'", fb.Comment)
	}
}

func TestAddFeedback_InvalidRating(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	taskID := createTestTask(t, "Test Task", "Desc", "2026-03-01", "high", "chinmai")
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
	taskID := createTestTask(t, "Open Task", "Desc", "2026-03-01", "high", "chinmai")

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
	taskID := createTestTask(t, "Test Task", "Desc", "2026-03-01", "high", "chinmai")
	completeTask(t, taskID, "chinmai", "alice")

	// First feedback
	body := `{"rating":5,"comment":"Great!"}`
	req := httptest.NewRequest("POST", "/api/tasks/1/feedback?username=alice", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(taskID)})
	rr := httptest.NewRecorder()
	AddFeedback(rr, req)

	// Duplicate feedback should fail
	body2 := `{"rating":3,"comment":"Changed my mind"}`
	req2 := httptest.NewRequest("POST", "/api/tasks/1/feedback?username=alice", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2 = mux.SetURLVars(req2, map[string]string{"id": strconv.Itoa(taskID)})
	rr2 := httptest.NewRecorder()
	AddFeedback(rr2, req2)

	if rr2.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", rr2.Code)
	}
}

func TestAddFeedback_Forbidden(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	createTestUser(t, "bob", "bob@ufl.edu", "pass123")
	taskID := createTestTask(t, "Test Task", "Desc", "2026-03-01", "high", "chinmai")
	completeTask(t, taskID, "chinmai", "alice")

	// Bob is not creator or claimer
	body := `{"rating":4,"comment":"Not my task"}`
	req := httptest.NewRequest("POST", "/api/tasks/1/feedback?username=bob", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(taskID)})
	rr := httptest.NewRecorder()

	AddFeedback(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rr.Code)
	}
}

func TestGetFeedback_Success(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	createTestUser(t, "alice", "alice@ufl.edu", "pass123")
	taskID := createTestTask(t, "Test Task", "Desc", "2026-03-01", "high", "chinmai")
	completeTask(t, taskID, "chinmai", "alice")

	// Add feedback
	body := `{"rating":5,"comment":"Excellent!"}`
	addReq := httptest.NewRequest("POST", "/api/tasks/1/feedback?username=alice", bytes.NewBufferString(body))
	addReq.Header.Set("Content-Type", "application/json")
	addReq = mux.SetURLVars(addReq, map[string]string{"id": strconv.Itoa(taskID)})
	addRR := httptest.NewRecorder()
	AddFeedback(addRR, addReq)

	// Get feedback
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
	if len(feedbacks) > 0 && feedbacks[0].Rating != 5 {
		t.Errorf("Expected rating 5, got %d", feedbacks[0].Rating)
	}
}

func TestGetFeedback_Empty(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "chinmai", "chinmai@ufl.edu", "pass123")
	taskID := createTestTask(t, "No Feedback Task", "Desc", "2026-03-01", "", "chinmai")

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

func TestGetFeedback_TaskNotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest("GET", "/api/tasks/999/feedback", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()

	GetFeedback(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}