package dto

import (
	"sync/atomic"
	"time"
)

type Metrics struct {
	UnSavedChangesQueueCount    int64         `json:"unSavedChangesCount"`
	MaxUnSavedChangesQueueCount int64         `json:"maxUnSavedChangesCount"`
	DbErrorCount                int64         `json:"dbErrorCount"`
	RecoveredPanicsCount        int64         `json:"recoveredPanicsCount"`
	InstanceStartTime           string        `json:"instanceStartTime"`
	IsGracefullyStopped         bool          `json:"isGracefullyStopped"`
	RequestsCount               RequestsCount `json:"requestsCount"`
}

type RequestsCount struct {
	GET__static_                         atomic.Int64 `json:"GET /static/"`
	GET__metrics                         atomic.Int64 `json:"GET /metrics"`
	GET__generate_id                     atomic.Int64 `json:"GET /generate/id"`
	POST__login                          atomic.Int64 `json:"POST /login"`
	GET__logout_me                       atomic.Int64 `json:"GET /logout/me"`
	GET__logout_all                      atomic.Int64 `json:"GET /logout/all"`
	POST__registration                   atomic.Int64 `json:"POST /registration"`
	PUT__password                        atomic.Int64 `json:"PUT /password"`
	PUT__user                            atomic.Int64 `json:"PUT /user"`
	GET__adv__advId_                     atomic.Int64 `json:"GET /adv/{advId}"`
	GET__adv                             atomic.Int64 `json:"GET /adv"`
	GET__user_adv__advId_                atomic.Int64 `json:"GET /user/adv/{advId}"`
	GET__user_adv                        atomic.Int64 `json:"GET /user/adv"`
	POST__user_adv                       atomic.Int64 `json:"POST /user/adv"`
	POST__adv                            atomic.Int64 `json:"POST /adv"`
	PUT__adv__advId_                     atomic.Int64 `json:"PUT /adv/{advId}"`
	DELETE__adv__advId_                  atomic.Int64 `json:"DELETE /adv/{advId}"`
	POST__adv__advId__photos             atomic.Int64 `json:"POST /adv/{advId}/photos"`
	DELETE__adv__advId__photos__photoId_ atomic.Int64 `json:"DELETE /adv/{advId}/photos/{photoId}"`
	UNKNOWN                              atomic.Int64 `json:"UNKNOWN"`
}

type GenerateIdResponse struct {
	Id int64 `json:"id"`
}

type Result struct {
	Result string `json:"result"`
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
	Price        int64   `json:"price,omitempty"`
	Latitude     float64 `json:"latitude,omitempty"`
	Longitude    float64 `json:"longitude,omitempty"`
	TranslatedTo string  `json:"translatedTo,omitempty"`
	Title        string  `json:"title,omitempty"`
	Description  string  `json:"description,omitempty"`
	Photos       string  `json:"photos,omitempty"`
	Currency     string  `json:"currency,omitempty"`
	Country      string  `json:"country,omitempty"`
	City         string  `json:"city,omitempty"`
	Address      string  `json:"address,omitempty"`
	UserComment  string  `json:"userComment,omitempty"`
}

type CreateAdvResponse struct {
	AdvId     int64 `json:"advId"`
	RequestId int64 `json:"requestId"`
}

type AddPhotoRequest struct {
	Filename string `json:"filename"`
}

type GetAdvListRequest struct {
	FirstNew     bool    `json:"firstNew,omitempty"`
	Page         int     `json:"page,omitempty"`
	MinPrice     int64   `json:"minPrice,omitempty"`
	MaxPrice     int64   `json:"maxPrice,omitempty"`
	MinLongitude float64 `json:"minLongitude,omitempty"`
	MaxLongitude float64 `json:"maxLongitude,omitempty"`
	MinLatitude  float64 `json:"minLatitude,omitempty"`
	MaxLatitude  float64 `json:"maxLatitude,omitempty"`
	Currency     string  `json:"currency,omitempty"`
	CountryCode  string  `json:"countryCode,omitempty"`
	Location     string  `json:"location,omitempty"`
}

type GetUserAdvListRequest struct {
	Page     int  `json:"page,omitempty"`
	FirstNew bool `json:"firstNew,omitempty"`
}

type GetAdvResponseItem struct {
	Id           int64     `json:"id,omitempty"`
	Price        int64     `json:"price,omitempty"`
	DollarPrice  int64     `json:"dollarPrice,omitempty"` //не хранится в БД
	Watches      int64     `json:"watches,omitempty"`
	Latitude     float64   `json:"latitude,omitempty"`
	Longitude    float64   `json:"longitude,omitempty"`
	Approved     bool      `json:"approved,omitempty"`
	SeVisible    bool      `json:"seVisible,omitempty"`
	Lang         int8      `json:"lang,omitempty"`
	OriginLang   int8      `json:"originLang,omitempty"`
	TranslatedBy int8      `json:"translatedBy,omitempty"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	UserEmail    string    `json:"userEmail,omitempty"`
	UserName     string    `json:"userName,omitempty"`
	Title        string    `json:"title,omitempty"`
	Description  string    `json:"description,omitempty"`
	Currency     string    `json:"currency,omitempty"`
	Country      string    `json:"country,omitempty"`
	City         string    `json:"city,omitempty"`
	Address      string    `json:"address,omitempty"`
	UserComment  string    `json:"userComment,omitempty"`
	Photos       []string  `json:"photos,omitempty"`
}

type GetAdvListResponse struct {
	Count int                   `json:"count"`
	List  []*GetAdvResponseItem `json:"list"`
}

type UpdateAdvRequest struct {
	OriginLang   int8    `json:"originLang,omitempty"`
	TranslatedBy int8    `json:"translatedBy,omitempty"`
	Price        int64   `json:"price,omitempty"`
	Latitude     float64 `json:"latitude,omitempty"`
	Longitude    float64 `json:"longitude,omitempty"`
	TranslatedTo string  `json:"translatedTo,omitempty"`
	Title        string  `json:"title,omitempty"`
	Description  string  `json:"description,omitempty"`
	Photos       string  `json:"photos,omitempty"`
	Currency     string  `json:"currency,omitempty"`
	Country      string  `json:"country,omitempty"`
	City         string  `json:"city,omitempty"`
	Address      string  `json:"address,omitempty"`
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
	RequestId  int64  `json:"requestId,omitempty"`
	ErrMessage string `json:"errMessage,omitempty"`
}
