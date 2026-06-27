package models

import (
	"time"
)

type User struct {
	ID int `json:"id"`
	Role UserRole `json:"role"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRole string

const (
	Customer UserRole = "customer"
	Merchant UserRole = "merchant"
	Admin UserRole = "admin"
)