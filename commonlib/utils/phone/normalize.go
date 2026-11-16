package utils

import (
	"github.com/JiagnNan/apd_New/commonlib/errorx"
	"strings"
)

// 规范化手机号输入，去掉空格、短横线和加号。
func Normalize(phoneNumber string) string {
	phoneNumber = strings.TrimSpace(phoneNumber)
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")
	phoneNumber = strings.TrimPrefix(phoneNumber, "+")
	return phoneNumber
}

// 规范化登录账号，邮箱直接返回，手机号解析后返回真实号码。
func NormalizeAccount(account string) (string, bool, error) {
	if account == "" {
		return "", false, errorx.ERR_PARAM_INVALID
	}
	if strings.Contains(account, "@") {
		return account, false, nil
	}

	_, phoneNumber, err := ParsePhone(account)
	if err != nil {
		return "", false, errorx.ERR_PHONE_NUMBER_INVALID
	}

	return phoneNumber, true, nil
}
