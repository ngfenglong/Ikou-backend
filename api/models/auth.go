package models

import (
	"time"
)

type Token struct {
	ID     string
	UserID string
	Token  string
	Expiry time.Time
}

type User struct {
	ID        string `json:"id"`
	UserName  string `json:"username"`
	FirstName string `json:"firstname"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Role      int    `json:"role"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type LoginUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
