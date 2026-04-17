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

func SplitToInt(s, sep string) []int64 {
	items := strings.Split(s, sep)
	values := make([]int64, 0)
	for _, item := range items {
		values = append(values, ToInt64(item))
	}
	return values
}

func SplitToFloat(s, sep string) []float64 {
	items := strings.Split(s, sep)
	values := make([]float64, 0)
	for _, item := range items {
		values = append(values, ToFloat(item))
	}
	return values
}
