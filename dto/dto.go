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
	Id           int64     `json:"id,omitempty"`
	UserEmail    string    `json:"userEmail,omitempty"`
	UserName     string    `json:"userName,omitempty"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	Approved     bool      `json:"approved,omitempty"`
	Lang         int8      `json:"lang,omitempty"`
	OriginLang   int8      `json:"originLang,omitempty"`
	TranslatedBy int8      `json:"translatedBy,omitempty"`
	Title        string    `json:"title,omitempty"`
	Description  string    `json:"description,omitempty"`
	Photos       string    `json:"photos,omitempty"`
	Price        int64     `json:"price,omitempty"`
	Currency     string    `json:"currency,omitempty"`
	DollarPrice  int64     `json:"dollarPrice,omitempty"` //не хранится в БД
	Country      string    `json:"country,omitempty"`
	City         string    `json:"city,omitempty"`
	Address      string    `json:"address,omitempty"`
	Latitude     float64   `json:"latitude,omitempty"`
	Longitude    float64   `json:"longitude,omitempty"`
	Watches      int64     `json:"watches,omitempty"`
	SeVisible    bool      `json:"seVisible,omitempty"`
	UserComment  string    `json:"userComment,omitempty"`
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
