package svc

import (
	"github.com/JiagnNan/apd_New/commonlib/dbx"
	"github.com/JiagnNan/apd_New/commonlib/redisx"
	auth2 "github.com/JiagnNan/apd_New/server/dao/auth"
	"github.com/JiagnNan/apd_New/server/rpc/auth/internal/config"
	"github.com/JiagnNan/apd_New/server/service/auth"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	Redis   *redisx.Client
	AuthSvc *auth.AuthSvc
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := dbx.MustNewMysql(c.MysqlConf)
	authDao := auth2.NewAuthDao(db)
	authSvc := auth.NewAuthSvc(authDao)
	redisClient := redisx.MustNewRedis(c.RedisConf)
	return &ServiceContext{
		Config:  c,
		DB:      db,
		Redis:   redisClient,
		AuthSvc: authSvc,
	}
}
