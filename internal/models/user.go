package models

import "time"

type User struct {
	ID        uint `gorm:"primary_key"`
	Name      string
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
}
