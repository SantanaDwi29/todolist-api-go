package controllers

import (
	"net/http"
	"todolist-api/config"
	"todolist-api/models"

	"github.com/gin-gonic/gin"
)

func GetNextMilestone(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var milestone models.Milestone
	// Find the closest active milestone by target date
	result := config.DB.Where("user_id = ? AND is_completed = ?", userID, false).Order("target_date asc").First(&milestone)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"milestone": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"milestone": milestone})
}

func CreateMilestone(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input models.MilestoneInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	milestone := models.Milestone{
		UserID:     userID.(uint),
		Title:      input.Title,
		TargetDate: input.TargetDate,
	}

	config.DB.Create(&milestone)
	c.JSON(http.StatusCreated, milestone)
}
