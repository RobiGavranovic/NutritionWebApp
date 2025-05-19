package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID uint
	Email  string
	jwt.RegisteredClaims
}
