package models

import "time"

type User struct {
	Id            uint64
	Email         string
	Name          string
	PasswordHash  []byte
	SessionSecret [24]byte //нужно перегенерить для выхода из всех устройств
	Invite        *Invite
	Balance       float64
	Trusted       bool
	Created       time.Time
	Enabled       bool
}

type Invite struct {
	Id      int64
	Company string
	Code    int64
	Used    bool
}

type Photo struct {
	Name uint64
	Ext  byte
}

type Adv struct {
	Id                      uint64
	UserId                  uint64
	User                    *User
	Created                 time.Time
	Updated                 time.Time
	Approved                bool
	Lang                    int8
	OriginLang              int8
	TranslatedBy            int8
	TranslatedToLangs       []int8
	Title                   string
	Description             string
	Photos                  []Photo
	Price                   uint64
	Currency                string
	Country                 string
	City                    string
	Address                 string
	Latitude                float64
	Longitude               float64
	Watches                 uint64
	PaidAdv                 uint64
	VisibleForSearchEngines bool
	UserComment             string
	AdminComment            string
}

type UserAction struct {
	Type int
	Time time.Time
	User *User
}

type CurrencyRate struct {
	Currency   string
	DollarRate float64
	EuroRate   float64
}
