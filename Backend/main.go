package main

import (
	"log"
	"net/http"
	"studenttaskhub/database"
	"studenttaskhub/handlers"
	"studenttaskhub/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize SQLite database
	database.InitDB("studenttaskhub.db")

	// Create router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// User routes (Sprint 2)
	api.HandleFunc("/register", handlers.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/login", handlers.Login).Methods("POST", "OPTIONS")

	// Profile routes (Sprint 3 - NEW)
	api.HandleFunc("/profile/{username}", handlers.GetProfile).Methods("GET")
	api.HandleFunc("/profile/{username}", handlers.UpdateProfile).Methods("PUT", "OPTIONS")

	// Task routes (Sprint 1)
	api.HandleFunc("/tasks", handlers.CreateTask).Methods("POST", "OPTIONS")
	api.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	api.HandleFunc("/tasks/{id}", handlers.GetTask).Methods("GET")
	api.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT", "OPTIONS")
	api.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/tasks/{id}/claim", handlers.ClaimTask).Methods("POST", "OPTIONS")
	api.HandleFunc("/tasks/{id}/status", handlers.UpdateTaskStatus).Methods("PATCH", "OPTIONS")

	// Feedback routes (Sprint 3 - NEW)
	api.HandleFunc("/tasks/{id}/feedback", handlers.AddFeedback).Methods("POST", "OPTIONS")
	api.HandleFunc("/tasks/{id}/feedback", handlers.GetFeedback).Methods("GET")

	// Apply CORS middleware
	handler := middleware.CORS(r)

	// Start server
	port := ":8080"
	log.Printf("StudentTaskHub API server starting on http://localhost%s", port)
	log.Println("")
	log.Println("=== Sprint 3 Endpoints (NEW) ===")
	log.Println("  GET    /api/profile/{username}")
	log.Println("  PUT    /api/profile/{username}?username=xxx")
	log.Println("  POST   /api/tasks/{id}/feedback?username=xxx")
	log.Println("  GET    /api/tasks/{id}/feedback")
	log.Println("")
	log.Println("=== Sprint 1 & 2 Endpoints ===")
	log.Println("  POST   /api/register")
	log.Println("  POST   /api/login")
	log.Println("  POST   /api/tasks")
	log.Println("  GET    /api/tasks              (query: ?status=open&sort=deadline&search=keyword)")
	log.Println("  GET    /api/tasks/{id}")
	log.Println("  PUT    /api/tasks/{id}?username=xxx")
	log.Println("  DELETE /api/tasks/{id}?username=xxx")
	log.Println("  POST   /api/tasks/{id}/claim")
	log.Println("  PATCH  /api/tasks/{id}/status?username=xxx")

	log.Fatal(http.ListenAndServe(port, handler))
}
