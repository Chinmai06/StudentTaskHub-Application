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

	// Task routes
	api.HandleFunc("/tasks", handlers.CreateTask).Methods("POST", "OPTIONS")
	api.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	api.HandleFunc("/tasks/{id}", handlers.GetTask).Methods("GET")
	api.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT", "OPTIONS")
	api.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/tasks/{id}/claim", handlers.ClaimTask).Methods("POST", "OPTIONS")
	api.HandleFunc("/tasks/{id}/status", handlers.UpdateTaskStatus).Methods("PATCH", "OPTIONS")

	// Apply CORS middleware
	handler := middleware.CORS(r)

	// Start server
	port := ":8080"
	log.Printf("StudentTaskHub API server starting on http://localhost%s", port)
	log.Println("Available endpoints:")
	log.Println("  POST   /api/tasks")
	log.Println("  GET    /api/tasks              (query: ?status=open&sort=deadline)")
	log.Println("  GET    /api/tasks/{id}")
	log.Println("  PUT    /api/tasks/{id}?username=xxx")
	log.Println("  DELETE /api/tasks/{id}?username=xxx")
	log.Println("  POST   /api/tasks/{id}/claim")
	log.Println("  PATCH  /api/tasks/{id}/status?username=xxx")

	log.Fatal(http.ListenAndServe(port, handler))
}
