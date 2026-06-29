package database

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDB runs SQL migrations using golang-migrate
func MigrateDB() {
	if DB == nil {
		log.Fatal("Database is not initialized. Call ConnectDB first.")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB from gorm:", err)
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatal("Failed to create migration driver:", err)
	}

	// Assuming migrations folder is at the root of the project
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql", 
		driver,
	)
	if err != nil {
		log.Fatal("Failed to initialize migration:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run database migrations:", err)
	}

	log.Println("Database SQL migration completed successfully!")
}
