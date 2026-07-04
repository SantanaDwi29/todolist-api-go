package repository

import (
	"todolist-api/internal/models"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	FindAll(userID uint) ([]models.Project, error)
	Create(project *models.Project) error
	FindByID(id string, userID uint) (models.Project, error)
	Update(project *models.Project) error
	Delete(project *models.Project) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db}
}

func (r *projectRepository) FindAll(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Where("user_id = ?", userID).
		Preload("Todos", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority ASC, deadline ASC")
		}).
		Find(&projects).Error
	return projects, err
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) FindByID(id string, userID uint) (models.Project, error) {
	var project models.Project
	// Preload Todos, ordered by priority (enum index: high=1, medium=2, easy=3) and deadline ascending
	err := r.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Todos", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority ASC, deadline ASC")
		}).
		First(&project).Error
	return project, err
}

func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) Delete(project *models.Project) error {
	return r.db.Delete(project).Error
}
