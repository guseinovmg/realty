package dto

import "time"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	OK bool `json:"OK"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	InviteId string `json:"inviteId"`
}

type RegisterResponse struct {
	OK bool `json:"OK"`
}

type CreateAdvRequest struct {
	OriginLang   int8
	TranslatedBy int8
	TranslatedTo string
	Title        string
	Description  string
	Photos       string
	Price        int64
	Currency     string
	Country      string
	City         string
	Address      string
	Latitude     float64
	Longitude    float64
	UserComment  string
}

type GetAdvResponse struct {
	Id           int64
	UserEmail    string
	UserName     string
	Created      time.Time
	Updated      time.Time
	Approved     bool
	Lang         int8
	OriginLang   int8
	TranslatedBy int8
	Title        string
	Description  string
	Photos       string
	Price        int64
	Currency     string
	DollarPrice  int64 //не хранится в БД
	Country      string
	City         string
	Address      string
	Latitude     float64
	Longitude    float64
	Watches      int64
	SeVisible    bool
}

type UpdateAdvRequest struct {
	OriginLang   int8
	TranslatedBy int8
	TranslatedTo string
	Title        string
	Description  string
	Photos       string
	Price        int64
	Currency     string
	Country      string
	City         string
	Address      string
	Latitude     float64
	Longitude    float64
	UserComment  string
}

type UpdateUserRequest struct {
	Name        string
	Description string
}

type UpdatePasswordRequest struct {
	OldPassword string
	NewPassword string
}

type Err struct {
	ErrMessage string
}
