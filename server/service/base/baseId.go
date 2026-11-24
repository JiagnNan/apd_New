package base

import (
	"fmt"
	"github.com/JiagnNan/apd_New/commonlib/utils/convert"
	"strconv"

	"github.com/JiagnNan/apd_New/commonlib/errorx"
	"github.com/JiagnNan/apd_New/commonlib/redisx"
	"golang.org/x/net/context"
)

const (
	userCodeKey                     = "auth:user_code:available"
	userCodeCursorKey               = "auth:user_code:cursor"
	userCodePublishBatchLimit       = 500
	targetStock               int64 = 10_000
)

func (b *BaseSvc) NextUserId(ctx context.Context) (int64, error) {
	const key = "userIdx"
	return b.redisClient.Incr(ctx, key)
}

func (b *BaseSvc) UserCode(ctx context.Context) (int64, error) {
	pop, err := b.redisClient.LPop(ctx, userCodeKey)

	if err != nil {
		if redisx.IsNil(err) {
			return 0, errorx.ERR_TOO_MANY_REQUESTS
		}
		return 0, err
	}
	num, err := strconv.ParseInt(pop, 10, 64)
	if err != nil {
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

// 号码池补库存
func (b *BaseSvc) SupplementUserCodes(ctx context.Context) error {
	var cursor int64
	lLen, err := b.redisClient.LLen(ctx, userCodeKey)
	if err != nil {
		return err
	}
	if lLen >= 2000 {
		return nil
	}
	expectedCursor, err := b.redisClient.Get(ctx, userCodeCursorKey)

	switch {
	case err == nil:
		cursor, err = convert.ParseInt64(expectedCursor)
	case redisx.IsNil(err):
		if lLen > 0 {
			return nil
		}
		expectedCursor = ""
		cursor, err = b.userDao.GetMaxUserCode(ctx)
		if err == nil && cursor == 0 {
			cursor = minUserCode
		}
	}
	if err != nil {
		return err
	}
	for lLen < targetStock {
		limit := int(min(targetStock-lLen, userCodePublishBatchLimit))
		codes := buildUserCodes(cursor, limit)
		if len(codes) == 0 {
			return nil
		}

		published, err := b.publishUserCodes(ctx, expectedCursor, codes)
		if err != nil {
			return err
		}
		if !published {
			break
		}
		lLen += int64(len(codes))
		cursor = codes[len(codes)-1]
		expectedCursor = strconv.FormatInt(cursor, 10)
	}

	return nil
}
