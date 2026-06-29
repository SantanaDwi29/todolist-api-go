package repository

import (
	"time"
	"todolist-api/internal/models"

	"gorm.io/gorm"
)

type TodoRepository interface {
	FindAll(userID uint, categoryID, priority, status string) ([]models.Todo, error)
	FindByID(id string, userID uint) (models.Todo, error)
	FindCompletedSince(userID uint, since time.Time) ([]models.Todo, error)
	Create(todo *models.Todo) error
	Update(todo *models.Todo) error
	Delete(todo *models.Todo) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db}
}

func (r *todoRepository) FindAll(userID uint, categoryID, priority, status string) ([]models.Todo, error) {
	query := r.db.Where("user_id = ?", userID).Preload("Category")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var todos []models.Todo
	if err := query.Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) FindByID(id string, userID uint) (models.Todo, error) {
	var todo models.Todo
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).Preload("Category").First(&todo).Error; err != nil {
		return todo, err
	}
	return todo, nil
}

func (r *todoRepository) FindCompletedSince(userID uint, since time.Time) ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.Where("user_id = ? AND status = ? AND updated_at >= ?", userID, models.StatusDone, since).Find(&todos).Error
	return todos, err
}

func (r *todoRepository) Create(todo *models.Todo) error {
	if err := r.db.Create(todo).Error; err != nil {
		return err
	}
	return r.db.Preload("Category").First(todo, todo.ID).Error
}

func (r *todoRepository) Update(todo *models.Todo) error {
	if err := r.db.Save(todo).Error; err != nil {
		return err
	}
	return r.db.Preload("Category").First(todo, todo.ID).Error
}

func (r *todoRepository) Delete(todo *models.Todo) error {
	return r.db.Delete(todo).Error
}
