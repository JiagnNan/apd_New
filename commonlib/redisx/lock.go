package redisx

import (
	"context"
	"time"
)

const unlockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`

// 尝试获取分布式锁，value 应使用调用方生成的唯一持有者标识。
func (c *Client) TryLock(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	return c.SetNX(ctx, key, value, ttl)
}

// 释放分布式锁，仅当 value 与持有者标识一致时才删除锁。
func (c *Client) Unlock(ctx context.Context, key, value string) (bool, error) {
	n, err := c.rdb.Eval(ctx, unlockScript, []string{key}, value).Int64()
	return n > 0, err
}
