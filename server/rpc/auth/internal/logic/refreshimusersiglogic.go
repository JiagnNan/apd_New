package logic

import (
	"context"
	"crypto/subtle"
	"errors"
	"strings"

	"github.com/JiagnNan/apd_New/commonlib/jwt"
	"github.com/JiagnNan/apd_New/commonlib/redisx"
	usermodel "github.com/JiagnNan/apd_New/server/model/user"
	"github.com/JiagnNan/apd_New/server/rpc/auth/authpb"
	"github.com/JiagnNan/apd_New/server/rpc/auth/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type RefreshImUserSigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshImUserSigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshImUserSigLogic {
	return &RefreshImUserSigLogic{ctx: ctx, svcCtx: svcCtx}
}

// RefreshImUserSig issues credentials for the current business session.
func (s *RefreshImUserSigLogic) RefreshImUserSig(_ *authpb.ImUserSigReq) (*authpb.ImUserSigResp, error) {
	ctx := s.ctx
	values := metadata.ValueFromIncomingContext(ctx, "authorization")
	if len(values) != 1 {
		return nil, status.Error(codes.Unauthenticated, "access token required")
	}
	parts := strings.Fields(values[0])
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, status.Error(codes.Unauthenticated, "access token required")
	}
	raw := parts[1]
	claims, err := s.svcCtx.IMTokens.ParseAccessToken(raw)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	if !s.svcCtx.IM.Enabled() {
		logx.WithContext(ctx).Errorf("refresh IM UserSig: stage=availability user_id=%d reason=IM disabled or initialization failed", claims.UserID)
		return nil, status.Error(codes.Unavailable, "IM is unavailable")
	}
	session, err := s.svcCtx.IMSessions.GetSession(ctx, claims.UserID, claims.Platform)
	if redisx.IsNil(err) {
		return nil, status.Error(codes.Unauthenticated, "session expired")
	}
	if err != nil {
		logx.WithContext(ctx).Errorf("refresh IM UserSig: stage=session_lookup user_id=%d err=%v", claims.UserID, err)
		return nil, status.Error(codes.Unavailable, "session check unavailable")
	}
	if session.SessionId != claims.SessionID ||
		subtle.ConstantTimeCompare([]byte(session.AccessHash), []byte(jwt.TokenHash(raw))) != 1 {
		return nil, status.Error(codes.Unauthenticated, "invalid session")
	}
	var user usermodel.UserInfo
	err = s.svcCtx.DB.WithContext(ctx).Select("id", "is_logoff", "black_time").
		Where("id = ?", claims.UserID).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.Unauthenticated, "invalid account")
	}
	if err != nil {
		logx.WithContext(ctx).Errorf("refresh IM UserSig: stage=account_lookup user_id=%d err=%v", claims.UserID, err)
		return nil, status.Error(codes.Unavailable, "account check unavailable")
	}
	if user.IsLogoff != 0 || user.BlackTime != 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid account session")
	}
	info, err := s.svcCtx.IM.IssueUserSig(ctx, user.Id)
	if err != nil {
		// IssueUserSig returns local signing errors without including credentials.
		logx.WithContext(ctx).Errorf("refresh IM UserSig: stage=sign user_id=%d err=%v", user.Id, err)
		return nil, status.Error(codes.Unavailable, "IM credentials unavailable")
	}
	return &authpb.ImUserSigResp{
		SdkAppId:  info.SDKAppID,
		UserId:    info.UserID,
		UserSig:   info.UserSig,
		ExpiresAt: info.ExpiresAt,
	}, nil
}
