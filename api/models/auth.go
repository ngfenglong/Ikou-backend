package models

import (
	"time"
)

type Token struct {
	ID           string
	UserID       string
	Token        string
	ProfileImage string
	FirstName    string
	LastName     string `json:"firstname"`
	Expiry       time.Time
}

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"firstname"`
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
