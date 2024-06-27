package cache

import (
	"realty/db"
	"realty/dto"
	"realty/models"
	"realty/utils"
	"strings"
	"sync"
	"time"
)

type AdvCache struct {
	CurrentAdv  models.Adv
	OldAdv      models.Adv
	UpdateCount int64
	ToDelete    bool
	Deleted     bool
	mu          sync.RWMutex
}

type UserCache struct {
	CurrentUser models.User
	OldUser     models.User
	UpdateCount int64
	ToDelete    bool
	Deleted     bool
	mu          sync.RWMutex
}

var users []*UserCache
var advs []*AdvCache

func Initialize() {
	db.ReadDb()
	go func() {
		for {
			time.Sleep(time.Second)
			var adv *AdvCache
			for i := 0; i < len(advs); i++ {
				adv = advs[i]
				adv.mu.Lock()
				if adv.ToDelete {
					err := db.DeleteAdv(adv.CurrentAdv.Id)
					if err == nil {
						adv.Deleted = true
						adv.ToDelete = false
					} else {
						//todo
					}
					continue
				}
				if adv.OldAdv.Id == 0 {
					err := db.CreateAdv(&adv.CurrentAdv)
					if err == nil {
						adv.OldAdv = adv.CurrentAdv
					} else {
						//todo
					}
					continue
				}
				if adv.UpdateCount > 0 {
					err := db.UpdateAdvChanges(&adv.OldAdv, &adv.CurrentAdv)
					if err == nil {
						adv.OldAdv = adv.CurrentAdv
						adv.UpdateCount = 0
					} else {
						//todo
					}
				}
				adv.mu.Unlock()
			}
			time.Sleep(time.Second)
			var user *UserCache
			for i := 0; i < len(users); i++ {
				user = users[i]
				user.mu.Lock()
				if user.ToDelete {
					err := db.DeleteAdv(user.CurrentUser.Id)
					if err == nil {
						user.Deleted = true
						user.ToDelete = false
					} else {
						//todo
					}
					continue
				}
				if user.OldUser.Id == 0 {
					err := db.CreateUser(&user.CurrentUser)
					if err == nil {
						user.OldUser = user.CurrentUser
					} else {
						//todo
					}
					continue
				}

				if user.UpdateCount > 0 {
					err := db.UpdateUserChanges(&user.OldUser, &user.CurrentUser)
					if err == nil {
						user.OldUser = user.CurrentUser
						user.UpdateCount++
					} else {
						//todo
					}
				}
				user.mu.Unlock()
			}
		}
	}()
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
			if users[i].ToDelete || users[i].Deleted {
				return nil
			}
			return users[i]
		}
	}
	return nil
}

func FindUserCacheByLogin(email string) *UserCache {
	for i := 0; i < len(users); i++ {
		if users[i].CurrentUser.Email == email {
			if users[i].ToDelete || users[i].Deleted {
				return nil
			}
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
			if advs[i].ToDelete || advs[i].Deleted {
				return nil
			}
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
		if advs[i].ToDelete || advs[i].Deleted {
			continue
		}
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
			result = append(result, adv)
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
		Watches:      0,
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
	adv.mu.Lock()
	defer adv.mu.Unlock()
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
	adv.UpdateCount++
}

func DeleteAdv(adv *AdvCache) {
	adv.mu.Lock()
	defer adv.mu.Unlock()
	adv.ToDelete = true
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

func UpdateUser(userCache *UserCache, request *dto.UpdateUserRequest) {
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
	userCache.CurrentUser.Name = request.Name
	userCache.CurrentUser.Description = request.Description
	userCache.UpdateCount++
}

func UpdatePassword(userCache *UserCache, request *dto.UpdatePasswordRequest) {
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
	userCache.CurrentUser.PasswordHash = utils.GeneratePasswordHash(request.NewPassword)
	userCache.CurrentUser.SessionSecret = utils.GenerateSessionsSecret()
	userCache.UpdateCount++
}

func DeleteUser(user *UserCache) {
	user.mu.Lock()
	defer user.mu.Unlock()
	user.ToDelete = true
}
