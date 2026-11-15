package auth

import (
	"github.com/JiagnNan/apd_New/commonlib/errorx"
	"github.com/JiagnNan/apd_New/commonlib/jwt"
	"github.com/JiagnNan/apd_New/commonlib/mdata"
	"github.com/JiagnNan/apd_New/commonlib/pkg/log"
	utils "github.com/JiagnNan/apd_New/commonlib/utils/phone"
	"github.com/JiagnNan/apd_New/server/dao/auth"
	model "github.com/JiagnNan/apd_New/server/model/auth"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

type AuthSvc struct {
	userDao    AuthCodeRepository
	sessionDao *auth.SessionDao
	jManager   *jwt.Manager
}

func NewAuthSvc(userDao AuthCodeRepository, sessionDao *auth.SessionDao, jManager *jwt.Manager) *AuthSvc {
	return &AuthSvc{userDao: userDao, sessionDao: sessionDao, jManager: jManager}

}

func (s *AuthSvc) CheckPhoneCode(ctx context.Context, phoneNumber, phoneCode string, codeType model.CodeType) (*model.AuthCode, error) {
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
func (s *AuthSvc) DeletedAuthCode(ctx context.Context, phoneNumber string, codeType model.CodeType) error {
	err := s.userDao.DeleteAuthCode(ctx, phoneNumber, codeType)
	if err != nil {
		return errorx.ERR_DELETED_PHONE_CODE_FAILED
	}
	return nil
}
func (s *AuthSvc) ConsumeAuthCode(ctx context.Context, code *model.AuthCode) error {
	if ok := s.userDao.ConsumeAuthCode(ctx, code); ok {
		return nil
	}
	return errorx.ERR_PHONE_CODE_INVALID
}

func (s *AuthSvc) Login(ctx context.Context, userId int64, actionType model.ActionType, loginType model.LoginType) error {
	var v *model.AuthLogin

	md := mdata.FromContext(ctx)

	v = &model.AuthLogin{
		UserId:        userId,
		Ip:            md.ClientIp,
		DeviceId:      md.DeviceId,
		DeviceModel:   md.DeviceModel,
		SystemVersion: md.SystemVersion,
		AppVersion:    md.AppVersion,
		Platform:      md.Platform,
		Language:      "en",
		LoginType:     uint8(loginType),
		ActionType:    uint8(actionType),
		ActionTime:    time.Now().Unix(),
		CountryCode:   strconv.FormatInt(md.CountryCode, 10),
		Deleted:       0,
	}
	err := s.userDao.Login(ctx, v)
	if err != nil {
		return errorx.ERR_CREATE_LOGIN_ERROR
	}
	return nil
}
