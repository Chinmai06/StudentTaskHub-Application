package Database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable()
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		owner_name TEXT NOT NULL,
		description TEXT,
		deadline TEXT,
		priority TEXT,
		category TEXT
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
