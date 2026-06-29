package controllers

import (
	"net/http"
	"time"
	"todolist-api/config"
	"todolist-api/models"

	"github.com/gin-gonic/gin"
)

func GetCurrentFocusSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var session models.FocusSession
	result := config.DB.Where("user_id = ? AND status != ?", userID, models.SessionCompleted).First(&session)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"session": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"session": session})
}

func StartFocusSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Complete any existing active sessions
	config.DB.Model(&models.FocusSession{}).
		Where("user_id = ? AND status != ?", userID, models.SessionCompleted).
		Update("status", models.SessionCompleted)

	var input struct {
		DurationMinutes int `json:"duration_minutes"`
	}
	c.ShouldBindJSON(&input)

	duration := 45
	if input.DurationMinutes > 0 {
		duration = input.DurationMinutes
	}

	// Create new session
	session := models.FocusSession{
		UserID:          userID.(uint),
		StartTime:       time.Now(),
		Status:          models.SessionActive,
		DurationMinutes: duration,
	}

	config.DB.Create(&session)
	c.JSON(http.StatusOK, session)
}

func PauseFocusSession(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var session models.FocusSession
	
	if err := config.DB.Where("user_id = ? AND status = ?", userID, models.SessionActive).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found"})
		return
	}

	now := time.Now()
	// Calculate elapsed since start or last resume
	var lastActive time.Time
	if session.PausedAt != nil {
		// Wait, if it's active, PausedAt should be nil.
		// So elapsed since StartTime or last resume. 
		// For simplicity, we just use the difference from start time if it was never paused. 
		// Actually a better way: keep track of ElapsedSeconds.
		lastActive = session.StartTime
	} else {
		lastActive = session.StartTime
	}

	// This logic is simple: we update ElapsedSeconds and set status to paused
	elapsed := int(now.Sub(lastActive).Seconds())
	
	// Wait, if it was paused and resumed, start time doesn't reflect the whole time.
	// Let's just track elapsed correctly.
	session.ElapsedSeconds += elapsed
	session.Status = models.SessionPaused
	session.PausedAt = &now

	config.DB.Save(&session)
	c.JSON(http.StatusOK, session)
}

func ResumeFocusSession(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var session models.FocusSession
	
	if err := config.DB.Where("user_id = ? AND status = ?", userID, models.SessionPaused).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No paused session found"})
		return
	}

	// Reset start time to now so we can measure elapsed from now
	session.StartTime = time.Now()
	session.Status = models.SessionActive
	session.PausedAt = nil

	config.DB.Save(&session)
	c.JSON(http.StatusOK, session)
}

func StopFocusSession(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var session models.FocusSession
	if err := config.DB.Where("user_id = ? AND status != ?", userID, models.SessionCompleted).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found"})
		return
	}

	now := time.Now()
	session.Status = models.SessionCompleted
	session.EndTime = &now

	config.DB.Save(&session)
	c.JSON(http.StatusOK, session)
}
