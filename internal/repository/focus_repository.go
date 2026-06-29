package repository

import (
	"todolist-api/internal/models"

	"gorm.io/gorm"
)

type FocusRepository interface {
	FindCurrent(userID uint) (models.FocusSession, error)
	FindByStatus(userID uint, status models.SessionStatus) (models.FocusSession, error)
	CompleteAllActive(userID uint) error
	Create(session *models.FocusSession) error
	Update(session *models.FocusSession) error
}

type focusRepository struct {
	db *gorm.DB
}

func NewFocusRepository(db *gorm.DB) FocusRepository {
	return &focusRepository{db}
}

func (r *focusRepository) FindCurrent(userID uint) (models.FocusSession, error) {
	var session models.FocusSession
	err := r.db.Where("user_id = ? AND status != ?", userID, models.SessionCompleted).First(&session).Error
	return session, err
}

func (r *focusRepository) FindByStatus(userID uint, status models.SessionStatus) (models.FocusSession, error) {
	var session models.FocusSession
	err := r.db.Where("user_id = ? AND status = ?", userID, status).First(&session).Error
	return session, err
}

func (r *focusRepository) CompleteAllActive(userID uint) error {
	return r.db.Model(&models.FocusSession{}).
		Where("user_id = ? AND status != ?", userID, models.SessionCompleted).
		Update("status", models.SessionCompleted).Error
}

func (r *focusRepository) Create(session *models.FocusSession) error {
	return r.db.Create(session).Error
}

func (r *focusRepository) Update(session *models.FocusSession) error {
	return r.db.Save(session).Error
}
