package dto

type SuccessResponseDto struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
