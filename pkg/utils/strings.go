package utils

import (
	"regexp"
	"strings"
)


func CapitalizeFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func KebabToCamel(s string) string {
	s = strings.ToLower(s)
	words := strings.Split(s, "-")
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}
func CamelToKebab(s string) string {
	// Regular expression to find uppercase letters
	re := regexp.MustCompile("([a-z])([A-Z])")
	// Insert a hyphen between lowercase-uppercase pairs and convert to lowercase
	kebab := re.ReplaceAllString(s, "${1}-${2}")
	return strings.ToLower(kebab)
}
