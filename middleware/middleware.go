package middleware

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"realty/cache"
	"realty/config"
	"realty/dto"
	"realty/parsing_input"
	"realty/render"
	"realty/utils"
	"realty/validator"
	"strconv"
	"time"
)

func Auth(rd *RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	cookie, err := request.Cookie("auth_token")
	if err != nil {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{RequestId: rd.RequestId, ErrMessage: "ошибка авторизации 1"})
	}
	if cookie.Value == "" {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{RequestId: rd.RequestId, ErrMessage: "ошибка авторизации 2 "})
	}
	tokenBytes, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{RequestId: rd.RequestId, ErrMessage: "неверный формат токена авторизации"})
	}
	if len(tokenBytes) != 36 {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{RequestId: rd.RequestId, ErrMessage: "неверная длина токена авторизации"})
	}
	tokenBytesArr := [36]byte(tokenBytes)
	tokenBytesArr = UnShuffle(tokenBytesArr)
	userId, expireTime := UnpackToken(tokenBytesArr)
	if time.Now().UnixNano() > expireTime {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{RequestId: rd.RequestId, ErrMessage: "ошибка авторизации 3"})
	}
	if time.Now().Add(time.Hour*24*30).UnixNano() < expireTime {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{RequestId: rd.RequestId, ErrMessage: "ошибка авторизации 4"})
	}
	if !validator.IsValidUnixNanoId(userId) {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{RequestId: rd.RequestId, ErrMessage: "ошибка авторизации 5"})
	}
	userCache := cache.FindUserCacheById(userId)
	if userCache == nil {
		return render.Json(writer, http.StatusNotFound, &dto.Err{RequestId: rd.RequestId, ErrMessage: "пользователь не найден"})
	}
	if userCache.Deleted {
		return render.Json(writer, http.StatusNotFound, &dto.Err{RequestId: rd.RequestId, ErrMessage: "пользователь удален"})
	}
	if !userCache.CurrentUser.Enabled {
		return render.Json(writer, http.StatusForbidden, &dto.Err{RequestId: rd.RequestId, ErrMessage: "пользователь заблокирован"})
	}
	if !IsValidToken(tokenBytesArr, userCache.CurrentUser.SessionSecret) {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{RequestId: rd.RequestId, ErrMessage: "неверный токен"})
	}
	rd.User = userCache
	return render.Next()
}

func Login(rd *RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	requestDto := &dto.LoginRequest{}
	if err := parsing_input.ParseRawJson(request, requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{RequestId: rd.RequestId, ErrMessage: err.Error()})
	}
	if err := validator.ValidateLoginRequest(requestDto); err != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{RequestId: rd.RequestId, ErrMessage: err.Error()})
	}
	userCache := cache.FindUserCacheByLogin(requestDto.Email)
	if userCache == nil {
		return render.Json(writer, http.StatusNotFound, &dto.Err{RequestId: rd.RequestId, ErrMessage: "пользователь не найден"})
	}
	if !bytes.Equal(utils.GeneratePasswordHash(requestDto.Password), userCache.CurrentUser.PasswordHash) {
		return render.Json(writer, http.StatusUnauthorized, &dto.Err{RequestId: rd.RequestId, ErrMessage: "неверный пароль"})
	}
	rd.User = userCache
	return render.Next()
}

func CheckIsAdmin(rd *RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	if rd.User == nil || rd.User.CurrentUser.Id != config.GetAdminId() {
		return render.Json(writer, http.StatusForbidden, &dto.Err{RequestId: rd.RequestId, ErrMessage: "пользователь не админ"})
	}
	return render.Next()
}

func CheckGracefullyStop(rd *RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	if cache.IsGracefullyStopped() {
		return render.Json(writer, http.StatusServiceUnavailable, &dto.Err{RequestId: rd.RequestId, ErrMessage: "сервис временно недоступен"})
	}
	return render.Next()
}

func SetAuthCookie(rd *RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	cookieDuration := time.Hour * 24 * 3
	newTokenBytes := CreateToken(rd.User.CurrentUser.Id, time.Now().Add(cookieDuration).UnixNano(), rd.User.CurrentUser.SessionSecret)
	newTokenBytes = Shuffle(newTokenBytes)
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
	return render.Next()
}

func FindAdv(rd *RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	advIdStr := request.PathValue("advId")
	advId, errConv := strconv.ParseInt(advIdStr, 10, 64)
	if errConv != nil {
		return render.Json(writer, http.StatusBadRequest, &dto.Err{RequestId: rd.RequestId, ErrMessage: errConv.Error()})
	}
	if !validator.IsValidUnixNanoId(advId) {
		return render.Json(writer, http.StatusNotFound, &dto.Err{RequestId: rd.RequestId, ErrMessage: "объявление не найдено"})
	}
	advCache := cache.FindAdvCacheById(advId)
	if advCache == nil {
		return render.Json(writer, http.StatusNotFound, &dto.Err{RequestId: rd.RequestId, ErrMessage: "объявление не найдено"})
	}
	rd.Adv = advCache
	return render.Next()
}

func CheckAdvOwner(rd *RequestData, writer http.ResponseWriter, request *http.Request) render.Result {
	if rd.Adv.CurrentAdv.UserId != rd.User.CurrentUser.Id {
		return render.Json(writer, http.StatusNotFound, &dto.Err{RequestId: rd.RequestId, ErrMessage: "объявление не принадлежит текущему пользователю"})
	}
	return render.Next()
}

func CreateToken(userId int64, nanoseconds int64, sessionSecret [24]byte) [36]byte {
	userIdBytes, expireTimeBytes := make([]byte, 8), make([]byte, 8)
	binary.LittleEndian.PutUint64(userIdBytes, uint64(userId))
	binary.LittleEndian.PutUint64(expireTimeBytes, uint64(nanoseconds))
	resultBytes := [36]byte{}
	for i := 0; i < 8; i++ {
		resultBytes[i] = userIdBytes[i]
	}
	for i := 8; i < 16; i++ {
		resultBytes[i] = expireTimeBytes[i-8]
	}
	hash := sha1.New()
	hash.Write(sessionSecret[:])
	hash.Write(resultBytes[:16])
	hashBytes := hash.Sum(nil)
	for i := 16; i < 36; i++ {
		resultBytes[i] = hashBytes[i-16]
	}
	return resultBytes
}

func UnpackToken(inputBytes [36]byte) (userId int64, nanoseconds int64) {
	userId = int64(binary.LittleEndian.Uint64(inputBytes[0:8]))
	nanoseconds = int64(binary.LittleEndian.Uint64(inputBytes[8:16]))
	return userId, nanoseconds
}

func IsValidToken(inputBytes [36]byte, sessionSecret [24]byte) bool {
	hash := sha1.New()
	hash.Write(sessionSecret[:])
	hash.Write(inputBytes[:16])
	hashBytes := hash.Sum(nil)[:]
	return bytes.Equal(inputBytes[16:], hashBytes)
}

func Shuffle(arr [36]byte) [36]byte {
	return [36]byte{
		arr[35],
		arr[10],
		arr[6],
		arr[8],
		arr[9],
		arr[7],
		arr[3],
		arr[2],
		arr[5],
		arr[4],
		arr[27],
		arr[26],
		arr[29],
		arr[30],
		arr[21],
		arr[22],
		arr[24],
		arr[25],
		arr[23],
		arr[15],
		arr[17],
		arr[28],
		arr[18],
		arr[20],
		arr[19],
		arr[13],
		arr[14],
		arr[11],
		arr[12],
		arr[33],
		arr[0],
		arr[34],
		arr[31],
		arr[32],
		arr[16],
		arr[1],
	}
}

func UnShuffle(arr [36]byte) [36]byte {
	return [36]byte{
		arr[30],
		arr[35],
		arr[7],
		arr[6],
		arr[9],
		arr[8],
		arr[2],
		arr[5],
		arr[3],
		arr[4],
		arr[1],
		arr[27],
		arr[28],
		arr[25],
		arr[26],
		arr[19],
		arr[34],
		arr[20],
		arr[22],
		arr[24],
		arr[23],
		arr[14],
		arr[15],
		arr[18],
		arr[16],
		arr[17],
		arr[11],
		arr[10],
		arr[21],
		arr[12],
		arr[13],
		arr[32],
		arr[33],
		arr[29],
		arr[31],
		arr[0],
	}
}
