package service

import (
	"github.com/JiagnNan/apd_New/commonlib/errorx"
	"github.com/JiagnNan/apd_New/commonlib/pkg/log"
	utils "github.com/JiagnNan/apd_New/commonlib/utils/phone"
	model "github.com/JiagnNan/apd_New/server/model/auth"
	"golang.org/x/net/context"
	"time"
)

const (
	CodeTypeLogin    uint8 = 1
	CodeTypeRegister uint8 = 2
)

type AuthSvc struct {
	userDao AuthCodeRepository
}

func NewAuthSvc(userDao AuthCodeRepository) *AuthSvc {
	return &AuthSvc{userDao: userDao}
}

func (s *AuthSvc) CheckPhoneCode(ctx context.Context, phoneNumber, phoneCode string, codeType uint8) (*model.AuthCode, error) {
	number, err := utils.NormalizeAccount(phoneNumber)
	if err != nil {
		return nil, err
	}
	code, err := s.userDao.GetAuthCode(ctx, number, codeType)
	if err != nil {
		log.Warnx(ctx, "GetAuthCode err:%v", err)
		return nil, errorx.ERR_PHONE_CODE_INVALID
	}
	if code.ExpiresTime <= time.Now().Unix() || code.Deleted == 1 {
		return nil, errorx.ERR_PHONE_CODE_INVALID
	}
	if code.Code != phoneCode {
		return nil, errorx.ERR_PHONE_CODE_INVALID
	}

	return code, nil
}
