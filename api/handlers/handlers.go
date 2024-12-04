package handlers

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"realty/application"
	"realty/cache"
	"realty/chain"
	"realty/config"
	"realty/currency"
	"realty/dto"
	"realty/models"
	"realty/moderation"
	"realty/parsing_input"
	"realty/render"
	"realty/utils"
	"realty/validator"
	"strconv"
	"strings"
	"time"
)

func TextError(recovered any, rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusInternalServerError)
	_, _ = writer.Write(utils.UnsafeStringToBytes("Internal error, requestId=" + strconv.FormatInt(rc.RequestId, 10)))
}

func JsonError(recovered any, rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) {
	render.Json(writer, http.StatusInternalServerError, &dto.Err{ErrMessage: "Internal error", RequestId: rc.RequestId})
}

func JsonOK(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func GetMetrics(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	count := cache.GetToSaveCount()
	if count > application.GetMaxUnSavedChangesQueueCount() {
		application.SetMaxUnSavedChangesQueueCount(count)
	}
	//todo надо еще добавить метрики из пакетов runtime и metrics
	m := dto.Metrics{
		InstanceStartTime:           application.GetInstanceStartTime().Format("2006/01/02 15:04:05"),
		UnSavedChangesQueueCount:    count,
		DbErrorCount:                application.GetDbErrorsCount(),
		RecoveredPanicsCount:        application.GetRecoveredPanicsCount(),
		MaxUnSavedChangesQueueCount: application.GetMaxUnSavedChangesQueueCount(),
	}
	return render.Json(writer, http.StatusOK, &m)
}

func GenerateId(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	return render.Json(writer, http.StatusOK, &dto.GenerateIdResponse{Id: utils.GenerateId()})
}

func LogoutMe(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	http.SetCookie(writer, &http.Cookie{
		SameSite: http.SameSiteStrictMode,
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Domain:   config.GetDomain(),
		MaxAge:   1,
		Secure:   true, // only sent over HTTPS
		HttpOnly: true, // not accessible via JavaScript
	})
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func LogoutAll(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	cache.UpdateSessionSecret(rc.RequestId, rc.User)
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func Registration(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	requestDto := &dto.RegisterRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateRegisterRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	cache.CreateUser(rc.RequestId, requestDto)
	return render.Json(writer, http.StatusOK, render.ResultOK)

}

func UpdatePassword(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	requestDto := &dto.UpdatePasswordRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateUpdatePasswordRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if !bytes.Equal(rc.User.CurrentUser.PasswordHash, utils.GeneratePasswordHash(requestDto.OldPassword)) {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверный пароль"})
	}
	cache.UpdatePassword(rc.RequestId, rc.User, requestDto)
	return render.Json(writer, http.StatusOK, render.ResultOK)

}

func UpdateUser(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	requestDto := &dto.UpdateUserRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateUpdateUserRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	cache.UpdateUser(rc.RequestId, rc.User, requestDto)
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func CreateAdv(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	requestDto := &dto.CreateAdvRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateCreateAdvRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if badWords := moderation.SearchBadWord(requestDto.Description); len(badWords) != 0 {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: fmt.Sprintf("в описании присутствуют запрещенные слова: %v", badWords)})
	}
	advId := cache.CreateAdv(rc.RequestId, &rc.User.CurrentUser, requestDto)
	return render.Json(writer, http.StatusOK, &dto.CreateAdvResponse{RequestId: rc.RequestId, AdvId: advId})
}

func GetAdv(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	adv := rc.Adv.CurrentAdv
	if !adv.Approved {
		return render.Json(writer, http.StatusLocked, &dto.Err{ErrMessage: "объявление на проверке"})
	}
	if !adv.SeVisible {
		//todo
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
		Photos:       rc.Adv.GetPhotosFilenames(),
		Price:        adv.Price,
		Currency:     adv.Currency,
		DollarPrice:  adv.DollarPrice,
		Country:      adv.Country,
		City:         adv.City,
		Address:      adv.Address,
		Latitude:     adv.Latitude,
		Longitude:    adv.Longitude,
		Watches:      rc.Adv.Watches.Watches.Count,
		SeVisible:    adv.SeVisible,
		UserComment:  adv.UserComment,
	}
	cache.IncAdvWatches(rc.Adv.Watches)
	return render.Json(writer, http.StatusOK, response)
}

func GetAdvList(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	var (
		minDollarPrice int64
		maxDollarPrice int64   = math.MaxInt64
		minLongitude   float64 = -180
		maxLongitude   float64 = 180
		minLatitude    float64 = -180
		maxLatitude    float64 = 180
		offset         int
		limit          int = 20
	)
	requestDto := &dto.GetAdvListRequest{}
	if err := parsing_input.Parse(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateGetAdvListRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if !currency.IsValidCurrency(requestDto.Currency) {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "неверный currency"})
	}
	if requestDto.MinPrice > 0 {
		minDollarPrice = currency.CalcDollarPrice(requestDto.Currency, requestDto.MinPrice)
	}
	if requestDto.MaxPrice > 0 {
		maxDollarPrice = currency.CalcDollarPrice(requestDto.Currency, requestDto.MaxPrice)
	}
	if requestDto.MinLongitude != 0 {
		minLongitude = requestDto.MinLongitude
	}
	if requestDto.MaxLongitude != 0 {
		maxLongitude = requestDto.MaxLongitude
	}
	if requestDto.MinLatitude != 0 {
		minLatitude = requestDto.MinLatitude
	}
	if requestDto.MaxLatitude != 0 {
		maxLatitude = requestDto.MaxLatitude
	}
	offset = (requestDto.Page - 1) * limit
	advs, count := cache.FindAdvs(
		minDollarPrice,
		maxDollarPrice,
		minLongitude,
		maxLongitude,
		minLatitude,
		maxLatitude,
		requestDto.CountryCode,
		requestDto.Location,
		offset,
		limit,
		requestDto.FirstNew)
	return render.Json(writer, http.StatusOK, &dto.GetAdvListResponse{List: advs, Count: count})
}

func GetUsersAdv(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	return render.Json(writer, http.StatusOK, rc.Adv.CurrentAdv)
}

func GetUsersAdvList(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	var (
		offset   int
		limit    int = 20
		firstNew bool
	)
	requestDto := &dto.GetUserAdvListRequest{}
	if err := parsing_input.Parse(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateGetUserAdvListRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	offset = (int(requestDto.Page) - 1) * limit
	advs, count := cache.FindUsersAdvs(rc.User.CurrentUser.Id, offset, limit, firstNew)
	return render.Json(writer, http.StatusOK, &dto.GetAdvListResponse{List: advs, Count: count})
}

func UpdateAdv(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	requestDto := &dto.UpdateAdvRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateUpdateAdvRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if badWords := moderation.SearchBadWord(requestDto.Description); len(badWords) != 0 {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: fmt.Sprintf("в описании присутствуют запрещенные слова: %v", badWords)})
	}
	cache.UpdateAdv(rc.RequestId, rc.Adv, requestDto)
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func DeleteAdv(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	cache.DeleteAdv(rc.RequestId, rc.Adv)
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func AddAdvPhoto(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	requestDto := &dto.AddPhotoRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateAddPhotoRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	splits := strings.Split(requestDto.Filename, ".")
	ext := splits[1]
	id, err := strconv.ParseInt(splits[0], 10, 64)
	if err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if !validator.IsValidUnixNanoId(id) {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "неверный id файла"})
	}
	photo := &models.Photo{
		AdvId: rc.Adv.CurrentAdv.Id,
		Id:    id,
	}
	switch ext {
	case "jpg":
		photo.Ext = 1
	case "png":
		photo.Ext = 2
	case "gif":
		photo.Ext = 3
	default:
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "не поддерживается тип изображения"})
	}
	cache.CreatePhoto(rc.RequestId, rc.Adv, photo)
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func DeleteAdvPhoto(rc *chain.RequestContext, writer http.ResponseWriter, request *http.Request) chain.Result {
	photoIdStr := request.PathValue("photoId")
	photoId, errConv := strconv.ParseInt(photoIdStr, 10, 64)
	if errConv != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: errConv.Error()})
	}
	if !validator.IsValidUnixNanoId(photoId) {
		return render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "фото не найдено"})
	}
	photoCache := cache.FindPhotoCacheById(photoId)
	if photoCache == nil {
		return render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "фото не найдено"})
	}
	if rc.Adv.CurrentAdv.Id != photoCache.Photo.AdvId {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "фото принадлежит другому объявлению"})
	}
	cache.DeletePhoto(rc.RequestId, rc.Adv, photoCache)
	return render.Json(writer, http.StatusOK, render.ResultOK)
}
