package main

import (
	"log"

	"todolist-api/internal/config"
	"todolist-api/internal/database"
)

func main() {
	// Load config
	config.LoadConfig()

	// Connect Database
	database.ConnectDB()

	// Run Migrations
	log.Println("Starting standalone database migration...")
	database.MigrateDB()
}
