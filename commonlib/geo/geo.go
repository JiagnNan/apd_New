package geo

import (
	"context"
	_ "embed"
	"github.com/JiagnNan/apd_New/commonlib/pkg/geolite"
	"net"
	"strings"
)

//go:embed GeoLite2-Country.mmdb
var data []byte

var mmdb *geolite.Reader

func init() {
	var err error
	mmdb, err = geolite.FromBytes(data)
	if err != nil {
		panic(err)
	}
}

func parseIp(ip string) net.IP {
	ip = strings.TrimSpace(ip)

	// 没有端口的 IPv4 或 IPv6
	if ip := net.ParseIP(ip); ip != nil {
		return ip
	}

	// 带端口的 IPv4 或 IPv6，例如：
	// 192.168.1.10:8080
	// [2001:db8::1]:8080
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		return nil
	}

	return net.ParseIP(host)
}

func IsCN(ip string) bool {
	country, err := mmdb.Country(parseIp(ip))
	if err != nil {
		return false
	}
	return country.Country.IsoCode == "CN"
}

// GetCountryCode 返回国家码，如：CN、US等
func GetCountryCode(ip string) string {
	country, err := mmdb.Country(parseIp(ip))
	if err != nil {
		return ""
	}
	return country.Country.IsoCode
}

// GetCountryName 返回国家名称，如：中国、美国等
func GetCountryName(ip string) string {
	country, err := mmdb.Country(parseIp(ip))
	if err != nil {
		return ""
	}
	return country.Country.Names["zh-CN"]
}

// GetContinentCode 返回大洲代码，如：AS - 亚洲,AF - 非洲,EU - 欧洲,NA - 北美洲,SA - 南美洲,OC - 大洋洲,AN - 南极洲
func GetContinentCode(ip string) string {
	country, err := mmdb.Country(parseIp(ip))
	if err != nil {
		return ""
	}
	return country.Continent.Code
}

func GetCountry(ctx context.Context, ip string) *geolite.Country {
	country, err := mmdb.Country(parseIp(ip))
	if err != nil {
		//log.Errorx(ctx, "GetCountry Err:%s", err)
		return &geolite.Country{}
	}
	return country
}
