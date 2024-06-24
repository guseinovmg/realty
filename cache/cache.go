package cache

import (
	"bytes"
	"realty/db"
	"realty/dto"
	"realty/models"
	"realty/utils"
	"strings"
	"sync/atomic"
	"time"
)

type AdvCache struct {
	currentAdv models.Adv
	oldAdv     models.Adv
}

type UserCache struct {
	CurrentUser models.User
	OldUser     models.User
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
				if users[i].OldUser.Id == 0 {
					err := db.CreateUser(&users[i].CurrentUser)
					if err == nil {
						users[i].OldUser = users[i].CurrentUser
					} else {
						//todo
					}
					continue
				}
				if !usersAreEqual(&users[i].OldUser, &users[i].CurrentUser) {
					err := db.UpdateUserChanges(&users[i].OldUser, &users[i].CurrentUser)
					if err == nil {
						users[i].OldUser = users[i].CurrentUser
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
		if users[i].CurrentUser.Id == id {
			return &users[i].CurrentUser
		}
	}
	return nil
}

func FindUserCacheById(id int64) *UserCache {
	for i := 0; i < len(users); i++ {
		if users[i].CurrentUser.Id == id {
			return users[i]
		}
	}
	return nil
}

func FindUserCacheByLogin(email string) *UserCache {
	for i := 0; i < len(users); i++ {
		if users[i].CurrentUser.Email == email {
			return users[i]
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

func FindAdvCacheById(id int64) *AdvCache {
	for i := 0; i < len(advs); i++ {
		if advs[i].currentAdv.Id == id {
			return advs[i]
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

func UpdateAdv(advId int64, request *dto.UpdateAdvRequest) {
	adv := FindAdvCacheById(advId)
	if adv == nil {
		return //todo
	}
	adv.currentAdv.OriginLang = request.OriginLang
	adv.currentAdv.TranslatedBy = request.TranslatedBy
	adv.currentAdv.TranslatedTo = request.TranslatedTo
	adv.currentAdv.Title = request.Title
	adv.currentAdv.Description = request.Description
	adv.currentAdv.Photos = request.Photos
	adv.currentAdv.Price = request.Price
	adv.currentAdv.Currency = request.Currency
	adv.currentAdv.Country = request.Country
	adv.currentAdv.City = request.City
	adv.currentAdv.Address = request.Address
	adv.currentAdv.Latitude = request.Latitude
	adv.currentAdv.Longitude = request.Longitude
	adv.currentAdv.UserComment = request.UserComment
}

func CreateUser(request *dto.RegisterRequest) {
	newUser := &models.User{
		Id:            time.Now().UnixMicro(),
		Email:         request.Email,
		Name:          request.Name,
		PasswordHash:  utils.GeneratePasswordHash(request.Password),
		SessionSecret: utils.GenerateSessionsSecret(),
		InviteId:      request.InviteId,
		Balance:       0,
		Trusted:       false,
		Created:       time.Now(),
		Enabled:       true,
		Description:   "",
	}
	userCache := &UserCache{
		CurrentUser: *newUser,
		OldUser:     models.User{},
	}
	users = append(users, userCache)
}
