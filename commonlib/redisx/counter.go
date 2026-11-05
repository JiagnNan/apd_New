package redisx

import "context"

// 对 key 做自增 1。
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.rdb.Incr(ctx, key).Result()
}

// 对 key 按指定步长自增。
func (c *Client) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.rdb.IncrBy(ctx, key, value).Result()
}

// 对 key 做自减 1。
func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	return c.rdb.Decr(ctx, key).Result()
}
