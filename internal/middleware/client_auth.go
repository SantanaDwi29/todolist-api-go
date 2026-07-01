package middleware

import (
	"net/http"
	"todolist-api/internal/database"
	"todolist-api/internal/models"

	"github.com/gin-gonic/gin"
)

// ClientAuthMiddleware verifies X-Client-Id and X-Client-Secret headers against the database
func ClientAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.GetHeader("X-Client-Id")
		clientSecret := c.GetHeader("X-Client-Secret")

		if clientID == "" || clientSecret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Akses ditolak: Header X-Client-Id dan X-Client-Secret wajib disertakan",
				"data":    nil,
			})
			c.Abort()
			return
		}

		var client models.OAuthClient
		// Check the credentials in the database
		if err := database.DB.Where("client_id = ? AND client_secret = ?", clientID, clientSecret).First(&client).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Akses ditolak: Client ID atau Client Secret tidak valid",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// You can set the client name into context if later handlers need to know which client made the request
		c.Set("oauth_client_name", client.Name)

		c.Next()
	}
}
