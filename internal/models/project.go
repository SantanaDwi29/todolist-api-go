package models

import "time"

type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusCompleted ProjectStatus = "completed"
)

type Project struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	UserID      uint          `json:"user_id" gorm:"not null"`
	Name        string        `json:"name" gorm:"not null"`
	Description string        `json:"description"`
	Status      ProjectStatus `json:"status" gorm:"type:enum('active','completed');default:'active'"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Todos       []Todo        `json:"todos,omitempty" gorm:"foreignKey:ProjectID"`
}

type ProjectInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
