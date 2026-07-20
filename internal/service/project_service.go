package service

import (
	"todolist-api/internal/models"
	"todolist-api/internal/repository"
)

type ProjectService interface {
	GetAllProjects(userID uint) ([]models.Project, error)
	CreateProject(userID uint, input models.ProjectInput) (models.Project, error)
	GetProjectByID(id string, userID uint) (models.Project, error)
	UpdateProject(id string, userID uint, input models.ProjectInput) (models.Project, error)
	UpdateProjectStatus(id string, userID uint, status models.ProjectStatus) (models.Project, error)
	DeleteProject(id string, userID uint) error
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo}
}

func (s *projectService) GetAllProjects(userID uint) ([]models.Project, error) {
	projects, err := s.repo.FindAll(userID)
	if err != nil {
		return nil, err
	}

	for i := range projects {
		expectedStatus := models.ProjectStatusActive
		if len(projects[i].Todos) > 0 {
			allDone := true
			for _, t := range projects[i].Todos {
				if t.Status != models.StatusDone {
					allDone = false
					break
				}
			}
			if allDone {
				expectedStatus = models.ProjectStatusCompleted
			}
		}

		if projects[i].Status != expectedStatus {
			projects[i].Status = expectedStatus
			_ = s.repo.Update(&projects[i])
		}
	}

	return projects, nil
}

func (s *projectService) CreateProject(userID uint, input models.ProjectInput) (models.Project, error) {
	project := models.Project{
		UserID:      userID,
		Name:        input.Name,
		Description: input.Description,
		Status:      models.ProjectStatusActive,
	}

	err := s.repo.Create(&project)
	return project, err
}

func (s *projectService) GetProjectByID(id string, userID uint) (models.Project, error) {
	project, err := s.repo.FindByID(id, userID)
	if err != nil {
		return project, err
	}

	expectedStatus := models.ProjectStatusActive
	if len(project.Todos) > 0 {
		allDone := true
		for _, t := range project.Todos {
			if t.Status != models.StatusDone {
				allDone = false
				break
			}
		}
		if allDone {
			expectedStatus = models.ProjectStatusCompleted
		}
	}

	if project.Status != expectedStatus {
		project.Status = expectedStatus
		_ = s.repo.Update(&project)
	}

	return project, nil
}

func (s *projectService) UpdateProject(id string, userID uint, input models.ProjectInput) (models.Project, error) {
	project, err := s.repo.FindByID(id, userID)
	if err != nil {
		return project, err
	}

	project.Name = input.Name
	project.Description = input.Description

	err = s.repo.Update(&project)
	return project, err
}

func (s *projectService) UpdateProjectStatus(id string, userID uint, status models.ProjectStatus) (models.Project, error) {
	project, err := s.repo.FindByID(id, userID)
	if err != nil {
		return project, err
	}

	project.Status = status
	err = s.repo.Update(&project)
	return project, err
}

func (s *projectService) DeleteProject(id string, userID uint) error {
	project, err := s.repo.FindByID(id, userID)
	if err != nil {
		return err
	}
	return s.repo.Delete(&project)
}
