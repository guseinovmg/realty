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
	login := &dto.LoginRequest{}
	if err := parsing_input.ParseRawJson(request, login); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	if err := validator.ValidateLoginRequest(login); err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	userCache := cache.FindUserCacheByLogin(login.Email)
	if userCache == nil {
		rd.Stop()
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "пользователь не найден"})
		return
	}
	if !bytes.Equal(utils.GeneratePasswordHash(login.Password), userCache.CurrentUser.PasswordHash) {
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
	if !validator.IsValidUnixMicroId(advId) {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return
	}
	adv := cache.FindAdvById(advId)
	if adv == nil {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return
	}
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

}

func GetAdvList(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	var (
		minDollarPrice int64
		maxDollarPrice int64   = math.MaxInt64
		minLongitude   float64 = -math.MaxFloat64
		maxLongitude   float64 = math.MaxFloat64
		minLatitude    float64 = -math.MaxFloat64
		maxLatitude    float64 = math.MaxFloat64
		countryCode    string
		location       string
		offset         int
		limit          int  = 20
		firstNew       bool = true
	)
	currencyStr := request.URL.Query().Get("currency")
	if !currency.IsValidCurrency(currencyStr) {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "неверный currency"})
		return
	}

	minPriceStr := request.URL.Query().Get("minPrice")
	if minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "minPrice должен быть числом"})
			return
		}
		minDollarPrice, err = currency.CalcDollarPrice(currencyStr, minPrice)
	}

	maxPriceStr := request.URL.Query().Get("maxPrice")
	if maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "maxPrice должен быть числом"})
			return
		}
		maxDollarPrice, err = currency.CalcDollarPrice(currencyStr, maxPrice)
	}

	minLongitudeStr := request.URL.Query().Get("minLongitude")
	if minLongitudeStr != "" {
		minLong, err := strconv.ParseFloat(minLongitudeStr, 64)
		if err != nil {
			_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "minLongitude должен быть числом"})
			return
		}
		minLongitude = minLong
	}

	maxLongitudeStr := request.URL.Query().Get("maxLongitude")
	if maxLongitudeStr != "" {
		maxLong, err := strconv.ParseFloat(maxLongitudeStr, 64)
		if err != nil {
			_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "maxLongitude должен быть числом"})
			return
		}
		maxLongitude = maxLong
	}

	minLatitudeStr := request.URL.Query().Get("minLatitude")
	if minLatitudeStr != "" {
		minLong, err := strconv.ParseFloat(minLatitudeStr, 64)
		if err != nil {
			_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "minLatitude должен быть числом"})
			return
		}
		minLatitude = minLong
	}

	maxLatitudeStr := request.URL.Query().Get("maxLatitude")
	if maxLatitudeStr != "" {
		maxLong, err := strconv.ParseFloat(maxLatitudeStr, 64)
		if err != nil {
			_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "maxLatitude должен быть числом"})
			return
		}
		maxLatitude = maxLong
	}

	countryCode = request.URL.Query().Get("countryCode")
	location = request.URL.Query().Get("location")
	firstNewStr := request.URL.Query().Get("firstNew")
	if firstNewStr != "" {
		first, err := strconv.ParseBool(firstNewStr)
		if err != nil {
			_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "firstNew должен быть иметь значение 0 или 1"})
			return
		}
		firstNew = first
	}

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
	advs := cache.FindAdvs(minDollarPrice,
		maxDollarPrice,
		minLongitude,
		maxLongitude,
		minLatitude,
		maxLatitude,
		countryCode,
		location,
		offset,
		limit,
		firstNew)
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
	if !validator.IsValidUnixMicroId(advId) {
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
	if !validator.IsValidUnixMicroId(advId) {
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
	if !validator.IsValidUnixMicroId(advId) {
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
