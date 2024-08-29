package cache

import (
	"log/slog"
	"os"
	"realty/db"
	"realty/dto"
	"realty/models"
	"realty/utils"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type SaveCache interface {
	Save() error
}

type SaveTask struct {
	Cache     SaveCache
	RequestId int64
}

var users []*UserCache
var advs []*AdvCache
var photos []*PhotoCache
var watches []*WatchesCache

var usersRWMutex sync.RWMutex
var advsRWMutex sync.RWMutex
var photosRWMutex sync.RWMutex
var watchesRWMutex sync.RWMutex

var toSave chan SaveTask

var gracefullyStop atomic.Bool

func GracefullyStopAndExitApp() {
	if !gracefullyStop.Load() {
		gracefullyStop.Store(true)
	}
}

func IsGracefullyStopped() bool {
	return gracefullyStop.Load()
}

func Initialize() {
	users_, advs_, photos_, watches_, errDb := db.ReadDb()
	if errDb != nil {
		panic(errDb)
	}

	photos = make([]*PhotoCache, len(photos_), len(photos_)+500)
	for i := range len(photos_) {
		photos[i] = &PhotoCache{
			Photo:    *photos_[i],
			ToDelete: false,
			Deleted:  false,
			mu:       sync.RWMutex{},
		}
	}
	watches = make([]*WatchesCache, len(watches_), len(watches_)+500)
	for i := range len(watches_) {
		watches[i] = &WatchesCache{
			Watches:  *watches_[i],
			ToCreate: false,
			ToUpdate: false,
			ToDelete: false,
			Deleted:  false,
			mu:       sync.RWMutex{},
		}
	}

	users = make([]*UserCache, len(users_), len(users_)+100)
	for i := range len(users_) {
		user := users_[i]
		users[i] = &UserCache{
			CurrentUser: *user,
			OldUser:     *user,
			ToUpdate:    false,
			ToDelete:    false,
			Deleted:     false,
			mu:          sync.RWMutex{},
		}
	}
	advs = make([]*AdvCache, len(advs_), len(advs_)+500)
	for i := range len(advs_) {
		adv := advs_[i]
		advs[i] = &AdvCache{
			CurrentAdv: *adv,
			OldAdv:     *adv,
			Photos:     GetPhotosByAdvId(adv.Id),
			Watches:    FindWatchesCacheById(adv.Id),
			ToUpdate:   false,
			ToDelete:   false,
			Deleted:    false,
			mu:         sync.RWMutex{},
		}
	}

	//todo надо просмотры и фото в adv добавить
	toSave = make(chan SaveTask, 1000)

	go func() {
		for saveCache := range toSave {
			for range 2 {
				if errSave := saveCache.Cache.Save(); errSave == nil {
					slog.Debug("saving", "requestId", saveCache.RequestId, "msg", "ok")
					break
				} else {
					slog.Error("saving", "requestId", saveCache.RequestId, "msg", errSave.Error())
				}
				time.Sleep(time.Millisecond * 100)
			}
			if gracefullyStop.Load() && len(toSave) == 0 {
				break
			}
		}
		os.Exit(1)
	}()

	go func() {
		for {
			time.Sleep(time.Minute)
			for i := range len(advs) {
				if gracefullyStop.Load() {
					return
				}
				if err := advs[i].Save(); err != nil {
					slog.Error("saving", "msg", err.Error())
					time.Sleep(time.Millisecond * 100)
				}
			}
			time.Sleep(time.Second)
			for i := range len(users) {
				if gracefullyStop.Load() {
					return
				}
				if err := users[i].Save(); err != nil {
					slog.Error("saving", "msg", err.Error())
					time.Sleep(time.Millisecond * 100)
				}
			}
			time.Sleep(time.Second)
			for i := range len(photos) {
				if gracefullyStop.Load() {
					return
				}
				if err := photos[i].Save(); err != nil {
					slog.Error("saving", "msg", err.Error())
					time.Sleep(time.Millisecond * 100)
				}
			}
			time.Sleep(time.Second)
			for i := range len(watches) {
				if gracefullyStop.Load() {
					return
				}
				if err := watches[i].Save(); err != nil {
					slog.Error("saving", "msg", err.Error())
					time.Sleep(time.Millisecond * 100)
				}
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
	usersRWMutex.RLock()
	defer usersRWMutex.RUnlock()
	low := 0
	high := len(users) - 1

	for low <= high {
		mid := (low + high) / 2
		if users[mid].CurrentUser.Id == id {
			if users[mid].ToDelete || users[mid].Deleted {
				return nil
			}
			return users[mid]
		} else if users[mid].CurrentUser.Id < id {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return nil
}

func FindUserCacheByLogin(email string) *UserCache {
	usersRWMutex.RLock()
	defer usersRWMutex.RUnlock()
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
	advsRWMutex.RLock()
	defer advsRWMutex.RUnlock()

	low := 0
	high := len(advs) - 1

	for low <= high {
		mid := (low + high) / 2
		if advs[mid].CurrentAdv.Id == id {
			if advs[mid].ToDelete || advs[mid].Deleted {
				return nil
			}
			return advs[mid]
		} else if advs[mid].CurrentAdv.Id < id {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return nil
}

func FindPhotoCacheById(id int64) *PhotoCache {
	photosRWMutex.RLock()
	defer photosRWMutex.RUnlock()
	low := 0
	high := len(photos) - 1

	for low <= high {
		mid := (low + high) / 2
		if photos[mid].Photo.Id == id {
			if photos[mid].ToDelete || photos[mid].Deleted {
				return nil
			}
			return photos[mid]
		} else if photos[mid].Photo.Id < id {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return nil
}

func FindWatchesCacheById(advId int64) *WatchesCache {
	watchesRWMutex.RLock()
	defer watchesRWMutex.RUnlock()
	low := 0
	high := len(watches) - 1

	for low <= high {
		mid := (low + high) / 2
		if watches[mid].Watches.AdvId == advId {
			if watches[mid].ToDelete || watches[mid].Deleted {
				return nil
			}
			return watches[mid]
		} else if watches[mid].Watches.AdvId < advId {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return nil
}

func FindAdvs(minDollarPrice int64, maxDollarPrice int64, minLongitude float64,
	maxLongitude float64, minLatitude float64, maxLatitude float64, countryCode string,
	location string, offset int, limit int, firstNew bool) ([]*dto.GetAdvResponseItem, int) {
	result := make([]*dto.GetAdvResponseItem, 0, limit)
	advsRWMutex.RLock()
	defer advsRWMutex.RUnlock()
	var i, step, count int
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
			count++
			if offset > 0 {
				offset--
				continue
			}
			if limit > 0 {
				limit--
			} else {
				continue
			}
			response := &dto.GetAdvResponseItem{
				Id:           adv.Id,
				UserEmail:    adv.User.Email,
				UserName:     adv.User.Name,
				Created:      time.UnixMicro(adv.Id / 1000),
				Updated:      adv.Updated,
				Approved:     adv.Approved,
				Lang:         adv.Lang,
				OriginLang:   adv.OriginLang,
				TranslatedBy: adv.TranslatedBy,
				Title:        adv.Title,
				Description:  adv.Description,
				Photos:       advs[i].GetPhotosFilenames(),
				Price:        adv.Price,
				Currency:     adv.Currency,
				DollarPrice:  adv.DollarPrice,
				Country:      adv.Country,
				City:         adv.City,
				Address:      adv.Address,
				Latitude:     adv.Latitude,
				Longitude:    adv.Longitude,
				Watches:      advs[i].Watches.Watches.Count,
				SeVisible:    adv.SeVisible,
			}
			result = append(result, response)
		}
	}
	return result, count
}

func FindUsersAdvs(userId int64, offset, limit int, firstNew bool) ([]*dto.GetAdvResponseItem, int) {
	advsRWMutex.RLock()
	defer advsRWMutex.RUnlock()
	length := len(advs)
	if offset > length {
		return []*dto.GetAdvResponseItem{}, 0 //todo ?
	}
	result := make([]*dto.GetAdvResponseItem, 0, limit)
	var i, step, count int
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
			count++
			if offset > 0 {
				offset--
				continue
			}
			if limit > 0 {
				limit--
			} else {
				continue
			}
			response := &dto.GetAdvResponseItem{
				Id:           adv.Id,
				UserEmail:    adv.User.Email,
				UserName:     adv.User.Name,
				Created:      time.UnixMicro(adv.Id / 1000),
				Updated:      adv.Updated,
				Approved:     adv.Approved,
				Lang:         adv.Lang,
				OriginLang:   adv.OriginLang,
				TranslatedBy: adv.TranslatedBy,
				Title:        adv.Title,
				Description:  adv.Description,
				Photos:       advs[i].GetPhotosFilenames(),
				Price:        adv.Price,
				Currency:     adv.Currency,
				DollarPrice:  adv.DollarPrice,
				Country:      adv.Country,
				City:         adv.City,
				Address:      adv.Address,
				Latitude:     adv.Latitude,
				Longitude:    adv.Longitude,
				Watches:      advs[i].Watches.Watches.Count,
				SeVisible:    adv.SeVisible,
			}
			result = append(result, response)
		}
	}
	return result, count
}

func CreateAdv(requestId int64, user *models.User, request *dto.CreateAdvRequest) {
	id := utils.GenerateId()
	newAdv := &models.Adv{
		Id:           id,
		UserId:       user.Id,
		User:         user,
		Updated:      time.Now(),
		Approved:     false,
		Lang:         request.OriginLang,
		OriginLang:   request.OriginLang,
		TranslatedBy: request.TranslatedBy,
		TranslatedTo: request.TranslatedTo,
		Title:        request.Title,
		Description:  request.Description,
		Price:        request.Price,
		Currency:     request.Currency,
		DollarPrice:  0, //todo
		Country:      request.Country,
		City:         request.City,
		Address:      request.Address,
		Latitude:     request.Latitude,
		Longitude:    request.Longitude,
		PaidAdv:      0,
		SeVisible:    true,
		UserComment:  request.UserComment,
		AdminComment: "",
	}
	advCache := &AdvCache{
		CurrentAdv: *newAdv,
		OldAdv:     models.Adv{},
		Watches: &WatchesCache{
			Watches: models.Watches{
				AdvId: id,
				Count: 0,
			},
			ToCreate: true,
			ToUpdate: false,
			ToDelete: false,
			Deleted:  false,
			mu:       sync.RWMutex{},
		},
		ToCreate: true,
	}

	advCache.mu.Lock() //todo нужно ли это
	defer advCache.mu.Unlock()

	advsRWMutex.Lock()
	advs = append(advs, advCache)
	advsRWMutex.Unlock()

	watchesRWMutex.Lock()
	watches = append(watches, advCache.Watches)
	watchesRWMutex.Unlock()

	toSave <- SaveTask{Cache: advCache, RequestId: requestId}
}

func UpdateAdv(requestId int64, adv *AdvCache, request *dto.UpdateAdvRequest) {
	adv.mu.Lock()
	defer adv.mu.Unlock()
	adv.CurrentAdv.OriginLang = request.OriginLang
	adv.CurrentAdv.TranslatedBy = request.TranslatedBy
	adv.CurrentAdv.TranslatedTo = request.TranslatedTo
	adv.CurrentAdv.Title = request.Title
	adv.CurrentAdv.Description = request.Description
	adv.CurrentAdv.Price = request.Price
	adv.CurrentAdv.Currency = request.Currency
	adv.CurrentAdv.Country = request.Country
	adv.CurrentAdv.City = request.City
	adv.CurrentAdv.Address = request.Address
	adv.CurrentAdv.Latitude = request.Latitude
	adv.CurrentAdv.Longitude = request.Longitude
	adv.CurrentAdv.UserComment = request.UserComment
	adv.ToUpdate = true
	toSave <- SaveTask{Cache: adv, RequestId: requestId}
}

func IncAdvWatches(watch *WatchesCache) {
	watch.mu.Lock()
	defer watch.mu.Unlock()
	watch.Watches.Count++
	watch.ToUpdate = true
	//toSave <- watch мы специально не отправляем в канал
}

func DeleteAdv(requestId int64, adv *AdvCache) {
	adv.mu.Lock()
	defer adv.mu.Unlock()
	if !adv.Deleted {
		adv.ToDelete = true
	}
	toSave <- SaveTask{Cache: adv, RequestId: requestId}
}

func CreateUser(requestId int64, request *dto.RegisterRequest) {
	passwordHash := utils.GeneratePasswordHash(request.Password)
	newUser := &models.User{
		Id:            utils.GenerateId(),
		Email:         request.Email,
		Name:          request.Name,
		PasswordHash:  passwordHash,
		SessionSecret: utils.GenerateSessionsSecret(passwordHash),
		InviteId:      request.InviteId,
		Balance:       0,
		Trusted:       false,
		Enabled:       true,
		Description:   "",
	}
	userCache := &UserCache{
		CurrentUser: *newUser,
		OldUser:     models.User{},
		ToCreate:    true,
	}
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
	users = append(users, userCache)
	toSave <- SaveTask{Cache: userCache, RequestId: requestId}
}

func UpdateUser(requestId int64, userCache *UserCache, request *dto.UpdateUserRequest) {
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
	userCache.CurrentUser.Name = request.Name
	userCache.CurrentUser.Description = request.Description
	userCache.ToUpdate = true
	toSave <- SaveTask{Cache: userCache, RequestId: requestId}
}

func UpdatePassword(requestId int64, userCache *UserCache, request *dto.UpdatePasswordRequest) {
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
	userCache.CurrentUser.PasswordHash = utils.GeneratePasswordHash(request.NewPassword)
	userCache.CurrentUser.SessionSecret = utils.GenerateSessionsSecret(userCache.CurrentUser.SessionSecret[:])
	userCache.ToUpdate = true
	toSave <- SaveTask{Cache: userCache, RequestId: requestId}
}

func UpdateSessionSecret(requestId int64, userCache *UserCache) {
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
	userCache.CurrentUser.SessionSecret = utils.GenerateSessionsSecret(userCache.CurrentUser.SessionSecret[:])
	userCache.ToUpdate = true
	toSave <- SaveTask{Cache: userCache, RequestId: requestId}
}

func DeleteUser(requestId int64, userCache *UserCache) {
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
	if !userCache.Deleted {
		userCache.ToDelete = true
	}
	toSave <- SaveTask{Cache: userCache, RequestId: requestId}
}

func CreatePhoto(requestId int64, adv *AdvCache, photo *models.Photo) {
	photoCache := &PhotoCache{
		Photo:    *photo,
		ToCreate: true,
	}
	photoCache.mu.Lock() //todo нужно ли тут вообще лочить
	toSave <- SaveTask{Cache: photoCache, RequestId: requestId}
	photoCache.mu.Unlock()

	photosRWMutex.Lock()
	photos = append(photos, photoCache)
	photosRWMutex.Unlock()

	adv.photoMu.Lock()
	adv.Photos = append(adv.Photos, photoCache)
	adv.photoMu.Unlock()

}

func DeletePhoto(requestId int64, adv *AdvCache, photoCache *PhotoCache) {
	photoCache.mu.Lock()
	if !photoCache.Deleted {
		photoCache.ToDelete = true
	}
	toSave <- SaveTask{Cache: photoCache, RequestId: requestId}
	photoCache.mu.Unlock()
}

func GetPhotosByAdvId(advId int64) []*PhotoCache {
	result := make([]*PhotoCache, 0, 15)
	photosRWMutex.RLock()
	for i := 0; i < len(photos); i++ {
		if photos[i].Photo.AdvId == advId && !photos[i].Deleted && !photos[i].ToDelete {
			result = append(result, photos[i])
		}
	}
	photosRWMutex.RUnlock()
	return result
}
