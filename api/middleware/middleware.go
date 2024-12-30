package middleware

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"realty/application"
	"realty/auth_token"
	"realty/cache"
	"realty/chain"
	"realty/config"
	"realty/dto"
	"realty/parsing_input"
	"realty/render"
	"realty/utils"
	"realty/validator"
	"strconv"
	"time"
)

func Auth(rd *chain.RequestData, writer http.ResponseWriter, request *http.Request) chain.Result {
	cookie, err := request.Cookie("auth_token")
	if err != nil {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации 1"})
	}
	if cookie.Value == "" {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации 2 "})
	}
	tokenBytes, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверный формат токена авторизации"})
	}
	if len(tokenBytes) != 36 {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверная длина токена авторизации"})
	}
	tokenBytesArr := [36]byte(tokenBytes)
	tokenBytesArr = auth_token.UnShuffle(tokenBytesArr)
	userId, expireTime := auth_token.UnpackToken(tokenBytesArr)
	if time.Now().UnixNano() > expireTime {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации 3"})
	}
	if time.Now().Add(time.Hour*24*30).UnixNano() < expireTime {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации 4"})
	}
	if !validator.IsValidUnixNanoId(userId) {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации 5"})
	}
	userCache := cache.FindUserCacheById(userId)
	if userCache == nil {
		return render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "пользователь не найден"})
	}
	if userCache.Deleted {
		return render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "пользователь удален"})
	}
	if !userCache.CurrentUser.Enabled {
		return render.Json(writer, http.StatusForbidden, &dto.Err{ErrMessage: "пользователь заблокирован"})
	}
	if !auth_token.IsValidToken(tokenBytesArr, userCache.CurrentUser.SessionSecret) {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "неверный токен"})
	}
	rd.User = userCache
	return chain.Next()
}

func Login(rd *chain.RequestData, writer http.ResponseWriter, request *http.Request) chain.Result {
	requestDto := &dto.LoginRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	if err := validator.ValidateLoginRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: err.Error()})
	}
	userCache := cache.FindUserCacheByLogin(requestDto.Email)
	if userCache == nil {
		return render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "пользователь не найден"})
	}
	if !bytes.Equal(utils.GeneratePasswordHash(requestDto.Password), userCache.CurrentUser.PasswordHash) {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверный пароль"})
	}
	rd.User = userCache
	return chain.Next()
}

func CheckIsAdmin(rd *chain.RequestData, writer http.ResponseWriter, request *http.Request) chain.Result {
	if rd.User == nil || rd.User.CurrentUser.Id != config.GetAdminId() {
		return render.Json(writer, http.StatusForbidden, &dto.Err{ErrMessage: "пользователь не админ"})
	}
	return chain.Next()
}

func CheckGracefullyStop(rd *chain.RequestData, writer http.ResponseWriter, request *http.Request) chain.Result {
	if application.IsGracefullyStopped() {
		return render.Json(writer, http.StatusServiceUnavailable, &dto.Err{ErrMessage: "сервис временно недоступен"})
	}
	return chain.Next()
}

func SetAuthCookie(rd *chain.RequestData, writer http.ResponseWriter, request *http.Request) chain.Result {
	cookieDuration := time.Hour * 24 * 3
	newTokenBytes := auth_token.CreateToken(rd.User.CurrentUser.Id, time.Now().Add(cookieDuration).UnixNano(), rd.User.CurrentUser.SessionSecret)
	newTokenBytes = auth_token.Shuffle(newTokenBytes)
	newTokenStr := base64.StdEncoding.EncodeToString(newTokenBytes[:])
	http.SetCookie(writer, &http.Cookie{
		SameSite: http.SameSiteStrictMode,
		Name:     "auth_token",
		Value:    newTokenStr,
		Path:     "/",
		Domain:   config.GetDomain(),
		MaxAge:   24 * 3600 * 3,
		Secure:   true, // only sent over HTTPS
		HttpOnly: true, // not accessible via JavaScript
	})
	return chain.Next()
}

func FindAdv(rd *chain.RequestData, writer http.ResponseWriter, request *http.Request) chain.Result {
	advIdStr := request.PathValue("advId")
	advId, errConv := strconv.ParseInt(advIdStr, 10, 64)
	if errConv != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: errConv.Error()})
	}
	if !validator.IsValidUnixNanoId(advId) {
		return render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
	}
	advCache := cache.FindAdvCacheById(advId)
	if advCache == nil {
		return render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
	}
	rd.Adv = advCache
	return chain.Next()
}

func CheckAdvOwner(rd *chain.RequestData, writer http.ResponseWriter, request *http.Request) chain.Result {
	if rd.Adv.CurrentAdv.UserId != rd.User.CurrentUser.Id {
		return render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не принадлежит текущему пользователю"})
	}
	return chain.Next()
}

func StopIfUnsavedMoreThan(count int64) chain.HandlerFunction {
	return func(rd *chain.RequestData, writer http.ResponseWriter, request *http.Request) chain.Result {
		if cache.GetToSaveCount() >= count {
			return render.Json(writer, http.StatusTooManyRequests, &dto.Err{ErrMessage: "попробуйте позже"})
		}
		return chain.Next()
	}
}

func CheckConnectionAndTimeout(rd *chain.RequestData, writer http.ResponseWriter, request *http.Request) chain.Result {
	if err := request.Context().Err(); err != nil {
		return chain.Result{WriteErr: err}
	}
	if rd.Timeout() {
		return rd.GetOnTimeout()(rd, writer, request)
	}
	return chain.Next()
}
