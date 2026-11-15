// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/JiagnNan/apd_New/api-gateway/internal/config"
	"github.com/JiagnNan/apd_New/server/rpc/auth/authpb"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	IMAuth authpb.AuthClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{Config: c}
	if c.IMAuthRPC != nil {
		svc.IMAuth = authpb.NewAuthClient(zrpc.MustNewClient(*c.IMAuthRPC).Conn())
	}
	return svc
}
