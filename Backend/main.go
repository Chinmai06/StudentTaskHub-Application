package main

import (
	"studenttaskhub/Database"
	"studenttaskhub/routes"
)

func main() {
	Database.ConnectDB()
	r := routes.SetupRoutes()
	r.Run(":8080")
}
