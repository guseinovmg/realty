package cache

import (
	"realty/db"
	"realty/models"
	"strings"
)

type AdvCache struct {
	currentAdv *models.Adv
	oldAdv     *models.Adv
}

var users []models.User
var advs []*AdvCache

func Initialize() {
	db.ReadDb()
}

func FindUserById(id int64) *models.User {
	for i := 0; i < len(users); i++ {
		if users[i].Id == id {
			return &users[i]
		}
	}
	return nil
}

func FindUserByLogin(email string) *models.User {
	for i := 0; i < len(users); i++ {
		if users[i].Email == email {
			return &users[i]
		}
	}
	return nil
}

func FindAdvById(id int64) *models.Adv {
	for i := 0; i < len(advs); i++ {
		if advs[i].currentAdv.Id == id {
			return advs[i].currentAdv
		}
	}
	return nil
}

func FindAdvs(minPrice uint64, maxPrice uint64, currency string, minLongitude float64,
	maxLongitude float64, minLatitude float64, maxLatitude float64, countryCode string,
	location string, offset int, limit int, firstCheap bool) []*models.Adv {
	result := make([]*models.Adv, 0, limit)
	var i, step int
	length := len(advs)
	if firstCheap {
		i = 0
		step = 1
	} else {
		i = length - 1
		step = -1
	}
	var adv *models.Adv
	for ; i < length && i >= 0; i += step {
		adv = advs[i].currentAdv
		price := CalcPrice(adv.Price, adv.Currency, currency)
		if price >= minPrice && price <= maxPrice &&
			adv.Longitude > minLongitude && adv.Longitude < maxLongitude &&
			adv.Latitude > minLatitude && adv.Latitude < maxLatitude &&
			(countryCode == "" || adv.Country == countryCode) &&
			(location == "" || strings.Contains(adv.Address, location)) {
			if offset > 0 {
				offset--
				continue
			}
			result = append(result, advs[i].currentAdv)
			if limit > 0 {
				limit--
			} else {
				break
			}
		}
	}
	return result
}

func CalcPrice(price int64, fromCurrency, toCurrency string) uint64 {
	return 1
}

func ReloadAdvFromDb(id int64) error {
	return nil
}

func ReloadUserFromDb(id int64) error {
	return nil
}
