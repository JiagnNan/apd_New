package user

import (
	"github.com/JiagnNan/apd_New/commonlib/dbx"
	"github.com/JiagnNan/apd_New/commonlib/errorx"
	model "github.com/JiagnNan/apd_New/server/model/user"
	"golang.org/x/net/context"
)

type UserSvc struct {
	UserDao UserRepository
}

func NewUserSvc(userRepository UserRepository) *UserSvc {
	return &UserSvc{
		UserDao: userRepository,
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
	// TODO: 写入用户记录时保存 countryCode。
	return nil, errorx.ERR_SERVER_INTERNAL
}
