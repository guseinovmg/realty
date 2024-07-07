package middleware

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"realty/cache"
	"realty/dto"
	"realty/render"
	"realty/validator"
	"time"
)

func Auth(rd *RequestData, writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("auth_token")
	if err != nil {
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации"})
		return
	}
	if cookie.Value == "" {
		rd.Stop()
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации"})
		return
	}
	tokenBytes, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		rd.Stop()
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверный формат токена авторизации"})
		return
	}
	if len(tokenBytes) != 36 {
		rd.Stop()
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "неверная длина токена авторизации"})
		return
	}
	tokenBytesArr := [36]byte(tokenBytes)
	userId, expireTime := UnpackToken(UnShuffle(tokenBytesArr))
	if time.Now().UnixNano() > expireTime {
		rd.Stop()
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации"})
		return
	}
	if time.Now().Add(time.Hour*24*30).UnixNano() < expireTime {
		rd.Stop()
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации"})
		return
	}
	if validator.IsValidUnixNanoId(userId) {
		rd.Stop()
		_ = render.Json(writer, http.StatusUnauthorized, &dto.Err{ErrMessage: "ошибка авторизации"})
		return
	}
	userCache := cache.FindUserCacheById(userId)
	if userCache == nil {
		rd.Stop()
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "пользователь не найден"})
		return
	}
	if userCache.Deleted {
		rd.Stop()
		_ = render.Json(writer, http.StatusNotFound, &dto.Err{ErrMessage: "пользователь удален"})
		return
	}
	if !userCache.CurrentUser.Enabled {
		rd.Stop()
		_ = render.Json(writer, http.StatusForbidden, &dto.Err{ErrMessage: "пользователь заблокирован"})
		return
	}
	if !IsValidToken(tokenBytesArr, userCache.CurrentUser.SessionSecret) {
		rd.Stop()
		_ = render.Json(writer, http.StatusBadRequest, &dto.Err{ErrMessage: "неверный токен"})
		return
	}
	rd.User = userCache
}

func CheckIsAdmin(rd *RequestData, writer http.ResponseWriter, request *http.Request) {
	if rd.User == nil || rd.User.CurrentUser.Id != 4446456464 {
		rd.Stop()
		writer.WriteHeader(http.StatusForbidden)
		return
	}
}

func SetAuthCookie(rd *RequestData, writer http.ResponseWriter, request *http.Request) {
	cookieDuration := time.Hour * 24 * 3
	newTokenBytes := CreateToken(rd.User.CurrentUser.Id, time.Now().Add(cookieDuration).UnixNano(), rd.User.CurrentUser.SessionSecret)
	newTokenBytes = Shuffle(newTokenBytes)
	newTokenStr := base64.StdEncoding.EncodeToString(newTokenBytes[:])
	http.SetCookie(writer, &http.Cookie{
		SameSite: http.SameSiteStrictMode, //todo разобраться какой нужен
		Name:     "auth_token",
		Value:    newTokenStr,
		Path:     "/",
		Domain:   "example.com", //todo
		Expires:  time.Now().Add(cookieDuration),
		Secure:   true, // only sent over HTTPS
		HttpOnly: true, // not accessible via JavaScript
	})
}

func CreateToken(userId int64, microseconds int64, sessionSecret [24]byte) [36]byte {
	userIdBytes, expireTimeBytes := make([]byte, 8), make([]byte, 8)
	binary.LittleEndian.PutUint64(userIdBytes, uint64(userId))
	binary.LittleEndian.PutUint64(expireTimeBytes, uint64(microseconds))
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

func UnpackToken(inputBytes [36]byte) (userId int64, microseconds int64) {
	userId = int64(binary.LittleEndian.Uint64(inputBytes[0:8]))
	microseconds = int64(binary.LittleEndian.Uint64(inputBytes[8:16]))
	return userId, microseconds
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
