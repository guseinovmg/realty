package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	OK bool `json:"OK"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
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
