package handler

import (
	"net/http"

	"todolist-api/internal/models"
	"todolist-api/internal/service"
	"todolist-api/utils"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	categories, err := h.service.GetCategories(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input models.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "message": utils.FormatValidationError(err)})
		return
	}

	category, err := h.service.CreateCategory(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	categoryID := c.Param("id")

	err := h.service.DeleteCategory(categoryID, userID)
	if err != nil {
		// Differentiate between not found and internal error if necessary,
		// but typically keeping it simple for now based on original code.
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
