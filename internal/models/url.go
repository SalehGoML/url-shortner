package models

import "time"

type URL struct {
	ID           uint `gorm:"primarykey"`
	LongURL      string
	shortCode    string `gorm:"unigue"`
	UserID       uint
	Clicks       int
	PasswordHash string
	Expiry       *time.Time
	CreatedAt    time.Time
}
