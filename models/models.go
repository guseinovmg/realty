package models

import (
	"time"
)

type User struct {
	Id            int64
	Balance       float64
	Trusted       bool
	Enabled       bool
	Email         string
	Name          string
	InviteId      string
	Description   string
	PasswordHash  []byte   `json:"-"`
	SessionSecret [24]byte `json:"-"` //нужно перегенерить для выхода из всех устройств
}

type Invite struct {
	Used    bool
	Id      string
	Company string
}

type Photo struct {
	AdvId int64
	Id    int64
	Ext   byte
}

type Watches struct {
	AdvId int64
	Count int64
}

type Adv struct {
	Id           int64
	UserId       int64
	Price        int64
	DollarPrice  int64 //не хранится в БД
	PaidAdv      int64
	Latitude     float64
	Longitude    float64
	Approved     bool
	SeVisible    bool
	Lang         int8
	OriginLang   int8
	TranslatedBy int8
	Updated      time.Time
	TranslatedTo string
	Title        string
	Description  string
	Currency     string
	Country      string
	City         string
	Address      string
	UserComment  string
	AdminComment string
	User         *User
}

type CurrencyRate struct {
	Currency   string
	DollarRate float64
}
