package controllers

import (
	"net/http"
	"todolist-api/config"
	"todolist-api/models"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var categories []models.Category
	if err := config.DB.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input models.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		UserID: userID,
		Name:   input.Name,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}
func DeleteCategory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	categoryID := c.Param("id")

	var category models.Category
	if err := config.DB.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
