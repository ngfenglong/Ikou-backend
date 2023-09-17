package dto

import "time"

type SignUpUserDTO struct{}

type LoginCredentialInputDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Error        bool      `json:"error"`
	Message      string    `json:"message"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
	User         UserDTO   `json:"user"`
}

type UserDTO struct {
	UserName     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Country      string `json:"country"`
	ProfileImage string `json:"profile_image"`
}

type RegisterFormInputDTO struct {
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
	Country      string `json:"country"`
}

type RefreshTokenResponseDTO struct {
	Error        bool      `json:"error"`
	Message      string    `json:"message"`
	AccessToken  string    `json:"access_token"`
	Expiry       time.Time `json:"expiry"`
}