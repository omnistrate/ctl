package utils

import (
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
)

func CheckIfEmpty(parameter string) error {
	if strings.TrimSpace(parameter) == "" {
		return errors.New("parameter is empty")
	}
	return nil
}

func CheckIfNilOrEmpty(parameter *string) bool {
	if parameter == nil {
		return true
	}
	if strings.TrimSpace(*parameter) == "" {
		return true
	}
	return false
}

func TruncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func CutString(str string, length int) string {
	if length <= 0 {
		return ""
	}

	if utf8.RuneCountInString(str) < length {
		return str
	}

	return string([]rune(str)[:length])
}
