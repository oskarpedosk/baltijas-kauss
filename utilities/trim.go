package utilities

import "strings"

// Trim space in a string
func TrimSpace(str string) string {
	str = strings.TrimSpace(str)
	return str
}

// Trim condition in a string
func Trim(str string, condition string) string {
	str = strings.Trim(str, condition)
	return str
}
