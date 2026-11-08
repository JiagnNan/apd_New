package utils

import (
	"github.com/JiagnNan/apd_New/commonlib/errorx"
	"strings"
)
// 规范化手机号输入，去掉空格、短横线和加号。
func normalize(phoneNumber string) string {
	phoneNumber = strings.TrimSpace(phoneNumber)
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")
	phoneNumber = strings.TrimPrefix(phoneNumber, "+")
	return phoneNumber
}

// 规范化登录账号，邮箱直接返回，手机号解析后返回真实号码。
func normalizeAccount(account string) (string, error) {
	if account == "" {
		return "", errorx.ERR_PARAM_INVALID
	}
	if strings.Contains(account, "@") {
		return account, nil
	}

	_, phoneNumber, err := parsePhone(account)
	if err != nil {
		return "", errorx.ERR_PHONE_NUMBER_INVALID
	}

	return phoneNumber, nil
}