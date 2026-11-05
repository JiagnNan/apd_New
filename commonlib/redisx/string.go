package redisx

import (
	"context"
	"time"
)

// 读取字符串 key。
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

// 写入字符串 key，ttl 为 0 时表示不过期。
func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

// 仅当 key 不存在时写入，常用于分布式锁或幂等占位。
func (c *Client) SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	return c.rdb.SetNX(ctx, key, value, ttl).Result()
}

// 删除一个或多个 key。
func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

// 判断一个或多个 key 是否存在，任意存在即返回 true。
func (c *Client) Exists(ctx context.Context, keys ...string) (bool, error) {
	n, err := c.rdb.Exists(ctx, keys...).Result()
	return n > 0, err
}

// 设置 key 的过期时间。
func (c *Client) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.rdb.Expire(ctx, key, ttl).Err()
}

// TTL 读取 key 的剩余过期时间，用于只延长不缩短缓存生命周期。
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.rdb.TTL(ctx, key).Result()
}

// 读取哈希字段。
func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	return c.rdb.HGet(ctx, key, field).Result()
}

// 写入哈希字段。
func (c *Client) HSet(ctx context.Context, key string, values ...any) error {
	return c.rdb.HSet(ctx, key, values...).Err()
}

// 删除哈希字段。
func (c *Client) HDel(ctx context.Context, key string, fields ...string) error {
	return c.rdb.HDel(ctx, key, fields...).Err()
}

// 读取哈希中的全部字段。
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.rdb.HGetAll(ctx, key).Result()
}

// 判断哈希字段是否存在。
func (c *Client) HExists(ctx context.Context, key, field string) (bool, error) {
	return c.rdb.HExists(ctx, key, field).Result()
}
