package controllers

import (
	"net/http"
	"time"
	"todolist-api/config"
	"todolist-api/models"

	"github.com/gin-gonic/gin"
)

func GetAnalytics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Calculate current week's Mon-Sun tasks completed
	now := time.Now()
	// Find Monday of the current week
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	monday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, offset)

	var todos []models.Todo
	config.DB.Where("user_id = ? AND status = ? AND updated_at >= ?", userID, models.StatusDone, monday).Find(&todos)

	chartData := []int{0, 0, 0, 0, 0, 0, 0}
	for _, todo := range todos {
		// Calculate the day of the week, Monday = 0, Sunday = 6
		dayIndex := int(todo.UpdatedAt.Weekday() - time.Monday)
		if dayIndex < 0 {
			dayIndex = 6 // Sunday
		}
		if dayIndex >= 0 && dayIndex < 7 {
			chartData[dayIndex]++
		}
	}

	// Calculate focus score based on number of tasks completed this week vs a target
	// Or simply sum the completed tasks
	totalCompleted := 0
	for _, count := range chartData {
		totalCompleted += count
	}

	focusScore := totalCompleted * 2 // Arbitrary multiplier for "score"

	c.JSON(http.StatusOK, gin.H{
		"chartData":  chartData,
		"focusScore": focusScore,
	})
}
