package service

import (
	"fmt"
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
	repo        repository.TodoRepository
	projectRepo repository.ProjectRepository
}

func NewTodoService(repo repository.TodoRepository, projectRepo repository.ProjectRepository) TodoService {
	return &todoService{repo, projectRepo}
}

func (s *todoService) GetTodos(userID uint, categoryID, priority, status string) ([]models.Todo, error) {
	return s.repo.FindAll(userID, categoryID, priority, status)
}

func (s *todoService) syncProjectStatus(userID uint, projectID *uint) {
	if projectID == nil || *projectID == 0 || s.projectRepo == nil {
		return
	}
	projectIDStr := fmt.Sprintf("%d", *projectID)
	project, err := s.projectRepo.FindByID(projectIDStr, userID)
	if err != nil {
		return
	}

	newProjectStatus := models.ProjectStatusActive
	if len(project.Todos) > 0 {
		allDone := true
		for _, t := range project.Todos {
			if t.Status != models.StatusDone {
				allDone = false
				break
			}
		}
		if allDone {
			newProjectStatus = models.ProjectStatusCompleted
		}
	}

	if project.Status != newProjectStatus {
		project.Status = newProjectStatus
		_ = s.projectRepo.Update(&project)
	}
}

func (s *todoService) CreateTodo(userID uint, input models.TodoInput) (models.Todo, error) {
	if input.ProjectID != nil && *input.ProjectID != 0 {
		projectIDStr := fmt.Sprintf("%d", *input.ProjectID)
		_, err := s.projectRepo.FindByID(projectIDStr, userID)
		if err != nil {
			return models.Todo{}, fmt.Errorf("invalid project_id or project does not belong to user")
		}
	} else {
		input.ProjectID = nil
	}

	todo := models.Todo{
		UserID:      userID,
		CategoryID:  input.CategoryID,
		ProjectID:   input.ProjectID,
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Deadline:    input.Deadline,
	}
	err := s.repo.Create(&todo)
	if err == nil && todo.ProjectID != nil {
		s.syncProjectStatus(userID, todo.ProjectID)
	}
	return todo, err
}

func (s *todoService) UpdateTodo(id string, userID uint, input models.TodoInput) (models.Todo, error) {
	if input.ProjectID != nil && *input.ProjectID != 0 {
		projectIDStr := fmt.Sprintf("%d", *input.ProjectID)
		_, err := s.projectRepo.FindByID(projectIDStr, userID)
		if err != nil {
			return models.Todo{}, fmt.Errorf("invalid project_id or project does not belong to user")
		}
	} else {
		input.ProjectID = nil
	}

	todo, err := s.repo.FindByID(id, userID)
	if err != nil {
		return todo, err
	}

	oldProjectID := todo.ProjectID

	todo.CategoryID = input.CategoryID
	todo.ProjectID = input.ProjectID
	todo.Title = input.Title
	todo.Description = input.Description
	todo.Priority = input.Priority
	todo.Deadline = input.Deadline

	err = s.repo.Update(&todo)
	if err == nil {
		if oldProjectID != nil {
			s.syncProjectStatus(userID, oldProjectID)
		}
		if todo.ProjectID != nil {
			s.syncProjectStatus(userID, todo.ProjectID)
		}
	}
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
	if err != nil {
		return newStatus, err
	}

	s.syncProjectStatus(userID, todo.ProjectID)
	return newStatus, err
}

func (s *todoService) DeleteTodo(id string, userID uint) error {
	todo, err := s.repo.FindByID(id, userID)
	if err != nil {
		return err
	}
	err = s.repo.Delete(&todo)
	if err == nil {
		s.syncProjectStatus(userID, todo.ProjectID)
	}
	return err
}
