package helpers

import (
	"strconv"
)

//Int64ToString function convert a int64 number to a string
func Int64ToString(input int64) string {
	return strconv.Itoa(int(input))
}

//StringToInt64 function convert a string to a int64 number
func StringToInt64(input string) (int64, error) {
	intvalue, err := strconv.Atoi(input)
	if err != nil {
		return int64(intvalue), err
	}
	return int64(intvalue), err
}
