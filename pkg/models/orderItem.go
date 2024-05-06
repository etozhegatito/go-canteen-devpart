package models

// Struct для нашей корзины, т.к в одном заказе может быть несколько одинаковых блюд и все это надо сохранять
type OrderItem struct {
	ID       uint  `gorm:"primaryKey"`
	OrderID  uint  `gorm:"not null"`
	DishID   uint  `gorm:"not null"`
	Quantity int   `gorm:"not null"`
	Order    Order `gorm:"foreignKey:OrderID"`
	Dish     Dish  `gorm:"foreignKey:DishID"`
}
