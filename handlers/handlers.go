package handlers

import (
	"bytes"
	"math"
	"net/http"
	"realty/cache"
	"realty/currency"
	"realty/dto"
	"realty/middleware"
	"realty/parsing_input"
	"realty/render"
	"realty/utils"
	"realty/validator"
	"time"
)

func TextError(recovered any, rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(500)
	_, _ = writer.Write(utils.UnsafeStringToBytes("Internal error"))
}

func JsonError(recovered any, rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	_ = render.Json(writer, http.StatusInternalServerError, &dto.Err{ErrMessage: "Internal error"})
}

func Login(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	requestDto := &dto.LoginRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if err := validator.ValidateLoginRequest(requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	userCache := cache.FindUserCacheByLogin(requestDto.Email)
	if userCache == nil {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "пользователь не найден"})
		return
	}
	if !bytes.Equal(utils.GeneratePasswordHash(requestDto.Password), userCache.CurrentUser.PasswordHash) {
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверный пароль"})
		return
	}
	rd.User = userCache
	middleware.SetAuthCookie(rd, writer, request) //todo надо подумать
	_ = render.JsonOK(writer, http.StatusOK)
	return
}

func LogoutMe(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	return false
}

func LogoutAll(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	cache.UpdateSessionSecret(rd.User)
	_ = render.JsonOK(writer, http.StatusOK)
	return false
}

func Registration(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	requestDto := &dto.RegisterRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if err := validator.ValidateRegisterRequest(requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	cache.CreateUser(requestDto)
	_ = render.JsonOK(writer, http.StatusOK)
	return
}

func UpdatePassword(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	requestDto := &dto.UpdatePasswordRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if err := validator.ValidateUpdatePasswordRequest(requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if !bytes.Equal(rd.User.CurrentUser.PasswordHash, utils.GeneratePasswordHash(requestDto.OldPassword)) {
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверный пароль"})
		return
	}
	cache.UpdatePassword(rd.User, requestDto)
	_ = render.JsonOK(writer, http.StatusOK)
	return
}

func UpdateUser(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	requestDto := &dto.UpdateUserRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if err := validator.ValidateUpdateUserRequest(requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	cache.UpdateUser(rd.User, requestDto)
	_ = render.JsonOK(writer, http.StatusOK)
	return
}

func CreateAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	requestDto := &dto.CreateAdvRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if err := validator.ValidateCreateAdvRequest(requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	cache.CreateAdv(&rd.User.CurrentUser, requestDto)
	_ = render.JsonOK(writer, http.StatusOK)
	return
}

func GetAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	adv := rd.Adv.CurrentAdv
	if !adv.Approved {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление на проверке"})
		return
	}
	if !adv.SeVisible {
		//todo
	}
	response := &dto.GetAdvResponse{
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
		UserComment:  adv.UserComment,
	}
	_ = render.Json(writer, http.StatusOK, response)
	cache.IncAdvWatches(rd.Adv)
	return
}

func GetAdvList(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	var (
		minDollarPrice int64
		maxDollarPrice int64 = math.MaxInt64
		offset         int
		limit          int = 20
	)
	requestDto := &dto.GetAdvListRequest{}
	if err := parsing_input.Parse(request, requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if err := validator.ValidateGetAdvListRequest(requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if !currency.IsValidCurrency(requestDto.Currency) {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "неверный currency"})
		return
	}
	if requestDto.MinPrice > 0 {
		minDollarPrice = currency.CalcDollarPrice(requestDto.Currency, requestDto.MinPrice)
	}
	if requestDto.MaxPrice > 0 {
		maxDollarPrice = currency.CalcDollarPrice(requestDto.Currency, requestDto.MaxPrice)
	}
	offset = (requestDto.Page - 1) * limit
	advs := cache.FindAdvs(
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
	_ = render.Json(writer, http.StatusOK, advs)
	return
}

func GetUsersAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	_ = render.Json(writer, http.StatusOK, rd.Adv.CurrentAdv)
	return
}

func GetUsersAdvList(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	var (
		offset   int
		limit    int = 20
		firstNew bool
	)
	requestDto := &dto.GetUserAdvListRequest{}
	if err := parsing_input.Parse(request, requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if err := validator.ValidateGetUserAdvListRequest(requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	offset = (int(requestDto.Page) - 1) * limit
	advs := cache.FindUsersAdvs(rd.User.CurrentUser.Id, offset, limit, firstNew)
	_ = render.Json(writer, http.StatusOK, advs)
	return
}

func UpdateAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	requestDto := &dto.UpdateAdvRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if err := validator.ValidateUpdateAdvRequest(requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	cache.UpdateAdv(rd.Adv, requestDto)
	_ = render.JsonOK(writer, http.StatusOK)
	return
}

func DeleteAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	cache.DeleteAdv(rd.Adv)
	_ = render.JsonOK(writer, http.StatusOK)
	return
}

func AddAdvPhoto(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	//todo
	return
}

func DeleteAdvPhoto(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) (next bool) {
	//todo
	_ = render.JsonOK(writer, http.StatusOK)
	return
}
