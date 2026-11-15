package auth

import (
	model "github.com/JiagnNan/apd_New/server/model/auth"
	"golang.org/x/net/context"
)

type AuthCodeRepository interface {
	GetAuthCode(ctx context.Context, phoneNumber string, codeType model.CodeType) (*model.AuthCode, error)
	SaveAuthCode(ctx context.Context, authCode *model.AuthCode) error
	DeleteAuthCode(ctx context.Context, phoneNumber string, codeType model.CodeType) error
	Login(ctx context.Context, authLogin *model.AuthLogin) error
}
