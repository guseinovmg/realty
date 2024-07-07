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
	"strconv"
)

func TextError(recovered any, rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(500)
	_, _ = writer.Write(utils.UnsafeStringToBytes("Internal error"))
}

func JsonError(recovered any, rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	_ = render.Json(writer, http.StatusInternalServerError, &dto.Err{ErrMessage: "Internal error"})
}

func Login(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
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
		rd.Stop()
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "пользователь не найден"})
		return
	}
	if !bytes.Equal(utils.GeneratePasswordHash(requestDto.Password), userCache.CurrentUser.PasswordHash) {
		rd.Stop()
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверный пароль"})
		return
	}

	rd.User = userCache

	_ = render.JsonOK(writer, http.StatusOK)

}

func LogoutMe(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func LogoutAll(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	cache.UpdateSessionSecret(rd.User)
	_ = render.JsonOK(writer, http.StatusOK)
}

func Registration(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
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
}

func UpdatePassword(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
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
}

func UpdateUser(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
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
}

func CreateAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
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
}

func GetAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	advIdStr := request.PathValue("id")
	advId, errConv := strconv.ParseInt(advIdStr, 10, 64)
	if errConv != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: errConv.Error()})
		return
	}
	if !validator.IsValidUnixNanoId(advId) {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return
	}
	advCache := cache.FindAdvCacheById(advId)
	if advCache == nil {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return
	}
	adv := &advCache.CurrentAdv
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
		UserComment:  adv.UserComment,
	}
	_ = render.Json(writer, http.StatusOK, response)
	cache.IncAdvWatches(advCache)
}

func GetAdvList(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	var (
		minDollarPrice int64
		maxDollarPrice int64 = math.MaxInt64
		offset         int
		limit          int = 20
	)
	requestDto := &dto.GetAdvListRequest{}
	if err := parsing_input.ParseQueryToGetAdvListRequest(request.URL.Query(), requestDto); err != nil {
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

}

func GetUsersAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	advIdStr := request.PathValue("id")
	advId, errConv := strconv.ParseInt(advIdStr, 10, 64)
	if errConv != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: errConv.Error()})
		return
	}
	adv := cache.FindAdvById(advId)
	if !validator.IsValidUnixNanoId(advId) {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return
	}
	if adv == nil {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return
	}
	if adv.UserId != rd.User.CurrentUser.Id {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не принадлежит текущему пользвателю"})
		return
	}
	_ = render.Json(writer, http.StatusOK, adv)
}

func GetUsersAdvList(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	var (
		offset   int
		limit    int = 20
		firstNew bool
	)
	pageStr := request.URL.Query().Get("page")
	if pageStr != "" {
		page, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "page должен быть целым числом"})
			return
		}
		if page < 1 {
			_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "page должен быть целым числом больше 1"})
			return
		}
		offset = (int(page) - 1) * limit
	}
	firstNewStr := request.URL.Query().Get("firstNew")
	if firstNewStr != "" {
		first, err := strconv.ParseBool(firstNewStr)
		if err != nil {
			_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "firstNew должен быть иметь значение 0 или 1"})
			return
		}
		firstNew = first
	}
	advs := cache.FindUsersAdvs(rd.User.CurrentUser.Id, offset, limit, firstNew)
	_ = render.Json(writer, http.StatusOK, advs)
}

func UpdateAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	advIdStr := request.PathValue("id")
	advId, errConv := strconv.ParseInt(advIdStr, 10, 64)
	if errConv != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: errConv.Error()})
		return
	}
	if !validator.IsValidUnixNanoId(advId) {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return
	}
	requestDto := &dto.UpdateAdvRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if err := validator.ValidateUpdateAdvRequest(requestDto); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	adv := cache.FindAdvCacheById(advId)
	if adv == nil {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return
	}
	if adv.CurrentAdv.UserId != rd.User.CurrentUser.Id {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не принадлежит текущему пользвателю"})
		return
	}
	cache.UpdateAdv(adv, requestDto)
	_ = render.JsonOK(writer, http.StatusOK)
}

func DeleteAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	advIdStr := request.PathValue("id")
	advId, errConv := strconv.ParseInt(advIdStr, 10, 64)
	if errConv != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: errConv.Error()})
		return
	}
	if !validator.IsValidUnixNanoId(advId) {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return
	}
	adv := cache.FindAdvCacheById(advId)
	if adv == nil {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return
	}
	if adv.CurrentAdv.UserId != rd.User.CurrentUser.Id {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не принадлежит текущему пользвателю"})
		return
	}
	cache.DeleteAdv(adv)
	_ = render.JsonOK(writer, http.StatusOK)
}
