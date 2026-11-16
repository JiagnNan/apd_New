package user

import (
	model "github.com/JiagnNan/apd_New/server/model/user"
	"golang.org/x/net/context"
)

type UserRepository interface {
	GetUserIdByPhoneOrEmail(ctx context.Context, phoneNumber string) (*model.UserInfo, error)
	CreateUserInfo(ctx context.Context, v *model.UserInfo) error
	UpdateUserNickName(ctx context.Context, userId int64, nickname string) error
	InsertRegisterInfo(ctx context.Context, v *model.RegisterInfo) error
}
