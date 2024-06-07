package dto

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
