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
	database.InitDB("studenttaskhub.db")

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	// User routes (Sprint 2)
	api.HandleFunc("/register", handlers.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/login", handlers.Login).Methods("POST", "OPTIONS")

	// Profile routes (Sprint 3)
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

	// Feedback routes (Sprint 3)
	api.HandleFunc("/tasks/{id}/feedback", handlers.AddFeedback).Methods("POST", "OPTIONS")
	api.HandleFunc("/tasks/{id}/feedback", handlers.GetFeedback).Methods("GET")

	// Notification routes (Sprint 4 - NEW)
	api.HandleFunc("/notifications", handlers.GetNotifications).Methods("GET")
	api.HandleFunc("/notifications/read-all", handlers.MarkAllNotificationsRead).Methods("PATCH", "OPTIONS")
	api.HandleFunc("/notifications/unread-count", handlers.GetUnreadCount).Methods("GET")
	api.HandleFunc("/notifications/{id}/read", handlers.MarkNotificationRead).Methods("PATCH", "OPTIONS")

	handler := middleware.CORS(r)

	port := ":8080"
	log.Printf("StudentTaskHub API server starting on http://localhost%s", port)
	log.Println("")
	log.Println("=== Sprint 4 Endpoints (NEW) ===")
	log.Println("  GET    /api/notifications?username=xxx")
	log.Println("  GET    /api/notifications/unread-count?username=xxx")
	log.Println("  PATCH  /api/notifications/{id}/read?username=xxx")
	log.Println("  PATCH  /api/notifications/read-all?username=xxx")
	log.Println("")
	log.Println("=== All Other Endpoints ===")
	log.Println("  POST   /api/register")
	log.Println("  POST   /api/login")
	log.Println("  GET    /api/profile/{username}")
	log.Println("  PUT    /api/profile/{username}?username=xxx")
	log.Println("  POST   /api/tasks")
	log.Println("  GET    /api/tasks")
	log.Println("  GET    /api/tasks/{id}")
	log.Println("  PUT    /api/tasks/{id}?username=xxx")
	log.Println("  DELETE /api/tasks/{id}?username=xxx")
	log.Println("  POST   /api/tasks/{id}/claim")
	log.Println("  PATCH  /api/tasks/{id}/status?username=xxx")
	log.Println("  POST   /api/tasks/{id}/feedback?username=xxx")
	log.Println("  GET    /api/tasks/{id}/feedback")

	log.Fatal(http.ListenAndServe(port, handler))
}
