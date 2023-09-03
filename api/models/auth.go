package models

import (
	"time"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"-"`
	Country      string `json:"country"`
	ProfileImage string `json:"profile_image"`
	Role         int    `json:"role"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type LoginUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	CreatedAt string
	UpdatedAt string
	ExpiresAt time.Time
}
