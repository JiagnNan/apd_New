package dbx

import "time"

type Config struct {
	DataSource    string
	DSN           string
	MaxIdleConn   int
	MaxOpenConn   int
	MaxLifetime   int
	SlowThreshold int
}

// 返回最终使用的 MySQL DSN。
func (c Config) dsn() string {
	if c.DSN != "" {
		return c.DSN
	}
	return c.DataSource
}

// 返回数据库连接最大生命周期，未配置时默认 10 分钟。
func (c Config) maxLifetime() time.Duration {
	if c.MaxLifetime <= 0 {
		return time.Minute * 10
	}
	return time.Second * time.Duration(c.MaxLifetime)
}

// 返回慢 SQL 阈值，未配置时默认 200 毫秒。
func (c Config) slowThreshold() time.Duration {
	if c.SlowThreshold <= 0 {
		return time.Millisecond * 200
	}
	return time.Millisecond * time.Duration(c.SlowThreshold)
}
