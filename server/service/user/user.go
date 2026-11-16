package user

import (
	"github.com/JiagnNan/apd_New/commonlib/dbx"
	"github.com/JiagnNan/apd_New/commonlib/errorx"
	model "github.com/JiagnNan/apd_New/server/model/user"
	"github.com/JiagnNan/apd_New/server/service/base"
	"golang.org/x/net/context"
	"time"
)

type UserSvc struct {
	UserDao UserRepository
	BaseSvc *base.BaseSvc
}

func NewUserSvc(userRepository UserRepository, baseSvc *base.BaseSvc) *UserSvc {
	return &UserSvc{
		UserDao: userRepository,
		BaseSvc: baseSvc,
	}
}

func (svc *UserSvc) GetUserIdByPhoneOrEmail(ctx context.Context, phoneNumber string) (*model.UserInfo, error) {
	user, err := svc.UserDao.GetUserIdByPhoneOrEmail(ctx, phoneNumber)
	if err != nil {
		if dbx.IsNotFound(err) {
			return nil, errorx.ERR_USER_NOT_EXISTED
		}
		return nil, errorx.ERR_SERVER_INTERNAL
	}
	if user.IsLogoff != 0 || user.BlackTime != 0 {
		return nil, errorx.ERR_USER_STATUS_ABNORMAL
	}
	return user, nil
}
func (svc *UserSvc) CreateUserInfo(ctx context.Context, phoneNumber, countryCode string, isPhone bool) (*model.UserInfo, error) {
	//默认初始化字段
	//随机code,自增id
	id, _ := svc.BaseSvc.NextUserId(ctx)
	code, _ := svc.BaseSvc.UserCode(ctx)
	user := &model.UserInfo{
		Id:          id,
		Code:        code,
		Birthday:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		CountryCode: countryCode,
	}
	registerInfo := &model.RegisterInfo{}
	if isPhone {
		registerInfo.Phone = phoneNumber
		user.Phone = phoneNumber
	} else {
		user.Email = phoneNumber
		registerInfo.Email = phoneNumber
	}
	user.Nickname = svc.BaseSvc.MakeRandomUsername(user.Id)
	if err := svc.UserDao.CreateRegistration(ctx, user, registerInfo); err != nil {
		return nil, err
	}
	return user, nil
}
