package service

import (
	"todolist-api/internal/models"
	"todolist-api/internal/repository"
)

type CategoryService interface {
	GetCategories(userID uint) ([]models.Category, error)
	CreateCategory(userID uint, input models.CategoryInput) (models.Category, error)
	DeleteCategory(id string, userID uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo}
}

func (s *categoryService) GetCategories(userID uint) ([]models.Category, error) {
	return s.repo.FindAll(userID)
}

func (s *categoryService) CreateCategory(userID uint, input models.CategoryInput) (models.Category, error) {
	category := models.Category{
		UserID: userID,
		Name:   input.Name,
	}
	err := s.repo.Create(&category)
	return category, err
}

func (s *categoryService) DeleteCategory(id string, userID uint) error {
	category, err := s.repo.FindByID(id, userID)
	if err != nil {
		return err
	}
	return s.repo.Delete(&category)
}
