package utils

import (
	"crypto/sha1"
	"reflect"
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
