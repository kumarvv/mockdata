package utils

import "strings"

func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func IsBlankPtr(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}
