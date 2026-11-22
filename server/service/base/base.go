package base

import (
	"github.com/JiagnNan/apd_New/commonlib/redisx"
	"github.com/JiagnNan/apd_New/server/dao/user"
)

type BaseSvc struct {
	redisClient *redisx.Client
	userDao     *user.UserDao
}

func NewBaseSvc(redis *redisx.Client, userDao *user.UserDao) *BaseSvc {
	return &BaseSvc{
		redisClient: redis,
		userDao:     userDao,
	}
}
