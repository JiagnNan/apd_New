// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"github.com/JiagnNan/apd_New/api-gateway/internal/svc"
	"github.com/JiagnNan/apd_New/api-gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NumberLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNumberLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NumberLoginLogic {
	return &NumberLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NumberLoginLogic) NumberLogin(req *types.NumberLoginReq) (resp *types.TokenResp, err error) {
	// todo: add your logic here and delete this line

	return
}
