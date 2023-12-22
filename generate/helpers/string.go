package helpers

import (
	"unicode"
)

func ToTitle(s string) string {

	if len(s) == 0 {
		return ""
	}

	newS := string(unicode.ToUpper(rune(s[0])))

	if len(s) == 1 {
		return newS
	}

	return newS + s[1:]
}

func ToCamel(s string) string {

	if len(s) == 0 {
		return ""
	}

	newS := string(unicode.ToLower(rune(s[0])))

	if len(s) == 1 {
		return newS
	}

	return newS + s[1:]
}
