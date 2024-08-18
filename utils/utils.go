package utils

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

type H map[string]string

// UnsafeStringToBytes converts string to byte slice without a memory allocation.
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func UnsafeStringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeBytesToString converts byte slice to string without a memory allocation.
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func UnsafeBytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func Iff[T any](cond bool, a T, b T) T {
	if cond {
		return a
	} else {
		return b
	}
}

func Is(value any) bool {
	if value == nil {
		return true
	}

	refVal := reflect.ValueOf(value)
	switch refVal.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return refVal.IsNil()
	default:
		return false
	}
}

// Normalize converts typed nils (e.g. []byte(nil)) into untyped nil. Other values are returned unmodified.
func Normalize(v any) any {
	if Is(v) {
		return nil
	}
	return v
}

// NormalizeSlice converts all typed nils (e.g. []byte(nil)) in s into untyped nils. Other values are unmodified. s is
// mutated in place.
func NormalizeSlice(s []any) {
	for i := range s {
		if Is(s[i]) {
			s[i] = nil
		}
	}
}

func GeneratePasswordHash(password string) []byte {
	hash := sha1.New()
	hash.Write(UnsafeStringToBytes(password))
	hash.Write(UnsafeStringToBytes("bla bla secret 36464663464"))
	return hash.Sum(nil)
}

func GenerateSessionsSecret(randomBytes []byte) [24]byte {
	hash := sha256.New()
	hash.Write(UnsafeStringToBytes(time.Now().Add(time.Hour * 24 * 100).String()))
	hash.Write(randomBytes)
	return [24]byte(hash.Sum(nil))
}

func NewRequest(method string, headers H, url string, pathParams H, queryParams H, body any) (*http.Request, error) {
	if pathParams != nil {
		for k, v := range pathParams {
			url = strings.Replace(url, k, "{"+v+"}", 1)
		}
	}

	if strings.Contains(url, "{") || strings.Contains(url, "}") {
		return nil, errors.New("неправильно сформирован путь " + url)
	}
	var buf io.Reader
	if body != nil {
		bytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = strings.NewReader(string(bytes))
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("fail to create request: %s", err.Error()))
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	if queryParams != nil {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	return req, nil
}
