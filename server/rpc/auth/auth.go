package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/JiagnNan/apd_New/server/rpc/auth/authpb"
	"github.com/JiagnNan/apd_New/server/rpc/auth/internal/config"
	"github.com/JiagnNan/apd_New/server/rpc/auth/internal/server"
	"github.com/JiagnNan/apd_New/server/rpc/auth/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const userCodeRefillTimeout = 5 * time.Second

var configFile = flag.String("f", "etc/auth.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	bootCtx, cancel := context.WithTimeout(context.Background(), userCodeRefillTimeout)
	err := ctx.BaseSvc.SupplementUserCodes(bootCtx)
	cancel()
	if err != nil {
		logx.Errorf("启动补号失败: %v", err)
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		authpb.RegisterAuthServer(grpcServer, server.NewAuthServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	refillCtx, stopRefill := context.WithCancel(context.Background())
	defer stopRefill()
	go runUserCodeRefill(refillCtx, ctx)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

func runUserCodeRefill(ctx context.Context, svcCtx *svc.ServiceContext) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			roundCtx, cancel := context.WithTimeout(ctx, userCodeRefillTimeout)
			err := svcCtx.BaseSvc.SupplementUserCodes(roundCtx)
			cancel()
			if err != nil {
				logx.Errorf("定时补号失败: %v", err)
			}
		}
	}
}
