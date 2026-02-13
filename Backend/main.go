package main

import (
	"student-task-hub/database"
	"student-task-hub/routes"
)

func main() {
	database.ConnectDB()
	r := routes.SetupRoutes()
	r.Run(":8080")
}
