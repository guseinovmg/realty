package auth_token

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
)

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
