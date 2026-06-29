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
	repo repository.MilestoneRepository
}

func NewMilestoneService(repo repository.MilestoneRepository) MilestoneService {
	return &milestoneService{repo}
}

func (s *milestoneService) GetNextMilestone(userID uint) (*models.Milestone, error) {
	milestone, err := s.repo.FindNextActive(userID)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

func (s *milestoneService) CreateMilestone(userID uint, input models.MilestoneInput) (models.Milestone, error) {
	milestone := models.Milestone{
		UserID:     userID,
		Title:      input.Title,
		TargetDate: input.TargetDate,
	}
	err := s.repo.Create(&milestone)
	return milestone, err
}
