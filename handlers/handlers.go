package handlers

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"os"
	"realty/cache"
	"realty/config"
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
	render.Json(rd.RequestId, writer, http.StatusInternalServerError, &dto.Err{ErrMessage: "Internal error"})
}

func JsonOK(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	render.Json(rd.RequestId, writer, http.StatusOK, render.ResultOK)
	return false
}

func GetMetrics(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	m := dto.Metrics{
		InstanceStartTime:           metrics.GetInstanceStartTime(),
		FreeRAM:                     metrics.GetFreeRAM(),
		CPUTemp:                     metrics.GetCPUTemp(),
		CPUConsumption:              metrics.GetCPUConsumption(),
		UnSavedChangesQueueCount:    metrics.GetUnSavedChangesQueueCount(),
		DiskUsagePercent:            metrics.GetDiskUsagePercent(),
		RecoveredPanicsCount:        metrics.GetRecoveredPanicsCount(),
		MaxRAMConsumptions:          metrics.GetMaxRAMConsumptions(),
		MaxCPUConsumptions:          metrics.GetMaxCPUConsumptions(),
		MaxRPS:                      metrics.GetMaxRPS(),
		MaxUnSavedChangesQueueCount: metrics.GetMaxUnSavedChangesQueueCount(),
	}
	render.Json(rd.RequestId, writer, http.StatusOK, &m)
	return false
}

func LogoutMe(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	return false
}

func LogoutAll(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	cache.UpdateSessionSecret(rd.RequestId, rd.User)
	render.Json(rd.RequestId, writer, http.StatusOK, render.ResultOK)
	return false
}

func Registration(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	requestDto := &dto.RegisterRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	if err := validator.ValidateRegisterRequest(requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	cache.CreateUser(rd.RequestId, requestDto)
	render.Json(rd.RequestId, writer, http.StatusOK, render.ResultOK)
	return false
}

func UpdatePassword(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	requestDto := &dto.UpdatePasswordRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	if err := validator.ValidateUpdatePasswordRequest(requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	if !bytes.Equal(rd.User.CurrentUser.PasswordHash, utils.GeneratePasswordHash(requestDto.OldPassword)) {
		render.Json(rd.RequestId, writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверный пароль"})
		return false
	}
	cache.UpdatePassword(rd.RequestId, rd.User, requestDto)
	render.Json(rd.RequestId, writer, http.StatusOK, render.ResultOK)
	return false
}

func UpdateUser(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	requestDto := &dto.UpdateUserRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	if err := validator.ValidateUpdateUserRequest(requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	cache.UpdateUser(rd.RequestId, rd.User, requestDto)
	render.Json(rd.RequestId, writer, http.StatusOK, render.ResultOK)
	return false
}

func CreateAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	requestDto := &dto.CreateAdvRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	if err := validator.ValidateCreateAdvRequest(requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	cache.CreateAdv(rd.RequestId, &rd.User.CurrentUser, requestDto)
	render.Json(rd.RequestId, writer, http.StatusOK, render.ResultOK)
	return false
}

func GetAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	adv := rd.Adv.CurrentAdv
	if !adv.Approved {
		render.Json(rd.RequestId, writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление на проверке"})
		return false
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
	render.Json(rd.RequestId, writer, http.StatusOK, response)
	cache.IncAdvWatches(rd.Adv.Watches)
	return false
}

func GetAdvList(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	var (
		minDollarPrice int64
		maxDollarPrice int64 = math.MaxInt64
		offset         int
		limit          int = 20
	)
	requestDto := &dto.GetAdvListRequest{}
	if err := parsing_input.Parse(request, requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	if err := validator.ValidateGetAdvListRequest(requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	if !currency.IsValidCurrency(requestDto.Currency) {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: "неверный currency"})
		return false
	}
	if requestDto.MinPrice > 0 {
		minDollarPrice = currency.CalcDollarPrice(requestDto.Currency, requestDto.MinPrice)
	}
	if requestDto.MaxPrice > 0 {
		maxDollarPrice = currency.CalcDollarPrice(requestDto.Currency, requestDto.MaxPrice)
	}
	offset = (requestDto.Page - 1) * limit
	advs, count := cache.FindAdvs(
		minDollarPrice,
		maxDollarPrice,
		requestDto.MinLongitude,
		requestDto.MaxLongitude,
		requestDto.MinLatitude,
		requestDto.MaxLatitude,
		requestDto.CountryCode,
		requestDto.Location,
		offset,
		limit,
		requestDto.FirstNew)
	render.Json(rd.RequestId, writer, http.StatusOK, &dto.GetAdvListResponse{List: advs, Count: count})
	return false
}

func GetUsersAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	render.Json(rd.RequestId, writer, http.StatusOK, rd.Adv.CurrentAdv)
	return false
}

func GetUsersAdvList(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	var (
		offset   int
		limit    int = 20
		firstNew bool
	)
	requestDto := &dto.GetUserAdvListRequest{}
	if err := parsing_input.Parse(request, requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	if err := validator.ValidateGetUserAdvListRequest(requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	offset = (int(requestDto.Page) - 1) * limit
	advs, count := cache.FindUsersAdvs(rd.User.CurrentUser.Id, offset, limit, firstNew)
	render.Json(rd.RequestId, writer, http.StatusOK, &dto.GetAdvListResponse{List: advs, Count: count})
	return false
}

func UpdateAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	requestDto := &dto.UpdateAdvRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	if err := validator.ValidateUpdateAdvRequest(requestDto); err != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	cache.UpdateAdv(rd.RequestId, rd.Adv, requestDto)
	render.Json(rd.RequestId, writer, http.StatusOK, render.ResultOK)
	return false
}

func DeleteAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	cache.DeleteAdv(rd.RequestId, rd.Adv)
	render.Json(rd.RequestId, writer, http.StatusOK, render.ResultOK)
	return false
}

func AddAdvPhoto(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	// ParseMultipartForm parses a request body as multipart/form-data
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		render.Json(rd.RequestId, writer, http.StatusInternalServerError, &dto.Err{ErrMessage: err.Error()})
		return false
	}

	file, header, err := request.FormFile("file") // Retrieve the file from form data

	if err != nil {
		render.Json(rd.RequestId, writer, http.StatusInternalServerError, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	defer file.Close() // Close the file when we finish
	splits := strings.Split(header.Filename, ".")
	ext := splits[len(splits)-1]
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		render.Json(rd.RequestId, writer, http.StatusInternalServerError, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	photo := &models.Photo{
		AdvId: rd.Adv.CurrentAdv.Id,
		Id:    utils.GenerateId(),
	}

	switch ext {
	case ".jpg":
		photo.Ext = 1
	case ".png":
		photo.Ext = 2
	case ".gif":
		photo.Ext = 3
	default:
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: "не поддерживается тип изображения"})
		return false
	}

	f, err := os.Create(config.GetUploadedFilesPath() + strconv.FormatInt(photo.Id, 10) + ext)
	if err != nil {
		render.Json(rd.RequestId, writer, http.StatusInternalServerError, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		render.Json(rd.RequestId, writer, http.StatusInternalServerError, &dto.Err{ErrMessage: err.Error()})
		return false
	}
	cache.CreatePhoto(rd.RequestId, rd.Adv, photo)
	render.Json(rd.RequestId, writer, http.StatusOK, render.ResultOK)
	return false
}

func DeleteAdvPhoto(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) bool {
	photoIdStr := request.PathValue("photoId")
	photoId, errConv := strconv.ParseInt(photoIdStr, 10, 64)
	if errConv != nil {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: errConv.Error()})
		return false
	}
	if !validator.IsValidUnixNanoId(photoId) {
		render.Json(rd.RequestId, writer, http.StatusNotFound, &dto.Err{ErrMessage: "фото не найдено"})
		return false
	}
	photoCache := cache.FindPhotoCacheById(photoId)
	if photoCache == nil {
		render.Json(rd.RequestId, writer, http.StatusNotFound, &dto.Err{ErrMessage: "фото не найдено"})
		return false
	}
	if rd.Adv.CurrentAdv.Id != photoCache.Photo.AdvId {
		render.Json(rd.RequestId, writer, http.StatusBadRequest, &dto.Err{ErrMessage: "фото принадлежит другому объявлению"})
		return false
	}
	cache.DeletePhoto(rd.RequestId, rd.Adv, photoCache)
	render.Json(rd.RequestId, writer, http.StatusOK, render.ResultOK)
	return false
}
