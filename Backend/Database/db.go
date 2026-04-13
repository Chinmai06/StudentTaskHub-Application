package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB creates the SQLite database and tables
func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	createTables()
	log.Println("Database initialized successfully")
}

func createTables() {
	// Sprint 2: Added users table for registration and login
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Sprint 3: Added profiles table for user profiles
	profilesTable := `
	CREATE TABLE IF NOT EXISTS profiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		full_name TEXT DEFAULT '',
		bio TEXT DEFAULT '',
		major TEXT DEFAULT '',
		year TEXT DEFAULT '',
		skills TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (username) REFERENCES users(username)
	);`

	tasksTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		deadline TEXT NOT NULL,
		priority TEXT NOT NULL DEFAULT 'normal' CHECK(priority IN ('high', 'normal')),
		status TEXT NOT NULL DEFAULT 'open' CHECK(status IN ('open', 'claimed', 'in_progress', 'done')),
		created_by TEXT NOT NULL,
		claimed_by TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (created_by) REFERENCES users(username),
		FOREIGN KEY (claimed_by) REFERENCES users(username)
	);`

	// Sprint 3: Added feedback table for task reviews
	feedbackTable := `
	CREATE TABLE IF NOT EXISTS feedback (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		rating INTEGER NOT NULL CHECK(rating >= 1 AND rating <= 5),
		comment TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (task_id) REFERENCES tasks(id),
		FOREIGN KEY (username) REFERENCES users(username)
	);`

	tables := []struct {
		name string
		sql  string
	}{
		{"users", usersTable},
		{"profiles", profilesTable},
		{"tasks", tasksTable},
		{"feedback", feedbackTable},
	}

	for _, table := range tables {
		_, err := DB.Exec(table.sql)
		if err != nil {
			log.Fatalf("Failed to create %s table: %v", table.name, err)
		}
	}
}
