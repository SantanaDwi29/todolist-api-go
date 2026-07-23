package repository

import (
	"todolist-api/internal/models"

	"gorm.io/gorm"
)

type MilestoneRepository interface {
	FindNextActive(userID uint) (models.Milestone, error)
	Create(milestone *models.Milestone) error
	Update(milestone *models.Milestone) error
}

type milestoneRepository struct {
	db *gorm.DB
}

func NewMilestoneRepository(db *gorm.DB) MilestoneRepository {
	return &milestoneRepository{db}
}

func (r *milestoneRepository) FindNextActive(userID uint) (models.Milestone, error) {
	var milestone models.Milestone
	err := r.db.Where("user_id = ? AND is_completed = ?", userID, false).Order("target_date asc").First(&milestone).Error
	return milestone, err
}

func (r *milestoneRepository) Create(milestone *models.Milestone) error {
	return r.db.Create(milestone).Error
}

func (r *milestoneRepository) Update(milestone *models.Milestone) error {
	return r.db.Save(milestone).Error
}
