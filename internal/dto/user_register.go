package dto

import (
	"auth-service/internal/models"
	"errors"
	"regexp"
)

type UserRegister struct {
	Email string `json:"email" binding:"required,email"`
	Role models.UserRole `json:"role" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (u *UserRegister) Validate() error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}