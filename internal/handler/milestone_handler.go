package handler

import (
	"net/http"

	"todolist-api/internal/models"
	"todolist-api/internal/service"
	"todolist-api/utils"

	"github.com/gin-gonic/gin"
)

type MilestoneHandler struct {
	service service.MilestoneService
}

func NewMilestoneHandler(service service.MilestoneService) *MilestoneHandler {
	return &MilestoneHandler{service}
}

func (h *MilestoneHandler) GetNextMilestone(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	milestone, _ := h.service.GetNextMilestone(userID)
	// If err != nil (like record not found), it will just return nil milestone, which matches old behavior

	c.JSON(http.StatusOK, gin.H{"milestone": milestone})
}

func (h *MilestoneHandler) CreateMilestone(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input models.MilestoneInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "message": utils.FormatValidationError(err)})
		return
	}

	milestone, err := h.service.CreateMilestone(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create milestone"})
		return
	}

	c.JSON(http.StatusCreated, milestone)
}
