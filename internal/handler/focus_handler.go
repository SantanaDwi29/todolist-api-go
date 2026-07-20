package handler

import (
	"net/http"

	"todolist-api/internal/service"

	"github.com/gin-gonic/gin"
)

type FocusHandler struct {
	service service.FocusService
}

func NewFocusHandler(service service.FocusService) *FocusHandler {
	return &FocusHandler{service}
}

func (h *FocusHandler) GetCurrentFocusSession(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	session, _ := h.service.GetCurrentFocusSession(userID)

	c.JSON(http.StatusOK, gin.H{"session": session})
}

func (h *FocusHandler) StartFocusSession(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input struct {
		DurationMinutes int  `json:"duration_minutes"`
		DurationSeconds *int `json:"duration_seconds"`
	}
	c.ShouldBindJSON(&input)

	session, err := h.service.StartFocusSession(userID, input.DurationMinutes, input.DurationSeconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *FocusHandler) PauseFocusSession(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	session, err := h.service.PauseFocusSession(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *FocusHandler) ResumeFocusSession(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	session, err := h.service.ResumeFocusSession(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No paused session found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *FocusHandler) StopFocusSession(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	session, err := h.service.StopFocusSession(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found"})
		return
	}

	c.JSON(http.StatusOK, session)
}
