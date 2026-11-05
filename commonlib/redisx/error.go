package redisx

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

// 判断 Redis 错误是否表示 key 或 field 不存在。
func IsNil(err error) bool {
	return errors.Is(err, redis.Nil)
}
