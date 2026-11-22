package convert

import "strconv"

func ParseInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
