package redisx

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client 是项目内部 Redis 封装，避免业务层直接依赖 go-redis。
type Client struct {
	rdb redis.UniversalClient
}

// 创建 Redis 客户端并立即 Ping，确保启动时发现配置或网络问题。
func NewRedis(c Config) (*Client, error) {
	client := NewLazyRedis(c)
	if err := client.rdb.Ping(context.Background()).Err(); err != nil {
		_ = client.rdb.Close()
		return nil, err
	}
	return client, nil
}

// NewLazyRedis 保留可重连客户端，不把首次网络故障变成永久nil；仅用于允许降级的业务。
func NewLazyRedis(c Config) *Client {
	opts := &redis.UniversalOptions{
		Addrs:                 strings.Split(c.addr(), ","),
		Username:              c.Username,
		Password:              c.password(),
		DB:                    c.DB,
		PoolSize:              c.PoolSize,
		MinIdleConns:          c.MinIdleConns,
		DialTimeout:           seconds(c.DialTimeout, 5*time.Second),
		ReadTimeout:           seconds(c.ReadTimeout, 3*time.Second),
		WriteTimeout:          seconds(c.WriteTimeout, 3*time.Second),
		ContextTimeoutEnabled: true,
	}

	return &Client{rdb: redis.NewUniversalClient(opts)}
}

// 创建 Redis 客户端，失败时直接 panic，用于服务启动阶段强依赖 Redis 的场景。
func MustNewRedis(c Config) *Client {
	client, err := NewRedis(c)
	if err != nil {
		panic(err)
	}
	return client
}

// 关闭 Redis 客户端连接。
func (c *Client) Close() error {
	if c == nil || c.rdb == nil {
		return nil
	}
	return c.rdb.Close()
}
