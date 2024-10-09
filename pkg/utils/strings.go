package utils

import "strings"


func CapitalizeFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func KebabToCamel(s string) string {
	words := strings.Split(s, "-")
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}
