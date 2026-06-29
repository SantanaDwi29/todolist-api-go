package service

import (
	"time"

	"todolist-api/internal/models"
	"todolist-api/internal/repository"
)

type FocusService interface {
	GetCurrentFocusSession(userID uint) (*models.FocusSession, error)
	StartFocusSession(userID uint, durationMinutes int) (models.FocusSession, error)
	PauseFocusSession(userID uint) (models.FocusSession, error)
	ResumeFocusSession(userID uint) (models.FocusSession, error)
	StopFocusSession(userID uint) (models.FocusSession, error)
}

type focusService struct {
	repo repository.FocusRepository
}

func NewFocusService(repo repository.FocusRepository) FocusService {
	return &focusService{repo}
}

func (s *focusService) GetCurrentFocusSession(userID uint) (*models.FocusSession, error) {
	session, err := s.repo.FindCurrent(userID)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *focusService) StartFocusSession(userID uint, durationMinutes int) (models.FocusSession, error) {
	s.repo.CompleteAllActive(userID)

	duration := 45
	if durationMinutes > 0 {
		duration = durationMinutes
	}

	session := models.FocusSession{
		UserID:          userID,
		StartTime:       time.Now(),
		Status:          models.SessionActive,
		DurationMinutes: duration,
	}

	err := s.repo.Create(&session)
	return session, err
}

func (s *focusService) PauseFocusSession(userID uint) (models.FocusSession, error) {
	session, err := s.repo.FindByStatus(userID, models.SessionActive)
	if err != nil {
		return session, err
	}

	now := time.Now()
	lastActive := session.StartTime
	
	elapsed := int(now.Sub(lastActive).Seconds())
	
	session.ElapsedSeconds += elapsed
	session.Status = models.SessionPaused
	session.PausedAt = &now

	err = s.repo.Update(&session)
	return session, err
}

func (s *focusService) ResumeFocusSession(userID uint) (models.FocusSession, error) {
	session, err := s.repo.FindByStatus(userID, models.SessionPaused)
	if err != nil {
		return session, err
	}

	session.StartTime = time.Now()
	session.Status = models.SessionActive
	session.PausedAt = nil

	err = s.repo.Update(&session)
	return session, err
}

func (s *focusService) StopFocusSession(userID uint) (models.FocusSession, error) {
	session, err := s.repo.FindCurrent(userID)
	if err != nil {
		return session, err
	}

	now := time.Now()
	session.Status = models.SessionCompleted
	session.EndTime = &now

	err = s.repo.Update(&session)
	return session, err
}
