package database

import (
	"log"
	"os"
	"path/filepath"

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

	// Dynamically locate the migrations folder relative to the current working directory
	migrationsPath := "migrations"
	if wd, err := os.Getwd(); err == nil {
		dir := wd
		for {
			target := filepath.Join(dir, "migrations")
			if info, err := os.Stat(target); err == nil && info.IsDir() {
				if rel, err := filepath.Rel(wd, target); err == nil {
					migrationsPath = filepath.ToSlash(rel)
				}
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	// Assuming migrations folder is at the root of the project
	m, err := migrate.NewWithDatabaseInstance(
		"file://" + migrationsPath,
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
