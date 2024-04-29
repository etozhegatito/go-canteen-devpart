package db

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model            // Включает поля ID, CreatedAt, UpdatedAt, DeletedAt
	Name        string    `gorm:"not null"`
	Surname     string    `gorm:"not null"`
	Email       string    `gorm:"uniqueIndex;not null"`
	Password    string    `gorm:"not null"`
	Phone       string    `gorm:"not null"`
	DateOfBirth time.Time `gorm:"type:date"`
	IsAdmin     bool      `gorm:"default:false"`
}
