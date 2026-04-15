package utils

import (
	"math/rand"
	"strings"
)

func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

//func IsBlankPtr(s *string) bool {
//	return s == nil || strings.TrimSpace(*s) == ""
//}

func Includes(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func RandomOneOf[T any](values ...T) T {
	l := len(values)
	i := rand.Intn(l)
	return values[i]
}
