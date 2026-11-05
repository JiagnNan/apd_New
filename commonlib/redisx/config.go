package redisx

import "time"

// Config 兼容 go-zero 和自定义 Redis 配置字段。
type Config struct {
	Addr         string `json:",optional"`
	Host         string `json:",optional"`
	Username     string `json:",optional"`
	Password     string `json:",optional"`
	Pass         string `json:",optional"`
	DB           int    `json:",optional"`
	Type         string `json:",optional"`
	PoolSize     int    `json:",optional"`
	MinIdleConns int    `json:",optional"`
	DialTimeout  int64  `json:",optional"`
	ReadTimeout  int64  `json:",optional"`
	WriteTimeout int64  `json:",optional"`
}

// 返回 Redis 地址，优先使用 Addr，兼容旧配置 Host。
func (c Config) addr() string {
	if c.Addr != "" {
		return c.Addr
	}
	return c.Host
}

// 返回 Redis 密码，优先使用 Password，兼容旧配置 Pass。
func (c Config) password() string {
	if c.Password != "" {
		return c.Password
	}
	return c.Pass
}

// 将秒级配置转换为 time.Duration，未配置时使用默认值。
func seconds(v int64, fallback time.Duration) time.Duration {
	if v <= 0 {
		return fallback
	}
	return time.Duration(v) * time.Second
}
