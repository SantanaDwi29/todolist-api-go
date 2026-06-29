package handler

import (
	"net/http"

	"todolist-api/internal/models"
	"todolist-api/internal/service"
	"todolist-api/utils"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	service service.TodoService
}

func NewTodoHandler(service service.TodoService) *TodoHandler {
	return &TodoHandler{service}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	categoryID := c.Query("category_id")
	priority := c.Query("priority")
	status := c.Query("status")

	todos, err := h.service.GetTodos(userID, categoryID, priority, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input models.TodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "message": utils.FormatValidationError(err)})
		return
	}

	todo, err := h.service.CreateTodo(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	todoID := c.Param("id")

	var input models.TodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "message": utils.FormatValidationError(err)})
		return
	}

	todo, err := h.service.UpdateTodo(todoID, userID, input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) ToggleTodoStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	todoID := c.Param("id")

	status, err := h.service.ToggleTodoStatus(todoID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated", "status": status})
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	todoID := c.Param("id")

	err := h.service.DeleteTodo(todoID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
