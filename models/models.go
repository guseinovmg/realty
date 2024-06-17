package models

import "time"

type User struct {
	Id            int64
	Email         string
	Name          string
	PasswordHash  []byte
	SessionSecret [24]byte //нужно перегенерить для выхода из всех устройств
	Invite        *Invite
	InviteId      int64
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
	Name int64
	Ext  byte
}

type Adv struct {
	Id                      int64
	UserId                  int64
	User                    *User
	Created                 time.Time
	Updated                 time.Time
	Approved                bool
	Lang                    int8
	OriginLang              int8
	TranslatedBy            int8
	TranslatedTo            string
	Title                   string
	Description             string
	Photos                  string
	Price                   int64
	Currency                string
	Country                 string
	City                    string
	Address                 string
	Latitude                float64
	Longitude               float64
	Watches                 int64
	PaidAdv                 int64
	VisibleForSearchEngines bool
	UserComment             string
	AdminComment            string
}

type CurrencyRate struct {
	Currency   string
	DollarRate float64
	EuroRate   float64
}
