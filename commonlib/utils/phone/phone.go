package utils

import (
	"errors"
	"strings"
)

// ErrPhoneNumberInvalid 表示手机号格式或国家码无法识别。
var ErrPhoneNumberInvalid = errors.New("phone number invalid")

func parsePhone(phoneNumber string) (*Country, string, error) {
	phoneNumber = normalize(phoneNumber)
	if !IsAllDigits(phoneNumber) {
		return nil, "", ErrPhoneNumberInvalid
	}
	if len(phoneNumber) < 7 || len(phoneNumber) > 16 {
		return nil, "", ErrPhoneNumberInvalid
	}
	var c *Country
	var cLen int

	for i := 4; i >= 1; i-- {
		v, ok := GetCountry(phoneNumber[:i])
		if ok {
			cLen = i
			c = &v
			break
		}
	}
	if c == nil {
		return nil, "", ErrPhoneNumberInvalid
	}

	//国内格式号码首位为0时，去掉该数字0例如：840912113260更改为84912113260
	national := strings.TrimLeft(phoneNumber[cLen:], "0")
	if c.Length != len(national) && c.Length > 0 && len(national) != 0 {
		return nil, "", ErrPhoneNumberInvalid
	}

	noZeroNumber := phoneNumber[:cLen] + national
	return c, noZeroNumber, nil
}

func IsAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
