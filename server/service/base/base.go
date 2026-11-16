package base

import (
	"github.com/JiagnNan/apd_New/commonlib/redisx"
)

type BaseSvc struct {
	redisClient *redisx.Client
}

func NewBaseSvc(redis *redisx.Client) *BaseSvc {
	return &BaseSvc{
		redisClient: redis,
	}
}
