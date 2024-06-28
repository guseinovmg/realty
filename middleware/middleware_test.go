package middleware

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

var sessionSecret = [24]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
var userId int64 = 1657567
var expireTime = time.Now().Add(24 * time.Hour).UnixMicro()

func TestToken(t *testing.T) {
	fmt.Println("544747")
	token := CreateToken(userId, expireTime, sessionSecret)
	unpackedUserId, unpackedExpireTime := UnpackToken(token)
	if unpackedUserId != userId {
		t.Fatal("userId!=1")
	}
	if expireTime != unpackedExpireTime {
		fmt.Println(expireTime)
		fmt.Println(unpackedExpireTime)
		t.Fatal("expireTime")
	}
	if !IsValidToken(token, sessionSecret) {
		t.Fatal("IsValidToken")
	}
}

func TestShuffle(t *testing.T) {
	var arr = [36]byte([]byte("fsdfsbjfhsjhfvsgefyefiw73wg72rgwehjgrtyu"))
	shuffled := Shuffle(arr)
	fmt.Println(string(shuffled[:]))
	unShuffled := UnShuffle(shuffled)
	fmt.Println(string(unShuffled[:]))
	if !bytes.Equal(arr[:], unShuffled[:]) {
		t.Fatal("Shuffle error")
	}
}
