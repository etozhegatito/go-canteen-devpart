package models

type OrderItem struct {
	ID       uint  `gorm:"primaryKey"`
	OrderID  uint  `gorm:"not null"`
	DishID   uint  `gorm:"not null"`
	Quantity int   `gorm:"not null"`
	Order    Order `gorm:"foreignKey:OrderID"`
	Dish     Dish  `gorm:"foreignKey:DishID"`
}
