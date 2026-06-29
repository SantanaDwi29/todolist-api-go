package handler

import (
	"net/http"

	"todolist-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service service.AnalyticsService
}

func NewAnalyticsHandler(service service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service}
}

func (h *AnalyticsHandler) GetAnalytics(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	analyticsData, err := h.service.GetAnalytics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate analytics"})
		return
	}

	c.JSON(http.StatusOK, analyticsData)
}
