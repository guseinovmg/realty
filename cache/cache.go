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
	CurrentAdv  models.Adv
	OldAdv      models.Adv
	UpdateCount atomic.Int64
}

type UserCache struct {
	CurrentUser models.User
	OldUser     models.User
	UpdateCount atomic.Int64
}

var users []*UserCache
var advs []*AdvCache

func Initialize() {
	db.ReadDb()
	go func() {
		for {
			time.Sleep(time.Second)
			for i := 0; i < len(advs); i++ {
				if advs[i].OldAdv.Id == 0 {
					err := db.CreateAdv(&advs[i].CurrentAdv)
					if err == nil {
						advs[i].OldAdv = advs[i].CurrentAdv
					} else {
						//todo
					}
					continue
				}
				if advs[i].UpdateCount.Load() > 0 {
					err := db.UpdateAdvChanges(&advs[i].OldAdv, &advs[i].CurrentAdv)
					if err == nil {
						advs[i].OldAdv = advs[i].CurrentAdv
						advs[i].UpdateCount.Store(0)
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
				if users[i].UpdateCount.Load() > 0 {
					err := db.UpdateUserChanges(&users[i].OldUser, &users[i].CurrentUser)
					if err == nil {
						users[i].OldUser = users[i].CurrentUser
						users[i].UpdateCount.Store(0)
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
	userCache := FindUserCacheById(id)
	if userCache == nil {
		return nil
	}
	return &userCache.CurrentUser
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
	advCache := FindAdvCacheById(id)
	if advCache == nil {
		return nil
	}
	return &advCache.CurrentAdv
}

func FindAdvCacheById(id int64) *AdvCache {
	for i := 0; i < len(advs); i++ {
		if advs[i].CurrentAdv.Id == id {
			return advs[i]
		}
	}
	return nil
}

func FindAdvs(minDollarPrice int64, maxDollarPrice int64, minLongitude float64,
	maxLongitude float64, minLatitude float64, maxLatitude float64, countryCode string,
	location string, offset int, limit int, firstNew bool) []*models.Adv {
	result := make([]*models.Adv, 0, limit)
	var i, step int
	length := len(advs)
	if firstNew {
		i = length - 1
		step = -1
	} else {
		i = 0
		step = 1
	}
	var adv *models.Adv
	for ; i < length && i >= 0; i += step {
		adv = &advs[i].CurrentAdv
		if adv.DollarPrice >= minDollarPrice && adv.DollarPrice <= maxDollarPrice &&
			adv.Longitude > minLongitude && adv.Longitude < maxLongitude &&
			adv.Latitude > minLatitude && adv.Latitude < maxLatitude &&
			(countryCode == "" || adv.Country == countryCode) &&
			(location == "" || strings.Contains(adv.Address, location)) {
			if offset > 0 {
				offset--
				continue
			}
			result = append(result, &advs[i].CurrentAdv)
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
		CurrentAdv: *newAdv,
		OldAdv:     models.Adv{},
	}
	advs = append(advs, advCache)

}

func UpdateAdv(adv *AdvCache, request *dto.UpdateAdvRequest) {
	adv.CurrentAdv.OriginLang = request.OriginLang
	adv.CurrentAdv.TranslatedBy = request.TranslatedBy
	adv.CurrentAdv.TranslatedTo = request.TranslatedTo
	adv.CurrentAdv.Title = request.Title
	adv.CurrentAdv.Description = request.Description
	adv.CurrentAdv.Photos = request.Photos
	adv.CurrentAdv.Price = request.Price
	adv.CurrentAdv.Currency = request.Currency
	adv.CurrentAdv.Country = request.Country
	adv.CurrentAdv.City = request.City
	adv.CurrentAdv.Address = request.Address
	adv.CurrentAdv.Latitude = request.Latitude
	adv.CurrentAdv.Longitude = request.Longitude
	adv.CurrentAdv.UserComment = request.UserComment
	adv.UpdateCount.Add(1)
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
