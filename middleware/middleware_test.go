package middleware

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

var sessionSecret = [24]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
var userId uint64 = 1657567
var expireTime = time.Now().Add(24 * time.Hour)

func TestToken(t *testing.T) {
	fmt.Println("544747")
	token := CreateToken(userId, expireTime, sessionSecret)
	unpackedUserId, unpackedExpireTime, err := UnpackToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if unpackedUserId != userId {
		t.Fatal("userId!=1")
	}
	if !expireTime.Equal(unpackedExpireTime) {
		t.Fatal("expireTime")
	}
	if !IsValidToken(token, sessionSecret) {
		t.Fatal("IsValidToken")
	}
}

func TestShuffle(t *testing.T) {
	var arr = [40]byte([]byte("fsdfsbjfhsjhfvsgefyefiw73wg72rgwehjgrtyu"))
	shuffled := Shuffle(arr)
	fmt.Println(string(shuffled[:]))
	unShaffled := UnShuffle(shuffled)
	fmt.Println(string(unShaffled[:]))
	if !bytes.Equal(arr[:], unShaffled[:]) {
		t.Fatal("Shuffle error")
	}
}
