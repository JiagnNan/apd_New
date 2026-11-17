package redisx

import "golang.org/x/net/context"

func (c *Client) LPop(ctx context.Context, key string) (string, error) {
	return c.rdb.LPop(ctx, key).Result()
}
