package utils

import (
	"strconv"
)

func StringMustInt64(str string) int64 {
	res, _ := strconv.ParseInt(str, 10, 64)
	return res
}
