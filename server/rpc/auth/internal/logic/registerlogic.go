package logic

import (
	"context"
	"errors"
	"github.com/JiagnNan/apd_New/commonlib/mdata"
	modelAuth "github.com/JiagnNan/apd_New/server/model/auth"
	"strings"

	"github.com/JiagnNan/apd_New/commonlib/errorx"
	utils "github.com/JiagnNan/apd_New/commonlib/utils/phone"

	"github.com/JiagnNan/apd_New/server/rpc/auth/authpb"
	"github.com/JiagnNan/apd_New/server/rpc/auth/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *authpb.RegisterReq) (*authpb.TokenResp, error) {
	var phoneNumber string
	var isPhone bool
	var err error
	countryCode := "84"
	if strings.Contains(in.Number, "@") {
		phoneNumber = in.Number
		isPhone = false
	} else {
		var c *utils.Country
		c, phoneNumber, err = utils.ParsePhone(in.Number)
		if err != nil {
			return nil, errorx.ERR_PARAM_INVALID
		}
		isPhone = true
		countryCode = c.CountryCode
	}
	platform := mdata.FromContext(l.ctx).Platform
	switch platform {
	case "web", "android", "ios":
	default:
		return nil, errorx.ERR_PARAM_INVALID
	}
	checkedCode, err := l.svcCtx.AuthSvc.CheckPhoneCode(l.ctx, phoneNumber, in.Code, modelAuth.CodeTypeRegister)
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.UserSvc.GetUserIdByPhoneOrEmail(l.ctx, phoneNumber)
	if err != nil && !errors.Is(err, errorx.ERR_USER_NOT_EXISTED) {
		return nil, err
	}

	if user != nil {
		return nil, errorx.ERR_DUPLICATE_REGISTERATION
	}
	//软删除验证码
	if err := l.svcCtx.AuthSvc.ConsumeAuthCode(l.ctx, checkedCode); err != nil {
		return nil, err
	}

	user, err = l.svcCtx.UserSvc.CreateUserInfo(l.ctx, phoneNumber, countryCode, isPhone)
	if err != nil {
		return nil, err
	}

	var t modelAuth.LoginType
	if !isPhone {
		t = modelAuth.EmailLoginType
	} else {
		t = modelAuth.PhoneLoginType
	}

	err = l.svcCtx.AuthSvc.Login(l.ctx, user.Id, modelAuth.ActionTypeRegister, t, user.CountryCode)
	if err != nil {
		return nil, err
	}

	pair, err := l.svcCtx.AuthSvc.CreateSession(l.ctx, user.Id, platform)
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
