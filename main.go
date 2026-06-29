package main

import (
	"log"
	"os"

	"todolist-api/config"
	"todolist-api/models"
	"todolist-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Initialize Database
	config.ConnectDB()

	// Auto Migrate Models
	config.DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Todo{})

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust this in production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
	}))

	// Setup Routes
	routes.SetupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}