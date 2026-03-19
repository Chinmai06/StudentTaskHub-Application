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

	_, err := DB.Exec(usersTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	_, err = DB.Exec(tasksTable)
	if err != nil {
		log.Fatal("Failed to create tasks table:", err)
	}
}
