package service

import (
	"todolist-api/internal/models"
	"todolist-api/internal/repository"
)

type TodoService interface {
	GetTodos(userID uint, categoryID, priority, status string) ([]models.Todo, error)
	CreateTodo(userID uint, input models.TodoInput) (models.Todo, error)
	UpdateTodo(id string, userID uint, input models.TodoInput) (models.Todo, error)
	ToggleTodoStatus(id string, userID uint) (models.Status, error)
	DeleteTodo(id string, userID uint) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo}
}

func (s *todoService) GetTodos(userID uint, categoryID, priority, status string) ([]models.Todo, error) {
	return s.repo.FindAll(userID, categoryID, priority, status)
}

func (s *todoService) CreateTodo(userID uint, input models.TodoInput) (models.Todo, error) {
	todo := models.Todo{
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Deadline:    input.Deadline,
	}
	err := s.repo.Create(&todo)
	return todo, err
}

func (s *todoService) UpdateTodo(id string, userID uint, input models.TodoInput) (models.Todo, error) {
	todo, err := s.repo.FindByID(id, userID)
	if err != nil {
		return todo, err
	}

	todo.CategoryID = input.CategoryID
	todo.Title = input.Title
	todo.Description = input.Description
	todo.Priority = input.Priority
	todo.Deadline = input.Deadline

	err = s.repo.Update(&todo)
	return todo, err
}

func (s *todoService) ToggleTodoStatus(id string, userID uint) (models.Status, error) {
	todo, err := s.repo.FindByID(id, userID)
	if err != nil {
		return "", err
	}

	newStatus := models.StatusDone
	if todo.Status == models.StatusDone {
		newStatus = models.StatusUndone
	}
	todo.Status = newStatus

	err = s.repo.Update(&todo)
	return newStatus, err
}

func (s *todoService) DeleteTodo(id string, userID uint) error {
	todo, err := s.repo.FindByID(id, userID)
	if err != nil {
		return err
	}
	return s.repo.Delete(&todo)
}
