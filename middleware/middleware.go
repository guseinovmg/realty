package middleware

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"realty/cache"
	"realty/utils"
	"time"
)

func Auth(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {
	tokenHeader := request.Header.Get("Authorization")
	if tokenHeader == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		rd.Stop()
		return
	}
	tokenBytes, err := base64.StdEncoding.DecodeString(tokenHeader)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		rd.Stop()
		return
	}
	if len(tokenBytes) != 40 {
		writer.WriteHeader(http.StatusUnauthorized)
		rd.Stop()
		return
	}
	tokenBytesArr := [40]byte(tokenBytes)
	userId, expireTime, err := UnpackToken(UnShuffle(tokenBytesArr)) //todo проверить как работает такое приведение
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		rd.Stop()
		return
	}
	if time.Now().After(expireTime) {
		writer.WriteHeader(http.StatusUnauthorized)
		rd.Stop()
		return
	}
	if time.Now().Add(time.Hour * 24 * 30).Before(expireTime) {
		writer.WriteHeader(http.StatusUnauthorized)
		rd.Stop()
		return
	}
	user := cache.FindUserById(userId)
	if user == nil {
		writer.WriteHeader(http.StatusNotFound)
		rd.Stop()
		return
	}
	if !user.Enabled {
		writer.WriteHeader(http.StatusForbidden)
		rd.Stop()
		return
	}
	if !IsValidToken(tokenBytesArr, user.SessionSecret) {
		writer.WriteHeader(http.StatusBadRequest)
		rd.Stop()
		return
	}
	rd.User = user
}

func SetAuthCookie(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {
	newTokenBytes := CreateToken(rd.User.Id, time.Now().Add(time.Hour*24*3), rd.User.SessionSecret)
	newTokenBytes = Shuffle(newTokenBytes)
	newTokenStr := base64.StdEncoding.EncodeToString(newTokenBytes[:])
	writer.Header().Set("Cookie", newTokenStr)
}

func CreateToken(userId int64, expireTime time.Time, sessionSecret [24]byte) [40]byte {
	userIdBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(userIdBytes, uint64(userId))
	expireTimeBytes, _ := expireTime.MarshalBinary()
	resultBytes := [40]byte{}
	for i := 0; i < 8; i++ {
		resultBytes[i] = userIdBytes[i]
	}
	for i := 8; i < 24 && i < len(expireTimeBytes)+8; i++ {
		resultBytes[i] = expireTimeBytes[i-8]
	}
	hash := sha1.New()
	hash.Write(sessionSecret[:])
	hash.Write(resultBytes[:24])
	hashBytes := hash.Sum(nil)
	for i := 24; i < 40; i++ {
		resultBytes[i] = hashBytes[i-24]
	}
	return resultBytes
}

func UnpackToken(inputBytes [40]byte) (userId int64, expireTime time.Time, err error) {
	userId = int64(binary.LittleEndian.Uint64(inputBytes[0:8]))
	var timeBuf []byte
	if inputBytes[8] == 1 { //todo надо решить по таймзоне
		timeBuf = inputBytes[8:23]
	} else {
		timeBuf = inputBytes[8:24]
	}
	err = expireTime.UnmarshalBinary(timeBuf)
	if err != nil {
		return 0, time.Time{}, err
	}
	return userId, expireTime, nil
}

func IsValidToken(inputBytes [40]byte, sessionSecret [24]byte) bool {
	hash := sha1.New()
	hash.Write(sessionSecret[:])
	hash.Write(inputBytes[:24])
	hashBytes := hash.Sum(nil)[:16]
	return bytes.Equal(inputBytes[24:], hashBytes)
}

func Shuffle(arr [40]byte) [40]byte {
	return [40]byte{
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
		arr[36],
		arr[37],
		arr[38],
		arr[39],
	}
}

func UnShuffle(arr [40]byte) [40]byte {
	return [40]byte{
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
		arr[36],
		arr[37],
		arr[38],
		arr[39],
	}
}
