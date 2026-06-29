package models

import "time"

type SessionStatus string

const (
	SessionActive    SessionStatus = "active"
	SessionPaused    SessionStatus = "paused"
	SessionCompleted SessionStatus = "completed"
)

type FocusSession struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	UserID          uint          `json:"user_id" gorm:"not null"`
	StartTime       time.Time     `json:"start_time"`
	EndTime         *time.Time    `json:"end_time"`
	Status          SessionStatus `json:"status" gorm:"type:enum('active','paused','completed');default:'active'"`
	DurationMinutes int           `json:"duration_minutes" gorm:"default:45"` // Target duration
	PausedAt        *time.Time    `json:"paused_at"`
	ElapsedSeconds  int           `json:"elapsed_seconds" gorm:"default:0"`   // How many seconds elapsed before pause
}
