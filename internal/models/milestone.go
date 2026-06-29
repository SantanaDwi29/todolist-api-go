package models

import "time"

type Milestone struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Title      string    `json:"title" gorm:"not null"`
	TargetDate time.Time `json:"target_date"`
	Progress   int       `json:"progress" gorm:"default:0"` // 0 to 100
	IsCompleted bool     `json:"is_completed" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
}

type MilestoneInput struct {
	Title      string    `json:"title" binding:"required"`
	TargetDate time.Time `json:"target_date" binding:"required"`
}
