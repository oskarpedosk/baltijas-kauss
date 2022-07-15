package utilities

import "strings"

func TrimSpace(str string) string {
	str = strings.TrimSpace(str)
	return str
}

func Trim(str string, condition string) string {
	str = strings.Trim(str, condition)
	return str
}
