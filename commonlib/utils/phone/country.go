package utils

import (
	_ "embed"
	"strconv"
	"strings"
)

//go:embed country.txt
var countryTxt string

type Country struct {
	CountryCode string
	ISO2        string
	DefaultName string
	Patterns    string
	Length      int
}

var countryList []Country

var countryMap = make(map[string]Country)

func GetCountry(countryCode string) (Country, bool) {
	v, ok := countryMap[countryCode]
	return v, ok
}

// 返回全部国家码配置。
func GetCountries() []Country {
	return countryList
}

func init() {
	x := strings.Split(countryTxt, "\n")
	for _, v := range x {
		v = strings.TrimSpace(v) //去掉首尾 空格、\t、\r、\n
		// 过滤国家码配置中的注释行。
		if strings.HasPrefix(v, "//") || strings.HasPrefix(v, "#") {
			continue
		}
		m := strings.Split(v, ";")
		if len(m) == 0 {
			continue
		}

		c := Country{CountryCode: m[0]}
		//第一段必须存在，作为 map key  如果后面段缺失，不 panic，只留空串或 0
		if len(m) > 1 {
			c.ISO2 = m[1]
		}
		if len(m) > 2 {
			c.DefaultName = m[2]
		}
		//去掉国家码前缀
		if len(m) > 3 {
			c.Patterns = strings.TrimPrefix(m[3], c.CountryCode)
		}
		//长度转整数
		if len(m) > 5 {
			c.Length, _ = strconv.Atoi(m[4])
		}
		//双容器存储
		countryList = append(countryList, c)
		countryMap[c.CountryCode] = c
	}
}
