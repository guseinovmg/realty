package dto

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
