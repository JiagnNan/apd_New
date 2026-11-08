package service

import (
	model "github.com/JiagnNan/apd_New/server/model/auth"
	"golang.org/x/net/context"
)

type AuthSvc struct {
	userDao AuthCodeRepository
}

func NewAuthSvc(userDao AuthCodeRepository) *AuthSvc {
	return &AuthSvc{userDao: userDao}
}

func (s *AuthSvc) GetAuthCode(ctx context.Context, phoneNumber string, codeType uint8) (*model.AuthCode, error) {
	//
	return nil, nil
}

func (s *AuthSvc) SaveAuthCode(ctx context.Context, authCode *model.AuthCode) error {
	return nil
}
func (s *AuthSvc) DeleteAuthCode(ctx context.Context, phoneNumber string, codeType uint8) error {
	return nil
}
