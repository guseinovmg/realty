package cache

import (
	"bytes"
	"realty/db"
	"realty/dto"
	"realty/models"
	"strings"
	"sync/atomic"
	"time"
)

type AdvCache struct {
	currentAdv models.Adv
	oldAdv     models.Adv
}

type UserCache struct {
	currentUser models.User
	oldUser     models.User
}

var users []*UserCache
var advs []*AdvCache

func Initialize() {
	db.ReadDb()
	go func() {
		for {
			time.Sleep(time.Second)
			for i := 0; i < len(advs); i++ {
				if advs[i].oldAdv.Id == 0 {
					err := db.CreateAdv(&advs[i].currentAdv)
					if err == nil {
						advs[i].oldAdv = advs[i].currentAdv
					} else {
						//todo
					}
					continue
				}
				if advs[i].oldAdv != advs[i].currentAdv {
					err := db.UpdateAdvChanges(&advs[i].oldAdv, &advs[i].currentAdv)
					if err == nil {
						advs[i].oldAdv = advs[i].currentAdv
					} else {
						//todo
					}
				}
			}
			time.Sleep(time.Second)
			for i := 0; i < len(users); i++ {
				if !usersAreEqual(&users[i].oldUser, &users[i].currentUser) {
					err := db.UpdateUserChanges(&users[i].oldUser, &users[i].currentUser)
					if err == nil {
						users[i].oldUser = users[i].currentUser
					} else {
						//todo
					}
				}
			}

		}
	}()
}

func usersAreEqual(u1, u2 *models.User) bool {
	if u1.Id != u2.Id {
		return false
	}
	if u1.Email != u2.Email {
		return false
	}
	if u1.Name != u2.Name {
		return false
	}
	if !bytes.Equal(u1.PasswordHash, u2.PasswordHash) {
		return false
	}
	if !bytes.Equal(u1.SessionSecret[:], u2.SessionSecret[:]) {
		return false
	}
	if u1.InviteId != u2.InviteId {
		return false
	}
	if u1.Trusted != u2.Trusted {
		return false
	}
	if u1.Enabled != u2.Enabled {
		return false
	}
	if u1.Balance != u2.Balance {
		return false
	}
	if !u1.Created.Equal(u2.Created) {
		return false
	}
	if u1.Description != u2.Description {
		return false
	}

	return true
}

func FindUserById(id int64) *models.User {
	for i := 0; i < len(users); i++ {
		if users[i].currentUser.Id == id {
			return &users[i].currentUser
		}
	}
	return nil
}

func FindUserByLogin(email string) *models.User {
	for i := 0; i < len(users); i++ {
		if users[i].currentUser.Email == email {
			return &users[i].currentUser
		}
	}
	return nil
}

func FindAdvById(id int64) *models.Adv {
	for i := 0; i < len(advs); i++ {
		if advs[i].currentAdv.Id == id {
			return &advs[i].currentAdv
		}
	}
	return nil
}

func FindAdvs(minDollarPrice int64, maxDollarPrice int64, minLongitude float64,
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
		adv = &advs[i].currentAdv
		if adv.DollarPrice >= minDollarPrice && adv.DollarPrice <= maxDollarPrice &&
			adv.Longitude > minLongitude && adv.Longitude < maxLongitude &&
			adv.Latitude > minLatitude && adv.Latitude < maxLatitude &&
			(countryCode == "" || adv.Country == countryCode) &&
			(location == "" || strings.Contains(adv.Address, location)) {
			if offset > 0 {
				offset--
				continue
			}
			result = append(result, &advs[i].currentAdv)
			if limit > 0 {
				limit--
			} else {
				break
			}
		}
	}
	return result
}

func CreateAdv(user *models.User, request *dto.CreateAdvRequest) {
	newAdv := &models.Adv{
		Id:           time.Now().UnixMicro(),
		UserId:       user.Id,
		User:         user,
		Created:      time.Now(),
		Updated:      time.Now(),
		Approved:     false,
		Lang:         request.OriginLang,
		OriginLang:   request.OriginLang,
		TranslatedBy: request.TranslatedBy,
		TranslatedTo: request.TranslatedTo,
		Title:        request.Title,
		Description:  request.Description,
		Photos:       request.Photos,
		Price:        request.Price,
		Currency:     request.Currency,
		DollarPrice:  0, //todo
		Country:      request.Country,
		City:         request.City,
		Address:      request.Address,
		Latitude:     request.Latitude,
		Longitude:    request.Longitude,
		Watches:      atomic.Int64{},
		PaidAdv:      0,
		SeVisible:    true,
		UserComment:  request.UserComment,
		AdminComment: "",
	}
	advCache := &AdvCache{
		currentAdv: *newAdv,
		oldAdv:     models.Adv{},
	}
	advs = append(advs, advCache)

}

func UpdateAdv(request *dto.UpdateAdvRequest) {

}

func CreateUser(request *dto.RegisterRequest) {

}

func UpdateUser(request *dto.UpdateUserRequest) {

}

func UpdatePassword(request *dto.UpdatePasswordRequest) {

}
