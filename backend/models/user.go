package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GoogleID     string         `gorm:"unique;not null"`
	Email        string         `gorm:"unique;not null"`
	Username     string         `gorm:"unique;not null"`
	Allergens    pq.StringArray `gorm:"type:text[]"`
	Intolerances pq.StringArray `gorm:"type:text[]"`
}

type RegisterUser struct {
	Allergens     []string      `json:"allergens"`
	Intolerances  []string      `json:"intolerances"`
	TokenResponse TokenResponse `json:"tokenResponse"` // match exactly: tokenResponse
	Username      string        `json:"username"`
}

type GoogleUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}
