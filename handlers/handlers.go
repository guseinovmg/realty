package handlers

import (
	"bytes"
	"net/http"
	"realty/cache"
	"realty/dto"
	"realty/middleware"
	"realty/parsing_input"
	"realty/render"
	"realty/utils"
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
	err := parsing_input.Parse(request, login)
	if err != nil {
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
	err := parsing_input.Parse(request, requestDto)
	if err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	cache.CreateUser(requestDto)
	_ = render.JsonOK(writer, http.StatusOK)
}

func UpdatePassword(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	requestDto := &dto.UpdatePasswordRequest{}
	err := parsing_input.Parse(request, requestDto)
	if err != nil {
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
	err := parsing_input.Parse(request, requestDto)
	if err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	cache.UpdateUser(rd.User, requestDto)
	_ = render.JsonOK(writer, http.StatusOK)
}

func CreateAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	requestDto := &dto.CreateAdvRequest{}
	err := parsing_input.Parse(request, requestDto)
	if err != nil {
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
	adv := cache.FindAdvById(advId)
	if adv == nil {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return
	}
	if !adv.Approved {
		_ = render.Json(writer, http.StatusLocked, &dto.Err{ErrMessage: "объявление на модерации"})
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
	}
	_ = render.Json(writer, http.StatusOK, response)

}

func GetAdvList(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func GetUsersAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	advIdStr := request.PathValue("id")
	advId, errConv := strconv.ParseInt(advIdStr, 10, 64)
	if errConv != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: errConv.Error()})
		return
	}
	adv := cache.FindAdvById(advId)
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

}

func UpdateAdv(rd *middleware.RequestData, writer http.ResponseWriter, request *http.Request) {
	requestDto := &dto.UpdateAdvRequest{}
	err := parsing_input.Parse(request, requestDto)
	if err != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
		return
	}
	advIdStr := request.PathValue("id")
	advId, errConv := strconv.ParseInt(advIdStr, 10, 64)
	if errConv != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: errConv.Error()})
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
