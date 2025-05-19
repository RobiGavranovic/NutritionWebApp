package models

import "time"

type Ingredient struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"unique;not null"`
	KcalPer100g float64 `gorm:"not null"`
}

type ConsumptionRequest struct {
	TokenResponse TokenResponse `json:"tokenResponse"`
	Ingredient    string        `json:"ingredient"`
	Weight        float64       `json:"weight"`
}

type IngredientsNamesRequest struct {
	TokenResponse TokenResponse `json:"tokenResponse"`
}

type Consumption struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null"`
	IngredientID uint      `gorm:"not null"`
	Ingredient   string    `gorm:"not null"`
	Weight       float64   `gorm:"not null"`
	Calories     float64   `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
