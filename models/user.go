package models

import (
	"gorm.io/gorm"
	// "time"
)

type User struct {
	gorm.Model          
	Name         string `gorm:"size:255;not null;index"`
	Email        string `gorm:"size:255;unique;not null;index"`
	Password     string `gorm:"size:255;not null"` 
	Role         string `gorm:"size:20;default:user;index"`
	RefreshToken string `gorm:"size:500"`

	// For response (avoid exposing password)
	PasswordHash string `gorm:"-" json:"-"`
}

