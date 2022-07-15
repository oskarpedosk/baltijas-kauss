package utilities

import "strconv"

func Atoi(str string) int {
	integer, _ := strconv.Atoi(str)
	return integer
}
