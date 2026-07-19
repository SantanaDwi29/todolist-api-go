package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadConfig loads environment variables from .env file, searching parent directories if necessary
func LoadConfig() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting working directory:", err)
		return
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err != nil {
				log.Println("Error loading .env file:", err)
			}
			return
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	log.Println("No .env file found, relying on environment variables")
}
