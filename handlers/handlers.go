package handlers

import (
	"bytes"
	"math"
	"net/http"
	"realty/cache"
	"realty/currency"
	"realty/dto"
	"realty/metrics"
	"realty/middleware"
	"realty/models"
	"realty/parsing_input"
	"realty/render"
	"realty/utils"
	"realty/validator"
	"strconv"
	"strings"
	"time"
)

func TextError(recovered any, rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusInternalServerError)
	_, _ = writer.Write(utils.UnsafeStringToBytes("Internal error, requestId=" + strconv.FormatInt(rd.RequestId, 10)))
}

func JsonError(recovered any, rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	render.Json(writer, http.StatusInternalServerError, &dto.Err{ErrMessage: "Internal error", RequestId: rd.RequestId})
}

func JsonOK(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func GetMetrics(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	m := dto.Metrics{
		InstanceStartTime:           metrics.GetInstanceStartTime().Format("2006/01/02 15:04:05"),
		UnSavedChangesQueueCount:    metrics.GetUnSavedChangesQueueCount(),
		RecoveredPanicsCount:        metrics.GetRecoveredPanicsCount(),
		MaxUnSavedChangesQueueCount: metrics.GetMaxUnSavedChangesQueueCount(),
	}
	return render.Json(writer, http.StatusOK, &m)
}

func LogoutMe(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func LogoutAll(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	cache.UpdateSessionSecret(rd.RequestId, rd.User)
	return render.Json(writer, http.StatusOK, render.ResultOK)

}

func Registration(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	requestDto := &dto.RegisterRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateRegisterRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	cache.CreateUser(rd.RequestId, requestDto)
	return render.Json(writer, http.StatusOK, render.ResultOK)

}

func UpdatePassword(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	requestDto := &dto.UpdatePasswordRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateUpdatePasswordRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if !bytes.Equal(rd.User.CurrentUser.PasswordHash, utils.GeneratePasswordHash(requestDto.OldPassword)) {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверный пароль"})
	}
	cache.UpdatePassword(rd.RequestId, rd.User, requestDto)
	return render.Json(writer, http.StatusOK, render.ResultOK)

}

func UpdateUser(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	requestDto := &dto.UpdateUserRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateUpdateUserRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	cache.UpdateUser(rd.RequestId, rd.User, requestDto)
	return render.Json(writer, http.StatusOK, render.ResultOK)

}

func CreateAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	requestDto := &dto.CreateAdvRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateCreateAdvRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	advId := cache.CreateAdv(rd.RequestId, &rd.User.CurrentUser, requestDto)
	return render.Json(writer, http.StatusOK, &dto.CreateAdvResponse{RequestId: rd.RequestId, AdvId: advId})
}

func GetAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	adv := rd.Adv.CurrentAdv
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
		Photos:       rd.Adv.GetPhotosFilenames(),
		Price:        adv.Price,
		Currency:     adv.Currency,
		DollarPrice:  adv.DollarPrice,
		Country:      adv.Country,
		City:         adv.City,
		Address:      adv.Address,
		Latitude:     adv.Latitude,
		Longitude:    adv.Longitude,
		Watches:      rd.Adv.Watches.Watches.Count,
		SeVisible:    adv.SeVisible,
		UserComment:  adv.UserComment,
	}
	cache.IncAdvWatches(rd.Adv.Watches)
	return render.Json(writer, http.StatusOK, response)
}

func GetAdvList(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
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

func GetUsersAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	return render.Json(writer, http.StatusOK, rd.Adv.CurrentAdv)
}

func GetUsersAdvList(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
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
	advs, count := cache.FindUsersAdvs(rd.User.CurrentUser.Id, offset, limit, firstNew)
	return render.Json(writer, http.StatusOK, &dto.GetAdvListResponse{List: advs, Count: count})
}

func UpdateAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	requestDto := &dto.UpdateAdvRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateUpdateAdvRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	cache.UpdateAdv(rd.RequestId, rd.Adv, requestDto)
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func DeleteAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	cache.DeleteAdv(rd.RequestId, rd.Adv)
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func AddAdvPhoto(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
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
		AdvId: rd.Adv.CurrentAdv.Id,
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
	cache.CreatePhoto(rd.RequestId, rd.Adv, photo)
	return render.Json(writer, http.StatusOK, render.ResultOK)
}

func DeleteAdvPhoto(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
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
	if rd.Adv.CurrentAdv.Id != photoCache.Photo.AdvId {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "фото принадлежит другому объявлению"})
	}
	cache.DeletePhoto(rd.RequestId, rd.Adv, photoCache)
	return render.Json(writer, http.StatusOK, render.ResultOK)
}
