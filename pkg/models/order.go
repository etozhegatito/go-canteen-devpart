package models

type Order struct {
	ID       uint    `gorm:"primaryKey"`
	UserID   uint    `gorm:"not null"`
	TotalSum float64 `gorm:"not null"`
	User     User    `gorm:"foreignKey:UserID"`
}
