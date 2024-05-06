package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Включает поля ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `gorm:"not null"`
	Surname    string `gorm:"not null"`
	Email      string `gorm:"uniqueIndex;not null"`
	Password   string `gorm:"not null"`
	Phone      string `gorm:"not null"`
	Age        int    `gorm:"not null;default:18"` // Возраст пользователя
	IsAdmin    bool   `gorm:"default:false"`
}
