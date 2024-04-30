package db

import (
	"gorm.io/gorm"
)

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

type Dish struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Description string
	InStock     bool
	Weight      int
}

type Order struct {
	ID       uint    `gorm:"primaryKey"`
	UserID   uint    `gorm:"not null"`
	TotalSum float64 `gorm:"not null"`
	User     User    `gorm:"foreignKey:UserID"`
}

type OrderItem struct {
	ID       uint  `gorm:"primaryKey"`
	OrderID  uint  `gorm:"not null"`
	DishID   uint  `gorm:"not null"`
	Quantity int   `gorm:"not null"`
	Order    Order `gorm:"foreignKey:OrderID"`
	Dish     Dish  `gorm:"foreignKey:DishID"`
}

type OrderRequest struct {
	UserID    uint       `json:"user_id"`
	CartItems []CartItem `json:"cart_items"`
}

type CartItem struct {
	DishID   uint `json:"dish_id"`
	Quantity int  `json:"quantity"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
