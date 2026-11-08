package svc

import (
	"github.com/JiagnNan/apd_New/commonlib/redisx"
	"github.com/JiagnNan/apd_New/server/rpc/auth/internal/config"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redisx.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
