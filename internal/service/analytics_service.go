package service

import (
	"time"

	"todolist-api/internal/repository"
)

type AnalyticsService interface {
	GetAnalytics(userID uint) (map[string]interface{}, error)
}

type analyticsService struct {
	todoRepo repository.TodoRepository
}

func NewAnalyticsService(todoRepo repository.TodoRepository) AnalyticsService {
	return &analyticsService{todoRepo}
}

func (s *analyticsService) GetAnalytics(userID uint) (map[string]interface{}, error) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	monday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, offset)

	todos, err := s.todoRepo.FindCompletedSince(userID, monday)
	if err != nil {
		return nil, err
	}

	chartData := []int{0, 0, 0, 0, 0, 0, 0}
	for _, todo := range todos {
		dayIndex := int(todo.UpdatedAt.Weekday() - time.Monday)
		if dayIndex < 0 {
			dayIndex = 6
		}
		if dayIndex >= 0 && dayIndex < 7 {
			chartData[dayIndex]++
		}
	}

	totalCompleted := 0
	for _, count := range chartData {
		totalCompleted += count
	}

	focusScore := totalCompleted * 2

	return map[string]interface{}{
		"chartData":  chartData,
		"focusScore": focusScore,
	}, nil
}
