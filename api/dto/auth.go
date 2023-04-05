package dto

type SignUpUser struct{}

type LoginCredentialInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterFormInput struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Country   string `json:"country"`
}
