package service

import (
	"todolist-api/internal/models"
	"todolist-api/internal/repository"
)

type MilestoneService interface {
	GetNextMilestone(userID uint) (*models.Milestone, error)
	CreateMilestone(userID uint, input models.MilestoneInput) (models.Milestone, error)
}

type milestoneService struct {
	repo     repository.MilestoneRepository
	todoRepo repository.TodoRepository
}

func NewMilestoneService(repo repository.MilestoneRepository, todoRepo repository.TodoRepository) MilestoneService {
	return &milestoneService{repo, todoRepo}
}

func (s *milestoneService) GetNextMilestone(userID uint) (*models.Milestone, error) {
	milestone, err := s.repo.FindNextActive(userID)
	if err != nil {
		return nil, err
	}

	// Calculate milestone progress based on todos
	todos, err := s.todoRepo.FindAll(userID, "", "", "")
	if err == nil {
		totalCount := 0
		completedCount := 0
		for _, todo := range todos {
			if todo.Deadline != nil && !todo.Deadline.After(milestone.TargetDate) {
				totalCount++
				if todo.Status == models.StatusDone {
					completedCount++
				}
			}
		}

		if totalCount > 0 {
			milestone.Progress = (completedCount * 100) / totalCount
		} else {
			// Fallback: use overall completion percentage if no tasks have specific deadlines before the target date
			allTotal := len(todos)
			if allTotal > 0 {
				allCompleted := 0
				for _, todo := range todos {
					if todo.Status == models.StatusDone {
						allCompleted++
					}
				}
				milestone.Progress = (allCompleted * 100) / allTotal
			} else {
				milestone.Progress = 0
			}
		}

		// If progress reaches 100% and there are indeed tasks, mark as completed
		// and recurse to get the next active milestone.
		if milestone.Progress == 100 && (totalCount > 0 || len(todos) > 0) {
			milestone.IsCompleted = true
			_ = s.repo.Update(&milestone)

			// Recurse to find the next active milestone
			return s.GetNextMilestone(userID)
		}
	}

	return &milestone, nil
}

func (s *milestoneService) CreateMilestone(userID uint, input models.MilestoneInput) (models.Milestone, error) {
	milestone := models.Milestone{
		UserID:     userID,
		Title:      input.Title,
		TargetDate: input.TargetDate,
		Progress:   0,
		IsCompleted: false,
	}
	err := s.repo.Create(&milestone)
	return milestone, err
}
