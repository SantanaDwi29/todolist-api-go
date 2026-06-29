package models

import "time"

type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	Todos     []Todo    `json:"todos,omitempty"`
}

type CategoryInput struct {
	Name string `json:"name" binding:"required"`
}
