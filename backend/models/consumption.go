package models

type Ingredient struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"unique;not null"`
	KcalPer100g float64 `gorm:"not null"`
}
