package models

type Dish struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Description string
	InStock     bool
	Weight      int
}
