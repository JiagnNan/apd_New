package logic

import (
	"context"

	"github.com/JiagnNan/apd_New/server/rpc/auth/authpb"
	"github.com/JiagnNan/apd_New/server/rpc/auth/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NumberLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNumberLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NumberLoginLogic {
	return &NumberLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NumberLoginLogic) NumberLogin(in *authpb.NumberLoginReq) (*authpb.TokenResp, error) {
	// todo: add your logic here and delete this line

	return &authpb.TokenResp{}, nil
}
