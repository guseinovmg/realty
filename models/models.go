package models

import (
	"time"
)

type User struct {
	Id            int64
	Email         string
	Name          string
	PasswordHash  []byte   `json:"-"`
	SessionSecret [24]byte `json:"-"` //нужно перегенерить для выхода из всех устройств
	InviteId      string
	Balance       float64
	Trusted       bool
	Created       time.Time
	Enabled       bool
	Description   string
}

type Invite struct {
	Id      string
	Company string
	Used    bool
}

type Photo struct {
	AdvId int64
	Id    int64
	Ext   byte
}

type Watches struct {
	AdvId   int64
	Watches int64
}

type Adv struct {
	Id           int64
	UserId       int64
	User         *User
	Updated      time.Time
	Approved     bool
	Lang         int8
	OriginLang   int8
	TranslatedBy int8
	TranslatedTo string
	Title        string
	Description  string
	Photos       []*Photo
	Price        int64
	Currency     string
	DollarPrice  int64 //не хранится в БД
	Country      string
	City         string
	Address      string
	Latitude     float64
	Longitude    float64
	Watches      *Watches
	PaidAdv      int64
	SeVisible    bool
	UserComment  string
	AdminComment string
}

type CurrencyRate struct {
	Currency   string
	DollarRate float64
}
