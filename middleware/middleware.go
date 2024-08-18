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
	"realty/render"
	"realty/validator"
	"strconv"
	"time"
)

func Auth(rd *RequestData, writer http.ResponseWriter, request *http.Request) bool {
	cookie, err := request.Cookie("auth_token")
	if err != nil {
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации"})
		return false
	}
	if cookie.Value == "" {
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации"})
		return false
	}
	tokenBytes, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверный формат токена авторизации"})
		return false
	}
	if len(tokenBytes) != 36 {
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверная длина токена авторизации"})
		return false
	}
	tokenBytesArr := [36]byte(tokenBytes)
	userId, expireTime := UnpackToken(UnShuffle(tokenBytesArr))
	if time.Now().UnixNano() > expireTime {
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации"})
		return false
	}
	if time.Now().Add(time.Hour*24*30).UnixNano() < expireTime {
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации"})
		return false
	}
	if validator.IsValidUnixNanoId(userId) {
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации"})
		return false
	}
	userCache := cache.FindUserCacheById(userId)
	if userCache == nil {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "пользователь не найден"})
		return false
	}
	if userCache.Deleted {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "пользователь удален"})
		return false
	}
	if !userCache.CurrentUser.Enabled {
		_ = render.Json(writer, http.StatusForbidden, &dto.Err{ErrMessage: "пользователь заблокирован"})
		return false
	}
	if !IsValidToken(tokenBytesArr, userCache.CurrentUser.SessionSecret) {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "неверный токен"})
		return false
	}
	rd.User = userCache
	return true
}

func CheckIsAdmin(rd *RequestData, writer http.ResponseWriter, request *http.Request) bool {
	if rd.User == nil || rd.User.CurrentUser.Id != config.GetAdminId() {
		_ = render.Json(writer, http.StatusForbidden, &dto.Err{ErrMessage: "пользователь не админ"})
		return false
	}
	return true
}

func CheckGracefullyStop(rd *RequestData, writer http.ResponseWriter, request *http.Request) bool {
	if cache.IsGracefullyStopped() {
		_ = render.Json(writer, http.StatusServiceUnavailable, &dto.Err{ErrMessage: "сервис временно недоступен"})
		return false
	}
	return true
}

func SetAuthCookie(rd *RequestData, writer http.ResponseWriter, request *http.Request) bool {
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
	return true
}

func FindAdv(rd *RequestData, writer http.ResponseWriter, request *http.Request) bool {
	advIdStr := request.PathValue("advId")
	advId, errConv := strconv.ParseInt(advIdStr, 10, 64)
	if errConv != nil {
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: errConv.Error()})
		return false
	}
	if !validator.IsValidUnixNanoId(advId) {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return false
	}
	advCache := cache.FindAdvCacheById(advId)
	if advCache == nil {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не найдено"})
		return false
	}
	rd.Adv = advCache
	return true
}

func CheckAdvOwner(rd *RequestData, writer http.ResponseWriter, request *http.Request) bool {
	if rd.Adv.CurrentAdv.UserId != rd.User.CurrentUser.Id {
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "объявление не принадлежит текущему пользвателю"})
		return false
	}
	return true
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
