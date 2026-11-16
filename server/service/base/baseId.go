package base

import (
	"golang.org/x/net/context"
)

func (b *BaseSvc) NextUserId(ctx context.Context) (int64, error) {
	const key = "userIdx"
	return b.redisClient.Incr(ctx, key)
}

func (b *BaseSvc) UserCode(ctx context.Context) (int64, error) {
	return 0, nil
}
