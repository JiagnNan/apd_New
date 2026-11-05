package redisx

import "context"

// Eval 原子执行 Lua 脚本，复杂队列操作仍由上层业务封装具体键和参数。
func (c *Client) Eval(ctx context.Context, script string, keys []string, args ...any) (any, error) {
	return c.rdb.Eval(ctx, script, keys, args...).Result()
}
