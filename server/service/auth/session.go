package auth

import (
	"github.com/JiagnNan/apd_New/commonlib/errorx"
	"github.com/JiagnNan/apd_New/commonlib/jwt"
	authDao "github.com/JiagnNan/apd_New/server/dao/auth"
	"golang.org/x/net/context"
)

func (s *AuthSvc) CreateSession(ctx context.Context, userId int64, platform, authVersion string) (*jwt.TokenPair, error) {
	sessionId, err := jwt.NewSessionID()
	if err != nil {
		return nil, errorx.ERR_SERVER_INTERNAL
	}
	pair, err := s.jManager.GenerateTokenPair(userId, sessionId, platform, authVersion)
	if err != nil {
		return nil, err
	}
	session := authDao.Session{
		SessionId:       sessionId,
		AccessHash:      jwt.TokenHash(pair.AccessToken),
		RefreshHash:     jwt.TokenHash(pair.RefreshToken),
		AuthVersion:     authVersion,
		RefreshExpireAt: pair.RefreshExpireAt.Unix(),
	}
	err = s.sessionDao.SaveSession(ctx, userId, platform, session)
	if err != nil {
		return nil, err
	}
	return &pair, nil
}
