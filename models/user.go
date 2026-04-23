package models

import "time"

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"` //Never send password to fronend
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
}
