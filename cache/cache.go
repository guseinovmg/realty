package cache

import (
	"realty/db"
	"realty/dto"
	"realty/models"
	"realty/utils"
	"slices"
	"strconv"
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
	photoMu    sync.RWMutex
}

func (adv *AdvCache) Save() error {
	adv.mu.Lock()
	defer adv.mu.Unlock()
	if adv.Deleted {
		return nil
	}
	if adv.ToDelete {
		err := db.DeleteAdv(adv.CurrentAdv.Id)
		if err != nil {
			return err
		}
		adv.Deleted = true
		adv.ToDelete = false
		adv.ToCreate = false
		adv.ToUpdate = false
	}
	if adv.ToCreate {
		err := db.CreateAdv(&adv.CurrentAdv)
		if err != nil {
			return err
		}
		adv.OldAdv = adv.CurrentAdv
		adv.ToCreate = false
		adv.ToUpdate = false
	}
	if adv.ToUpdate {
		err := db.UpdateAdvChanges(&adv.OldAdv, &adv.CurrentAdv)
		if err != nil {
			return err
		}
		adv.OldAdv = adv.CurrentAdv
		adv.ToUpdate = false
	}
	return nil
}

func (adv *AdvCache) GetPhotosFilenames() []string {
	result := make([]string, 0, len(adv.CurrentAdv.Photos))
	adv.photoMu.RLock()
	defer adv.photoMu.RUnlock()
	for _, v := range adv.CurrentAdv.Photos { //todo если бы тут был массив PhotoCache можно было бы сразу удаленные убрать
		ext := ""
		switch v.Ext {
		case 1:
			ext = ".jpg"
		case 2:
			ext = ".png"
		case 3:
			ext = ".gif"
		}
		name := strconv.FormatInt(v.Id, 10) + ext
		result = append(result, name)
	}
	return result
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
	if user.Deleted {
		return nil
	}
	if user.ToDelete {
		err := db.DeleteAdv(user.CurrentUser.Id)
		if err != nil {
			return err
		}
		user.Deleted = true
		user.ToDelete = false
		user.ToCreate = false
		user.ToUpdate = false
	}
	if user.ToCreate {
		err := db.CreateUser(&user.CurrentUser)
		if err != nil {
			return err
		}
		user.OldUser = user.CurrentUser
		user.ToCreate = false
		user.ToUpdate = false
	}
	if user.ToUpdate {
		err := db.UpdateUserChanges(&user.OldUser, &user.CurrentUser)
		if err != nil {
			return err
		}
		user.OldUser = user.CurrentUser
		user.ToUpdate = false
	}
	return nil
}

type PhotoCache struct {
	Photo    models.Photo
	ToCreate bool
	ToDelete bool
	Deleted  bool
	mu       sync.RWMutex
}

func (photo *PhotoCache) Save() error {
	photo.mu.Lock()
	defer photo.mu.Unlock()
	if photo.Deleted {
		return nil
	}
	if photo.ToDelete {
		err := db.DeletePhoto(photo.Photo.Id)
		if err != nil {
			return err
		}
		photo.Deleted = true
		photo.ToDelete = false
		photo.ToCreate = false
	}
	if photo.ToCreate {
		err := db.CreatePhoto(photo.Photo)
		if err != nil {
			return err
		}
		photo.ToCreate = false
	}
	return nil
}

type WatchesCache struct {
	Watches  models.Watches
	ToCreate bool
	ToUpdate bool
	ToDelete bool
	Deleted  bool
	mu       sync.RWMutex
}

func (watch *WatchesCache) Save() error {
	watch.mu.Lock()
	defer watch.mu.Unlock()
	if watch.Deleted {
		return nil
	}
	if watch.ToDelete {
		err := db.DeleteWatches(watch.Watches.AdvId)
		if err != nil {
			return err
		}
		watch.Deleted = true
		watch.ToDelete = false
		watch.ToCreate = false
	}
	if watch.ToCreate {
		err := db.CreateWatches(watch.Watches)
		if err != nil {
			return err
		}
		watch.ToCreate = false
	}
	if watch.ToUpdate {
		err := db.UpdateWatches(&watch.Watches)
		if err != nil {
			return err
		}
		watch.ToUpdate = false
	}
	return nil
}

var users []*UserCache
var advs []*AdvCache
var photos []*PhotoCache
var watches []*WatchesCache
var toSave chan SaveCache
var idGenerationMutex sync.Mutex

func GenerateId() int64 {
	idGenerationMutex.Lock()
	defer idGenerationMutex.Unlock()
	time.Sleep(time.Microsecond)
	return time.Now().UnixNano()
}

func Initialize() {
	users_, advs_, photos_, watches_, err := db.ReadDb()
	if err != nil {
		panic(err)
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
			CurrentUser: *user, //todo надо подумать, может непосредственно ссылку оптимальнее использовать
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
		adv.Photos = GetPhotosByAdvId(adv.Id)
		w := FindWatchesCacheById(adv.Id)
		adv.Watches = &w.Watches
		advs[i] = &AdvCache{
			CurrentAdv: *adv,
			OldAdv:     *adv,
			ToUpdate:   false,
			ToDelete:   false,
			Deleted:    false,
			mu:         sync.RWMutex{},
		}
	}

	//todo надо просмотры и фото в adv добавить
	toSave = make(chan SaveCache, 100)
	go func() {
		for saveCache := range toSave {
			for range 3 {
				err := saveCache.Save()
				if err == nil {
					break
				}
				time.Sleep(time.Second)
			}
		}
	}()
	go func() {
		for {
			time.Sleep(time.Minute)
			for i := range len(advs) {
				err := advs[i].Save()
				if err != nil {
					time.Sleep(time.Millisecond * 100)
				}
			}
			for i := range len(users) {
				err := users[i].Save()
				if err != nil {
					time.Sleep(time.Millisecond * 100)
				}
			}
			for i := range len(photos) {
				err := photos[i].Save()
				if err != nil {
					time.Sleep(time.Millisecond * 100)
				}
			}
			for i := range len(watches) {
				err := watches[i].Save()
				if err != nil {
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
				Watches:      adv.Watches.Watches,
				SeVisible:    adv.SeVisible,
			}
			result = append(result, response)
		}
	}
	return result, count
}

func FindUsersAdvs(userId int64, offset, limit int, firstNew bool) ([]*dto.GetAdvResponseItem, int) {
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
				Watches:      adv.Watches.Watches,
				SeVisible:    adv.SeVisible,
			}
			result = append(result, response)
		}
	}
	return result, count
}

func CreateAdv(user *models.User, request *dto.CreateAdvRequest) {
	id := GenerateId()
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
		Watches: &models.Watches{
			AdvId:   id,
			Watches: 0,
		},
		PaidAdv:      0,
		SeVisible:    true,
		UserComment:  request.UserComment,
		AdminComment: "",
	}
	advCache := &AdvCache{
		CurrentAdv: *newAdv,
		OldAdv:     models.Adv{},
		ToCreate:   true,
	}
	advCache.mu.Lock()
	defer advCache.mu.Unlock()
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

func IncAdvWatches(adv *AdvCache) {
	adv.mu.Lock()
	defer adv.mu.Unlock()
	adv.CurrentAdv.Watches.Watches++
	adv.ToUpdate = true
	//toSave <- adv мы специально не отправляем в канал
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
		Id:            GenerateId(),
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
		ToCreate:    true,
	}
	userCache.mu.Lock()
	defer userCache.mu.Unlock()
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

func CreatePhoto(adv *AdvCache, photo *models.Photo) {
	photoCache := &PhotoCache{
		Photo:    *photo,
		ToCreate: true,
	}
	photoCache.mu.Lock() //todo нужно ли тут вообще лочить
	toSave <- photoCache
	defer photoCache.mu.Unlock()

	photos = append(photos, photoCache) //todo тут видимо отдельный лок нужен для слайсов
	adv.photoMu.Lock()
	adv.CurrentAdv.Photos = append(adv.CurrentAdv.Photos, &photoCache.Photo)
	adv.photoMu.Unlock()

}

func DeletePhoto(adv *AdvCache, photoCache *PhotoCache) {
	photoCache.mu.Lock()
	if !photoCache.Deleted {
		photoCache.ToDelete = true
	}
	toSave <- photoCache
	photoCache.mu.Unlock()

	adv.photoMu.RLock()
	ind := slices.Index(adv.CurrentAdv.Photos, &photoCache.Photo)
	adv.photoMu.RUnlock()
	if ind != -1 {
		adv.photoMu.Lock()
		adv.CurrentAdv.Photos = slices.Delete(adv.CurrentAdv.Photos, ind, ind)
		adv.photoMu.Unlock()
	}
}

func GetPhotosByAdvId(advId int64) []*models.Photo {
	result := make([]*models.Photo, 0, 15)
	for i := 0; i < len(photos); i++ {
		if photos[i].Photo.AdvId == advId && !photos[i].Deleted && !photos[i].ToDelete {
			result = append(result, &photos[i].Photo)
		}
	}
	return result
}
