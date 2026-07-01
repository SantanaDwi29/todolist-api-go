package models

import "time"

type OAuthClient struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	ClientID     string    `json:"client_id" gorm:"unique;not null;size:100"`
	ClientSecret string    `json:"client_secret" gorm:"not null;size:100"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName overrides the default table name
func (OAuthClient) TableName() string {
	return "oauth_clients"
}

