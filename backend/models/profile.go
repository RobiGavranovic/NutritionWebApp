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
