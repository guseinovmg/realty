package dto

import "time"

type Metrics struct {
	InstanceStartTime           int64   `json:"instanceStartTime,omitempty"`
	FreeRAM                     int64   `json:"freeRAM,omitempty"`
	CPUTemp                     float64 `json:"cpuTemp,omitempty"`
	CPUConsumption              float64 `json:"cpuConsumption,omitempty"`
	UnSavedChangesQueueCount    int64   `json:"unSavedChangesCount,omitempty"`
	DiskUsagePercent            float64 `json:"diskUsagePercent,omitempty"`
	RecoveredPanicsCount        int64   `json:"recoveredPanicsCount,omitempty"`
	MaxRAMConsumptions          int64   `json:"maxRAMConsumptions,omitempty"`
	MaxCPUConsumptions          int64   `json:"maxCPUConsumptions,omitempty"`
	MaxRPS                      int64   `json:"maxRPS,omitempty"`
	MaxUnSavedChangesQueueCount int64   `json:"maxUnSavedChangesCount,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	InviteId string `json:"inviteId"`
}

type CreateAdvRequest struct {
	OriginLang   int8    `json:"originLang,omitempty"`
	TranslatedBy int8    `json:"translatedBy,omitempty"`
	TranslatedTo string  `json:"translatedTo,omitempty"`
	Title        string  `json:"title,omitempty"`
	Description  string  `json:"description,omitempty"`
	Photos       string  `json:"photos,omitempty"`
	Price        int64   `json:"price,omitempty"`
	Currency     string  `json:"currency,omitempty"`
	Country      string  `json:"country,omitempty"`
	City         string  `json:"city,omitempty"`
	Address      string  `json:"address,omitempty"`
	Latitude     float64 `json:"latitude,omitempty"`
	Longitude    float64 `json:"longitude,omitempty"`
	UserComment  string  `json:"userComment,omitempty"`
}

type GetAdvListRequest struct {
	Currency     string  `json:"currency,omitempty"`
	MinPrice     int64   `json:"minPrice,omitempty"`
	MaxPrice     int64   `json:"maxPrice,omitempty"`
	MinLongitude float64 `json:"minLongitude,omitempty"`
	MaxLongitude float64 `json:"maxLongitude,omitempty"`
	MinLatitude  float64 `json:"minLatitude,omitempty"`
	MaxLatitude  float64 `json:"maxLatitude,omitempty"`
	CountryCode  string  `json:"countryCode,omitempty"`
	Location     string  `json:"location,omitempty"`
	Page         int     `json:"page,omitempty"`
	FirstNew     bool    `json:"firstNew,omitempty"`
}

type GetUserAdvListRequest struct {
	Page     int  `json:"page,omitempty"`
	FirstNew bool `json:"firstNew,omitempty"`
}

type GetAdvResponseItem struct {
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
	Photos       []string  `json:"photos,omitempty"`
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

type GetAdvListResponse struct {
	List  []*GetAdvResponseItem `json:"list"`
	Count int                   `json:"count"`
}

type UpdateAdvRequest struct {
	OriginLang   int8    `json:"originLang,omitempty"`
	TranslatedBy int8    `json:"translatedBy,omitempty"`
	TranslatedTo string  `json:"translatedTo,omitempty"`
	Title        string  `json:"title,omitempty"`
	Description  string  `json:"description,omitempty"`
	Photos       string  `json:"photos,omitempty"`
	Price        int64   `json:"price,omitempty"`
	Currency     string  `json:"currency,omitempty"`
	Country      string  `json:"country,omitempty"`
	City         string  `json:"city,omitempty"`
	Address      string  `json:"address,omitempty"`
	Latitude     float64 `json:"latitude,omitempty"`
	Longitude    float64 `json:"longitude,omitempty"`
	UserComment  string  `json:"userComment,omitempty"`
}

type UpdateUserRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}

type Err struct {
	ErrMessage string `json:"errMessage,omitempty"`
}
