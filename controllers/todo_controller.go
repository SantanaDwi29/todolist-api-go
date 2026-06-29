package controllers

import (
	"net/http"
	"todolist-api/config"
	"todolist-api/models"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	categoryID := c.Query("category_id")
	priority := c.Query("priority")
	status := c.Query("status")

	query := config.DB.Where("user_id = ?", userID).Preload("Category")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var todos []models.Todo
	if err := query.Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func CreateTodo(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input models.TodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo := models.Todo{
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Deadline:    input.Deadline,
	}

	if err := config.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	// Fetch with category preloaded to return complete data
	config.DB.Preload("Category").First(&todo, todo.ID)

	c.JSON(http.StatusCreated, todo)
}

func UpdateTodo(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	todoID := c.Param("id")

	var todo models.Todo
	if err := config.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	var input models.TodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&todo).Updates(models.Todo{
		CategoryID:  input.CategoryID,
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Deadline:    input.Deadline,
	})

	config.DB.Preload("Category").First(&todo, todo.ID)
	c.JSON(http.StatusOK, todo)
}

func ToggleTodoStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	todoID := c.Param("id")

	var todo models.Todo
	if err := config.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	newStatus := models.StatusDone
	if todo.Status == models.StatusDone {
		newStatus = models.StatusUndone
	}

	config.DB.Model(&todo).Update("status", newStatus)

	c.JSON(http.StatusOK, gin.H{"message": "Status updated", "status": newStatus})
}

func DeleteTodo(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	todoID := c.Param("id")

	var todo models.Todo
	if err := config.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	if err := config.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
