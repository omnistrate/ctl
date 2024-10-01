package utils

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	KeyRegex = `^[a-zA-Z][a-zA-Z0-9_]*$`
)

var (
	camelCaseRegex = regexp.MustCompile("[^a-zA-Z0-9-]+")

	// Use a regular expression to detect non-alphanumeric characters
	nonAlphaNumericRegEx = regexp.MustCompile("[^a-zA-Z0-9]+")

	// by using a regex with Unicode categories, we also accept letters from alphabets other than English
	nonAlphaNumericAllowSpacesRegEx = regexp.MustCompile(`[^\p{L}\p{N} ]+`)
)

var (
	keyRegexParsed = regexp.MustCompile(KeyRegex)
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

func CheckIfNonAlphanumeric(parameter string) bool {
	return nonAlphaNumericAllowSpacesRegEx.MatchString(parameter)
}

func ValidateURL(urlTest string) (err error) {
	_, err = url.Parse(urlTest)
	if err != nil {
		err = errors.Wrap(err, "invalid URL")
		return
	}
	return
}

// convertToLowerCamelCase converts a slice of words into a camelCase string.
func convertToLowerCamelCase(words []string) string {
	var result string
	for i, word := range words {
		if i == 0 {
			result += strings.ToLower(word) // First word in lowercase
		} else {
			result += cases.Title(language.English).String(strings.ToLower(word)) // Capitalize the first letter of each subsequent word
		}
	}
	return result
}

// convertToCamelCase converts a slice of words into a CamelCase string.
func convertToCamelCase(words []string) string {
	var result string
	for _, word := range words {
		result += cases.Title(language.English).String(strings.ToLower(word)) // Capitalize the first letter of each subsequent word
	}
	return result
}

// ToLowerCamelCase converts a string to lowerCamelCase.
func ToLowerCamelCase(str string) string {
	if camelCaseRegex.MatchString(str) {
		// If the string contains non-alphanumeric characters, treat it as non-camelCase and convert.
		words := camelCaseRegex.Split(str, -1)
		return convertToLowerCamelCase(words)
	} else if strings.ToUpper(str) == str {
		// If the string is all uppercase, convert the entire string to lowercase.
		return strings.ToLower(str)
	}

	// For strings that don't match the above conditions, only ensure the first letter is lowercase,
	// assuming it's either already camelCase or a single word that needs minimal conversion.
	return convertFirstToLower(str)
}

// ToCamelCase converts a string to CamelCase.
func ToCamelCase(str string) string {
	if camelCaseRegex.MatchString(str) {
		// If the string contains non-alphanumeric characters, treat it as non-camelCase and convert.
		words := camelCaseRegex.Split(str, -1)
		return convertToCamelCase(words)
	} else if strings.ToUpper(str) == str {
		// If the string is all uppercase, convert the entire string to title case.
		return cases.Title(language.English).String(strings.ToLower(str))
	}

	// For strings that don't match the above conditions, only ensure the first letter is uppercase,
	// assuming it's either already camelCase or a single word that needs minimal conversion.
	return convertFirstToUpper(str)
}

// convertFirstToLower ensures the first letter of the string is lowercase,
// useful for already camelCased or single-word inputs.
func convertFirstToLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

// convertFirstToUpper ensures the first letter of the string is uppercase,
// useful for already CamelCased or single-word inputs.
func convertFirstToUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func ToSnakeCase(key string) (snakeCased string) {
	// Convert key in any format to snake_case
	keyParts := nonAlphaNumericRegEx.Split(key, -1)
	if len(keyParts) < 2 {
		// No key parts found
		snakeCased = strings.ToLower(key)
		return
	}

	caser := cases.Lower(language.English)

	// Convert first part to lowercase
	snakeCased = strings.ToLower(keyParts[0])

	// Convert remaining parts to lowercase
	for _, keyPart := range keyParts[1:] {
		snakeCased += "_" + caser.String(keyPart)
	}

	return
}

func ValidateKey(key string) (err error) {
	if err = CheckIfEmpty(key); err != nil {
		err = errors.Wrap(err, "key is empty")
		return
	}
	if !keyRegexParsed.MatchString(key) {
		err = errors.Errorf("invalid key: %s. expected to match regex: %s", key, KeyRegex)
		return
	}

	return
}

func NormalizeToPathSafeStringPreserveCase(key string) (result string) {
	result = filepath.Clean(key)
	result = strings.ReplaceAll(result, " ", "-")
	result = strings.ReplaceAll(result, ",", "-")
	result = strings.ReplaceAll(result, ":", "-")
	result = strings.ReplaceAll(result, ";", "-")
	result = strings.ReplaceAll(result, "/", "-")
	result = strings.ReplaceAll(result, "\\", "-")
	result = strings.ReplaceAll(result, "\"", "-")
	result = strings.ReplaceAll(result, "'", "-")
	result = strings.ReplaceAll(result, "(", "-")
	result = strings.ReplaceAll(result, ")", "-")
	result = strings.ReplaceAll(result, "[", "-")
	result = strings.ReplaceAll(result, "]", "-")
	result = strings.ReplaceAll(result, "{", "-")
	result = strings.ReplaceAll(result, "}", "-")
	result = strings.ReplaceAll(result, "<", "-")
	result = strings.ReplaceAll(result, ">", "-")
	result = strings.ReplaceAll(result, "?", "-")
	result = strings.ReplaceAll(result, "!", "-")
	result = strings.ReplaceAll(result, "@", "-")
	result = strings.ReplaceAll(result, "#", "-")
	result = strings.ReplaceAll(result, "$", "-")
	result = strings.ReplaceAll(result, "%", "-")
	result = strings.ReplaceAll(result, "^", "-")
	result = strings.ReplaceAll(result, "&", "-")
	result = strings.ReplaceAll(result, "*", "-")
	result = strings.ReplaceAll(result, "|", "-")
	result = strings.ReplaceAll(result, "`", "-")
	result = strings.ReplaceAll(result, "~", "-")
	result = strings.ReplaceAll(result, "=", "-")
	result = strings.ReplaceAll(result, "+", "-")
	result = strings.ReplaceAll(result, "_", "-")
	result = strings.ReplaceAll(result, ".", "-")

	return
}

func NormalizeToPathSafeString(key string) (result string) {
	result = NormalizeToPathSafeStringPreserveCase(ToSnakeCase(key))
	return
}

func ValidateIfIsAValidUnixPath(path string) (err error) {
	cleanedPath := filepath.Clean(path)

	if err = CheckIfEmpty(cleanedPath); err != nil {
		err = errors.Wrap(err, "path is empty")
		return
	}

	if !filepath.IsAbs(cleanedPath) {
		err = errors.Errorf("path is not absolute: %s", path)
		return
	}

	var subPath string
	if subPath, err = filepath.Rel(string(filepath.Separator), cleanedPath); err != nil {
		err = errors.Wrap(err, "invalid path")
		return
	}

	if subPath == ".." || strings.HasPrefix(subPath, ".."+string(filepath.Separator)) {
		err = errors.Errorf("path is not absolute: %s", path)
		return
	}
	return
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

func Keyify(in string) (out string) {
	return strings.ReplaceAll(NormalizeToPathSafeString(strings.ToLower(in)), "_", "")
}

func RemoveDashes(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "-", ""))
}

func ConvertSnakeToLowerCamelCase(key string) (camelCased string) {
	// Convert key in any format to snake_case
	keyParts := nonAlphaNumericRegEx.Split(key, -1)
	if len(keyParts) < 2 {
		// No key parts found
		camelCased = strings.ToLower(keyParts[0])
		return
	}

	caser := cases.Title(language.English)

	// Convert first part to lowercase
	camelCased = strings.ToLower(keyParts[0])

	// Convert remaining parts to lowercase
	for _, keyPart := range keyParts[1:] {
		camelCased += caser.String(keyPart)
	}

	return
}
