package base

import (
	"strconv"
	"strings"
)

const (
	minUserCode = 10000000
	maxUserCode = 99999999
)

func isValidUserCode(code int64) bool {
	if code < minUserCode || code > maxUserCode {
		return false
	}
	codeStr := strconv.FormatInt(code, 10)
	if codeStr[0] == codeStr[7] {
		return false
	}
	if strings.Contains(codeStr, "168") {
		return false
	}
	var counts [10]int
	oddCount := 0
	for i := 0; i < len(codeStr); i++ {
		//字符 '7' 对应 ASCII 编码值 55，字符 '0' 对应 ASCII 编码值 48。
		digit := codeStr[i] - '0'
		counts[digit]++
		if counts[digit] > 2 {
			return false
		}
		//计算奇数个数
		if digit%2 != 0 {
			oddCount++
		}
	}
	if oddCount < 3 || oddCount > 5 {
		return false
	}
	for i := 0; i < len(codeStr)-1; i++ {
		if codeStr[i] == codeStr[i+1] {
			return false
		}
		// 先检查隔一位，例如 434。
		if i+2 < len(codeStr) && codeStr[i] == codeStr[i+2] {
			return false
		}
		// 剩下不足四位时，跳过升降序检查。
		if i+3 >= len(codeStr) {
			continue
		}

		step := int(codeStr[i+1]) - int(codeStr[i])
		//判断是否升序，1234,4321
		if step == 1 || step == -1 {
			secondStep := int(codeStr[i+2]) - int(codeStr[i+1])
			thirdStep := int(codeStr[i+3]) - int(codeStr[i+2])
			if secondStep == step && thirdStep == step {
				return false
			}
		}
	}
	return true
}

func nextValidUserCode(cursor int64) (int64, bool) {
	if cursor >= maxUserCode {
		return 0, false
	}
	cursor += 1
	if cursor < minUserCode {
		cursor = minUserCode
	}
	for ; cursor <= maxUserCode; cursor++ {
		ok := isValidUserCode(cursor)
		if ok {
			return cursor, true
		}
	}
	return 0, false
}

func buildUserCodes(cursor int64, limit int) []int64 {
	if limit <= 0 {
		return nil
	}
	codes := make([]int64, 0, limit)
	for len(codes) < limit {
		next, ok := nextValidUserCode(cursor)
		if !ok {
			break
		}
		codes = append(codes, next)
		cursor = next
	}
	return codes
}
