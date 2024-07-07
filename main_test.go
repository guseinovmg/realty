package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	fmt.Println(math.MaxInt64)
	//18446744073709551615
	//9223372036854775807
	//1718200199829214
	fmt.Println(time.UnixMicro(math.MaxInt64))
	fmt.Println(time.Now().UnixMicro())

	userIdBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(userIdBytes, uint64(time.Now().UnixNano()))
	log.Println(userIdBytes)
}
