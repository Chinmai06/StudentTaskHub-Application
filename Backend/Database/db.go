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

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	createTables()
	log.Println("Database initialized successfully")
}

func createTables() {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		ufid TEXT DEFAULT '',
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

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
		category TEXT DEFAULT '',
		description TEXT DEFAULT '',
		location TEXT DEFAULT '',
		priority TEXT NOT NULL DEFAULT 'Medium' CHECK(priority IN ('High', 'Medium', 'Low')),
		deadline TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'open' CHECK(status IN ('open', 'claimed', 'in_progress', 'done')),
		created_by TEXT NOT NULL,
		claimed_by TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (created_by) REFERENCES users(username),
		FOREIGN KEY (claimed_by) REFERENCES users(username)
	);`

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
