package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID int `json:"user_id"`
	Email string `json:"email"`
	Role UserRole `json:"role"`
	jwt.RegisteredClaims
}