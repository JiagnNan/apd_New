package logic

import (
	"context"
	"github.com/JiagnNan/apd_New/commonlib/errorx"
	"github.com/JiagnNan/apd_New/commonlib/mdata"
	utils "github.com/JiagnNan/apd_New/commonlib/utils/phone"
	modelAuth "github.com/JiagnNan/apd_New/server/model/auth"
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
	var isPhone bool
	var err error
	phoneNumber, isPhone, err = utils.NormalizeAccount(in.Number)
	if err != nil {
		return nil, err
	}
	checkedCode, err := l.svcCtx.AuthSvc.CheckPhoneCode(l.ctx, phoneNumber, in.Code, modelAuth.CodeTypeLogin)
	if err != nil {
		return nil, err
	}
	user, err := l.svcCtx.UserSvc.GetUserIdByPhoneOrEmail(l.ctx, phoneNumber)
	if err != nil {
		return nil, err
	}
	platform := mdata.FromContext(l.ctx).Platform
	switch platform {
	case "web", "android", "ios":
	default:
		return nil, errorx.ERR_PARAM_INVALID
	}
	if user.AuthVersion == "" {
		l.Errorf("missing auth version for user %d", user.Id)
		return nil, errorx.ERR_SERVER_INTERNAL
	}
	var t modelAuth.LoginType
	if !isPhone {
		t = modelAuth.EmailLoginType
	} else {
		t = modelAuth.PhoneLoginType
	}

	//软删除验证码
	if err := l.svcCtx.AuthSvc.ConsumeAuthCode(l.ctx, checkedCode); err != nil {
		return nil, err
	}
	err = l.svcCtx.AuthSvc.Login(l.ctx, user.Id, modelAuth.ActionTypeLogin, t, user.CountryCode)
	if err != nil {
		return nil, err
	}

	pair, err := l.svcCtx.AuthSvc.CreateSession(l.ctx, user.Id, platform, user.AuthVersion)
	if err != nil {
		l.Errorf("create login session: %v", err)
		return nil, errorx.ERR_SERVER_INTERNAL
	}

	l.svcCtx.IM.SyncAccountAsync(l.ctx, user.Id)
	return &authpb.TokenResp{
		AccessToken:          pair.AccessToken,
		RefreshToken:         pair.RefreshToken,
		AccessTokenExpireAt:  pair.AccessExpireAt.Unix(),
		RefreshTokenExpireAt: pair.RefreshExpireAt.Unix(),
		Jti:                  pair.AccessJTI,
	}, nil
}
