package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"todolist-api/internal/config"
	"todolist-api/internal/database"
	"todolist-api/internal/models"
)

// generateRandomString generates a random hex string of given byte length
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func main() {
	// Parse CLI flags
	clientName := flag.String("name", "Default Client", "The name of the OAuth client")
	flag.Parse()

	// 1. Load configuration and connect to database
	config.LoadConfig()
	database.ConnectDB()

	// Note: Auto migration for the oauth_clients table should happen via the normal migration process
	// when the server is run. If it hasn't been run, this might fail, so ensure you ran `go run cmd/server/main.go` 
	// at least once to trigger the up.sql files.

	// 2. Generate Client ID and Client Secret
	// Usually 40 hex characters (20 bytes) is standard for Client ID and Secret
	clientID, err := generateRandomString(20)
	if err != nil {
		log.Fatalf("Failed to generate client_id: %v", err)
	}

	clientSecret, err := generateRandomString(20)
	if err != nil {
		log.Fatalf("Failed to generate client_secret: %v", err)
	}

	// 3. Save to database
	client := models.OAuthClient{
		Name:         *clientName,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	db := database.DB
	if db == nil {
		log.Fatal("Database connection is nil")
	}

	result := db.Create(&client)
	if result.Error != nil {
		log.Fatalf("Failed to create OAuth client in database: %v", result.Error)
	}

	// 4. Output the credentials for the user
	fmt.Println("=====================================================")
	fmt.Println("OAuth Client Created Successfully!")
	fmt.Println("=====================================================")
	fmt.Printf("Client ID     : %s\n", client.ClientID)
	fmt.Printf("Client Secret : %s\n", client.ClientSecret)
	fmt.Println("=====================================================")
	fmt.Println("Store this client secret safely. It will be used by")
	fmt.Println("your frontend application to authenticate with the API.")
}
