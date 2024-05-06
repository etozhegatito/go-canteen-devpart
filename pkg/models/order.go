package models

// Struct для заказов, связан с User struct т.к у каждого заказа должеть быть свой юзер
type Order struct {
	ID       uint    `gorm:"primaryKey"`
	UserID   uint    `gorm:"not null"`
	TotalSum float64 `gorm:"not null"`
	User     User    `gorm:"foreignKey:UserID"`
}
