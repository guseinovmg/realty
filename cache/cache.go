package cache

import (
	"realty/db"
	"realty/models"
	"strings"
)

var users []models.User
var advs []*models.Adv

func Initialize() {
	db.ReadDb()
}

func FindUserById(id uint64) *models.User {
	for i := 0; i < len(users); i++ {
		if users[i].Id == id {
			return &users[i]
		}
	}
	return nil
}

func FindAdvById(id uint64) *models.Adv {
	for i := 0; i < len(advs); i++ {
		if advs[i].Id == id {
			return advs[i]
		}
	}
	return nil
}

func FindAdvs(minPrice uint64, maxPrice uint64, currency string, minLongitude float64,
	maxLongitude float64, minLatitude float64, maxLatitude float64, countryCode string,
	location string, offset int, limit int, sortedArr *[]*models.Adv) []*models.Adv {
	result := make([]*models.Adv, 0, limit)
	var adv *models.Adv
	for i := 0; i < len(*sortedArr); i++ {
		adv = (*sortedArr)[i]
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
			result = append(result, advs[i])
			if limit > 0 {
				limit--
			} else {
				break
			}
		}
	}
	return result
}

func CalcPrice(price uint64, fromCurrency, toCurrency string) uint64 {
	return 1
}

func ReloadAdvFromDb(id int64) error {
	return nil
}

func ReloadUserFromDb(id int64) error {
	return nil
}
