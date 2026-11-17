package base

import (
	"fmt"
	"strconv"

	"github.com/JiagnNan/apd_New/commonlib/errorx"
	"github.com/JiagnNan/apd_New/commonlib/redisx"
	"golang.org/x/net/context"
)

const (
	userCodeKey               = "auth:user_code:available"
	userCodeCursorKey         = "auth:user_code:cursor"
	userCodePublishBatchLimit = 500
)

func (b *BaseSvc) NextUserId(ctx context.Context) (int64, error) {
	const key = "userIdx"
	return b.redisClient.Incr(ctx, key)
}

func (b *BaseSvc) UserCode(ctx context.Context) (int64, error) {
	pop, err := b.redisClient.LPop(ctx, userCodeKey)
	num, _ := strconv.ParseInt(pop, 10, 64)
	if err != nil {
		if redisx.IsNil(err) {
			return 0, errorx.ERR_TOO_MANY_REQUESTS
		}
		return 0, err
	}
	return num, nil
}

func (b *BaseSvc) publishUserCodes(ctx context.Context, expectedCursor string, codes []int64) (bool, error) {
	const script = `
local current = redis.call('GET', KEYS[2]) or ''
if current ~= ARGV[1] then
    return 0
end

redis.call('RPUSH', KEYS[1], unpack(ARGV, 3))
redis.call('SET', KEYS[2], ARGV[2])
return 1
`
	if len(codes) == 0 || len(codes) > userCodePublishBatchLimit {
		return false, fmt.Errorf("user code batch size %d is outside 1–%d", len(codes), userCodePublishBatchLimit)
	}
	keys := []string{userCodeKey, userCodeCursorKey}
	args := make([]any, 0, len(codes)+2)
	args = append(args, expectedCursor, codes[len(codes)-1])
	for _, code := range codes {
		args = append(args, code)
	}

	result, err := b.redisClient.Eval(ctx, script, keys, args...)
	if err != nil {
		return false, err
	}

	status, ok := result.(int64)
	if !ok {
		return false, fmt.Errorf("unexpected user code publish result type %T", result)
	}
	switch status {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("unexpected user code publish status %d", status)
	}
}
