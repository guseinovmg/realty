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

type SaveCache interface {
	Save() error
}

type AdvCache struct {
	CurrentAdv models.Adv
	OldAdv     models.Adv
	ToCreate   bool
	ToUpdate   bool
	ToDelete   bool
	Deleted    bool
	mu         sync.RWMutex
}

func (adv *AdvCache) Save() error {
	adv.mu.Lock()
	defer adv.mu.Unlock()
	if !adv.Deleted {
		if adv.ToDelete {
			err := db.DeleteAdv(adv.CurrentAdv.Id)
			if err != nil {
				return err
			}
			adv.Deleted = true
			adv.ToDelete = false
			adv.ToCreate = false
			adv.ToUpdate = false
		} else if adv.ToCreate {
			err := db.CreateAdv(&adv.CurrentAdv)
			if err != nil {
				return err
			}
			adv.OldAdv = adv.CurrentAdv
			adv.ToCreate = false
			adv.ToUpdate = false
		} else if adv.ToUpdate {
			err := db.UpdateAdvChanges(&adv.OldAdv, &adv.CurrentAdv)
			if err != nil {
				return err
			}
			adv.OldAdv = adv.CurrentAdv
			adv.ToUpdate = false
		}
	}
	return nil
}

type UserCache struct {
	CurrentUser models.User
	OldUser     models.User
	ToCreate    bool
	ToUpdate    bool
	ToDelete    bool
	Deleted     bool
	mu          sync.RWMutex
}

func (user *UserCache) Save() error {
	user.mu.Lock()
	defer user.mu.Unlock()
	if !user.Deleted {
		if user.ToDelete {
			err := db.DeleteAdv(user.CurrentUser.Id)
			if err != nil {
				return err
			}
			user.Deleted = true
			user.ToDelete = false
			user.ToCreate = false
			user.ToUpdate = false
		} else if user.ToCreate {
			err := db.CreateUser(&user.CurrentUser)
			if err != nil {
				return err
			}
			user.OldUser = user.CurrentUser
			user.ToCreate = false
			user.ToUpdate = false
		} else if user.ToUpdate {
			err := db.UpdateUserChanges(&user.OldUser, &user.CurrentUser)
			if err != nil {
				return err
			}
			user.OldUser = user.CurrentUser
			user.ToUpdate = false
		}
	}
	return nil
}

var users []*UserCache
var advs []*AdvCache
var toSave chan SaveCache
var idGenerationMutex sync.Mutex

func generateId() int64 {
	idGenerationMutex.Lock()
	defer idGenerationMutex.Unlock()
	time.Sleep(time.Millisecond)
	return time.Now().UnixMicro()
}

func Initialize() {
	users_, advs_, err := db.ReadDb()
	if err != nil {
		panic(err)
	}
	users = make([]*UserCache, len(users_), len(users_)+100)
	for i := range len(users_) {
		users[i] = &UserCache{
			CurrentUser: users_[i],
			OldUser:     users_[i],
			ToUpdate:    false,
			ToDelete:    false,
			Deleted:     false,
			mu:          sync.RWMutex{},
		}
	}
	advs = make([]*AdvCache, len(advs_), len(advs_)+500)
	for i := range len(advs_) {
		advs[i] = &AdvCache{
			CurrentAdv: advs_[i],
			OldAdv:     advs_[i],
			ToUpdate:   false,
			ToDelete:   false,
			Deleted:    false,
			mu:         sync.RWMutex{},
		}
	}
	toSave = make(chan SaveCache, 100)
	go func() {
		for saveCache := range toSave {
			for range 3 {
				err := saveCache.Save()
				if err != nil {
					time.Sleep(time.Second)
				} else {
					break
				}
			}
		}
	}()
	go func() {
		for {
			time.Sleep(time.Minute * 3)
			for i := range len(advs) {
				_ = advs[i].Save()
				time.Sleep(time.Millisecond * 100)
			}
			for i := range len(users) {
				_ = users[i].Save()
				time.Sleep(time.Millisecond * 100)
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
	for i := range len(users) {
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
	for i := range len(users) {
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
	for i := range len(advs) {
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
	location string, offset int, limit int, firstNew bool) []*dto.GetAdvResponse {
	result := make([]*dto.GetAdvResponse, 0, limit)
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
		if adv.Approved && adv.DollarPrice >= minDollarPrice && adv.DollarPrice <= maxDollarPrice &&
			adv.Longitude > minLongitude && adv.Longitude < maxLongitude &&
			adv.Latitude > minLatitude && adv.Latitude < maxLatitude &&
			(countryCode == "" || adv.Country == countryCode) &&
			(location == "" || strings.Contains(adv.Address, location)) {
			if offset > 0 {
				offset--
				continue
			}
			response := &dto.GetAdvResponse{
				Id:           adv.Id,
				UserEmail:    adv.User.Email,
				UserName:     adv.User.Name,
				Created:      adv.Created,
				Updated:      adv.Updated,
				Approved:     adv.Approved,
				Lang:         adv.Lang,
				OriginLang:   adv.OriginLang,
				TranslatedBy: adv.TranslatedBy,
				Title:        adv.Title,
				Description:  adv.Description,
				Photos:       adv.Photos,
				Price:        adv.Price,
				Currency:     adv.Currency,
				DollarPrice:  adv.DollarPrice,
				Country:      adv.Country,
				City:         adv.City,
				Address:      adv.Address,
				Latitude:     adv.Latitude,
				Longitude:    adv.Longitude,
				Watches:      adv.Watches,
				SeVisible:    adv.SeVisible,
			}
			result = append(result, response)
			if limit > 0 {
				limit--
			} else {
				break
			}
		}
	}
	return result
}

func FindUsersAdvs(userId int64, offset, limit int, firstNew bool) []*dto.GetAdvResponse {
	result := make([]*dto.GetAdvResponse, 0, limit)
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
		if adv.UserId == userId {
			if offset > 0 {
				offset--
				continue
			}
			response := &dto.GetAdvResponse{
				Id:           adv.Id,
				UserEmail:    adv.User.Email,
				UserName:     adv.User.Name,
				Created:      adv.Created,
				Updated:      adv.Updated,
				Approved:     adv.Approved,
				Lang:         adv.Lang,
				OriginLang:   adv.OriginLang,
				TranslatedBy: adv.TranslatedBy,
				Title:        adv.Title,
				Description:  adv.Description,
				Photos:       adv.Photos,
				Price:        adv.Price,
				Currency:     adv.Currency,
				DollarPrice:  adv.DollarPrice,
				Country:      adv.Country,
				City:         adv.City,
				Address:      adv.Address,
				Latitude:     adv.Latitude,
				Longitude:    adv.Longitude,
				Watches:      adv.Watches,
				SeVisible:    adv.SeVisible,
			}
			result = append(result, response)
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
		Id:           generateId(),
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
	toSave <- advCache
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
	adv.ToUpdate = true
	toSave <- adv
}

func DeleteAdv(adv *AdvCache) {
	adv.mu.Lock()
	defer adv.mu.Unlock()
	if !adv.Deleted {
		adv.ToDelete = true
	}
	toSave <- adv
}

func CreateUser(request *dto.RegisterRequest) {
	passwordHash := utils.GeneratePasswordHash(request.Password)
	newUser := &models.User{
		Id:            generateId(),
		Email:         request.Email,
		Name:          request.Name,
		PasswordHash:  passwordHash,
		SessionSecret: utils.GenerateSessionsSecret(passwordHash),
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
	toSave <- userCache
}

func UpdateUser(userCache *UserCache, request *dto.UpdateUserRequest) {
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
	userCache.CurrentUser.Name = request.Name
	userCache.CurrentUser.Description = request.Description
	userCache.ToUpdate = true
	toSave <- userCache
}

func UpdatePassword(userCache *UserCache, request *dto.UpdatePasswordRequest) {
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
	userCache.CurrentUser.PasswordHash = utils.GeneratePasswordHash(request.NewPassword)
	userCache.CurrentUser.SessionSecret = utils.GenerateSessionsSecret(userCache.CurrentUser.SessionSecret[:])
	userCache.ToUpdate = true
	toSave <- userCache
}

func UpdateSessionSecret(userCache *UserCache) {
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
	userCache.CurrentUser.SessionSecret = utils.GenerateSessionsSecret(userCache.CurrentUser.SessionSecret[:])
	userCache.ToUpdate = true
	toSave <- userCache
}

func DeleteUser(userCache *UserCache) {
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
	if !userCache.Deleted {
		userCache.ToDelete = true
	}
	toSave <- userCache

}
