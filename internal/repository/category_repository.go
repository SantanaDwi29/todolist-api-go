package repository

import (
	"todolist-api/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(userID uint) ([]models.Category, error)
	Create(category *models.Category) error
	FindByID(id string, userID uint) (models.Category, error)
	Delete(category *models.Category) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) FindAll(userID uint) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("user_id = ?", userID).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByID(id string, userID uint) (models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error
	return category, err
}

func (r *categoryRepository) Delete(category *models.Category) error {
	return r.db.Delete(category).Error
}
