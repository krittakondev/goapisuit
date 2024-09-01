package utils

import "strings"


func CapitalizeFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

