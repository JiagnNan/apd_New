package logic

import (
	"context"
	utils "github.com/JiagnNan/apd_New/commonlib/utils/phone"
	"github.com/JiagnNan/apd_New/server/service/auth"

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
	var phoneNumber string
	var err error
	phoneNumber, err = utils.NormalizeAccount(in.Number)
	if err != nil {
		return nil, err
	}
	code, err := l.svcCtx.AuthSvc.CheckPhoneCode(l.ctx, phoneNumber, in.Code, auth.CodeTypeLogin)
	if err != nil {
		return nil, err
	}

	return &authpb.TokenResp{}, nil
}
