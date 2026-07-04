package main

import (
	"fmt"
	"log"

	"todolist-api/internal/config"
	"todolist-api/internal/database"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()

	// 1. Force migration version to 6 (clean)
	err := database.DB.Exec("UPDATE schema_migrations SET version = 6, dirty = 0").Error
	if err != nil {
		log.Fatal("Failed to update schema_migrations:", err)
	}

	// 2. Clean up any tables that might have been partially created in migration 7
	database.DB.Exec("ALTER TABLE todos DROP FOREIGN KEY fk_todos_project")
	database.DB.Exec("ALTER TABLE todos DROP COLUMN project_id")
	database.DB.Exec("DROP TABLE IF EXISTS projects")

	fmt.Println("Successfully fixed dirty database and reverted migration 7!")
}
