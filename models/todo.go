package models

import "time"

type Priority string
type Status string

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityEasy   Priority = "easy"

	StatusDone   Status = "done"
	StatusUndone Status = "undone"
)

type Todo struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	CategoryID *uint     `json:"category_id"`
	Title      string    `json:"title" gorm:"not null"`
	Description string   `json:"description"`
	Priority   Priority  `json:"priority" gorm:"type:enum('high','medium','easy');default:'easy'"`
	Deadline   *time.Time`json:"deadline"`
	Status     Status    `json:"status" gorm:"type:enum('done','undone');default:'undone'"`
	CreatedAt  time.Time `json:"created_at"`
	Category   *Category `json:"category,omitempty"`
}

type TodoInput struct {
	CategoryID  *uint     `json:"category_id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority" binding:"required,oneof=high medium easy"`
	Deadline    *time.Time`json:"deadline"`
}
