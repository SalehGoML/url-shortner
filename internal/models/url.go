package models

import "time"

type URL struct {
	ID           uint   `gorm:"primarykey"`
	LongURL      string `gorm:"type:text;not null"`
	shortCode    string `gorm:"size:20;uniqueIndex;not null"`
	UserID       uint   `gorm:"index"`
	Clicks       int    `gorm:"default:0"`
	PasswordHash *string
	Expiry       *time.Time `gorm:"index"`
	IsActive     bool
	CreatedAt    time.Time
	updatedAt    time.Time
}
