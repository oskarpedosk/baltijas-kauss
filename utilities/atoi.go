package utilities

import "strconv"

// Convert string to integer
func Atoi(str string) int {
	integer, _ := strconv.Atoi(str)
	return integer
}
