package models

import "github.com/lib/pq"

type UpdateAllergens struct {
	TokenResponse TokenResponse  `json:"tokenResponse"`
	Allergens     pq.StringArray `json:"allergens"`
}

type UpdateIntolarences struct {
	TokenResponse TokenResponse  `json:"tokenResponse"`
	Intolerances  pq.StringArray `json:"intolerances"`
}

type UpdateUsername struct {
	TokenResponse TokenResponse `json:"tokenResponse"`
	Username      string        `json:"username"`
}

type UpdatePersonalInfo struct {
	TokenResponse TokenResponse `json:"tokenResponse"`
	Age           int           `json:"age"`
	Height        int           `json:"height"`
	Weight        int           `json:"weight"`
}

type UpdateDailyCalorieGoal struct {
	TokenResponse    TokenResponse `json:"tokenResponse"`
	DailyCalorieGoal int           `json:"dailyCalorieGoal"`
	DailyGoalType    string        `json:"dailyGoalType"`
}
