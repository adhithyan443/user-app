package domain

import (
	"gorm.io/gorm"
	// "time"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:255;not null;index"`
	Email    string `gorm:"size:255;unique;not null;index"`
	Password string `gorm:"size:255;not null"`
	Role     string `gorm:"size:20;default:user;index"`

	// For response (avoid exposing password)
	PasswordHash string `gorm:"-" json:"-"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}


type LoginInput struct{
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}