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

// Order представляет заказ в системе
type Order struct {
	gorm.Model
	UserID      uint          // Ссылка на пользователя
	TotalSum    float64       // Общая сумма заказа
	TotalWeight float64       // Общий вес заказа
	Details     []OrderDetail `gorm:"foreignKey:OrderID;references:ID"`
}

// OrderItem представляет элемент заказа в системе
type OrderDetail struct {
	gorm.Model
	OrderID uint  // ID заказа
	MenuID  uint  // ID меню
	Count   int   // Количество
	Menu    Menu  `gorm:"foreignKey:MenuID"`
	Order   Order `gorm:"foreignKey:OrderID;references:ID"`
}

// Menu представляет элемент меню в системе столовой университета
type Menu struct {
	gorm.Model
	Name        string
	Price       float64
	Description string
	InStock     bool
	Weight      int // Вес в граммах
}
