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
	if user.IsLogoff == 1 || user.BlackTime != 1 {
		return nil, errorx.ERR_USER_STATUS_ABNORMAL
	}
	return user, nil
}
