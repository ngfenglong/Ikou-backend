package dto

import "time"

type SignUpUser struct{}

type LoginCredentialInputDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDto struct {
	Error        bool      `json:"error"`
	Message      string    `json:"message"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
	User         UserDto   `json:"user"`
}

type UserDto struct {
	UserName     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Country      string `json:"country"`
	ProfileImage string `json:"profile_image"`
}

type RegisterFormInputDto struct {
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
	Country      string `json:"country"`
}
