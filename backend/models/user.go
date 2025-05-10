package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GoogleID         string         `gorm:"unique;not null"`
	Email            string         `gorm:"unique;not null"`
	Username         string         `gorm:"unique;not null"`
	Allergens        pq.StringArray `gorm:"type:text[]"`
	Intolerances     pq.StringArray `gorm:"type:text[]"`
	Gender           string         `gorm:"not null"`
	Age              int            `gorm:"default:null"`
	Height           int            `gorm:"default:null"`
	Weight           int            `gorm:"default:null"`
	DailyCalorieGoal int            `gorm:"default:null"`
	DailyGoalType    string         `gorm:"default:null"`
}

type RegisterUser struct {
	Allergens     []string      `json:"allergens"`
	Intolerances  []string      `json:"intolerances"`
	TokenResponse TokenResponse `json:"tokenResponse"` // match exactly: tokenResponse
	Username      string        `json:"username"`
	Gender        string        `json:"gender"`
}

type GoogleUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type Profile struct {
	Username         string
	Allergens        pq.StringArray
	Intolerances     pq.StringArray
	Age              int
	Height           int
	Weight           int
	DailyCalorieGoal int
	DailyGoalType    string
}
