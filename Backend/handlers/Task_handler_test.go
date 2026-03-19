package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"studenttaskhub/database"
)

func TestGetTasksByCreator(t *testing.T) {

	// Initialize test database
	database.InitDB("test.db")

	req, err := http.NewRequest("GET", "/api/tasks?created_by=alim", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(GetTasks)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 but got %v", rr.Code)
	}
}
