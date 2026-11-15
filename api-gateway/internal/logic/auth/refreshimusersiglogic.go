package auth

import (
	"context"
	"strings"

	"github.com/JiagnNan/apd_New/api-gateway/internal/svc"
	"github.com/JiagnNan/apd_New/api-gateway/internal/types"
	"github.com/JiagnNan/apd_New/server/rpc/auth/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type RefreshImUserSigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshImUserSigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshImUserSigLogic {
	return &RefreshImUserSigLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *RefreshImUserSigLogic) RefreshImUserSig(authorization string) (*types.ImUserSigResp, error) {
	parts := strings.Fields(authorization)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, status.Error(codes.Unauthenticated, "access token required")
	}
	if l.svcCtx.IMAuth == nil {
		return nil, status.Error(codes.Unavailable, "IM is unavailable")
	}
	out, _ := metadata.FromOutgoingContext(l.ctx)
	if out == nil {
		out = metadata.MD{}
	}
	out.Set("authorization", "Bearer "+parts[1])
	ctx := metadata.NewOutgoingContext(l.ctx, out)
	resp, err := l.svcCtx.IMAuth.RefreshImUserSig(ctx, &authpb.ImUserSigReq{})
	if err != nil {
		return nil, err
	}
	return &types.ImUserSigResp{
		ImEnabled:         true,
		ImSdkAppId:        resp.SdkAppId,
		ImUserId:          resp.UserId,
		ImUserSig:         resp.UserSig,
		ImUserSigExpireAt: resp.ExpiresAt,
	}, nil
}
