package config

import (
	"github.com/JiagnNan/apd_New/commonlib/dbx"
	"github.com/JiagnNan/apd_New/commonlib/pkg/log"
	"github.com/JiagnNan/apd_New/commonlib/redisx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MysqlConf dbx.Config
	RedisConf redisx.Config
	AuthKey   AuthKeyConf
	PhoneCode PhoneCodeConf
	AppLog    log.Config `json:"AppLog,optional"`
}

type AuthKeyConf struct {
	AccessSecret  string `json:",optional"`
	AccessExpire  int64
	RefreshSecret string `json:",optional"`
	RefreshExpire int64
}

type PhoneCodeConf struct {
	Expire         int64
	ResendInterval int64
}
