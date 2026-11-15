package svc

import (
	"fmt"
	"github.com/JiagnNan/apd_New/commonlib/dbx"
	"github.com/JiagnNan/apd_New/commonlib/jwt"
	"github.com/JiagnNan/apd_New/commonlib/redisx"
	auth2 "github.com/JiagnNan/apd_New/server/dao/auth"
	user2 "github.com/JiagnNan/apd_New/server/dao/user"
	"github.com/JiagnNan/apd_New/server/rpc/auth/internal/config"
	"github.com/JiagnNan/apd_New/server/service/auth"
	"github.com/JiagnNan/apd_New/server/service/user"
	"gorm.io/gorm"
	"time"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	Redis   *redisx.Client
	AuthSvc *auth.AuthSvc
	UserSvc *user.UserSvc
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := dbx.MustNewMysql(c.MysqlConf)
	redisClient := redisx.MustNewRedis(c.RedisConf)

	manager, err := jwt.NewManager(jwt.Config{
		AccessSecret:  c.AuthKey.AccessSecret,
		RefreshSecret: c.AuthKey.RefreshSecret,
		AccessTTL:     time.Duration(c.AuthKey.AccessExpire) * time.Second,
		RefreshTTL:    time.Duration(c.AuthKey.RefreshExpire) * time.Second,
	})
	if err != nil {
		fmt.Println("jwt.NewManager 初始化失败")
		panic(err)
	}
	
	authDao := auth2.NewAuthDao(db)
	sessionDao := auth2.NewSessionDao(redisClient)
	authSvc := auth.NewAuthSvc(authDao, sessionDao, manager)

	userDao := user2.NewUserDao(db)
	userSvc := user.NewUserSvc(userDao)

	return &ServiceContext{
		Config:  c,
		DB:      db,
		Redis:   redisClient,
		AuthSvc: authSvc,
		UserSvc: userSvc,
	}
}
